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
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	tools "github.com/MottainaiCI/simplestreams-builder/pkg/tools"
)

type FilterReport struct {
	FilterDate string   `json:"filter_date,omitempty"`
	FilterType string   `json:"filter_type,omitempty"`
	Matches    []string `json:"matches,omitempty"`
	NotMatches []string `json:"not_matches,omitempty"`
}

func NewFilterReport(filterType string) (*FilterReport, error) {
	if filterType == "" {
		return nil, errors.New("Invalid filter type")
	}
	ans := &FilterReport{
		FilterDate: fmt.Sprintf("%d", time.Now().Unix()),
		FilterType: filterType,
		Matches:    make([]string, 0),
		NotMatches: make([]string, 0),
	}

	return ans, nil
}

func (f *FilterReport) WriteReport(reportPrefix string) error {
	var reportFile string

	if reportPrefix == "" {
		return errors.New("Invalid report prefix")
	}

	dir := filepath.Dir(reportPrefix)
	_, err := tools.MkdirIfNotExist(dir, 0760)
	if err != nil {
		return err
	}

	if dir == filepath.Clean(reportPrefix) {
		// POST: prefix is a directory
		reportFile = filepath.Join(reportPrefix, "/report.filtered")
	} else {
		reportFile = fmt.Sprintf("%s-report.filtered", filepath.Clean(reportPrefix))
	}
	file, err := os.OpenFile(reportFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	enc := json.NewEncoder(w)
	err = enc.Encode(*f)
	if err != nil {
		return err
	}
	w.Flush()

	return nil
}

func (f *FilterReport) GetReport() (string, error) {
	bytes, err := json.Marshal(f)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
