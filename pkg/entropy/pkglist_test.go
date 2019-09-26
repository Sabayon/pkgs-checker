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
package entropy_test

import (
	. "github.com/Sabayon/pkgs-checker/pkg/entropy"
	. "github.com/Sabayon/pkgs-checker/pkg/gentoo"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Entropy Pkglist", func() {

	Describe("Parse String1", func() {

		ep, err := NewEntropyPackage("sys-fs/udftools-2.1~1")

		Context("Check processing phase", func() {
			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check element", func() {
				Expect(ep).Should(Equal(&EntropyPackage{
					GentooPackage: &GentooPackage{
						Category:  "sys-fs",
						Name:      "udftools",
						Version:   "2.1",
						Slot:      "0",
						Condition: PkgCondEqual,
					},
					Revision: 1,
				}))
			})

		})

	})

	Describe("Parse String2 without revision", func() {

		ep, err := NewEntropyPackage("sys-fs/udftools-2.1")

		Context("Check processing phase", func() {
			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check element", func() {
				Expect(ep).Should(Equal(&EntropyPackage{
					GentooPackage: &GentooPackage{
						Category:  "sys-fs",
						Name:      "udftools",
						Version:   "2.1",
						Slot:      "0",
						Condition: PkgCondEqual,
					},
					Revision: 0,
				}))
			})

		})

	})

})
