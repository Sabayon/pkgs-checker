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
	. "github.com/Sabayon/pkgs-checker/commons"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FilterRule", func() {

	Describe("NewSarkFilterRuleConf", func() {

		f := NewSarkFilterRuleConf("test")
		f.AddPackage("net-misc/ntp")
		f.AddCategory("net-misc")
		f.AddFile("pkglist|/tmp/Pkglist")
		f.AddUrl("build|https://mynode.it/build.yaml")

		Context("Check rule", func() {
			It("Check attributes", func() {
				Expect(f.Descr).To(Equal("test"))
				Expect(f.Packages[0]).To(Equal("net-misc/ntp"))
				Expect(f.Categories[0]).Should(Equal("net-misc"))
				Expect(f.Files[0]).To(Equal("pkglist|/tmp/Pkglist"))
				Expect(f.Urls[0]).To(Equal("build|https://mynode.it/build.yaml"))
			})

		})

	})

})
