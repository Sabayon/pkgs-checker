/*

Copyright (C) 2017-2021  Daniele Rondina <geaaru@sabayonlinux.org>

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
package gentoo

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type PortageMetaData struct {
	*GentooPackage `json:"package,omitempty"`
	IUse           []string `json:"iuse,omitempty"`
	IUseEffective  []string `json:"iuse_effective,omitempty"`
	Use            []string `json:"use,omitempty"`
	Eapi           string   `json:"eapi,omitempty"`
	CxxFlags       string   `json:"cxxflags,omitempty"`
	LdFlags        string   `json:"ldflags,omitempty"`
	CHost          string   `json:"chost,omitempty"`
	BDEPEND        string   `json:"bdepend,omitempty"`
	RDEPEND        string   `json:"rdepend,omitempty"`
	DEPEND         string   `json:"depend,omitempty"`
	REQUIRES       string   `json:"requires,omitempty"`
	KEYWORDS       string   `json:"keywords,omitempty"`
	PROVIDES       string   `json:"provides,omitempty"`
	SIZE           string   `json:"size,omitempty"`
}

type PortageUseParseOpts struct {
	UseFilters []string `json:"use_filters,omitempty" yaml:"use_filters,omitempty"`
	Categories []string `json:"categories,omitempty" yaml:"categories,omitempty"`
	Packages   []string `json:"pkgs_filters,omitempty" yaml:"pkgs_filters,omitempty"`
}

func NewPortageMetaData(pkg *GentooPackage) *PortageMetaData {
	return &PortageMetaData{
		GentooPackage: pkg,
		IUse:          make([]string, 0),
		IUseEffective: make([]string, 0),
		Use:           make([]string, 0),
		Eapi:          "",
		CxxFlags:      "",
		LdFlags:       "",
		BDEPEND:       "",
		RDEPEND:       "",
		DEPEND:        "",
	}
}

func (o *PortageUseParseOpts) IsCatAdmit(cat string) bool {
	ans := false

	if len(o.Categories) > 0 {
		for _, c := range o.Categories {
			if c == cat {
				ans = true
				break
			}
		}
	} else {
		ans = true
	}

	return ans
}

func (o *PortageUseParseOpts) IsPkgAdmit(pkg string) bool {
	ans := false

	// Prepare regex
	if len(o.Packages) > 0 {
		for _, f := range o.Packages {
			r := regexp.MustCompile(f)
			if r != nil {
				if r.MatchString(pkg) {
					ans = true
					break
				}
			} else {
				fmt.Println("WARNING: Regex " + f + " not compiled.")
			}
		}
	} else {
		ans = true
	}

	return ans
}

func ParseMetadataDir(dir string, opts PortageUseParseOpts) ([]*PortageMetaData, error) {
	ans := make([]*PortageMetaData, 0)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return ans, err
	}

	for _, file := range files {
		if file.IsDir() && opts.IsCatAdmit(file.Name()) {
			pkgs, err := ParseMetadataCatDir(filepath.Join(dir, file.Name()), opts)
			if err != nil {
				return ans, errors.New(
					fmt.Sprintf("Error on parse directory %s: %s",
						file.Name(), err.Error()))
			}

			ans = append(ans, pkgs...)
		}
	}

	return ans, nil
}

func ParseMetadataCatDir(dir string, opts PortageUseParseOpts) ([]*PortageMetaData, error) {
	ans := make([]*PortageMetaData, 0)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return ans, err
	}

	for _, file := range files {
		if file.IsDir() {
			pm, err := ParsePackageMetadataDir(filepath.Join(dir, file.Name()), opts)
			if err != nil {
				return ans, errors.New(
					fmt.Sprintf("Error on parse directory %s/%s: %s",
						dir, file.Name(), err.Error()))
			}

			if opts.IsPkgAdmit(pm.GetPackageNameWithSlot()) {
				ans = append(ans, pm)
			}
		}
	}

	return ans, nil
}

func ParsePackageMetadataDir(dir string, opts PortageUseParseOpts) (*PortageMetaData, error) {
	var ans *PortageMetaData = nil

	// Check if the directory is valid
	fi, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}

	if !fi.IsDir() {
		return nil, errors.New("Path " + dir + " is not a directory!")
	}

	metaDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	words := strings.Split(metaDir, "/")

	if len(words) <= 2 {
		return nil, errors.New("Path " + dir + " is invalid!")
	}

	pkgname := fmt.Sprintf("%s/%s",
		words[len(words)-2],
		words[len(words)-1],
	)

	gp, err := ParsePackageStr(pkgname)
	if err != nil {
		return nil, errors.New("Error on parse pkgname " + err.Error())
	}

	ans = NewPortageMetaData(gp)

	ans.BDEPEND, err = parseMetaFile(filepath.Join(metaDir, "BDEPEND"), true)
	if err != nil {
		return nil, err
	}

	ans.RDEPEND, err = parseMetaFile(filepath.Join(metaDir, "RDEPEND"), true)
	if err != nil {
		return nil, err
	}

	ans.DEPEND, err = parseMetaFile(filepath.Join(metaDir, "DEPEND"), true)
	if err != nil {
		return nil, err
	}

	ans.GentooPackage.Slot, err = parseMetaFile(
		filepath.Join(metaDir, "SLOT"), true,
	)
	if err != nil {
		return nil, err
	}

	ans.Eapi, err = parseMetaFile(
		filepath.Join(metaDir, "EAPI"), true,
	)
	if err != nil {
		return nil, err
	}

	ans.CxxFlags, err = parseMetaFile(
		filepath.Join(metaDir, "CXXFLAGS"), true,
	)
	if err != nil {
		return nil, err
	}

	ans.LdFlags, err = parseMetaFile(
		filepath.Join(metaDir, "LDFLAGS"), true,
	)
	if err != nil {
		return nil, err
	}

	ans.CHost, err = parseMetaFile(
		filepath.Join(metaDir, "CHOST"), true,
	)
	if err != nil {
		return nil, err
	}

	ans.GentooPackage.License, err = parseMetaFile(
		filepath.Join(metaDir, "LICENSE"), true,
	)
	if err != nil {
		return nil, err
	}

	ans.GentooPackage.Repository, err = parseMetaFile(
		filepath.Join(metaDir, "repository"), true,
	)
	if err != nil {
		return nil, err
	}

	ans.REQUIRES, err = parseMetaFile(
		filepath.Join(metaDir, "REQUIRES"), true,
	)
	if err != nil {
		return nil, err
	}

	ans.KEYWORDS, err = parseMetaFile(
		filepath.Join(metaDir, "KEYWORDS"), true,
	)
	if err != nil {
		return nil, err
	}

	ans.PROVIDES, err = parseMetaFile(
		filepath.Join(metaDir, "PROVIDES"), true,
	)
	if err != nil {
		return nil, err
	}

	ans.SIZE, err = parseMetaFile(
		filepath.Join(metaDir, "SIZE"), true,
	)
	if err != nil {
		return nil, err
	}

	iuse, err := parseMetaFile(
		filepath.Join(metaDir, "IUSE"), true,
	)
	if err != nil {
		return nil, err
	}
	if iuse != "" {
		ans.IUse = strings.Split(iuse, " ")
	}

	iuseEffective, err := parseMetaFile(
		filepath.Join(metaDir, "IUSE_EFFECTIVE"), true,
	)
	if err != nil {
		return nil, err
	}
	if iuseEffective != "" {
		ans.IUseEffective = strings.Split(iuseEffective, " ")
	}

	use, err := parseMetaFile(
		filepath.Join(metaDir, "USE"), true,
	)
	if err != nil {
		return nil, err
	}
	if use != "" {
		ans.Use = strings.Split(use, " ")
	}

	if len(ans.IUseEffective) > 0 {
		ans.GentooPackage.UseFlags = elaborateUses(ans.IUseEffective, ans.Use, opts)
	}

	return ans, nil
}

func useInArray(use string, arr []string) bool {
	ans := false
	for _, u := range arr {
		if use == u {
			ans = true
			break
		}
	}
	return ans
}

func elaborateUses(iuse, use []string, opts PortageUseParseOpts) []string {
	ans := []string{}

	// Prepare regex
	listRegex := []*regexp.Regexp{}
	for _, f := range opts.UseFilters {
		r := regexp.MustCompile(f)
		if r != nil {
			listRegex = append(listRegex, r)
		} else {
			fmt.Println("WARNING: Regex " + f + " not compiled.")
		}
	}

	for _, u := range iuse {

		toSkip := false

		if strings.HasPrefix(u, "+") {
			u = u[1:]
		}

		// Check if use flags is filtered
		if len(listRegex) > 0 {
			for _, r := range listRegex {
				if r.MatchString(u) {
					toSkip = true
					//fmt.Println("MATCHED FILTER ", u)
					break
				}
			}
		}

		if toSkip {
			continue
		}

		if useInArray(u, use) {
			ans = append(ans, u)
		} else {
			ans = append(ans, "-"+u)
		}
	}

	return ans
}

func parseMetaFile(file string, dropLn bool) (string, error) {
	_, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	ans := string(data)
	if dropLn {
		ans = strings.TrimRight(ans, "\n")
	}

	return ans, nil
}
