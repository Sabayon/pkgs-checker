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

package commons_test

import (
	"fmt"

	. "github.com/Sabayon/pkgs-checker/commons"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gentoo Packages", func() {

	Describe("Parse package strings", func() {

		Context("Matches any version of a package", func() {

			pkg, err := ParsePackageStr("x11-libs/gtk+")
			fmt.Println(fmt.Sprintf("pkg %s", pkg))

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check pkgName", func() {
				Expect((*pkg).Name).Should(Equal("gtk+"))
			})

			It("Check category", func() {
				Expect((*pkg).Category).Should(Equal("x11-libs"))
			})

			//It("Check cond", func() {
			// TODO: check how use PkgCondInvalid
			//		Expect((*pkg).Condition).Should(Equal(PkgCondInvalid.(PackageCond)))
			//		})
		})

	})

})
