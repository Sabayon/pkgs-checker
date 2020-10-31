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
package pkg

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Sabayon/pkgs-checker/pkg/gentoo"
)

type PkgInfoReport struct {
	Packages []gentoo.GentooPackage `json:"packages,omitempty"`
}

func newPkgInfoCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "info [OPTIONS]",
		Short: "Parse package string and print detail.",
		Args:  cobra.OnlyValidArgs,
		Example: `
$> pkgs-checker pkg info app/foo-3.30
`,
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("No packages availables.")
				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			jsonOut, _ := cmd.Flags().GetBool("json")

			pkgs := []gentoo.GentooPackage{}

			for _, pkg := range args {
				gp, err := gentoo.ParsePackageStr(pkg)
				if err != nil {
					fmt.Println(fmt.Sprintf("Invalid package %s: %s", pkg, err))
					os.Exit(1)
				}

				pkgs = append(pkgs, *gp)

			}

			if len(pkgs) == 0 {
				os.Exit(0)
			}

			if jsonOut {

				var err error
				var out []byte

				if len(pkgs) > 1 {
					preport := PkgInfoReport{
						Packages: pkgs,
					}

					out, err = json.Marshal(preport)
				} else {
					out, err = json.Marshal(pkgs[0])
				}
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				fmt.Println(string(out))

			} else {

				for _, gp := range pkgs {
					fmt.Println("name:", gp.Name)
					fmt.Println("category:", gp.Category)
					fmt.Println("version:", gp.Version)
					fmt.Println("version_suffix:", gp.VersionSuffix)
					fmt.Println("version_build:", gp.VersionBuild)
					fmt.Println("slot:", gp.Slot)
					fmt.Println("condition:", gp.Condition)
					fmt.Println("repository:", gp.Repository)
					fmt.Println("uses:", gp.UseFlags)
				}

			}
		},
	}

	cmd.Flags().BoolP("json", "j", false, "Enable json output on stdout.")

	return cmd
}
