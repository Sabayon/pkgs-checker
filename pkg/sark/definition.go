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
package sark

import (
	v "github.com/spf13/viper"
)

type SarkConfig struct {
	Viper *v.Viper

	Id         string
	Repository SarkRepository   `mapstructure:"repository" yaml:"repository,omitempty"`
	Build      SarkBuild        `mapstructure:"build" yaml:"build,omitempty"`
	Injector   SarkInjectConfig `mapstructure:"injector" yaml:"injector,omitempty"`
}

type SarkBuild struct {
	Script     SarkBuildScript `mapstructure:"script" yaml:"script,omitempty"`
	Verbose    int             `mapstructure:"verbose" yaml:"verbose,omitempty"`
	QA_Checks  int             `mapstructure:"qa_checks" yaml:"qa_checks,omitempty"`
	Overlays   []string        `mapstructure:"overlays" yaml:"overlays,omitempty"`
	TargetPkgs []string        `mapstructure:"target" yaml:"target,omitempty"`
	Equo       SarkBuildEquo   `mapstructure:"equo" yaml:"equo,omitempty"`
	Emerge     SarkBuildEmerge `mapstructure:"emerge" yaml:"emerge,omitempty"`
}

type SarkBuildEmerge struct {
	DefaultArgs       string `mapstructure:"default_args" yaml:"default_args,omitempty"`
	SplitInstall      int    `mapstructure:"split_install" yaml:"split_install,omitempty"`
	Features          string `mapstructure:"features" yaml:"features,omitempty"`
	Profile           string `mapstructure:"profile" yaml:"profile,omitempty"`
	Jobs              int    `mapstructure:"jobs" yaml:"jobs,omitempty"`
	PreserverdRebuild int    `mapstructure:"preserved_rebuild" yaml:"preserved_rebuild,omitempty"`
	SkipSync          int    `mapstructure:"skip_sync" yaml:"skip_sync,omitempty"`
	WebRsync          int    `mapstructure:"webrsync" yaml:"webrsync,omitempty"`

	RemoteOverlay       []string `mapstructure:"remote_overlay" yaml:"remote_overlay,omitempty"`
	RemoveRemoveOverlay []string `mapstructure:"remove_remote_overlay" yaml:"remove_remote_overlay,omitempty"`
	RemoveLaymanOverlay []string `mapstructure:"remove_layman_overlay" yaml:"remove_layman_overlay,omitempty"`
	RemovePkgs          []string `mapstructure:"remove" yaml:"remove,omitempty"`
}

type SarkBuildEquo struct {
	// no_cache is needed?

	EnmanAddRepositories []string `mapstructure:"repositories" yaml:"repositories,omitempty"`
	EnmanDelRepositories []string `mapstructure:"remove_repositories" yaml:"remove_repositories,omitempty"`
	EnmanSelf            int      `mapstructure:"enman_self" yaml:"enman_self,omitempty"`

	Packages          SarkBuildEquoPackage     `mapstructure:"package" yaml:"package,omitempty"`
	Repository        string                   `mapstructure:"repository" yaml:"repository,omitempty"`
	DependencyInstall SarkBuildEquoDepsInstall `mapstructure:"dependency_install" yaml:"dependency_install,omitempty"`
}

type SarkBuildEquoDepsInstall struct {
	Enable                   int `mapstructure:"enable" yaml:"enable,omitempty"`
	InstallAtoms             int `mapstructure:"install_atoms" yaml:"install_atoms,omitempty"`
	DependencyScanDepth      int `mapstructure:"dependency_scan_depth" yaml:"dependency_scan_depth,omitempty"`
	DependencyIgnoreVersions int `mapstructure:"dependency_ignore_versions" yaml:"dependency_ignore_versions,omitempty"`
	PruneVirtuals            int `mapstructure:"prune_virtuals" yaml:"prune_virtuals,omitempty"`
	InstallVersion           int `mapstructure:"install_version" yaml:"install_version,omitempty"`
	SplitInstall             int `mapstructure:"split_install" yaml:"split_install,omitempty"`
}

type SarkBuildEquoPackage struct {
	Install []string `mapstructure:"install" yaml:"install,omitempty"`
	Remove  []string `mapstructure:"remove" yaml:"remove,omitempty"`
	Mask    []string `mapstructure:"mask" yaml:"mask,omitempty"`
	Unmask  []string `mapstructure:"unmask yaml:"unmask,omitempty"`
}

type SarkBuildScript struct {
	PreScripts  []string `mapstructure:"pre" yaml:"pre,omitempty"`
	PostScripts []string `mapstructure:"post" yaml:"post,omitempty"`
}

type SarkRepository struct {
	Description string                    `mapstructure:"description" yaml:"description,omitempty"`
	Maintenance SarkRepositoryMaintenance `mapstructure:"maintenance" yaml:"maintenance,omitempty"`
}

type SarkRepositoryMaintenance struct {
	CheckDiffs           int      `mapstructure:"check_diffs" yaml:"check_diffs,omitempty"`
	CleanCache           int      `mapstructure:"clean_cache" yaml:"clean_cache,omitempty"`
	KeepPreviousVersions int      `mapstructure:"keep_previous_versions" yaml:"keep_previous_versions,omitempty"`
	RemovePkgs           []string `mapstructure:"remove" yaml:"remove,omitempty"`
}

type SarkInjectConfig struct {
	Filter SarkInjectFilterConfig `mapstructure:"filter" yaml:"filter,omitempty"`
}

// https://github.com/mitchellh/mapstructure/pull/145 (omitempty is not yet supported)
type SarkInjectFilterConfig struct {
	FilterType string               `mapstructure:"type" yaml:"type,omitempty"` // values whitelist|blacklist
	Rules      []SarkFilterRuleConf `mapstructure:"rules" yaml:"rules,omitempty"`
}

type SarkFilterRuleConf struct {
	Descr      string   `mapstructure:"description" yaml:"description,omitempty"`
	Packages   []string `mapstructure:"pkgs" yaml:"pkgs,omitempty"`
	Categories []string `mapstructure:"categories" yaml:"categories,omitempty"`
	Files      []string `mapstructure:"files" yaml:"files,omitempty"`
	Urls       []string `mapstructure:"urls" yaml:"urls,omitempty"`
}
