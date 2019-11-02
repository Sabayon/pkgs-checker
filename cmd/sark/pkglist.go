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
	"sort"

	"github.com/spf13/cobra"
	settings "github.com/spf13/viper"

	"github.com/Sabayon/pkgs-checker/pkg/commons"
	"github.com/Sabayon/pkgs-checker/pkg/pkglist"
	"github.com/Sabayon/pkgs-checker/pkg/sark"
)

func newSarkPkglistCommand() *cobra.Command {
	var sark_files []string

	var cmd = &cobra.Command{
		Use:   "pkglist [OPTIONS]",
		Short: "Show list of targets defined between build files.",
		Args:  cobra.OnlyValidArgs,
		Example: `
Show targets defined
$> pkgs-checker sark pkglist -s core-staging1-build.yaml -s core-staging2-build.yaml
`,
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(sark_files) == 0 {
				fmt.Fprintln(os.Stderr, "No sark config resources defined")
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

			sark_targets, err = pkglist.PkgListWithSlot(sark_targets, false)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error on process sark targets: %s\n", err.Error())
				os.Exit(1)
			}

			sort.Strings(sark_targets)
			for _, pkg := range sark_targets {
				fmt.Println(pkg)
			}
		},
	}

	var flags = cmd.Flags()

	flags.StringSliceVarP(&sark_files, "sark-files", "s", []string{},
		"Path or URL of sark config resources.")

	return cmd
}
