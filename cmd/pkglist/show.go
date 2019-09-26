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
package pkglist

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
	settings "github.com/spf13/viper"

	"github.com/Sabayon/pkgs-checker/pkg/commons"
	"github.com/Sabayon/pkgs-checker/pkg/pkglist"
)

func newPkglistShowCommand() *cobra.Command {
	var resources []string
	var quiet bool
	var parse bool

	var cmd = &cobra.Command{
		Use:     "show [OPTIONS]",
		Short:   "Show pkglist from one or multiple resources.",
		Args:    cobra.OnlyValidArgs,
		Example: `$> pkgs-checker pkglist show -r https://server1/sbi/namespace/base-arm/base-arm-binhost/base-arm.pkglist,https://server2/sbi/namespace/core-arm/core-arm-binhost/core-arm.pkglist`,

		PreRun: func(cmd *cobra.Command, args []string) {
			if len(resources) == 0 {
				fmt.Fprintln(os.Stderr, "No pkglist resource defined")
				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			opts := commons.NewHttpClientDefaultOpts()
			if settings.GetBool("insecure_skipverify") {
				opts.InsecureSkipVerify = true
			}
			apiKey := settings.GetString("apikey")

			plist := make([]string, 0)

			// TODO: Improve compare algorithm
			for _, r1 := range resources {

				list, err := pkglist.PkgListLoadResource(r1, apiKey, opts)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error on load resource %s\n", r1)
					os.Exit(1)
				}

				plist = append(plist, list...)
			}

			// Print results
			if len(plist) > 0 {

				if parse {
					emap, err := pkglist.PkgListConvertToMap(plist)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Error on parse package name: %s\n", err.Error())
						os.Exit(1)
					}

					plist = make([]string, 0)
					for _, pkgs := range emap {
						for _, p := range pkgs {
							plist = append(plist, p.String())
						}
					}
				} else {
					// Create map to avoid duplicate.
					pmap := make(map[string]bool, 0)
					for _, p := range plist {
						pmap[p] = true
					}

					plist = make([]string, 0)
					for pkg, _ := range pmap {
						plist = append(plist, pkg)
					}
				}

				sort.Strings(plist)

				for _, pkg := range plist {
					fmt.Println(pkg)
				}
			} else if !quiet {
				fmt.Println("No packages available.")
			}
		},
	}

	var flags = cmd.Flags()

	flags.StringSliceVarP(&resources, "pkglist", "r", []string{}, "Path or URL of pkglist resource.")
	flags.BoolVarP(&quiet, "quiet", "q", false, "Quiet output.")
	flags.BoolVarP(&parse, "parse-pkgname", "p", false, "Parse package version string and hide entropy revision.")

	return cmd
}
