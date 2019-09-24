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

package binhostdir

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"

	logger "github.com/sirupsen/logrus"
)

func ProcessCategoryDir(dir string, log *logger.Logger, tree *map[string][]string) error {
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
		log.WithFields(logger.Fields{
			"file":     file.Name(),
			"category": cat,
		}).Debugf("Processing file...")

		if file.IsDir() {
			continue
		}

		// Check only directory of categories.
		if !regexCat.MatchString(file.Name()) {
			log.WithFields(logger.Fields{
				"file":     file.Name(),
				"category": cat,
			}).Debugf("File skipped.")
			continue
		}

		pkgFiles = append(pkgFiles, path.Join(dir, file.Name()))
	}

	// Write to handle in mutual exclusion
	if len(pkgFiles) > 0 {
		(*tree)[cat] = pkgFiles
	}

	log.WithFields(logger.Fields{
		"category": cat,
		"files":    len(pkgFiles),
	}).Debugf("Complete navigation of directory.")

	return nil
}

// Parse binhost Directory
func AnalyzeBinHostDirectory(binhostDir string, log *logger.Logger, tree *map[string][]string) error {
	var err error
	var files []os.FileInfo
	var categoryDirs []string = make([]string, 0)

	files, err = ioutil.ReadDir(binhostDir)
	if err != nil {
		return errors.New(fmt.Sprintf("Error on read directory %s: %s",
			binhostDir, err.Error()))
	}

	var regexCat = regexp.MustCompile(`^[a-z]+[-][a-z]+$`)
	for _, file := range files {
		log.WithFields(logger.Fields{
			"file": file.Name(),
		}).Debugf("Processing file...")

		if !file.IsDir() {
			continue
		}

		// Check only directory of categories.
		if !regexCat.MatchString(file.Name()) {
			log.WithFields(logger.Fields{
				"file": file.Name(),
			}).Debugf("Is not a category directory.")
			continue
		}

		categoryDirs = append(categoryDirs, path.Join(binhostDir, file.Name()))
	}

	if len(categoryDirs) == 0 {
		log.Infoln("No directory of categories found. Nothing to filter.")
		return nil
	}

	// TODO: handle this with concurrency
	for _, cat := range categoryDirs {
		_ = ProcessCategoryDir(cat, log, tree)
	}

	return nil
}
