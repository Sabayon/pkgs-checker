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

package hash

import (
	"archive/tar"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"sort"

	commons "github.com/Sabayon/pkgs-checker/pkg/commons"
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
	if pkg == "" {
		return nil, errors.New("Invalid pkg value")
	}

	log := l
	if log == nil {
		log = logger.StandardLogger()
	}

	return &Package{
		pkg:     pkg,
		skipped: 0,
		logger:  log,
		files:   make(map[string][]byte),
	}, nil
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
		`Package: { pkg: %s, basename: %s, abspath: %s, checksum: %s, files: %d, dirs: %d, skipped: %d }`,
		p.pkg,
		p.basename,
		p.abspath,
		p.checksum,
		len(p.files),
		len(p.dirs),
		p.skipped,
	)
}

func (p *Package) ProcessTarFile(tarReader *tar.Reader, name string) error {
	buf := make([]byte, BYTE_BUFFER_LEN)
	fmd5 := md5.New()

	for {
		nBytes, err := tarReader.Read(buf)

		if nBytes > 0 {
			if _, writeErr := fmd5.Write(buf[:nBytes]); writeErr != nil {
				return writeErr
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}
	}

	fileHash := fmd5.Sum(nil)

	p.logger.Debugf(
		"[%s] %s - MD5 (without dirs): %s",
		p.basename,
		name,
		hex.EncodeToString(fileHash),
	)

	p.AddFile(name, fileHash)

	return nil
}

func (p *Package) CalculateCRC() error {
	var pmd5 hash.Hash = md5.New()

	if len(p.files) == 0 && len(p.dirs) == 0 {
		p.logger.Warnf("[%s] No directories or files found for CRC.\n", p.pkg)
		p.checksum = commons.PKGS_CHECKER_EMPTY_PKGHASH
		return nil
	}

	if len(p.files) > 0 {
		files := make([]string, 0, len(p.files))

		for k := range p.files {
			files = append(files, k)
		}

		sort.Strings(files)

		for _, k := range files {
			if _, err := pmd5.Write(p.files[k]); err != nil {
				return err
			}
		}
	}

	if len(p.dirs) > 0 {
		sort.Strings(p.dirs)

		for _, k := range p.dirs {
			if _, err := pmd5.Write([]byte(k)); err != nil {
				return err
			}
		}
	}

	p.checksum = hex.EncodeToString(pmd5.Sum(nil))

	return nil
}

func (p PackageSorter) Len() int {
	return len(p)
}

func (p PackageSorter) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p PackageSorter) Less(i, j int) bool {
	return p[i].pkg < p[j].pkg
}
