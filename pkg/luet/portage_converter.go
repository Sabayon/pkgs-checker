/*

Copyright (C) 2021  Daniele Rondina <geaaru@sabayonlinux.org>

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
package luet

import (
	"strings"

	"github.com/Sabayon/pkgs-checker/pkg/gentoo"
)

type PortageConverterUseFlags struct {
	Disabled []string `json:"disabled,omitempty" yaml:"disabled,omitempty"`
	Enabled  []string `json:"enabled,omitempty" yaml:"enabled,omitempty"`
}

type PortageConverterArtefact struct {
	Tree     string                   `json:"tree" yaml:"tree"`
	Uses     PortageConverterUseFlags `json:"uses,omitempty" yaml:"uses,omitempty"`
	Packages []string                 `json:"packages" yaml:"packages"`
}

type PortageConverterArtefacts struct {
	Artefacts []PortageConverterArtefact `json:"artefacts" yaml:"artefacts"`
}

func ConvertPortageMeta2PortageConverter(pkgs []*gentoo.PortageMetaData, treePath string) PortageConverterArtefacts {
	ans := PortageConverterArtefacts{
		Artefacts: []PortageConverterArtefact{},
	}

	// Create a map to merge packages with same use flags
	mPkgs := make(map[string][]*gentoo.PortageMetaData, 0)

	for _, p := range pkgs {

		key := strings.Join(p.GentooPackage.UseFlags, "-")
		if _, ok := mPkgs[key]; ok {
			mPkgs[key] = append(mPkgs[key], p)
		} else {
			mPkgs[key] = []*gentoo.PortageMetaData{p}
		}
	}

	for _, pp := range mPkgs {
		enabledUses := []string{}
		disabledUses := []string{}

		artefact := PortageConverterArtefact{
			Tree:     treePath,
			Packages: []string{},
		}

		for _, useflag := range pp[0].UseFlags {
			if strings.HasPrefix(useflag, "-") {
				disabledUses = append(disabledUses, useflag[1:])
			} else {
				enabledUses = append(enabledUses, useflag)
			}
		}

		for _, p := range pp {
			// Ignoring sub slot
			if strings.Contains(p.Slot, "/") {
				p.Slot = p.Slot[0:strings.Index(p.Slot, "/")]
			}

			artefact.Packages = append(artefact.Packages, p.GetPackageNameWithSlot())
		}

		artefact.Uses = PortageConverterUseFlags{
			Enabled:  enabledUses,
			Disabled: disabledUses,
		}

		ans.Artefacts = append(ans.Artefacts, artefact)
	}

	return ans
}
