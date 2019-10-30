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

package filter_test

import (
	"fmt"

	. "github.com/Sabayon/pkgs-checker/pkg/filter"
	gentoo "github.com/Sabayon/pkgs-checker/pkg/gentoo"
	sark "github.com/Sabayon/pkgs-checker/pkg/sark"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FilterRule", func() {

	Describe("NewSarkFilterRuleConf", func() {

		f := sark.NewSarkFilterRuleConf("test")
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

	Describe("NewFilterMatrix1", func() {

		matrix, err := NewFilterMatrix("whitelist")

		Context("Check matrix", func() {

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check type", func() {
				Expect(matrix.FilterType).Should(Equal("whitelist"))
			})

		})

		pkgs := []string{"net-libs/gnutls", "dev-libs/mpc"}
		resource, _ := NewFilterResource("test", "buildfile", pkgs, nil)
		err = matrix.AddResource(resource)

		binHostTree := make(map[string][]string, 2)
		binHostTree["net-libs"] = []string{
			"/tmp/net-libs/gnutls-1.1.1.tbz2",
			"/tmp/net-libs/nodejs-9.11.1.tbz2",
		}
		binHostTree["dev-libs"] = []string{
			"/tmp/dev-libs/mpc-22.2.2.tbz",
		}

		fmt.Println(fmt.Sprintf("RESOURCE %s", resource))
		Context("Check resource", func() {

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

		})

		err = matrix.CreateBranches()
		Context("Check branch1", func() {

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			b, ok := matrix.Branches["net-libs"]
			It("Check if exist net-libs branch", func() {
				Expect(ok).Should(Equal(true))
			})

			It("Check if resource is set", func() {
				Expect(b.Resources[0]).Should(Equal(resource))
			})

			It("Check if package is set", func() {
				Expect(b.Packages[0].Name).Should(Equal("gnutls"))
			})

			It("Check if matrix is set", func() {
				Expect(b.Matrix).Should(Equal(matrix))
			})
		})

		Context("Check branch1", func() {
			b, ok := matrix.Branches["dev-libs"]
			It("Check if exist dev-libs branch", func() {
				Expect(ok).Should(Equal(true))
			})

			It("Check if resource is set", func() {
				Expect(b.Resources[0]).Should(Equal(resource))
			})

			It("Check if package is set", func() {
				Expect(b.Packages[0].Name).Should(Equal("mpc"))
			})

			It("Check if matrix is set", func() {
				Expect(b.Matrix).Should(Equal(matrix))
			})
		})

		Context("Check matches", func() {
			err := matrix.CheckMatches(binHostTree)
			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check matched len", func() {
				Expect(len(matrix.GetMatches())).Should(Equal(2))
			})

			It("Check not matched len", func() {
				Expect(len(matrix.GetNotMatches())).Should(Equal(1))
			})

			It("Check not matched leaf name", func() {
				Expect((*(matrix.GetNotMatches()[0])).Name).Should(Equal("nodejs"))
				Expect((*(matrix.GetNotMatches()[0])).Path).Should(Equal(
					"/tmp/net-libs/nodejs-9.11.1.tbz2",
				))

				Expect(*(*(matrix.GetNotMatches()[0])).Package).Should(Equal(
					gentoo.GentooPackage{
						Name:      "nodejs",
						Category:  "net-libs",
						Version:   "9.11.1",
						Condition: gentoo.PkgCondEqual,
						Slot:      "0",
					},
				))
			})

		})
	})

	// Check filter when there are both files and categories configured
	Describe("NewFilterMatrix2", func() {

		matrix, err := NewFilterMatrix("whitelist")

		Context("Check matrix", func() {

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check type", func() {
				Expect(matrix.FilterType).Should(Equal("whitelist"))
			})

		})

		pkgs := []string{"net-libs/gnutls", "dev-libs/mpc"}
		resource, _ := NewFilterResource("test", "buildfile", pkgs, nil)
		err = matrix.AddResource(resource)

		resource2, _ := NewFilterResource("test2", "buildfile", nil, []string{"dev-node"})
		err2 := matrix.AddResource(resource2)

		binHostTree := make(map[string][]string, 2)
		binHostTree["net-libs"] = []string{
			"/tmp/net-libs/gnutls-1.1.1.tbz2",
			"/tmp/net-libs/nodejs-9.11.1.tbz2",
		}
		binHostTree["dev-libs"] = []string{
			"/tmp/dev-libs/mpc-22.2.2.tbz",
		}
		binHostTree["dev-node"] = []string{
			"/tmp/dev-node/os-tmpdir-1.0.2.tbz2",
		}

		fmt.Println(fmt.Sprintf("RESOURCE %s", resource))
		Context("Check resource", func() {

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check error2", func() {
				Expect(err2).Should(BeNil())
			})
		})

		err = matrix.CreateBranches()
		Context("Check branch1", func() {

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			b, ok := matrix.Branches["net-libs"]
			It("Check if exist net-libs branch", func() {
				Expect(ok).Should(Equal(true))
			})

			It("Check if resource is set", func() {
				Expect(b.Resources[0]).Should(Equal(resource))
			})

			It("Check if package is set", func() {
				Expect(b.Packages[0].Name).Should(Equal("gnutls"))
			})

			It("Check if matrix is set", func() {
				Expect(b.Matrix).Should(Equal(matrix))
			})

		})

		Context("Check branch2", func() {
			b, ok := matrix.Branches["dev-libs"]
			It("Check if exist dev-libs branch", func() {
				Expect(ok).Should(Equal(true))
			})

			It("Check if resource is set", func() {
				Expect(b.Resources[0]).Should(Equal(resource))
			})

			It("Check if package is set", func() {
				Expect(b.Packages[0].Name).Should(Equal("mpc"))
			})

			It("Check if matrix is set", func() {
				Expect(b.Matrix).Should(Equal(matrix))
			})
		})

		Context("Check branch3", func() {
			b, ok := matrix.Branches["dev-node"]
			It("Check if exist dev-node branch", func() {
				Expect(ok).Should(Equal(true))
			})

			It("Check if resource is set", func() {
				Expect(b.Resources[0]).Should(Equal(resource2))
			})

			It("Check if category is set", func() {
				Expect(b.CategoryFiltered).Should(Equal(true))
			})

			It("Check if matrix is set", func() {
				Expect(b.Matrix).Should(Equal(matrix))
			})
		})

		Context("Check matches", func() {
			err := matrix.CheckMatches(binHostTree)
			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check matched len", func() {
				Expect(len(matrix.GetMatches())).Should(Equal(3))
			})

			It("Check not matched len", func() {
				Expect(len(matrix.GetNotMatches())).Should(Equal(1))
			})

			It("Check not matched leaf name", func() {
				Expect((*(matrix.GetNotMatches()[0])).Name).Should(Equal("nodejs"))
				Expect((*(matrix.GetNotMatches()[0])).Path).Should(Equal(
					"/tmp/net-libs/nodejs-9.11.1.tbz2",
				))

				Expect(*(*(matrix.GetNotMatches()[0])).Package).Should(Equal(
					gentoo.GentooPackage{
						Name:      "nodejs",
						Category:  "net-libs",
						Version:   "9.11.1",
						Condition: gentoo.PkgCondEqual,
						Slot:      "0",
					},
				))
			})

		})
	})

})
