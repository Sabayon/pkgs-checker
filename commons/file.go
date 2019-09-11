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
	"errors"
	"path/filepath"
)

func AbsPathFromBase(basepath, resourcepath string) (string, error) {
	if basepath == "" {
		return "", errors.New("Invalid basepath")
	}
	if resourcepath == "" {
		return "", errors.New("Invalid resourcepath")
	}

	if filepath.IsAbs(resourcepath) {
		// Nothing to do.
		return resourcepath, nil
	}

	var err error
	if !filepath.IsAbs(basepath) {
		basepath, err = filepath.Abs(basepath)
		if err != nil {
			return "", err
		}
	}

	return filepath.Clean(filepath.Join(basepath, resourcepath)), nil
}
