/*

Copyright (C) 2017-2018  Daniele Rondina <geaaru@sabayonlinux.org>

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
	"sort"

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
		return errors.New("No dirs and files found for CRC")
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
