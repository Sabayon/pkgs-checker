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
	"io"
	"os"
	"path/filepath"
	"strings"

	logger "github.com/sirupsen/logrus"
	viper "github.com/spf13/viper"
)

type Checker struct {
	settings *viper.Viper
	logger   *logger.Logger
	packages []Package
}

func New(settings *viper.Viper, l *logger.Logger) (*Checker, error) {

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

	return &Checker{settings: settings, logger: log, packages: []Package{}}, nil
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
	//var global_md5 = md5.New()
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

	c.packages = append(c.packages, *p)

	return err
}

func (c *Checker) Run() error {

	var err error
	var okCounter, n_pkgs int
	var pkgs = c.settings.GetStringSlice("package")

	okCounter = 0
	n_pkgs = len(c.settings.GetStringSlice("package"))

	for _, pkg := range pkgs {
		c.logger.Infof("[%s] Checking package...", filepath.Base(pkg))

		var absp, extension, pkgname string
		absp, err = checkPackage(pkg)
		if err != nil {
			c.logger.Errorf("[%s] Error: %s", filepath.Base(pkg), err)
			continue
		}
		pkgname = filepath.Base(filepath.Dir(absp)) + "/" + filepath.Base(pkg)

		c.logger.Debugf("[%s] File %s checking OK.", pkgname, absp)

		extension = filepath.Ext(absp)

		if extension == ".tbz2" || strings.HasSuffix(filepath.Base(pkg), ".tar.bz2") {
			err = c.processTarBz2(pkgname, absp)
			if err != nil {
				c.logger.Errorf("[%s] Error: %s", pkgname, err)
				continue
			}
		} else {
			c.logger.Errorf("[%s] File with extension %s not supported.", pkgname)
		}

		okCounter++
	}

	c.logger.Infof("For %d packages: %d OK, %d KO.",
		n_pkgs, okCounter, n_pkgs-okCounter)

	return err
}

func (c *Checker) GetPackages() []Package {
	return c.packages
}
