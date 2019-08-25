/*

Copyright (C) 2017-2019  Daniele Rondina <geaaru@sabayonlinux.org>

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
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"

	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Filter struct {
	settings *viper.Viper
	logger   *logger.Logger
	config   *SarkConfig
	Tree     map[string][]string
}

func NewFilter(settings *viper.Viper, l *logger.Logger, config *SarkConfig) (*Filter, error) {
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

	logger.Debug("Created new Filter object")

	var ans = &Filter{
		settings: settings,
		logger:   log,
		config:   config,
		Tree:     make(map[string][]string, 0),
	}

	return ans, nil
}

func (f *Filter) Run(binhostDir string) error {

	var files []os.FileInfo
	var categoryDirs []string = make([]string, 0)
	var err error

	if binhostDir == "" {
		return errors.New("Invalid binhost directory")
	}

	files, err = ioutil.ReadDir(binhostDir)
	if err != nil {
		return errors.New(fmt.Sprintf("Error on read directory %s: %s",
			binhostDir, err.Error()))
	}

	var regexCat = regexp.MustCompile(`^[a-z]+[-][a-z]+$`)
	for _, file := range files {
		f.logger.WithFields(logger.Fields{
			"file": file.Name(),
		}).Debugf("Processing file...")

		if !file.IsDir() {
			continue
		}

		// Check only directory of categories.
		if !regexCat.MatchString(file.Name()) {
			f.logger.WithFields(logger.Fields{
				"file": file.Name(),
			}).Debugf("Is not a category directory.")
			continue
		}

		categoryDirs = append(categoryDirs, path.Join(binhostDir, file.Name()))
	}

	if len(categoryDirs) == 0 {
		f.logger.Infoln("No directory of categories found. Nothing to filter.")
		return nil
	}

	fmt.Println(len(categoryDirs))

	for _, cat := range categoryDirs {
		_ = f.processCategoryDir(cat)
	}

	fmt.Printf("MAP %s\n", f.Tree)

	return nil
}

func (f *Filter) processCategoryDir(dir string) error {
	var files []os.FileInfo
	var pkgFiles []string = make([]string, 0)
	var err error
	cat := path.Base(dir)

	files, err = ioutil.ReadDir(dir)
	if err != nil {
		return errors.New(fmt.Sprintf("Error on read directory %s: %s",
			dir, err.Error()))
	}

	var regexCat = regexp.MustCompile(`.tbz2$`)
	for _, file := range files {
		f.logger.WithFields(logger.Fields{
			"file":     file.Name(),
			"category": cat,
		}).Debugf("Processing file...")

		if file.IsDir() {
			continue
		}

		// Check only directory of categories.
		if !regexCat.MatchString(file.Name()) {
			f.logger.WithFields(logger.Fields{
				"file":     file.Name(),
				"category": cat,
			}).Debugf("File skipped.")
			continue
		}

		pkgFiles = append(pkgFiles, path.Join(dir, file.Name()))
	}

	// Write to handle in mutual exclusion
	f.Tree[cat] = pkgFiles

	logger.WithFields(logger.Fields{
		"category": cat,
		"files":    len(pkgFiles),
	}).Debugf("Complete navigation of directory.")

	return nil
}
