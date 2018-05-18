/*

Copyright (C) 2017-2018  Daniele Rondina <geaaru@sabayonlinux.org>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.

*/

package commons

import (
	"archive/tar"
	"compress/bzip2"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	logger "github.com/sirupsen/logrus"
	viper "github.com/spf13/viper"
)

type CheckerExecutor interface {
	AddPackage(p *Package) error
	Run() error
	GetPackages() []Package
}

type Checker struct {
	settings     *viper.Viper
	logger       *logger.Logger
	packages     []Package
	mutex        sync.Mutex
	elabPackages func(pkgs []string) error
}

// I use anonymous field to override
// Run and other methods
type CheckerConcurrent struct {
	*Checker
}

func NewChecker(settings *viper.Viper, l *logger.Logger) (*Checker, error) {

	var log *logger.Logger = nil
	if settings == nil {
		return nil, errors.New("Invalid settings param")
	}
	if l == nil {
		// Use standard logger
		log = logger.StandardLogger()
	} else {
		log = l
	}

	logger.Debug("Created new Checker object")

	var c = &Checker{
		settings: settings,
		logger:   log,
		packages: []Package{},
	}

	c.elabPackages = c.processPackages

	return c, nil
}

func NewCheckerConcurrent(settings *viper.Viper, l *logger.Logger) (*CheckerConcurrent, error) {

	var c, err = NewChecker(settings, l)
	if err != nil {
		return nil, err
	}

	var cc = &CheckerConcurrent{c}
	cc.elabPackages = cc.processPackages

	return cc, nil
}

func checkPackage(pkg string) (string, error) {

	var abspath string
	var err error

	abspath, err = filepath.Abs(pkg)
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(abspath); os.IsNotExist(err) {
		return "", err
	}

	return abspath, nil
}

func (c *Checker) is2SkipFile(pkg string, file string) (bool, error) {
	var ans bool = false

	if file == "" {
		return ans, errors.New("Invalid file")
	}

	// Check if file has an extention to skip
	if len(c.settings.GetStringSlice("ignoreExt")) > 0 {

		var ext string = filepath.Ext(file)

		for _, e := range c.settings.GetStringSlice("ignoreExt") {

			if !strings.HasPrefix(e, ".") {
				e = "." + e
			}
			if strings.Compare(e, ext) == 0 {
				ans = true
				goto ret
			}

		}

	}

	// Check if file must be ignored
	if len(c.settings.GetStringSlice("ignoreFiles")) > 0 {

		var base string = filepath.Base(file)

		for _, f := range c.settings.GetStringSlice("ignoreFiles") {

			if strings.Count(f, "/") > 1 {
				// POST: it's a path not a single file. I try to match it.
				if file == f {
					ans = true
					goto ret
				}

			} else if base == filepath.Base(f) {
				ans = true
				goto ret
			}

		}

	}

ret:

	return ans, nil
}

func (c *Checker) processTarBz2(pkg string, abs string) error {

	var p *Package
	var tarbz2 io.Reader
	var f, err = os.Open(abs)
	if err != nil {
		return err
	}
	defer f.Close()

	tarbz2 = bzip2.NewReader(f)

	// Create Package object
	p, err = NewPackage(pkg, c.logger)
	if err != nil {
		return err
	}
	p.abspath = abs
	p.basename = filepath.Base(pkg)

	var tarReader = tar.NewReader(tarbz2)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			err = nil
			break
		}

		if err != nil {
			return err
		}

		//var fi = header.FileInfo()

		var isDir = false
		var toSkip = false

		if header.Typeflag == tar.TypeDir {
			isDir = true
			p.AddDir(header.Name)
		} else {
			toSkip, err = c.is2SkipFile(pkg, header.Name)

			if toSkip == false && err == nil {
				err = p.ProcessTarFile(tarReader, header.Name)
			} else if toSkip {
				p.skipped++
			}
		}
		if err != nil {
			return err
		}

		c.logger.Debugf("[%s] File %s (dir = %t, skip = %t).", pkg, header.Name,
			isDir, toSkip)

	}

	err = p.CalculateCRC()

	c.logger.Infof("[%s] %s", pkg, p)

	c.AddPackage(p)

	return err
}

func (c *Checker) AddPackage(p *Package) error {

	if p == nil {
		return errors.New("Invalid package")
	}

	c.mutex.Lock()

	c.packages = append(c.packages, *p)

	c.mutex.Unlock()

	return nil
}

func (c *Checker) processPackage(pkg string) error {
	var err error = nil

	c.logger.Debugf("[%s] Checking package...", filepath.Base(pkg))

	var absp, extension, pkgname string
	absp, err = checkPackage(pkg)
	if err != nil {
		c.logger.Errorf("[%s] Error: %s", filepath.Base(pkg), err)
		return err
	}
	pkgname = filepath.Base(filepath.Dir(absp)) + "/" + filepath.Base(pkg)

	c.logger.Debugf("[%s] File %s checking OK.", pkgname, absp)

	extension = filepath.Ext(absp)

	if extension == ".tbz2" || strings.HasSuffix(filepath.Base(pkg), ".tar.bz2") {
		err = c.processTarBz2(pkgname, absp)
		if err != nil {
			c.logger.Errorf("[%s] Error: %s", pkgname, err)
			return err
		}
	} else {
		c.logger.Errorf("[%s] File with extension %s not supported.", pkgname)
		err = errors.New("Extension not supported")
	}

	return err
}

func (c *Checker) processPackages(pkgs []string) error {
	var err error
	var okCounter, n_pkgs int

	okCounter = 0
	n_pkgs = len(pkgs)

	for _, pkg := range pkgs {
		err = c.processPackage(pkg)
		if err == nil {
			okCounter++
		}
	}

	c.logger.Infof("For %d packages: %d OK, %d KO.",
		n_pkgs, okCounter, n_pkgs-okCounter)

	return err
}

func (c *Checker) findPackages(dir string) ([]string, error) {

	var err error
	var ans []string
	var files []os.FileInfo

	files, err = ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		c.logger.Debugf("For dir %s found child %s",
			dir, f.Name())
		if f.IsDir() {
			var child []string

			child, err = c.findPackages(fmt.Sprintf("%s/%s", dir, f.Name()))
			if err != nil {
				return nil, err
			}
			for _, childFile := range child {
				ans = append(ans, childFile)
			}
		} else {
			var ext string = filepath.Ext(f.Name())
			if ext == ".tbz2" {
				ans = append(ans, fmt.Sprintf("%s/%s", dir, f.Name()))
			}
		}
	}

	return ans, err
}

func (c *Checker) processDirectory(dir string) error {
	var err error
	var files []string

	files, err = c.findPackages(dir)
	if err != nil {
		return err
	}
	err = c.elabPackages(files)

	return err
}

func (c *Checker) sortPackages() {
	c.mutex.Lock()

	sort.Sort(PackageSorter(c.packages))

	c.mutex.Unlock()
}

func (c *Checker) Run() error {

	var err error

	// Elaborate list of packages if present
	if len(c.settings.GetStringSlice("package")) > 0 {
		err = c.elabPackages(c.settings.GetStringSlice("package"))
		if err != nil {
			return err
		}
	}

	// Elaborate .tbz2 file under directory
	if c.settings.GetString("directory") != "" {
		err = c.processDirectory(c.settings.GetString("directory"))
	}

	// TODO: process stdin

	// Sort package list
	c.sortPackages()

	return err
}

// Override processPackages for use gorouting for CheckerConcurrent struct
func (c *CheckerConcurrent) processPackages(pkgs []string) error {
	var err error
	var i int
	var ch chan ChannelResp = make(chan ChannelResp, c.settings.GetInt("maxconcurrency"))
	var okCounter, n_pkgs int

	okCounter = 0
	n_pkgs = len(pkgs)

	for _, pkg := range pkgs {
		c.logger.Debugf("Starting gorouting for package %s...\n", pkg)
		go c.go2ProcessPackage(ch, pkg)
	}

	for i = 0; i < n_pkgs; i++ {
		resp := <-ch
		if resp.Error == nil {
			c.logger.Debugf("[%s] Received response: OK\n", resp.Result)
			okCounter++
		} else {
			c.logger.Debugf("[%s] Received response: KO\n", resp.Result)
		}
	}

	c.logger.Infof("For %d packages: %d OK, %d KO.",
		n_pkgs, okCounter, n_pkgs-okCounter)

	if okCounter != n_pkgs {
		err = errors.New("Something goes wrong")
	}

	return err
}

func (c *CheckerConcurrent) go2ProcessPackage(channel chan ChannelResp,
	pkg string) {
	var err error
	c.logger.Debugf("[%s] Starting goroutine", pkg)
	err = c.processPackage(pkg)
	channel <- NewChannelResp(pkg, err)
	c.logger.Debugf("[%s] End goroutine", pkg)
}

func (c *Checker) GetPackages() []Package {
	return c.packages
}
