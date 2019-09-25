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

	"github.com/spf13/cobra"
	settings "github.com/spf13/viper"

	"github.com/Sabayon/pkgs-checker/pkg/commons"
	"github.com/Sabayon/pkgs-checker/pkg/pkglist"
)

type CompareResult struct {
	Resource1    string
	Resource2    string
	Intersection []string
}

func newPkglistIntersectCommand() *cobra.Command {
	var resources []string
	var quiet bool

	var cmd = &cobra.Command{
		Use:     "intersect [OPTIONS]",
		Short:   "Search duplicate package between multiple pkglist.",
		Example: `$> pkgs-checker pkglist intersect -r https://server1/sbi/namespace/base-arm/base-arm-binhost/base-arm.pkglist,https://server2/sbi/namespace/core-arm/core-arm-binhost/core-arm.pkglist`,
		Args:    cobra.OnlyValidArgs,

		PreRun: func(cmd *cobra.Command, args []string) {
			if len(resources) == 0 {
				fmt.Fprintln(os.Stderr, "No pkglist resource defined")
				os.Exit(1)
			}
			if len(resources) == 1 {
				fmt.Fprintln(os.Stderr, "At least two resources needed")
				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {

			results := make([]CompareResult, 0)
			pkglist_data := make(map[string][]string)
			opts := commons.NewHttpClientDefaultOpts()
			if settings.GetBool("insecure_skipverify") {
				opts.InsecureSkipVerify = true
			}
			apiKey := settings.GetString("apikey")
			var list2, list1 []string
			var ok bool
			var err error

			// TODO: Improve compare algorithm
			for _, r1 := range resources {

				list1, err = pkglist.PkgListLoadResource(r1, apiKey, opts)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error on load resource %s\n", r1)
					os.Exit(1)
				}
				pkglist_data[r1] = list1

				for _, r2 := range resources {
					if r1 == r2 {
						continue
					}

					if list2, ok = pkglist_data[r2]; !ok {
						list2, err = pkglist.PkgListLoadResource(r2, apiKey, opts)
						if err != nil {
							fmt.Fprintf(os.Stderr, "Error on load resource %s\n", r2)
							os.Exit(1)
						}
						pkglist_data[r2] = list2
					}

					intersect, err := pkglist.PkgListIntersectFromLists(list1, list2)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Error on compare %s with %s\n", r1, r2)
						os.Exit(1)
					}

					if len(intersect) > 0 {
						results = append(results, CompareResult{
							Resource1:    r1,
							Resource2:    r2,
							Intersection: intersect,
						})
					}
				}

			}

			// Print results
			if len(results) > 0 {
				// Create map to avoid duplicate.
				intersectMap := make(map[string]bool, 0)
				for _, res := range results {
					for _, pkg := range res.Intersection {
						intersectMap[pkg] = true
					}
				}

				for pkg, _ := range intersectMap {
					fmt.Println(pkg)
				}

			} else if !quiet {
				fmt.Println("No intersection found.")
			}
		},
	}

	var flags = cmd.Flags()

	flags.StringSliceVarP(&resources, "pkglist", "r", []string{}, "Path or URL of pkglist resource.")
	flags.BoolVarP(&quiet, "quiet", "q", false, "Quiet output.")

	return cmd
}
