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
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpClientOpts struct {
	InsecureSkipVerify bool
	MaxIdleConns       int
	IdleConnTimeout    int
	ConnTimeout        int
}

func NewHttpClientDefaultOpts() HttpClientOpts {
	return HttpClientOpts{
		InsecureSkipVerify: false,
		MaxIdleConns:       5,
		IdleConnTimeout:    30,
		ConnTimeout:        60,
	}
}

func GetResource(url, apiKey string, opts HttpClientOpts) ([]byte, error) {
	var err error
	var req *http.Request = nil

	transport := &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		MaxIdleConns:    opts.MaxIdleConns,
		IdleConnTimeout: time.Duration(opts.IdleConnTimeout) * time.Second,
	}

	if opts.InsecureSkipVerify {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(opts.ConnTimeout) * time.Second,
	}

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if apiKey != "" {
		req.Header.Add("Authorization", "token "+apiKey)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Invalid response %d for url %s",
			resp.StatusCode, url)
	}

	byteValue, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// String conversion suffer of performance issue
	// for now i ignore it.
	return byteValue, nil
}
