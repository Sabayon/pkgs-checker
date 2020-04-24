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
package entropy

import (
	"errors"
	"regexp"
	"strconv"

	gentoo "github.com/Sabayon/pkgs-checker/pkg/gentoo"
)

type EntropyPackage struct {
	*gentoo.GentooPackage
	Revision int `json:"revision",omitempty`
}

func EntropyIsPkgWithRevision(pkgname string) (ans bool) {
	ans = false

	if pkgname != "" {
		regexRev := regexp.MustCompile(
			"[~][0-9]*$",
		)

		matches := regexRev.FindAllString(pkgname, -1)
		if len(matches) > 0 {
			ans = true
		}
	}

	return
}

func NewEntropyPackage(pkgname string) (*EntropyPackage, error) {
	var ans *EntropyPackage

	if pkgname == "" {
		return nil, errors.New("Invalid pkgname")
	}

	regexRev := regexp.MustCompile(
		"[~][0-9]*$",
	)

	matches := regexRev.FindAllString(pkgname, -1)
	if len(matches) > 0 {
		gPkgname := pkgname[:len(pkgname)-len(matches[0])]

		gp, err := gentoo.ParsePackageStr(gPkgname)
		if err != nil {
			return nil, err
		}

		rev, err := strconv.Atoi(matches[0][1:])
		if err != nil {
			return nil, err
		}
		ans = &EntropyPackage{
			GentooPackage: gp,
			Revision:      rev,
		}
	} else {
		gp, err := gentoo.ParsePackageStr(pkgname)
		if err != nil {
			return nil, err
		}
		ans = &EntropyPackage{
			GentooPackage: gp,
			Revision:      0,
		}
	}

	return ans, nil
}
