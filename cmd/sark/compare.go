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
	"fmt"
	"os"

	"github.com/spf13/cobra"
	settings "github.com/spf13/viper"

	"github.com/Sabayon/pkgs-checker/pkg/commons"
	"github.com/Sabayon/pkgs-checker/pkg/pkglist"
	"github.com/Sabayon/pkgs-checker/pkg/sark"
)

func newSarkCompareCommand() *cobra.Command {
	var pkglist_files []string
	var sark_files []string
	var targetsNotInList, pkgsNotInTarget bool

	var cmd = &cobra.Command{
		Use:   "compare [OPTIONS]",
		Short: "Compare sark targets with pkglist files",
		Args:  cobra.OnlyValidArgs,
		Example: `
Show targets not present on package list:
$> pkgs-checker sark compare -s core-staging1-build.yaml -r core-arm.pkglist -t -m

Show packages not present between SARK targets:
$> pkgs-checker sark compare -s core-staging1-build.yaml -r core-arm.pkglist -v -t
`,
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(pkglist_files) == 0 {
				fmt.Fprintln(os.Stderr, "No pkglist resources defined")
				os.Exit(1)
			}
			if len(sark_files) == 0 {
				fmt.Fprintln(os.Stderr, "No sark config resources defined")
				os.Exit(1)
			}

			if targetsNotInList && pkgsNotInTarget {
				fmt.Fprintln(os.Stderr,
					"Both missing-targets and missing-packages couldn't be enabled.")
				os.Exit(1)
			} else if !targetsNotInList && !pkgsNotInTarget {
				fmt.Fprintln(os.Stderr,
					"One of missing-targets and missing-packages is needed.")
				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			opts := commons.NewHttpClientDefaultOpts()
			if settings.GetBool("insecure_skipverify") {
				opts.InsecureSkipVerify = true
			}
			apiKey := settings.GetString("apikey")

			// Load pkglist resources
			plist := make([]string, 0)

			for _, r := range pkglist_files {
				list, err := pkglist.PkgListLoadResource(r, apiKey, opts)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error on load pkglist %s\n", r)
					os.Exit(1)
				}

				plist = append(plist, list...)
			}

			plist, err = pkglist.PkgListWithoutVersions(plist)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error on process pkglist %s\n", err.Error())
				os.Exit(1)
			}

			// Load sark resources
			sark_targets := make([]string, 0)

			for _, s := range sark_files {
				conf, err := sark.NewSarkConfigFromResource(nil, s, apiKey, opts)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error on load sark config %s\n", s)
					os.Exit(1)
				}

				sark_targets = append(sark_targets, conf.Build.TargetPkgs...)
			}

			sark_targets, err = pkglist.PkgListWithoutVersions(sark_targets)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error on process sark targets: %s\n", err.Error())
				os.Exit(1)
			}

			var pkgs []string
			if targetsNotInList {
				pkgs = pkglist.PkgListPkgsNotInList(sark_targets, plist)
			} else if pkgsNotInTarget {
				pkgs = pkglist.PkgListPkgsNotInList(plist, sark_targets)
			}

			for _, pkg := range pkgs {
				fmt.Println(pkg)
			}
		},
	}

	var flags = cmd.Flags()

	flags.StringSliceVarP(&pkglist_files, "pkglist-files", "p", []string{},
		"Path or URL of pkglist resources.")
	flags.StringSliceVarP(&sark_files, "sark-files", "s", []string{},
		"Path or URL of sark config resources.")
	flags.BoolVarP(&targetsNotInList, "missing-packages", "m", false,
		"Show targets not present on pkglist(s).")
	flags.BoolVarP(&pkgsNotInTarget, "missing-targets", "t", false,
		"Show packages not present on target.")

	return cmd
}
