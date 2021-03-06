/*

Copyright (C) 2017-2020  Daniele Rondina <geaaru@sabayonlinux.org>

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
package entropy

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	e "github.com/Sabayon/pkgs-checker/pkg/entropy"
	gentoo "github.com/Sabayon/pkgs-checker/pkg/gentoo"
)

func newEntropyInfoCommand() *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "info [OPTIONS]",
		Short: "Retrieve package detail.",
		Args:  cobra.OnlyValidArgs,
		Example: `
$> pkgs-checker entropy info app/foo-1
`,
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("No package supply.")
				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {

			pkgname := args[0]

			jsonOut, _ := cmd.Flags().GetBool("json")
			onlyDeps, _ := cmd.Flags().GetBool("onlydeps")
			dbPath, _ := cmd.Flags().GetString("entropy-db")
			if dbPath == "" {
				fmt.Println("Invalid entropy path")
				os.Exit(1)
			}

			pkg, err := e.NewEntropyPackage(pkgname)
			if err != nil {
				fmt.Println(fmt.Sprintf("Invalid package %s: %s", pkg, err))
				os.Exit(1)
			}

			detail, err := e.RetrievePackageData(pkg, dbPath)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if jsonOut {

				out, err := json.Marshal(detail)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				fmt.Println(string(out))

			} else {

				if !onlyDeps {
					fmt.Println("name:", detail.Package.Name)
					fmt.Println("category:", detail.Package.Category)
					fmt.Println("version:", detail.Package.Version)
					fmt.Println("version_suffix:", detail.Package.VersionSuffix)
					fmt.Println("version_build:", detail.Package.VersionBuild)
					fmt.Println("slot:", detail.Package.Slot)
					fmt.Println("condition:", detail.Package.Condition.String())
					fmt.Println("uses:", strings.Join(detail.Package.UseFlags, " "))
					fmt.Println("license:", detail.Package.License)
				}

				if len(detail.Dependencies) > 0 {
					if !onlyDeps {
						fmt.Println("\nDependencies:\n")
					}
					for _, dep := range detail.Dependencies {
						if onlyDeps {
							if dep.Condition == gentoo.PkgCondMatchVersion {
								fmt.Printf("=%s:%s*\n", dep.String(), dep.Slot)
							} else {
								fmt.Printf("%s%s:%s\n", dep.Condition.String(),
									dep.String(), dep.Slot)
							}
						} else {
							fmt.Println("\tname:", dep.Name)
							fmt.Println("\tcategory:", dep.Category)
							fmt.Println("\tversion:", dep.Version)
							fmt.Println("\tversion_suffix:", dep.VersionSuffix)
							fmt.Println("\tslot:", dep.Slot)
							fmt.Println("\tcondition:", dep.Condition.String())
							fmt.Println("")
						}
					}
				}

				if len(detail.Files) > 0 {

					if !onlyDeps {
						fmt.Println("\nFiles:\n")

						for _, f := range detail.Files {
							fmt.Println(fmt.Sprintf("\t%s", f))
						}
					}

				}

			}

		},
	}

	cmd.Flags().StringP("entropy-db", "d", "", "Path of the entropy database")
	cmd.Flags().Bool("onlydeps", false, "Print only deps in quiet mode.")
	cmd.Flags().BoolP("json", "j", false, "Enable json output on stdout.")

	return cmd
}
