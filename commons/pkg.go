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
package commons

import (
	"archive/tar"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"regexp"
	"sort"
	"strings"

	logger "github.com/sirupsen/logrus"
)

const BYTE_BUFFER_LEN = 100

type Package struct {
	pkg      string
	basename string
	abspath  string
	checksum string
	files    map[string][]byte
	dirs     []string
	skipped  int
	logger   *logger.Logger
}

type PackageSorter []Package

func NewPackage(pkg string, l *logger.Logger) (*Package, error) {

	var log *logger.Logger = nil

	if pkg == "" {
		return nil, errors.New("Invalid pkg value")
	}
	if l == nil {
		// Use standard logger
		log = logger.StandardLogger()
	} else {
		log = l
	}

	return &Package{pkg: pkg, skipped: 0,
		logger: log, files: make(map[string][]byte)}, nil
}

func (p *Package) AddFile(f string, hash []byte) {
	p.files[f] = hash
}

func (p *Package) AddDir(d string) {
	p.dirs = append(p.dirs, d)
}

func (p *Package) Name() string {
	return p.pkg
}

func (p *Package) CheckSum() string {
	return p.checksum
}

func (p *Package) String() string {
	return fmt.Sprintf(
		`Package: {
	pkg: %s,
	basename: %s,
	abspath: %s,
	checksum: %s,
	files: %d,
	dirs: %d,
	skipped: %d
}`,
		p.pkg, p.basename, p.abspath, p.checksum, len(p.files), len(p.dirs), p.skipped)
}

func (p *Package) ProcessTarFile(tarReader *tar.Reader, name string) error {

	// Read file
	var err error
	var buf []byte = make([]byte, BYTE_BUFFER_LEN)
	var pb []byte
	var n_bytes int
	var fmd5 hash.Hash = md5.New()

	for {
		n_bytes, err = tarReader.Read(buf)

		if n_bytes > 0 {
			//
			pb = buf[0:n_bytes]
			fmd5.Write(pb)
		}

		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			return err
		}

	}

	var h []byte = fmd5.Sum(nil)
	p.logger.Debugf("[%s] %s - MD5 (without dirs): %s",
		p.basename, name, hex.EncodeToString(h))

	p.AddFile(name, h)

	return err
}

func (p *Package) CalculateCRC() error {

	var pmd5 hash.Hash = md5.New()

	if len(p.files) == 0 && len(p.dirs) == 0 {
		p.logger.Warnf("[%s] No directories or files found for CRC.\n", p.pkg)
		p.checksum = PKGS_CHECKER_EMPTY_PKGHASH
		return nil
	}

	// Create MD5 with MD5 of all files
	if len(p.files) > 0 {
		var files []string

		for k, _ := range p.files {
			files = append(files, k)
		}

		sort.Strings(files)
		for _, k := range files {
			pmd5.Write(p.files[k])
		}
	}

	// Append to MD5 list of directories
	if len(p.dirs) > 0 {

		sort.Strings(p.dirs)
		for _, k := range p.dirs {
			pmd5.Write([]byte(k))
		}
	}

	var h []byte = pmd5.Sum(nil)
	p.checksum = hex.EncodeToString(h)

	return nil
}

func (p PackageSorter) Len() int           { return len(p) }
func (p PackageSorter) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PackageSorter) Less(i, j int) bool { return p[i].pkg < p[j].pkg }

// ----------------------------------
// Code to move and merge inside luet project
// ----------------------------------

// Package condition
type PackageCond int

const (
	PkgCondInvalid = 0
	// >
	PkgCondGreater = 1
	// >=
	PkgCondGreaterEqual = 2
	// <
	PkgCondLess = 3
	// <=
	PkgCondLessEqual = 4
	// =
	PkgCondEqual = 5
	// !
	PkgCondNot = 6
	// ~
	PkgCondAnyRevision = 7
	// =<pkg>*
	PkgCondMatchVersion = 8
)

type GentooPackage struct {
	Name          string
	Category      string
	Version       string
	VersionSuffix string
	Slot          string
	Condition     PackageCond
	Repository    string
}

func (p *GentooPackage) String() string {
	// TODO
	return fmt.Sprintf("%s/%s", p.Category, p.Name)
}

// return category, package, version, slot, condition
func ParsePackageStr(pkg string) (*GentooPackage, error) {
	if pkg == "" {
		return nil, errors.New("Invalid package string")
	}

	ans := GentooPackage{
		Slot:      "0",
		Condition: PkgCondInvalid,
	}

	if strings.HasPrefix(pkg, ">=") {
		pkg = pkg[2:]
		ans.Condition = PkgCondGreaterEqual
	} else if strings.HasPrefix(pkg, ">") {
		pkg = pkg[1:]
		ans.Condition = PkgCondGreater
	} else if strings.HasPrefix(pkg, "<=") {
		pkg = pkg[2:]
		ans.Condition = PkgCondLessEqual
	} else if strings.HasPrefix(pkg, "<") {
		pkg = pkg[1:]
		ans.Condition = PkgCondLess
	} else if strings.HasPrefix(pkg, "=") {
		pkg = pkg[1:]
		if strings.HasSuffix(pkg, "*") {
			ans.Condition = PkgCondMatchVersion
			pkg = pkg[0 : len(pkg)-1]
		} else {
			ans.Condition = PkgCondEqual
		}
	} else if strings.HasPrefix(pkg, "~") {
		pkg = pkg[1:]
		ans.Condition = PkgCondAnyRevision
	}

	words := strings.Split(pkg, "/")
	if len(words) != 2 {
		return nil, errors.New(fmt.Sprintf("Invalid package string %s", pkg))
	}
	ans.Category = words[0]
	pkgname := words[1]

	// Check if has repository
	if strings.Contains(pkgname, "::") {
		words = strings.Split(pkgname, "::")
		ans.Repository = words[1]
		pkgname = words[0]
	}

	// Check if has slot
	if strings.Contains(pkgname, ":") {
		words = strings.Split(pkgname, ":")
		ans.Slot = words[1]
		pkgname = words[0]
	}

	regexPkg := regexp.MustCompile(
		`([0-9]+[.][0-9]+|[0-9]+|[0-9]+[.][0-9]+[.][0-9]+|[0-9]+[.][0-9]+[.][0.9]+[.][0-9]+)(_p[0-9]+|_pre|_rc[0-9]+|_alpha|_beta)*$`,
	)
	matches := regexPkg.FindAllString(pkgname, -1)

	// NOTE: Now suffix comples like _alpha_rc1 are not supported.

	if len(matches) > 0 {
		// Check if there patch
		if strings.Contains(matches[0], "_p") {
			ans.Version = matches[0][:strings.Index(matches[0], "_p")]
			ans.VersionSuffix = matches[0][strings.Index(matches[0], "_p"):]
		} else if strings.Contains(matches[0], "_rc") {
			ans.Version = matches[0][:strings.Index(matches[0], "_rc")]
			ans.VersionSuffix = matches[0][strings.Index(matches[0], "_rc"):]
		} else if strings.Contains(matches[0], "_alpha") {
			ans.Version = matches[0][:strings.Index(matches[0], "_alpha")]
			ans.VersionSuffix = matches[0][strings.Index(matches[0], "_alpha"):]
		} else if strings.Contains(matches[0], "_beta") {
			ans.Version = matches[0][:strings.Index(matches[0], "_beta")]
			ans.VersionSuffix = matches[0][strings.Index(matches[0], "_beta"):]
		} else {
			ans.Version = matches[0]
		}
		ans.Name = pkgname[0 : len(pkgname)-len(ans.Version)-1-len(ans.VersionSuffix)]
	} else {
		ans.Name = pkgname
	}

	return &ans, nil
}
