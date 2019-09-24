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
	"os"
	"path"

	. "github.com/Sabayon/pkgs-checker/pkg/commons"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PKGLIST", func() {

	Describe("file commons", func() {

		Context("AbsPathFromBase: Example1", func() {
			ans, err := AbsPathFromBase(
				"/sabayon/sbi-tasks/next/arm/core",
				"../base/base-core-staging1.yaml",
			)
			It("Check err", func() {
				Expect(err).Should(BeNil())
			})
			It("Check ans", func() {
				Expect(ans).Should(
					Equal("/sabayon/sbi-tasks/next/arm/base/base-core-staging1.yaml"))
			})
		})

		Context("AbsPathFromBase: Example2", func() {
			ans, err := AbsPathFromBase(
				"/sabayon/sbi-tasks/next/arm/core",
				"/sabayon/sbi-tasks/next/arm/base/base-core-staging1.yaml",
			)
			It("Check err", func() {
				Expect(err).Should(BeNil())
			})
			It("Check ans", func() {
				Expect(ans).Should(
					Equal("/sabayon/sbi-tasks/next/arm/base/base-core-staging1.yaml"))
			})
		})

		Context("AbsPathFromBase: Example3", func() {
			// Override PWD for influence GetWd used by filepath.Abs method
			pwd, _ := os.Getwd()
			ans, err := AbsPathFromBase(
				"next/arm/core",
				"../base/base-core-staging1.yaml",
			)
			It("Check err", func() {
				Expect(err).Should(BeNil())
			})
			It("Check ans", func() {
				Expect(ans).Should(
					Equal(path.Clean(
						path.Join(pwd, "next/arm/", "base/base-core-staging1.yaml"))))
			})
		})

		Context("AbsPathFromBase: Example4", func() {
			ans, err := AbsPathFromBase(
				"/sabayon/sbi-tasks/next/arm/core",
				"../base",
			)
			It("Check err", func() {
				Expect(err).Should(BeNil())
			})
			It("Check ans", func() {
				Expect(ans).Should(
					Equal("/sabayon/sbi-tasks/next/arm/base"))
			})
		})

		Context("AbsPathFromBase: Example5", func() {
			ans, err := AbsPathFromBase(
				"/sabayon/sbi-tasks/next/arm/core",
				"../../../",
			)
			It("Check err", func() {
				Expect(err).Should(BeNil())
			})
			It("Check ans", func() {
				Expect(ans).Should(
					Equal("/sabayon/sbi-tasks"))
			})
		})

		Context("AbsPathFromBase: Example6", func() {
			ans, err := AbsPathFromBase(
				"/sabayon/sbi-tasks/next/arm/core",
				"../../../../../",
			)
			It("Check err", func() {
				Expect(err).Should(BeNil())
			})
			It("Check ans", func() {
				Expect(ans).Should(
					Equal("/"))
			})
		})

		Context("AbsPathFromBase: Example7", func() {
			ans, err := AbsPathFromBase(
				"/sabayon/sbi-tasks/next/arm/core",
				"../../../../../../../../etc/sabayon/release",
			)
			It("Check err", func() {
				Expect(err).Should(BeNil())
			})
			It("Check ans", func() {
				Expect(ans).Should(
					Equal("/etc/sabayon/release"))
			})
		})

	})

})
