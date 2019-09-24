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

package gentoo_test

import (
	"fmt"
	"sort"

	. "github.com/Sabayon/pkgs-checker/pkg/gentoo"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	version "github.com/hashicorp/go-version"
)

var _ = Describe("Gentoo Packages", func() {

	Describe("Parse package strings", func() {

		// Tests by: https://wiki.gentoo.org/wiki/Version_specifier

		Context("Matches any version of a package", func() {

			pkg, err := ParsePackageStr("x11-libs/gtk+")
			g := GentooPackage{
				Name:       "gtk+",
				Category:   "x11-libs",
				Condition:  PkgCondInvalid,
				Slot:       "0",
				Repository: "",
			}
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

			It("Check cond", func() {
				// TODO: check how use PkgCondInvalid
				Expect((*pkg).Condition).Should(Equal(g.Condition))
			})

			It("Check struct", func() {
				// TODO: check how use PkgCondInvalid
				Expect((*pkg)).Should(Equal(g))
			})
		})

		Context("Matches any version and any revision", func() {

			pkg, err := ParsePackageStr("~sys-devel/gdb-7.3")
			g := GentooPackage{
				Name:       "gdb",
				Category:   "sys-devel",
				Condition:  PkgCondAnyRevision,
				Version:    "7.3",
				Slot:       "0",
				Repository: "",
			}
			fmt.Println(fmt.Sprintf("pkg %s", pkg))

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check pkgName", func() {
				Expect((*pkg).Name).Should(Equal("gdb"))
			})

			It("Check category", func() {
				Expect((*pkg).Category).Should(Equal("sys-devel"))
			})

			It("Check cond", func() {
				// TODO: check how use PkgCondInvalid
				Expect((*pkg).Condition).Should(Equal(g.Condition))
			})

			It("Check struct", func() {
				// TODO: check how use PkgCondInvalid
				Expect((*pkg)).Should(Equal(g))
			})
		})

		Context("Matches a version by the version range.", func() {

			pkg, err := ParsePackageStr("=sys-devel/gdb-7.3*")
			g := GentooPackage{
				Name:       "gdb",
				Category:   "sys-devel",
				Condition:  PkgCondMatchVersion,
				Version:    "7.3",
				Slot:       "0",
				Repository: "",
			}
			fmt.Println(fmt.Sprintf("pkg %s", pkg))

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check pkgName", func() {
				Expect((*pkg).Name).Should(Equal("gdb"))
			})

			It("Check category", func() {
				Expect((*pkg).Category).Should(Equal("sys-devel"))
			})

			It("Check cond", func() {
				// TODO: check how use PkgCondInvalid
				Expect((*pkg).Condition).Should(Equal(g.Condition))
			})

			It("Check struct", func() {
				// TODO: check how use PkgCondInvalid
				Expect((*pkg)).Should(Equal(g))
			})
		})

		Context("Matches a version exactly.", func() {

			pkg, err := ParsePackageStr("=www-client/firefox-7.0")
			g := GentooPackage{
				Name:       "firefox",
				Category:   "www-client",
				Condition:  PkgCondEqual,
				Version:    "7.0",
				Slot:       "0",
				Repository: "",
			}
			fmt.Println(fmt.Sprintf("pkg %s", pkg))

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check pkgName", func() {
				Expect((*pkg).Name).Should(Equal("firefox"))
			})

			It("Check category", func() {
				Expect((*pkg).Category).Should(Equal("www-client"))
			})

			It("Check cond", func() {
				// TODO: check how use PkgCondInvalid
				Expect((*pkg).Condition).Should(Equal(g.Condition))
			})

			It("Check struct", func() {
				// TODO: check how use PkgCondInvalid
				Expect((*pkg)).Should(Equal(g))
			})
		})

		Context("Matches the specified version or any higher version.", func() {

			pkg, err := ParsePackageStr(">=dev-lang/python-2.7")
			g := GentooPackage{
				Name:       "python",
				Category:   "dev-lang",
				Condition:  PkgCondGreaterEqual,
				Version:    "2.7",
				Slot:       "0",
				Repository: "",
			}
			fmt.Println(fmt.Sprintf("pkg %s", pkg))

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check pkgName", func() {
				Expect((*pkg).Name).Should(Equal("python"))
			})

			It("Check category", func() {
				Expect((*pkg).Category).Should(Equal("dev-lang"))
			})

			It("Check cond", func() {
				// TODO: check how use PkgCondInvalid
				Expect((*pkg).Condition).Should(Equal(g.Condition))
			})

			It("Check struct", func() {
				// TODO: check how use PkgCondInvalid
				Expect((*pkg)).Should(Equal(g))
			})
		})

		Context("Matches a version strictly later then specified.", func() {

			pkg, err := ParsePackageStr(">dev-lang/python-2.7")
			g := GentooPackage{
				Name:       "python",
				Category:   "dev-lang",
				Condition:  PkgCondGreater,
				Version:    "2.7",
				Slot:       "0",
				Repository: "",
			}
			fmt.Println(fmt.Sprintf("pkg %s", pkg))

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check pkgName", func() {
				Expect((*pkg).Name).Should(Equal("python"))
			})

			It("Check category", func() {
				Expect((*pkg).Category).Should(Equal("dev-lang"))
			})

			It("Check cond", func() {
				// TODO: check how use PkgCondInvalid
				Expect((*pkg).Condition).Should(Equal(g.Condition))
			})

			It("Check struct", func() {
				// TODO: check how use PkgCondInvalid
				Expect((*pkg)).Should(Equal(g))
			})
		})

		Context("Matches a version strictly older than specified.", func() {

			pkg, err := ParsePackageStr("<dev-python/beautifulsoup-3.2")
			g := GentooPackage{
				Name:       "beautifulsoup",
				Category:   "dev-python",
				Condition:  PkgCondLess,
				Version:    "3.2",
				Slot:       "0",
				Repository: "",
			}
			fmt.Println(fmt.Sprintf("pkg %s", pkg))

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check pkgName", func() {
				Expect((*pkg).Name).Should(Equal("beautifulsoup"))
			})

			It("Check category", func() {
				Expect((*pkg).Category).Should(Equal("dev-python"))
			})

			It("Check cond", func() {
				// TODO: check how use PkgCondInvalid
				Expect((*pkg).Condition).Should(Equal(g.Condition))
			})

			It("Check struct", func() {
				// TODO: check how use PkgCondInvalid
				Expect((*pkg)).Should(Equal(g))
			})
		})

		Context("Matches the specified version or any older version.", func() {

			pkg, err := ParsePackageStr("<=sys-fs/udev-171")
			g := GentooPackage{
				Name:       "udev",
				Category:   "sys-fs",
				Condition:  PkgCondLessEqual,
				Version:    "171",
				Slot:       "0",
				Repository: "",
			}
			fmt.Println(fmt.Sprintf("pkg %s", pkg))

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check pkgName", func() {
				Expect((*pkg).Name).Should(Equal("udev"))
			})

			It("Check category", func() {
				Expect((*pkg).Category).Should(Equal("sys-fs"))
			})

			It("Check cond", func() {
				// TODO: check how use PkgCondInvalid
				Expect((*pkg).Condition).Should(Equal(g.Condition))
			})

			It("Check struct", func() {
				// TODO: check how use PkgCondInvalid
				Expect((*pkg)).Should(Equal(g))
			})
		})

		Context("Matches package in the specified package SLOT. Note that there is no prefix.", func() {

			pkg, err := ParsePackageStr("dev-db/sqlite:1")
			g := GentooPackage{
				Name:       "sqlite",
				Category:   "dev-db",
				Condition:  PkgCondInvalid,
				Version:    "",
				Slot:       "1",
				Repository: "",
			}
			fmt.Println(fmt.Sprintf("pkg %s", pkg))

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check pkgName", func() {
				Expect((*pkg).Name).Should(Equal("sqlite"))
			})

			It("Check category", func() {
				Expect((*pkg).Category).Should(Equal("dev-db"))
			})

			It("Check cond", func() {
				// TODO: check how use PkgCondInvalid
				Expect((*pkg).Condition).Should(Equal(g.Condition))
			})

			It("Check struct", func() {
				// TODO: check how use PkgCondInvalid
				Expect((*pkg)).Should(Equal(g))
			})
		})

		Context("Matches a package from a specific ebuild repository. This can be combined with other specifiers. The official Gentoo repository is ::gentoo.", func() {

			pkg, err := ParsePackageStr("=media-libs/mesa-9999::x11")
			g := GentooPackage{
				Name:       "mesa",
				Category:   "media-libs",
				Condition:  PkgCondEqual,
				Version:    "9999",
				Slot:       "0",
				Repository: "x11",
			}
			fmt.Println(fmt.Sprintf("pkg %s", pkg))

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check pkgName", func() {
				Expect((*pkg).Name).Should(Equal("mesa"))
			})

			It("Check category", func() {
				Expect((*pkg).Category).Should(Equal("media-libs"))
			})

			It("Check cond", func() {
				// TODO: check how use PkgCondInvalid
				Expect((*pkg).Condition).Should(Equal(g.Condition))
			})

			It("Check struct", func() {
				// TODO: check how use PkgCondInvalid
				Expect((*pkg)).Should(Equal(g))
			})
		})

		Context("Invalid package string1", func() {

			pkg, err := ParsePackageStr("foo")
			fmt.Println(fmt.Sprintf("pkg %s", pkg))

			It("Check error", func() {
				Expect(err).ShouldNot(BeNil())
			})
		})

		Context("Invalid package string2", func() {

			pkg, err := ParsePackageStr("")
			fmt.Println(fmt.Sprintf("pkg %s", pkg))

			It("Check error", func() {
				Expect(err).ShouldNot(BeNil())
			})
		})

		Context("Matches version with two dots.", func() {

			pkg, err := ParsePackageStr("=dev-python/docker-py-3.7.1")
			g := GentooPackage{
				Name:       "docker-py",
				Category:   "dev-python",
				Condition:  PkgCondEqual,
				Version:    "3.7.1",
				Slot:       "0",
				Repository: "",
			}
			fmt.Println(fmt.Sprintf("pkg %s", pkg))

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check pkgName", func() {
				Expect((*pkg).Name).Should(Equal("docker-py"))
			})

			It("Check category", func() {
				Expect((*pkg).Category).Should(Equal("dev-python"))
			})

			It("Check cond", func() {
				// TODO: check how use PkgCondInvalid
				Expect((*pkg).Condition).Should(Equal(g.Condition))
			})

			It("Check struct", func() {
				// TODO: check how use PkgCondInvalid
				Expect((*pkg)).Should(Equal(g))
			})
		})

		Context("Matches version with 4 numbers on version.", func() {

			pkg, err := ParsePackageStr("=dev-db/oracle-instantclient-sqlplus-12.1.0.2")
			g := GentooPackage{
				Name:       "oracle-instantclient-sqlplus",
				Category:   "dev-db",
				Condition:  PkgCondEqual,
				Version:    "12.1.0.2",
				Slot:       "0",
				Repository: "",
			}
			fmt.Println(fmt.Sprintf("pkg %s", pkg))

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check pkgName", func() {
				Expect((*pkg).Name).Should(Equal("oracle-instantclient-sqlplus"))
			})

			It("Check category", func() {
				Expect((*pkg).Category).Should(Equal("dev-db"))
			})

			It("Check cond", func() {
				// TODO: check how use PkgCondInvalid
				Expect((*pkg).Condition).Should(Equal(g.Condition))
			})

			It("Check struct", func() {
				// TODO: check how use PkgCondInvalid
				Expect((*pkg)).Should(Equal(g))
			})
		})

		Context("Matches version with 4 numbers on version and patch", func() {

			pkg, err := ParsePackageStr("=dev-db/oracle-instantclient-sqlplus-12.1.0.2_p1")
			g := GentooPackage{
				Name:          "oracle-instantclient-sqlplus",
				Category:      "dev-db",
				Condition:     PkgCondEqual,
				Version:       "12.1.0.2",
				Slot:          "0",
				VersionSuffix: "_p1",
				Repository:    "",
			}
			fmt.Println(fmt.Sprintf("pkg %s", pkg))

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check pkgName", func() {
				Expect((*pkg).Name).Should(Equal("oracle-instantclient-sqlplus"))
			})

			It("Check category", func() {
				Expect((*pkg).Category).Should(Equal("dev-db"))
			})

			It("Check cond", func() {
				// TODO: check how use PkgCondInvalid
				Expect((*pkg).Condition).Should(Equal(g.Condition))
			})

			It("Check struct", func() {
				// TODO: check how use PkgCondInvalid
				Expect((*pkg)).Should(Equal(g))
			})
		})

		Context("Matches version with 4 numbers", func() {

			pkg, err := ParsePackageStr("x11-libs/gtk+-2.1.0.1")
			g := GentooPackage{
				Name:       "gtk+",
				Category:   "x11-libs",
				Condition:  PkgCondEqual,
				Version:    "2.1.0.1",
				Slot:       "0",
				Repository: "",
			}
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

			It("Check cond", func() {
				// TODO: check how use PkgCondInvalid
				Expect((*pkg).Condition).Should(Equal(g.Condition))
			})

			It("Check struct", func() {
				// TODO: check how use PkgCondInvalid
				Expect((*pkg)).Should(Equal(g))
			})
		})

		Context("Matches version with 4 numbers (2)", func() {

			pkg, err := ParsePackageStr("x11-libs/gtk+-2.0.1.0")
			g := GentooPackage{
				Name:       "gtk+",
				Category:   "x11-libs",
				Condition:  PkgCondEqual,
				Version:    "2.0.1.0",
				Slot:       "0",
				Repository: "",
			}
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

			It("Check cond", func() {
				// TODO: check how use PkgCondInvalid
				Expect((*pkg).Condition).Should(Equal(g.Condition))
			})

			It("Check struct", func() {
				// TODO: check how use PkgCondInvalid
				Expect((*pkg)).Should(Equal(g))
			})
		})

		Context("Matches version with 4 numbers on version and release candidate", func() {

			pkg, err := ParsePackageStr("=dev-db/oracle-instantclient-sqlplus-12.1.0.2_rc1")
			g := GentooPackage{
				Name:          "oracle-instantclient-sqlplus",
				Category:      "dev-db",
				Condition:     PkgCondEqual,
				Version:       "12.1.0.2",
				Slot:          "0",
				VersionSuffix: "_rc1",
				Repository:    "",
			}
			fmt.Println(fmt.Sprintf("pkg %s", pkg))

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check pkgName", func() {
				Expect((*pkg).Name).Should(Equal("oracle-instantclient-sqlplus"))
			})

			It("Check category", func() {
				Expect((*pkg).Category).Should(Equal("dev-db"))
			})

			It("Check cond", func() {
				// TODO: check how use PkgCondInvalid
				Expect((*pkg).Condition).Should(Equal(g.Condition))
			})

			It("Check struct", func() {
				// TODO: check how use PkgCondInvalid
				Expect((*pkg)).Should(Equal(g))
			})
		})

		Context("Matches version with 4 numbers on version and alpha", func() {

			pkg, err := ParsePackageStr("=dev-db/oracle-instantclient-sqlplus-12.1.0.2_alpha")
			g := GentooPackage{
				Name:          "oracle-instantclient-sqlplus",
				Category:      "dev-db",
				Condition:     PkgCondEqual,
				Version:       "12.1.0.2",
				Slot:          "0",
				VersionSuffix: "_alpha",
				Repository:    "",
			}
			fmt.Println(fmt.Sprintf("pkg %s", pkg))

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check pkgName", func() {
				Expect((*pkg).Name).Should(Equal("oracle-instantclient-sqlplus"))
			})

			It("Check category", func() {
				Expect((*pkg).Category).Should(Equal("dev-db"))
			})

			It("Check cond", func() {
				// TODO: check how use PkgCondInvalid
				Expect((*pkg).Condition).Should(Equal(g.Condition))
			})

			It("Check struct", func() {
				// TODO: check how use PkgCondInvalid
				Expect((*pkg)).Should(Equal(g))
			})
		})

		Context("Matches version with 4 numbers on version and beta", func() {

			pkg, err := ParsePackageStr("=dev-db/oracle-instantclient-sqlplus-12.1.0.2_beta")
			g := GentooPackage{
				Name:          "oracle-instantclient-sqlplus",
				Category:      "dev-db",
				Condition:     PkgCondEqual,
				Version:       "12.1.0.2",
				Slot:          "0",
				VersionSuffix: "_beta",
				Repository:    "",
			}
			fmt.Println(fmt.Sprintf("pkg %s", pkg))

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check pkgName", func() {
				Expect((*pkg).Name).Should(Equal("oracle-instantclient-sqlplus"))
			})

			It("Check category", func() {
				Expect((*pkg).Category).Should(Equal("dev-db"))
			})

			It("Check cond", func() {
				// TODO: check how use PkgCondInvalid
				Expect((*pkg).Condition).Should(Equal(g.Condition))
			})

			It("Check struct", func() {
				// TODO: check how use PkgCondInvalid
				Expect((*pkg)).Should(Equal(g))
			})
		})

		Context("Check Admit() example1", func() {

			pkgA, err := ParsePackageStr("x11-libs/gtk+")
			pkgB, err := ParsePackageStr("x11-libs/gtk+-3.0.1")
			admitted, err := pkgA.Admit(pkgB)

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check Admit", func() {
				Expect(admitted).Should(Equal(true))
			})

		})

		Context("Check Admit() example2", func() {
			var pkgA, pkgB *GentooPackage
			var err error
			pkgA, err = ParsePackageStr("x11-libs/gtk+")
			pkgB, err = ParsePackageStr("www-servers/apache")
			_, err = pkgA.Admit(pkgB)

			It("Check error", func() {
				Expect(err).ShouldNot(BeNil())
			})

		})

		Context("Check Admit() example3", func() {
			var pkgA, pkgB *GentooPackage
			var err error

			pkgA, err = ParsePackageStr("x11-libs/gtk+")
			pkgB, err = ParsePackageStr("x11-libs/libX11")
			_, err = pkgA.Admit(pkgB)

			It("Check error", func() {
				Expect(err).ShouldNot(BeNil())
			})

		})

		Context("Check Admit() example4", func() {

			pkgA, err := ParsePackageStr("=x11-libs/gtk+-3.0.1")
			pkgB, err := ParsePackageStr("x11-libs/gtk+-3.0.1")
			admitted, err := pkgA.Admit(pkgB)

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check Admit", func() {
				Expect(admitted).Should(Equal(true))
			})

		})

		Context("Check Admit() example5", func() {

			pkgA, err := ParsePackageStr("x11-libs/gtk+-3.0.1")
			pkgB, err := ParsePackageStr("x11-libs/gtk+-3.0.1")
			admitted, err := pkgA.Admit(pkgB)

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check Admit", func() {
				Expect(admitted).Should(Equal(true))
			})

		})

		Context("Check Admit() example6", func() {

			pkgA, err := ParsePackageStr(">=x11-libs/gtk+-3.0.1")
			pkgB, err := ParsePackageStr("x11-libs/gtk+-3.0.1")
			admitted, err := pkgA.Admit(pkgB)

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check Admit", func() {
				Expect(admitted).Should(Equal(true))
			})

		})

		Context("Check Admit() example7", func() {

			pkgA, err := ParsePackageStr(">x11-libs/gtk+-3.0.1")
			pkgB, err := ParsePackageStr("x11-libs/gtk+-3.0.1")
			admitted, err := pkgA.Admit(pkgB)

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check Admit", func() {
				Expect(admitted).Should(Equal(false))
			})

		})

		Context("Check Admit() example8", func() {

			pkgA, err := ParsePackageStr(">x11-libs/gtk+-3.0.1")
			pkgB, err := ParsePackageStr("x11-libs/gtk+-3.0.2")
			admitted, err := pkgA.Admit(pkgB)

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check Admit", func() {
				Expect(admitted).Should(Equal(true))
			})

		})

		Context("Check Admit() example9", func() {

			pkgA, err := ParsePackageStr(">x11-libs/gtk+-3.0.1")
			pkgB, err := ParsePackageStr("x11-libs/gtk+-3.1.0")
			admitted, err := pkgA.Admit(pkgB)

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check Admit", func() {
				Expect(admitted).Should(Equal(true))
			})

		})

		Context("Check Admit() example10", func() {

			pkgA, err := ParsePackageStr(">x11-libs/gtk+-3.0.1")
			pkgB, err := ParsePackageStr("x11-libs/gtk+")
			admitted, err := pkgA.Admit(pkgB)

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check Admit", func() {
				Expect(admitted).Should(Equal(false))
			})

		})

		Context("Check Admit() example11", func() {

			pkgA, err := ParsePackageStr("x11-libs/gtk+")
			pkgB, err := ParsePackageStr("x11-libs/gtk+-1.3.4")
			admitted, err := pkgA.Admit(pkgB)

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check Admit", func() {
				Expect(admitted).Should(Equal(true))
			})

		})

		Context("Check Admit() example12", func() {

			pkgA, err := ParsePackageStr("<x11-libs/gtk+-2.0.0")
			pkgB, err := ParsePackageStr("x11-libs/gtk+")
			admitted, err := pkgA.Admit(pkgB)

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check Admit", func() {
				Expect(admitted).Should(Equal(false))
			})

		})

		Context("Check Admit() example13", func() {

			pkgA, err := ParsePackageStr("<x11-libs/gtk+-2.0.0")
			pkgB, err := ParsePackageStr("x11-libs/gtk+-1.0.0")
			admitted, err := pkgA.Admit(pkgB)

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check Admit", func() {
				Expect(admitted).Should(Equal(true))
			})

		})

		Context("Check Admit() example14", func() {

			pkgA, err := ParsePackageStr("<x11-libs/gtk+-2.0.0")
			pkgB, err := ParsePackageStr("x11-libs/gtk+-1.0.0")
			admitted, err := pkgA.Admit(pkgB)

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check Admit", func() {
				Expect(admitted).Should(Equal(true))
			})

		})

		Context("Check Admit() example15", func() {

			pkgA, err := ParsePackageStr("<=x11-libs/gtk+-2.0.0")
			pkgB, err := ParsePackageStr("x11-libs/gtk+-1.0.0")
			admitted, err := pkgA.Admit(pkgB)

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check Admit", func() {
				Expect(admitted).Should(Equal(true))
			})

		})

		Context("Check Admit() example16", func() {

			pkgA, err := ParsePackageStr("!x11-libs/gtk+-2.0.0")
			pkgB, err := ParsePackageStr("x11-libs/gtk+-1.0.0")
			admitted, err := pkgA.Admit(pkgB)

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check Admit", func() {
				Expect(admitted).Should(Equal(true))
			})

		})

		Context("Check Admit() example17", func() {

			pkgA, err := ParsePackageStr("!x11-libs/gtk+-2.0.0")
			pkgB, err := ParsePackageStr("x11-libs/gtk+-2.0.0")
			admitted, err := pkgA.Admit(pkgB)

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check Admit", func() {
				Expect(admitted).Should(Equal(false))
			})

		})

		Context("Check Admit() example18", func() {

			pkgA, err := ParsePackageStr("~x11-libs/gtk+-2.0")
			pkgB, err := ParsePackageStr("x11-libs/gtk+-2.0_rc1")
			admitted, err := pkgA.Admit(pkgB)

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check Admit", func() {
				Expect(admitted).Should(Equal(true))
			})

		})

		Context("Check Admit() example19", func() {

			pkgA, err := ParsePackageStr("~x11-libs/gtk+-2.0")
			pkgB, err := ParsePackageStr("x11-libs/gtk+-2.0.1")
			admitted, err := pkgA.Admit(pkgB)

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check Admit", func() {
				Expect(admitted).Should(Equal(false))
			})

		})

		Context("Check Admit() example20", func() {

			pkgA, err := ParsePackageStr("=x11-libs/gtk+-2.0*")
			pkgB, err := ParsePackageStr("x11-libs/gtk+-2.0.1")
			admitted, err := pkgA.Admit(pkgB)

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check Admit", func() {
				Expect(admitted).Should(Equal(true))
			})

		})

		Context("Check Admit() example21", func() {

			pkgA, err := ParsePackageStr("=x11-libs/gtk+-2.0.1*")
			pkgB, err := ParsePackageStr("x11-libs/gtk+-2.0.1.0")
			admitted, err := pkgA.Admit(pkgB)

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check Admit", func() {
				Expect(admitted).Should(Equal(true))
			})

		})

		Context("Check Admit() example22", func() {

			pkgA, err := ParsePackageStr("=x11-libs/gtk+-2.0.1*")
			pkgB, err := ParsePackageStr("x11-libs/gtk+-2.0.1-r1")
			admitted, err := pkgA.Admit(pkgB)

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check Admit", func() {
				Expect(admitted).Should(Equal(true))
			})

		})

		Context("Check Admit() example23", func() {

			pkgA, err := ParsePackageStr("=x11-libs/gtk+-2.0.1")
			pkgB, err := ParsePackageStr("x11-libs/gtk+-2.0.1-r1")
			admitted, err := pkgA.Admit(pkgB)

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check Admit", func() {
				Expect(admitted).Should(Equal(false))
			})

		})

		Context("Test go-version - example1", func() {
			v1, err := version.NewVersion("1.1.1.1")
			It("Check error", func() {
				Expect(err).Should(BeNil())
			})
			fmt.Println("VERSION = ", v1)
		})

		Context("Test go-version - example2", func() {
			v1, err := version.NewVersion("1.1.1.1")
			It("Check error", func() {
				Expect(err).Should(BeNil())
			})
			v2, err := version.NewVersion("1.1.1.2")
			It("Check error V2", func() {
				Expect(err).Should(BeNil())
			})

			res := v1.LessThan(v2)
			It("Check result", func() {
				Expect(res).Should(Equal(true))
			})

			res = v2.GreaterThanOrEqual(v1)
			It("Check result2", func() {
				Expect(res).Should(Equal(true))
			})
		})

		Context("Test go-version - example3", func() {
			v1, err := version.NewVersion("1.1.2_rc1")
			// NOTE: err SHOLD be with a value instead is null
			It("Check error", func() {
				Expect(err).Should(BeNil())
			})
			v2, err := version.NewVersion("1.1.1.2-alpha")
			It("Check error V2", func() {
				Expect(err).Should(BeNil())
			})
			fmt.Println("V1 = ", v1, err)
			fmt.Println("V2 = ", v2, err)

		})

		Context("Test go-version - example4", func() {
			v1, err := version.NewVersion("1.1.2-rc1")
			// NOTE: err SHOLD be with a value instead is null
			It("Check error", func() {
				Expect(err).Should(BeNil())
			})
			v2, err := version.NewVersion("1.1.1.2")
			It("Check error V2", func() {
				Expect(err).Should(BeNil())
			})
			v3, err := version.NewVersion("1.1.2")
			It("Check error V3", func() {
				Expect(err).Should(BeNil())
			})
			fmt.Println("V1 = ", v1, err)
			fmt.Println("V2 = ", v2, err)
			fmt.Println("V3 = ", v3, err)
			versions := make([]*version.Version, 3)
			versions[0] = v1
			versions[1] = v2
			versions[2] = v3
			sort.Sort(version.Collection(versions))
			fmt.Println("SORT = ", versions)

		})

		Context("GetPackageName", func() {
			gp, err := ParsePackageStr("sys-base/gcc-8.2.0")
			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check package name", func() {
				Expect(gp.GetPackageName()).Should(Equal("sys-base/gcc"))
			})
		})
	})

})
