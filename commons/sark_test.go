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

var _ = Describe("SarkConfig", func() {

	Describe("NewSarkConfigFromString", func() {

		conf := `---
injector:
  filter:
    type: "whitelist"
    rules:
      - categories:
        - "dev-libs"
        - "www-servers"
`
		sark, err := NewSarkConfigFromString(nil, conf)

		out, _ := sark.ToString()
		fmt.Println(fmt.Sprintf("SARK %s", out))

		Context("Check processing phase", func() {
			It("Check attributes", func() {
				Expect(err).Should(BeNil())
			})
		})

		Context("Check filter", func() {
			It("Check attributes", func() {
				Expect(sark.Injector.Filter.FilterType).To(Equal("whitelist"))
			})
		})

	})

	Describe("NewSarkConfigFromFile", func() {

		sark, err := NewSarkConfigFromFile(nil, "../tests/sark/inject_example1.yaml")

		Context("Check processing phase", func() {
			It("Check attributes", func() {
				Expect(err).Should(BeNil())
			})
		})

		Context("Check filter", func() {
			It("Check attributes", func() {
				Expect(sark.Injector.Filter.FilterType).To(Equal("blacklist"))
			})
		})

	})

})
