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

		// https://devmanual.gentoo.org/general-concepts/dependencies/
		Context("Parse dependency1", func() {

			pkg, err := ParsePackageStr(">=sys-libs/ncurses-5.2-r5:0=")
			g := GentooPackage{
				Name:          "ncurses",
				Category:      "sys-libs",
				Condition:     PkgCondGreaterEqual,
				Slot:          "0=",
				Version:       "5.2",
				VersionSuffix: "-r5",
				VersionBuild:  "",
				Repository:    "",
			}
			fmt.Println(fmt.Sprintf("pkg %s", pkg))

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check pkgName", func() {
				Expect((*pkg).Name).Should(Equal("ncurses"))
			})

			It("Check category", func() {
				Expect((*pkg).Category).Should(Equal("sys-libs"))
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

		Context("Parse dependency2", func() {

			pkg, err := ParsePackageStr(">=sys-libs/ncurses-5.2-r5:0*")
			g := GentooPackage{
				Name:          "ncurses",
				Category:      "sys-libs",
				Condition:     PkgCondGreaterEqual,
				Slot:          "0*",
				Version:       "5.2",
				VersionSuffix: "-r5",
				VersionBuild:  "",
				Repository:    "",
			}
			fmt.Println(fmt.Sprintf("pkg %s", pkg))

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check pkgName", func() {
				Expect((*pkg).Name).Should(Equal("ncurses"))
			})

			It("Check category", func() {
				Expect((*pkg).Category).Should(Equal("sys-libs"))
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

		Context("Parse dependency3", func() {

			pkg, err := ParsePackageStr(">=sys-libs/ncurses-5.2-r5:*")
			g := GentooPackage{
				Name:          "ncurses",
				Category:      "sys-libs",
				Condition:     PkgCondGreaterEqual,
				Slot:          "*",
				Version:       "5.2",
				VersionSuffix: "-r5",
				VersionBuild:  "",
				Repository:    "",
			}
			fmt.Println(fmt.Sprintf("pkg %s", pkg))

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check pkgName", func() {
				Expect((*pkg).Name).Should(Equal("ncurses"))
			})

			It("Check category", func() {
				Expect((*pkg).Category).Should(Equal("sys-libs"))
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

		Context("Parse dependency4", func() {

			pkg, err := ParsePackageStr("app-cdr/cdrtools-3.02_alpha09-r2")
			g := GentooPackage{
				Name:          "cdrtools",
				Category:      "app-cdr",
				Condition:     PkgCondEqual,
				Slot:          "0",
				Version:       "3.02",
				VersionSuffix: "_alpha09-r2",
				VersionBuild:  "",
				Repository:    "",
			}
			fmt.Println(fmt.Sprintf("pkg %s", pkg))

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check pkgName", func() {
				Expect((*pkg).Name).Should(Equal("cdrtools"))
			})

			It("Check category", func() {
				Expect((*pkg).Category).Should(Equal("app-cdr"))
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
				Name:         "gdb",
				Category:     "sys-devel",
				Condition:    PkgCondAnyRevision,
				Version:      "7.3",
				VersionBuild: "",
				Slot:         "0",
				Repository:   "",
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
				Name:         "gdb",
				Category:     "sys-devel",
				Condition:    PkgCondMatchVersion,
				Version:      "7.3",
				VersionBuild: "",
				Slot:         "0",
				Repository:   "",
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

		Context("Matches version with 4 numbers and build version", func() {

			pkg, err := ParsePackageStr("=dev-db/database-release-manager-0.1.0.1+AB")
			g := GentooPackage{
				Name:          "database-release-manager",
				Category:      "dev-db",
				Condition:     PkgCondEqual,
				Version:       "0.1.0.1",
				Slot:          "0",
				VersionSuffix: "",
				VersionBuild:  "AB",
				Repository:    "",
			}
			fmt.Println(fmt.Sprintf("pkg %s", pkg))

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check pkgName", func() {
				Expect((*pkg).Name).Should(Equal("database-release-manager"))
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

		Context("Matches version with 2 numbers and build version", func() {

			pkg, err := ParsePackageStr("=app-misc/c_rehash-1.7+r1")
			g := GentooPackage{
				Name:          "c_rehash",
				Category:      "app-misc",
				Condition:     PkgCondEqual,
				Version:       "1.7",
				Slot:          "0",
				VersionSuffix: "",
				VersionBuild:  "r1",
				Repository:    "",
			}
			fmt.Println(fmt.Sprintf("pkg %s", pkg))

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check pkgName", func() {
				Expect((*pkg).Name).Should(Equal("c_rehash"))
			})

			It("Check category", func() {
				Expect((*pkg).Category).Should(Equal("app-misc"))
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

		Context("Matches version with dot and numbers on pkgname and build version", func() {

			pkg, err := ParsePackageStr("=app-misc/geoclue-2.0-2.5.3+r2")
			g := GentooPackage{
				Name:          "geoclue-2.0",
				Category:      "app-misc",
				Condition:     PkgCondEqual,
				Version:       "2.5.3",
				Slot:          "0",
				VersionSuffix: "",
				VersionBuild:  "r2",
				Repository:    "",
			}
			fmt.Println(fmt.Sprintf("pkg %s", pkg))

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check pkgName", func() {
				Expect((*pkg).Name).Should(Equal("geoclue-2.0"))
			})

			It("Check category", func() {
				Expect((*pkg).Category).Should(Equal("app-misc"))
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

		Context("Matches version with pkgname that start with upper chars and build version", func() {

			pkg, err := ParsePackageStr("=dev-perl/WWW-RobotRules-6.20.0+r1")
			g := GentooPackage{
				Name:          "WWW-RobotRules",
				Category:      "dev-perl",
				Condition:     PkgCondEqual,
				Version:       "6.20.0",
				Slot:          "0",
				VersionSuffix: "",
				VersionBuild:  "r1",
				Repository:    "",
			}
			fmt.Println(fmt.Sprintf("pkg %s", pkg))

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check pkgName", func() {
				Expect((*pkg).Name).Should(Equal("WWW-RobotRules"))
			})

			It("Check category", func() {
				Expect((*pkg).Category).Should(Equal("dev-perl"))
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

		Context("Matches version with 3 numbers and build version with chars and number", func() {

			pkg, err := ParsePackageStr("=dev-db/mysql-8.1.0+0.dev")
			g := GentooPackage{
				Name:          "mysql",
				Category:      "dev-db",
				Condition:     PkgCondEqual,
				Version:       "8.1.0",
				Slot:          "0",
				VersionSuffix: "",
				VersionBuild:  "0.dev",
				Repository:    "",
			}
			fmt.Println(fmt.Sprintf("pkg %s", pkg))

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check pkgName", func() {
				Expect((*pkg).Name).Should(Equal("mysql"))
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

		Context("Parse dep 5", func() {

			pkg, err := ParsePackageStr("app/A-1.0+1")
			g := GentooPackage{
				Name:          "A",
				Category:      "app",
				Condition:     PkgCondEqual,
				Slot:          "0",
				Version:       "1.0",
				VersionSuffix: "",
				VersionBuild:  "1",
				Repository:    "",
			}
			fmt.Println(fmt.Sprintf("pkg %s", pkg))

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check pkgName", func() {
				Expect((*pkg).Name).Should(Equal("A"))
			})

			It("Check category", func() {
				Expect((*pkg).Category).Should(Equal("app"))
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

		Context("Parse dep 6", func() {

			pkg, err := ParsePackageStr("app/A-1.0+pre20200315.1")
			g := GentooPackage{
				Name:          "A",
				Category:      "app",
				Condition:     PkgCondEqual,
				Slot:          "0",
				Version:       "1.0",
				VersionSuffix: "",
				VersionBuild:  "pre20200315.1",
				Repository:    "",
			}
			fmt.Println(fmt.Sprintf("pkg %s", pkg))

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check pkgName", func() {
				Expect((*pkg).Name).Should(Equal("A"))
			})

			It("Check category", func() {
				Expect((*pkg).Category).Should(Equal("app"))
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

		Context("Parse dep 7", func() {

			pkg, err := ParsePackageStr("app/A-1.0_pre20200315+d1.1")
			g := GentooPackage{
				Name:          "A",
				Category:      "app",
				Condition:     PkgCondEqual,
				Slot:          "0",
				Version:       "1.0",
				VersionSuffix: "_pre20200315",
				VersionBuild:  "d1.1",
				Repository:    "",
			}
			fmt.Println(fmt.Sprintf("pkg %s", pkg))

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check pkgName", func() {
				Expect((*pkg).Name).Should(Equal("A"))
			})

			It("Check category", func() {
				Expect((*pkg).Category).Should(Equal("app"))
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

		Context("Parse dep 8", func() {

			pkg, err := ParsePackageStr("app/A-1.0_pre20200315")
			g := GentooPackage{
				Name:          "A",
				Category:      "app",
				Condition:     PkgCondEqual,
				Slot:          "0",
				Version:       "1.0",
				VersionSuffix: "_pre20200315",
				VersionBuild:  "",
				Repository:    "",
			}
			fmt.Println(fmt.Sprintf("pkg %s", pkg))

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check pkgName", func() {
				Expect((*pkg).Name).Should(Equal("A"))
			})

			It("Check category", func() {
				Expect((*pkg).Category).Should(Equal("app"))
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

		Context("Parse dep 9", func() {

			pkg, err := ParsePackageStr(">=dev-libs/libsigc++-2-2.3.2")
			g := GentooPackage{
				Name:          "libsigc++-2",
				Category:      "dev-libs",
				Condition:     PkgCondGreaterEqual,
				Slot:          "0",
				Version:       "2.3.2",
				VersionSuffix: "",
				VersionBuild:  "",
				Repository:    "",
			}
			fmt.Println(fmt.Sprintf("pkg %s", pkg))

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check pkgName", func() {
				Expect((*pkg).Name).Should(Equal("libsigc++-2"))
			})

			It("Check category", func() {
				Expect((*pkg).Category).Should(Equal("dev-libs"))
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

		Context("Parse dep 10", func() {

			pkg, err := ParsePackageStr(">=media-libs/libsndfile-1.0.29+pre2_p20191024.1")
			g := GentooPackage{
				Name:          "libsndfile",
				Category:      "media-libs",
				Condition:     PkgCondGreaterEqual,
				Slot:          "0",
				Version:       "1.0.29",
				VersionSuffix: "",
				VersionBuild:  "pre2_p20191024.1",
				Repository:    "",
			}
			fmt.Println(fmt.Sprintf("pkg %s", pkg))

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check pkgName", func() {
				Expect((*pkg).Name).Should(Equal("libsndfile"))
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

		Context("Parse dep 11", func() {

			pkg, err := ParsePackageStr(">=dev-libs/libsigc++-2-2.3.2+1")
			g := GentooPackage{
				Name:          "libsigc++-2",
				Category:      "dev-libs",
				Condition:     PkgCondGreaterEqual,
				Slot:          "0",
				Version:       "2.3.2",
				VersionSuffix: "",
				VersionBuild:  "1",
				Repository:    "",
			}
			fmt.Println(fmt.Sprintf("pkg %s", pkg))

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check pkgName", func() {
				Expect((*pkg).Name).Should(Equal("libsigc++-2"))
			})

			It("Check category", func() {
				Expect((*pkg).Category).Should(Equal("dev-libs"))
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

		Context("Parse dep 12", func() {

			pkg, err := ParsePackageStr(">=dev-libs/dbus-c++-0.9.0+r3")
			g := GentooPackage{
				Name:          "dbus-c++",
				Category:      "dev-libs",
				Condition:     PkgCondGreaterEqual,
				Slot:          "0",
				Version:       "0.9.0",
				VersionSuffix: "",
				VersionBuild:  "r3",
				Repository:    "",
			}
			fmt.Println(fmt.Sprintf("pkg %s", pkg))

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check pkgName", func() {
				Expect((*pkg).Name).Should(Equal("dbus-c++"))
			})

			It("Check category", func() {
				Expect((*pkg).Category).Should(Equal("dev-libs"))
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

		Context("Parse font package1", func() {

			pkg, err := ParsePackageStr(">=media-fonts/font-bitstream-100dpi-1.0.3-r2")
			g := GentooPackage{
				Name:          "font-bitstream-100dpi",
				Category:      "media-fonts",
				Condition:     PkgCondGreaterEqual,
				Slot:          "0",
				Version:       "1.0.3",
				VersionSuffix: "-r2",
				VersionBuild:  "",
				Repository:    "",
			}

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check pkgName", func() {
				Expect((*pkg).Name).Should(Equal("font-bitstream-100dpi"))
			})

			It("Check category", func() {
				Expect((*pkg).Category).Should(Equal("media-fonts"))
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

		Context("Parse font package2", func() {

			pkg, err := ParsePackageStr(">=media-fonts/font-bitstream-100dpi")
			g := GentooPackage{
				Name:          "font-bitstream-100dpi",
				Category:      "media-fonts",
				Condition:     PkgCondGreaterEqual,
				Slot:          "0",
				Version:       "",
				VersionSuffix: "",
				VersionBuild:  "",
				Repository:    "",
			}

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check pkgName", func() {
				Expect((*pkg).Name).Should(Equal("font-bitstream-100dpi"))
			})

			It("Check category", func() {
				Expect((*pkg).Category).Should(Equal("media-fonts"))
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

		Context("Parse dep with use flags", func() {

			pkg, err := ParsePackageStr(">=dev-libs/dbus-9.9.9[use1,use2,use3(+)]")
			g := GentooPackage{
				Name:          "dbus",
				Category:      "dev-libs",
				Condition:     PkgCondGreaterEqual,
				Slot:          "0",
				Version:       "9.9.9",
				VersionSuffix: "",
				VersionBuild:  "",
				Repository:    "",
				UseFlags: []string{
					"use1", "use2", "use3(+)",
				},
			}
			fmt.Println(fmt.Sprintf("pkg %s", pkg))

			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check pkgName", func() {
				Expect((*pkg).Name).Should(Equal("dbus"))
			})

			It("Check category", func() {
				Expect((*pkg).Category).Should(Equal("dev-libs"))
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

		Context("GetPackageName2", func() {
			gp, err := ParsePackageStr("app-arch/rpm2targz-9.0.0.5g")
			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check package name", func() {
				Expect(gp.GetPackageName()).Should(Equal("app-arch/rpm2targz"))
			})

			It("Check package version", func() {
				Expect(gp.Version).Should(Equal("9.0.0.5g"))
			})
		})

		Context("GetPackageName3", func() {
			gp, err := ParsePackageStr("dev-lang/spidermonkey-52.9.1_pre1")
			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check package name", func() {
				Expect(gp.GetPackageName()).Should(Equal("dev-lang/spidermonkey"))
			})

			It("Check package version", func() {
				Expect(gp.Version).Should(Equal("52.9.1"))
			})

			It("Check package version suffix", func() {
				Expect(gp.VersionSuffix).Should(Equal("_pre1"))
			})
		})

		Context("GetPackageName4", func() {
			gp, err := ParsePackageStr("sys-libs/timezone-data-2018i")
			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check package name", func() {
				Expect(gp.GetPackageName()).Should(Equal("sys-libs/timezone-data"))
			})

			It("Check package version", func() {
				Expect(gp.Version).Should(Equal("2018i"))
			})

			It("Check package version suffix", func() {
				Expect(gp.VersionSuffix).Should(Equal(""))
			})
		})

		Context("GetPackageName5", func() {
			gp, err := ParsePackageStr("net-im/zoom-bin-2.8.222599.0519")
			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check package name", func() {
				Expect(gp.GetPackageName()).Should(Equal("net-im/zoom-bin"))
			})

			It("Check package version", func() {
				Expect(gp.Version).Should(Equal("2.8.222599.0519"))
			})

			It("Check package version suffix", func() {
				Expect(gp.VersionSuffix).Should(Equal(""))
			})
		})

		Context("GetPackageName6", func() {
			gp, err := ParsePackageStr("dev-util/idea-community-2019.2.0.192.5728.98")
			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check package name", func() {
				Expect(gp.GetPackageName()).Should(Equal("dev-util/idea-community"))
			})

			It("Check package version", func() {
				Expect(gp.Version).Should(Equal("2019.2.0.192.5728.98"))
			})

			It("Check package version suffix", func() {
				Expect(gp.VersionSuffix).Should(Equal(""))
			})
		})

		Context("GetPackageName6", func() {
			gp, err := ParsePackageStr("kernel-5.4/sabayon-full-5.4+1")
			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check package name", func() {
				Expect(gp.GetPackageName()).Should(Equal("kernel-5.4/sabayon-full"))
			})

			It("Check package version", func() {
				Expect(gp.Version).Should(Equal("5.4"))
			})

			It("Check package version suffix", func() {
				Expect(gp.VersionBuild).Should(Equal("1"))
			})
		})
		Context("PkgWithUseFlags", func() {
			gp, err := ParsePackageStr("dev-util/mottainai-agent-0.0_pre20191012[lxd]")
			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check package name", func() {
				Expect(gp.GetPackageName()).Should(Equal("dev-util/mottainai-agent"))
			})

			It("Check package version", func() {
				Expect(gp.Version).Should(Equal("0.0"))
			})

			It("Check package version suffix", func() {
				Expect(gp.VersionSuffix).Should(Equal("_pre20191012"))
			})

			It("Check package use flags", func() {
				Expect(gp.UseFlags).Should(Equal([]string{"lxd"}))
			})
		})

		Context("PkgWithUseFlags2", func() {
			gp, err := ParsePackageStr("dev-util/mottainai-agent-0.0_pre20191012[lxd,zfs,3dnow]")
			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check package name", func() {
				Expect(gp.GetPackageName()).Should(Equal("dev-util/mottainai-agent"))
			})

			It("Check package version", func() {
				Expect(gp.Version).Should(Equal("0.0"))
			})

			It("Check package version suffix", func() {
				Expect(gp.VersionSuffix).Should(Equal("_pre20191012"))
			})

			It("Check package use flags", func() {
				Expect(gp.UseFlags).Should(Equal([]string{"lxd", "zfs", "3dnow"}))
			})
		})

		Context("PkgWithUseFlags3", func() {
			gp, err := ParsePackageStr("dev-util/mottainai-agent[lxd,zfs,3dnow]")
			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check package name", func() {
				Expect(gp.GetPackageName()).Should(Equal("dev-util/mottainai-agent"))
			})

			It("Check package use flags", func() {
				Expect(gp.UseFlags).Should(Equal([]string{"lxd", "zfs", "3dnow"}))
			})
		})

		Context("PerlVersion1", func() {
			gp, err := ParsePackageStr("virtual/perl-Storable-3.80.100_rc")
			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check package name", func() {
				Expect(gp.GetPackageName()).Should(Equal("virtual/perl-Storable"))
			})

			It("Check package version", func() {
				Expect(gp.Version).Should(Equal("3.80.100"))
			})

			It("Check package version suffix", func() {
				Expect(gp.VersionSuffix).Should(Equal("_rc"))
			})
		})

		Context("ConfrontVersions", func() {
			gp, err := ParsePackageStr(">=sys-power-5.3/acpi_call-3.17.5.3.2.1")
			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check package name", func() {
				Expect(gp.GetPackageName()).Should(Equal("sys-power-5.3/acpi_call"))
			})

			It("Check package category", func() {
				Expect(gp.Category).Should(Equal("sys-power-5.3"))
			})
			It("Check package version", func() {
				Expect(gp.Version).Should(Equal("3.17.5.3.2.1"))
			})

			gp2, err := ParsePackageStr("sys-power-5.3/acpi_call-3.17.5.3.10")
			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check package name", func() {
				Expect(gp2.GetPackageName()).Should(Equal("sys-power-5.3/acpi_call"))
			})

			It("Check package category", func() {
				Expect(gp2.Category).Should(Equal("sys-power-5.3"))
			})
			It("Check package version", func() {
				Expect(gp2.Version).Should(Equal("3.17.5.3.10"))
			})

			It("Check Admit", func() {
				Expect(gp.Admit(gp2)).Should(Equal(true))
			})

			gp3, err := ParsePackageStr("sys-power-5.3/acpi_call-3.17.5.4")
			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check package name", func() {
				Expect(gp3.GetPackageName()).Should(Equal("sys-power-5.3/acpi_call"))
			})

			It("Check package category", func() {
				Expect(gp3.Category).Should(Equal("sys-power-5.3"))
			})
			It("Check package version", func() {
				Expect(gp3.Version).Should(Equal("3.17.5.4"))
			})

			It("Check Admit", func() {
				Expect(gp.Admit(gp3)).Should(Equal(true))
			})
		})

		Context("ConfrontVersions2", func() {
			gp, err := ParsePackageStr(">=net-vpn-4.9/wireguard-0.0.1")
			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check package name", func() {
				Expect(gp.GetPackageName()).Should(Equal("net-vpn-4.9/wireguard"))
			})

			It("Check package category", func() {
				Expect(gp.Category).Should(Equal("net-vpn-4.9"))
			})
			It("Check package version", func() {
				Expect(gp.Version).Should(Equal("0.0.1"))
			})

			gp2, err := ParsePackageStr("net-vpn-4.9/wireguard-0.0.20190406.4.9.172-r1")
			It("Check error", func() {
				Expect(err).Should(BeNil())
			})

			It("Check package name", func() {
				Expect(gp2.GetPackageName()).Should(Equal("net-vpn-4.9/wireguard"))
			})

			It("Check package category", func() {
				Expect(gp2.Category).Should(Equal("net-vpn-4.9"))
			})
			It("Check package version", func() {
				Expect(gp2.Version).Should(Equal("0.0.20190406.4.9.172"))
			})

			It("Check Admit", func() {
				Expect(gp.Admit(gp2)).Should(Equal(true))
			})

		})

	})

	Context("Check Condition2Int", func() {
		gp, err := ParsePackageStr(">=net-vpn-4.9/wireguard-0.0.1")
		It("Check error", func() {
			Expect(err).Should(BeNil())
		})
		It("Check Int", func() {
			Expect(gp.Condition.Int()).Should(Equal(PkgCondGreaterEqual))
		})
	})

	Context("Check Package sorter", func() {
		gp1, err := ParsePackageStr("net-vpn/wireguard-0.6.0")
		It("Check error", func() {
			Expect(err).Should(BeNil())
		})
		gp2, err := ParsePackageStr("net-vpn/wireguard-0.1.0")
		It("Check error", func() {
			Expect(err).Should(BeNil())
		})
		gp3, err := ParsePackageStr("net-vpn/wireguard-0.4.0")
		It("Check error", func() {
			Expect(err).Should(BeNil())
		})

		pkgs := []GentooPackage{*gp1, *gp2, *gp3}

		sort.Sort(GentooPackageSorter(pkgs))

		It("Check order", func() {
			Expect(pkgs[0]).Should(Equal(*gp2))
			Expect(pkgs[1]).Should(Equal(*gp3))
			Expect(pkgs[2]).Should(Equal(*gp1))
		})
	})

	Context("Check Package sorter2", func() {
		gp1, err := ParsePackageStr("net-vpn/wireguard-0.6.0")
		It("Check error", func() {
			Expect(err).Should(BeNil())
		})
		gp2, err := ParsePackageStr("net-vpn/wireguard-0.6.0-r1")
		It("Check error", func() {
			Expect(err).Should(BeNil())
		})
		gp3, err := ParsePackageStr("net-vpn/wireguard-0.4.0")
		It("Check error", func() {
			Expect(err).Should(BeNil())
		})

		pkgs := []GentooPackage{*gp2, *gp3, *gp1}

		sort.Sort(GentooPackageSorter(pkgs))

		It("Check order", func() {
			Expect(pkgs[0]).Should(Equal(*gp3))
			Expect(pkgs[1]).Should(Equal(*gp1))
			Expect(pkgs[2]).Should(Equal(*gp2))
		})
	})

	Context("Check Package sorter3", func() {
		gp1, err := ParsePackageStr("net-vpn/wireguard-0.6.0+5")
		It("Check error", func() {
			Expect(err).Should(BeNil())
		})
		gp2, err := ParsePackageStr("net-vpn/wireguard-0.6.0+1")
		It("Check error", func() {
			Expect(err).Should(BeNil())
		})
		gp3, err := ParsePackageStr("net-vpn/wireguard-0.4.0")
		It("Check error", func() {
			Expect(err).Should(BeNil())
		})

		pkgs := []GentooPackage{*gp2, *gp3, *gp1}

		sort.Sort(GentooPackageSorter(pkgs))

		It("Check order", func() {
			Expect(pkgs[0]).Should(Equal(*gp3))
			Expect(pkgs[1]).Should(Equal(*gp2))
			Expect(pkgs[2]).Should(Equal(*gp1))
		})
	})

})
