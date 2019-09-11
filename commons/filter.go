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
	"strings"

	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Filter struct {
	settings    *viper.Viper
	logger      *logger.Logger
	Config      *SarkConfig
	BinHostTree map[string][]string
	RulesTree   *FilterMatrix
}

type FilterMatrix struct {
	FilterType string
	Father     *Filter
	Branches   map[string]FilterMatrixBranch
	Resources  []*FilterResource
}

type FilterMatrixBranch struct {
	Category         string
	CategoryFiltered bool
	Resources        []*FilterResource
	Matches          map[string]FilterMatrixLeaf
	NotMatches       map[string]FilterMatrixLeaf
}

type FilterMatrixLeaf struct {
	Name     string
	Path     string
	Father   *FilterMatrixBranch
	Resource *FilterResource
}

type FilterResource struct {
	Source     string
	Type       string
	Packages   []string
	Categories []string
}

func NewFilterResource(source string, rtype string, pkgs []string, categories []string) (*FilterResource, error) {
	var ans *FilterResource
	if source == "" {
		return nil, errors.New("Invalid source")
	}
	if rtype == "" {
		return nil, errors.New("Invalid resource type")
	}
	ans = &FilterResource{
		Source:     source,
		Type:       rtype,
		Packages:   make([]string, 0),
		Categories: make([]string, 0),
	}

	if len(pkgs) > 0 {
		ans.Packages = pkgs
	}
	if len(categories) > 0 {
		ans.Categories = categories
	}

	return ans, nil
}

func (r *FilterResource) AddCategory(category string) {
	r.Categories = append(r.Categories, category)
}

func (r *FilterResource) AddPackage(pkg string) {
	r.Packages = append(r.Packages, pkg)
}

func NewFilterMatrix(ftype string) (*FilterMatrix, error) {
	if ftype == "" {
		return nil, errors.New("Invalid filter type")
	}
	return &FilterMatrix{
		FilterType: ftype,
		Resources:  make([]*FilterResource, 0),
	}, nil
}

func (m *FilterMatrix) AddResource(r *FilterResource) error {
	if (*r).Source == "" {
		return errors.New("AddResource: Invalid source")
	}
	if (*r).Type == "" {
		return errors.New("AddResource: Invalid type")
	}
	m.Resources = append(m.Resources, r)
	return nil
}

func (m *FilterMatrix) GetResourceFilterBySource(source string) (*FilterResource, error) {
	if source == "" {
		return nil, errors.New("Invalid source")
	}
	for _, r := range m.Resources {
		if (*r).Source == source {
			return r, nil
		}
	}
	return nil, nil
}

func (m *FilterMatrix) processSarkBuildFile(conf *SarkConfig, level int, fromFile bool) error {
	if conf != nil {
		if r, _ := m.GetResourceFilterBySource(conf.Id); r == nil {
			if len(conf.Build.TargetPkgs) > 0 {
				br, _ := NewFilterResource(conf.Id, "buildfile", conf.Build.TargetPkgs, nil)
				m.AddResource(br)
				level++
				if fromFile {
					// For now handle this only for local sark file
					if len(conf.Injector.Filter.Rules) > 0 {
						for _, srule := range conf.Injector.Filter.Rules {
							err := m.LoadInjectRule(br, &srule, level)
							if err != nil {
								return errors.New("Error on load injection rule: " + err.Error())
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func (m *FilterMatrix) LoadInjectRule(r *FilterResource, rule *SarkFilterRuleConf, level int) error {
	if r == nil {
		return errors.New("LoadInjectRule: Invalid FilterResource")
	}
	if rule == nil {
		return errors.New("LoadInjectRule: Invalid rule")
	}

	if level >= 3 {
		// Avoid infinite loop
		return nil
	}

	if len((*rule).Categories) > 0 {
		for _, cat := range (*rule).Categories {
			(*r).AddCategory(cat)
		}
	}

	if len((*rule).Packages) > 0 {
		for _, p := range (*rule).Packages {
			(*r).AddPackage(p)
		}
	}

	// Handle Files Rules
	if len(rule.Files) > 0 {
		for _, f := range rule.Files {
			absfile, err := AbsPathFromBase((*r).Source, f)
			if err != nil {
				return errors.New(
					fmt.Sprintf("LoadInjectRule: Error on retrieve abs path for file %s: %s",
						f, err.Error()))
			}
			conf, err := NewSarkConfigFromFile(m.Father.settings, absfile)
			if err != nil {
				return errors.New(
					fmt.Sprintf("LoadInjectRule: Error on load file %s: %s",
						f, err.Error()))
			}

			err = m.processSarkBuildFile(conf, level, true)
			if err != nil {
				return errors.New(
					fmt.Sprintf("LoadInjectRule: Error on parse file %s: %s",
						f, err.Error()))
			}
		}
	}

	// Handle URL
	if len(rule.Urls) > 0 {
		opts := NewHttpClientDefaultOpts()
		if m.Father.settings.GetBool("insecure_skipverify") {
			opts.InsecureSkipVerify = true
		}
		apiKey := m.Father.settings.GetString("apikey")
		for _, u := range rule.Urls {

			if r, _ := m.GetResourceFilterBySource(u); r != nil {
				m.Father.logger.Warnf(fmt.Sprintf("Url %s duplicated.", u))
				continue
			}

			if !strings.HasPrefix(u, "buildfile") &&
				!strings.HasPrefix(u, "pkglist|") {
				return errors.New(
					fmt.Sprintf("LoadInjectRule: Invalid rule %s", u))
			}

			if strings.HasPrefix(u, "buildfile|") {
				resp, err := GetResource(u[10:], apiKey, opts)
				if err != nil {
					return errors.New(fmt.Sprintf("Error on fetch url %s: %s", u, err))
				}

				var remoteBuildfile *SarkConfig
				remoteBuildfile, err = NewSarkConfigFromBytes(m.Father.settings, resp)
				if err != nil {
					return errors.New(
						fmt.Sprintf("LoadInjectRule: Error on parse data from url %s: %s",
							u, err))
				}
				remoteBuildfile.Id = u
				err = m.processSarkBuildFile(remoteBuildfile, level, false)

			} else {
				resp, err := GetResource(u[9:], apiKey, opts)
				if err != nil {
					return errors.New(
						fmt.Sprintf("Error on fetch url %s: %s", u, err))
				}

				var pkgs []string
				pkgs, err = PkgListParser(resp)
				if len(pkgs) > 0 {
					br, _ := NewFilterResource(u, "pkglist", pkgs, nil)
					m.AddResource(br)
				}

			}

		}
	}

	return nil
}

func (m *FilterMatrix) LoadInjectRules(source, rtype string, rules []SarkFilterRuleConf) error {
	// NOTE: currently is not supported the inclusion of remote injection rules.

	r, err := m.GetResourceFilterBySource(source)
	if err != nil {
		return errors.New("Error on check for existing resource " + err.Error())
	}
	if r != nil {
		m.Father.logger.Warnf(
			"Resource %s already loaded. Skip rules to avoid circular injection.",
			source)
		return nil
	}

	r, err = NewFilterResource(source, rtype, nil, nil)

	for _, rule := range rules {
		err = m.LoadInjectRule(r, &rule, 1)
		if err != nil {
			return errors.New("Error on load injection rule: " + err.Error())
		}
	}

	return nil
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
		settings:    settings,
		logger:      log,
		Config:      config,
		BinHostTree: make(map[string][]string, 0),
	}

	return ans, nil
}

// Parse binhost Directory
func (f *Filter) analyzeBinHostDirectory(binhostDir string) error {
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

	// TODO: handle this with concurrency
	for _, cat := range categoryDirs {
		_ = f.processCategoryDir(cat)
	}

	return nil
}

func (f *Filter) Run(binhostDir string) error {
	var err error

	if binhostDir == "" {
		return errors.New("Invalid binhost directory")
	}

	// Phase1: Analyze binhost Directory
	err = f.analyzeBinHostDirectory(binhostDir)
	if err != nil {
		return err
	}

	if len(f.BinHostTree) > 0 {
		// Phase2: Create FilterMatrix
		err = f.createFilterMatrix()

	} else {
		f.logger.Infof("No files found to filter. Nothing to do.")
	}

	return nil
}

func (f *Filter) createFilterMatrix() error {
	if f.Config == nil {
		// Create an empty Sark Config where injection filter
		// has blacklist as filter type and no packages blocked
		f.Config, _ = NewSarkConfig(f.settings, "blacklist")
		f.Config.Id = "filter"
	} else if f.Config.Injector.Filter.FilterType == "" {
		// If there is a filter section I consider filter of type
		// whitelist where packages are all present on target section.
		f.Config.Injector = *NewSarkInjectConfig("whitelist")
		f.Config.Id = "filter"
	}

	f.RulesTree, _ = NewFilterMatrix(f.Config.Injector.Filter.FilterType)
	f.RulesTree.Father = f

	if len(f.Config.Injector.Filter.Rules) > 0 {
		// POST: Inject rules available. Load FilterResources
		err := f.RulesTree.LoadInjectRules(f.Config.Id, "buildfile", f.Config.Injector.Filter.Rules)
		if err != nil {
			return err
		}

	} else {
		if len(f.Config.Build.TargetPkgs) > 0 {
			// POST: No injection rules are present but on input configuration
			// file there are targets that will be used for filter.
			r, _ := NewFilterResource(f.Config.Id, "buildfile", f.Config.Build.TargetPkgs, nil)
			f.RulesTree.AddResource(r)
		} else if f.RulesTree.FilterType == "whitelist" {
			f.logger.Warnf("No packages defined and whitelist used. All packages will be filtered")
		}
	}

	// Elaborate loaded data

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
	if len(pkgFiles) > 0 {
		f.BinHostTree[cat] = pkgFiles
	}

	logger.WithFields(logger.Fields{
		"category": cat,
		"files":    len(pkgFiles),
	}).Debugf("Complete navigation of directory.")

	return nil
}
