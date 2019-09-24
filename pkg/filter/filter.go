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

package filter

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	binhostdir "github.com/Sabayon/pkgs-checker/pkg/binhostdir"
	commons "github.com/Sabayon/pkgs-checker/pkg/commons"
	gentoo "github.com/Sabayon/pkgs-checker/pkg/gentoo"
	pkglist "github.com/Sabayon/pkgs-checker/pkg/pkglist"
	sark "github.com/Sabayon/pkgs-checker/pkg/sark"
)

type Filter struct {
	settings    *viper.Viper
	logger      *logger.Logger
	Config      *sark.SarkConfig
	BinHostTree map[string][]string
	RulesTree   *FilterMatrix
}

type FilterMatrix struct {
	FilterType string
	Father     *Filter
	Branches   map[string]*FilterMatrixBranch
	Resources  []*FilterResource
}

type FilterMatrixBranch struct {
	Category         string
	CategoryFiltered bool
	Matrix           *FilterMatrix
	Resources        []*FilterResource
	Packages         []*gentoo.GentooPackage
	// The key of the map contains file path
	Matches map[string]*FilterMatrixLeaf
	// The key of the map contains file path
	NotMatches map[string]*FilterMatrixLeaf
}

type FilterMatrixLeaf struct {
	Name     string
	Path     string
	Package  *gentoo.GentooPackage
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

func NewFilterMatrixBranch(category string) (*FilterMatrixBranch, error) {
	if category == "" {
		return nil, errors.New("Invalid category param")
	}

	return &FilterMatrixBranch{
		Category:         category,
		CategoryFiltered: false,
		Resources:        make([]*FilterResource, 0),
		Packages:         make([]*gentoo.GentooPackage, 0),
		Matches:          make(map[string]*FilterMatrixLeaf, 0),
		NotMatches:       make(map[string]*FilterMatrixLeaf, 0),
	}, nil
}

func (b *FilterMatrixBranch) ContainsResource(resource *FilterResource) (bool, error) {
	var ans bool = false
	if resource == nil {
		return false, errors.New(
			fmt.Sprintf("Invalid resource for branch %s", b.Category))
	}

	for _, r := range b.Resources {
		if r == resource {
			ans = true
			break
		}
	}

	return ans, nil
}

func (b *FilterMatrixBranch) CheckPackages(files []string) error {
	var admitted bool
	var hasPkgRule bool

	for _, f := range files {
		admitted = false
		hasPkgRule = false

		pkgname := filepath.Base(f)
		pkgname = pkgname[:strings.Index(pkgname, filepath.Ext(pkgname))]

		gentooPkg, err := gentoo.ParsePackageStr(
			fmt.Sprintf("%s/%s", b.Category, pkgname))
		if err != nil {
			return err
		}

		// TODO: replace packages with a map
		for _, pkg := range b.Packages {
			if pkg.Name == gentooPkg.Name {
				hasPkgRule = true
				admitted, err = pkg.Admit(gentooPkg)
				if err != nil {
					return err
				}

				if !admitted {
					break
				}
			}
		}

		if !admitted && b.CategoryFiltered && !hasPkgRule {
			admitted = true
		}

		// TODO: store FilterResource
		_, err = b.AddPackage(f, admitted)
		if err != nil {
			return err
		}

	}

	return nil
}

func (b *FilterMatrixBranch) AddResource(resource *FilterResource) error {
	var present bool = false

	if resource == nil {
		return errors.New(
			fmt.Sprintf("Invalid resource to add in branch %s",
				b.Category))
	}

	// TODO: check if convert list of resources in map
	// Check resource is already present
	for _, r := range b.Resources {
		if r == resource {
			present = true
			break
		}
	}

	if !present {
		b.Resources = append(b.Resources, resource)
	}

	return nil
}

func (b *FilterMatrixBranch) AddPackage(file string, match bool) (*FilterMatrixLeaf, error) {
	pkgname := filepath.Base(file)
	pkgname = pkgname[:strings.Index(pkgname, filepath.Ext(pkgname))]

	gentooPkg, err := gentoo.ParsePackageStr(
		fmt.Sprintf("%s/%s", b.Category, pkgname))
	if err != nil {
		return nil, err
	}

	leaf := &FilterMatrixLeaf{
		Name:    gentooPkg.Name,
		Path:    file,
		Package: gentooPkg,
		Father:  b,
	}

	if match {
		b.Matches[file] = leaf
		b.Matrix.Log(logger.DebugLevel, "Branch %s: Add matched package %s (%s)",
			b.Category, gentooPkg, leaf.Path)
	} else {
		b.NotMatches[file] = leaf
		b.Matrix.Log(logger.DebugLevel, "Branch %s: Add not matched package %s (%s)",
			b.Category, gentooPkg, leaf.Path)
	}

	return leaf, nil
}

func NewFilterMatrix(ftype string) (*FilterMatrix, error) {
	if ftype == "" {
		return nil, errors.New("Invalid filter type")
	}
	return &FilterMatrix{
		FilterType: ftype,
		Resources:  make([]*FilterResource, 0),
		Branches:   make(map[string]*FilterMatrixBranch, 0),
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

func (m *FilterMatrix) Log(level logger.Level, msg string, args ...interface{}) {
	if m.Father != nil {
		m.Father.logger.Logf(level, msg, args...)
	} else {
		fmt.Printf("%s: %s\n", level.String(), fmt.Sprintf(msg, args...))
	}
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

func (m *FilterMatrix) processSarkBuildFile(conf *sark.SarkConfig, level int, fromFile bool) error {
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

func (m *FilterMatrix) LoadInjectRule(r *FilterResource, rule *sark.SarkFilterRuleConf, level int) error {
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
			absfile, err := commons.AbsPathFromBase(filepath.Dir((*r).Source), f)
			if err != nil {
				return errors.New(
					fmt.Sprintf("LoadInjectRule: Error on retrieve abs path for file %s: %s",
						f, err.Error()))
			}
			conf, err := sark.NewSarkConfigFromFile(m.Father.settings, absfile)
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
		opts := commons.NewHttpClientDefaultOpts()
		if m.Father.settings.GetBool("insecure_skipverify") {
			opts.InsecureSkipVerify = true
		}
		apiKey := m.Father.settings.GetString("apikey")
		for _, u := range rule.Urls {

			if r, _ := m.GetResourceFilterBySource(u); r != nil {
				m.Log(logger.WarnLevel, "Url %s duplicated.", u)
				continue
			}

			if !strings.HasPrefix(u, "buildfile") &&
				!strings.HasPrefix(u, "pkglist|") {
				return errors.New(
					fmt.Sprintf("LoadInjectRule: Invalid rule %s", u))
			}

			if strings.HasPrefix(u, "buildfile|") {
				remoteBuildfile, err := sark.NewSarkConfigFromResource(
					m.Father.settings,
					u[10:], apiKey, opts)
				if err != nil {
					return errors.New(fmt.Sprintf("Error on load resource url %s: %s", u, err))
				}
				remoteBuildfile.Id = u
				err = m.processSarkBuildFile(remoteBuildfile, level, false)

			} else {
				pkgs, err := pkglist.PkgListLoadResource(u[9:], apiKey, opts)
				if err != nil {
					return errors.New(
						fmt.Sprintf("Error on fetch url %s: %s", u, err))
				}
				if len(pkgs) > 0 {
					br, _ := NewFilterResource(u, "pkglist", pkgs, nil)
					m.AddResource(br)
				}

			}

		}
	}

	return nil
}

func (m *FilterMatrix) LoadInjectRules(source, rtype string, rules []sark.SarkFilterRuleConf) error {
	// NOTE: currently is not supported the inclusion of remote injection rules.

	r, err := m.GetResourceFilterBySource(source)
	if err != nil {
		return errors.New("Error on check for existing resource " + err.Error())
	}
	if r != nil {
		m.Log(logger.WarnLevel,
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

func (m *FilterMatrix) CreateBranches() error {
	for _, r := range m.Resources {
		// Check categories
		for _, category := range r.Categories {

			// Check if exists branch
			branch, present := m.Branches[category]
			if present {
				if !branch.CategoryFiltered {
					branch.CategoryFiltered = true
				}
				m.Log(logger.DebugLevel, "Branch for category %s already present.",
					category)
			} else {
				branch, _ := NewFilterMatrixBranch(category)
				branch.CategoryFiltered = true
				branch.Matrix = m
				branch.Resources = make([]*FilterResource, 0)
				m.Branches[category] = branch
				m.Log(logger.DebugLevel, "Added branch for category %s.", category)
			}
			branch.AddResource(r)
		}

		// Check packages
		for _, pkg := range r.Packages {
			gp, err := gentoo.ParsePackageStr(pkg)
			if err != nil {
				return errors.New(
					fmt.Sprintf("Invalid package string %s", pkg))
			}

			// Check if exists branch
			branch, present := m.Branches[gp.Category]
			if present {
				m.Log(logger.DebugLevel, "Add package %s for category %s.",
					pkg, gp.Category)
			} else {
				branch, err = NewFilterMatrixBranch(gp.Category)
				if err != nil {
					return errors.New(
						fmt.Sprintf("Error on create FilterMatrixBranch for category %s",
							gp.Category))
				}
				branch.Matrix = m
				m.Branches[gp.Category] = branch
				m.Log(logger.DebugLevel, "Added branch for category %s for package %s.",
					gp.Category, pkg)
			}
			branch.Packages = append(branch.Packages, gp)
			branch.AddResource(r)
		}

	}

	return nil
}

func (m *FilterMatrix) GetMatches() []*FilterMatrixLeaf {
	ans := make([]*FilterMatrixLeaf, 0)

	for _, branch := range m.Branches {
		for _, match := range branch.Matches {
			ans = append(ans, match)
		}
	}

	return ans
}

func (m *FilterMatrix) GetMatchesFiles() []string {
	ans := make([]string, 0)

	for _, branch := range m.Branches {
		for _, match := range branch.Matches {
			ans = append(ans, (*match).Path)
		}
	}

	return ans
}

func (m *FilterMatrix) GetNotMatches() []*FilterMatrixLeaf {
	ans := make([]*FilterMatrixLeaf, 0)

	for _, branch := range m.Branches {
		for _, notMatch := range branch.NotMatches {
			ans = append(ans, notMatch)
		}
	}

	return ans
}

func (m *FilterMatrix) GetNotMatchesFiles() []string {
	ans := make([]string, 0)

	for _, branch := range m.Branches {
		for _, notMatch := range branch.NotMatches {
			ans = append(ans, (*notMatch).Path)
		}
	}

	return ans
}

func (m *FilterMatrix) CheckMatches(binhost map[string][]string) error {
	for category, pkgs := range binhost {

		m.Log(logger.DebugLevel,
			"FilterMatrix/CheckMatches: Analyze category %s", category)
		if branch, ok := m.Branches[category]; ok {
			// POST: exists a branch for the category
			err := branch.CheckPackages(pkgs)
			if err != nil {
				return err
			}

		} else {
			// POST: doesn't exist a branch for the category
			branch, err := NewFilterMatrixBranch(category)
			if err != nil {
				return err
			}
			branch.Matrix = m
			for _, pkg := range pkgs {
				_, err := branch.AddPackage(pkg, false)
				if err != nil {
					return err
				}
			}
			m.Branches[category] = branch
		}

	}

	return nil
}

func NewFilter(settings *viper.Viper, l *logger.Logger, config *sark.SarkConfig) (*Filter, error) {
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

func (f *Filter) Run(binhostDir string) error {
	var err error

	if binhostDir == "" {
		return errors.New("Invalid binhost directory")
	}

	start := time.Now()
	// Phase1: Analyze binhost Directory
	err = binhostdir.AnalyzeBinHostDirectory(binhostDir, f.logger, &f.BinHostTree)
	if err != nil {
		return err
	}
	f.logger.Infoln(
		fmt.Sprintf("Analyze of binhost directory elapsed in %d µs.",
			time.Now().Sub(start).Nanoseconds()/1e3))

	if len(f.BinHostTree) > 0 {
		start = time.Now()
		// Phase2: Create FilterMatrix
		err = f.CreateFilterMatrix()
		if err != nil {
			return err
		}
		f.logger.Infoln(
			fmt.Sprintf("Creation of filter matrix (%s) elapsed in %d µs.",
				f.RulesTree.FilterType, time.Now().Sub(start).Nanoseconds()/1e3))

	} else {
		f.logger.Infof("No files found to filter. Nothing to do.")
		return nil
	}

	start = time.Now()
	// Elaborate matches/not matches
	err = f.RulesTree.CheckMatches(f.BinHostTree)
	if err != nil {
		return err
	}
	f.logger.Infoln(
		fmt.Sprintf("Check matches elapsed in %d µs.",
			time.Now().Sub(start).Nanoseconds()/1e3))

	matches := f.RulesTree.GetMatches()
	f.logger.Infof("Matches packages found %d.", len(matches))

	notMatches := f.RulesTree.GetNotMatches()
	f.logger.Infof("Not matches packages found %d.", len(notMatches))

	// Write report
	if f.settings.GetString("report-prefix-path") != "" {
		report, err := NewFilterReport(f.RulesTree.FilterType)
		if err != nil {
			return err
		}
		report.Matches = f.RulesTree.GetMatchesFiles()
		report.NotMatches = f.RulesTree.GetNotMatchesFiles()
		err = report.WriteReport(f.settings.GetString("report-prefix-path"))
		if err != nil {
			return err
		}
	}

	// Remove filtered files
	if !f.settings.GetBool("dry-run") {
		if f.RulesTree.FilterType == "whitelist" {
			// POST: Remove Not matched
			err = f.unlinkFiles(notMatches)
		} else {
			// POST: Remove matches with blacklist
			err = f.unlinkFiles(matches)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (f *Filter) unlinkFiles(files []*FilterMatrixLeaf) error {
	var inError = false
	for _, l := range files {
		f.logger.Infof("Removing file %s...", l.Path)
		err := os.Remove(l.Path)
		if err != nil {
			f.logger.Errorf("Error on remove file %s: %s",
				l.Path, err)
			inError = true
		}
	}

	if inError {
		return errors.New("Error on removing files")
	} else {
		f.logger.Infof("Removed %d files.", len(files))
	}

	return nil
}

func (f *Filter) CreateFilterMatrix() error {
	if f.Config == nil {
		// Create an empty Sark Config where injection filter
		// has blacklist as filter type and no packages blocked
		f.Config, _ = sark.NewSarkConfig(f.settings, "blacklist")
		f.Config.Id = "filter"
	} else if f.Config.Injector.Filter.FilterType == "" {
		// If there is a filter section I consider filter of type
		// whitelist where packages are all present on target section.
		f.Config.Injector = *sark.NewSarkInjectConfig("whitelist")
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
	err := f.RulesTree.CreateBranches()
	if err != nil {
		return err
	}

	return nil
}
