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
	"strings"

	yaml "github.com/go-yaml/yaml"
	v "github.com/spf13/viper"
)

type SarkConfig struct {
	Viper *v.Viper

	Injector SarkInjectConfig `mapstructure:"injector"`
}

type SarkInjectConfig struct {
	Filter SarkInjectFilterConfig `mapstructure:"filter"`
}

// https://github.com/mitchellh/mapstructure/pull/145 (omitempty is not yet supported)
type SarkInjectFilterConfig struct {
	FilterType string               `mapstructure:"type"` // values whitelist|blacklist
	Rules      []SarkFilterRuleConf `mapstructure:"rules"`
}

type SarkFilterRuleConf struct {
	Descr      string   `mapstructure:"description"`
	Packages   []string `mapstructure:"pkgs"`
	Categories []string `mapstructure:"categories"`
	Files      []string `mapstructure:"files"`
	Urls       []string `mapstructure:"urls"`
}

func (s *SarkConfig) unmarshalAndVerify() error {
	err := s.Viper.Unmarshal(&s)
	if err != nil {
		return err
	}

	filterType := s.Injector.Filter.FilterType
	if filterType != "whitelist" && filterType != "blacklist" {
		return errors.New("Invalid or missing filter type")
	}

	return nil
}

func NewSarkConfigFromString(viper *v.Viper, config string) (*SarkConfig, error) {
	var ans *SarkConfig
	var err error

	if config == "" {
		return nil, errors.New("Invalid configuration")
	}

	if viper == nil {
		viper = v.New()
	}

	viper.SetConfigType("yaml")
	err = viper.ReadConfig(strings.NewReader(config))
	if err != nil {
		return nil, err
	}

	ans = &SarkConfig{
		Viper: viper,
	}

	err = ans.unmarshalAndVerify()

	return ans, err
}

func NewSarkConfigFromFile(viper *v.Viper, file string) (*SarkConfig, error) {
	var ans *SarkConfig
	var err error

	if file == "" {
		return nil, errors.New("Invalid file path")
	}

	if viper == nil {
		viper = v.New()
	}

	viper.SetConfigFile(file)
	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	ans = &SarkConfig{
		Viper: viper,
	}

	err = ans.unmarshalAndVerify()

	return ans, err
}

func NewSarkConfig(viper *v.Viper, filterType string) (*SarkConfig, error) {
	if viper == nil {
		viper = v.New()
	}

	if filterType != "whitelist" && filterType != "blacklist" {
		return nil, errors.New("Invalid filter type")
	}

	return &SarkConfig{
		Injector: SarkInjectConfig{
			Filter: SarkInjectFilterConfig{
				FilterType: filterType,
				Rules:      make([]SarkFilterRuleConf, 0),
			},
		},
	}, nil
}

func NewSarkFilterRuleConf(desc string) *SarkFilterRuleConf {
	return &SarkFilterRuleConf{
		Descr:      desc,
		Packages:   make([]string, 0),
		Categories: make([]string, 0),
		Files:      make([]string, 0),
		Urls:       make([]string, 0),
	}
}

func (f *SarkFilterRuleConf) AddPackage(pkg string) {
	f.Packages = append(f.Packages, pkg)
}

func (f *SarkFilterRuleConf) AddCategory(category string) {
	f.Categories = append(f.Categories, category)
}

func (f *SarkFilterRuleConf) AddUrl(url string) {
	f.Urls = append(f.Urls, url)
}

func (f *SarkFilterRuleConf) AddFile(file string) {
	f.Files = append(f.Files, file)
}

func (s *SarkConfig) ToString() (string, error) {
	out, err := yaml.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
