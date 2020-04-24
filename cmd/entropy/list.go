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
	"sort"

	"github.com/spf13/cobra"
	settings "github.com/spf13/viper"

	e "github.com/Sabayon/pkgs-checker/pkg/entropy"
)

func newEntropyListCommand() *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "list [OPTIONS]",
		Short: "Retrieve packages available on database.",
		Args:  cobra.OnlyValidArgs,
		Example: `
$> pkgs-checker entropy list -d /var/lib/entropy/client/database/amd64/sabayon-weekly/standard/amd64/5/packages.db
`,
		Run: func(cmd *cobra.Command, args []string) {

			jsonOut, _ := cmd.Flags().GetBool("json")
			verbose := settings.GetBool("verbose")
			dbPath, _ := cmd.Flags().GetString("entropy-db")
			if dbPath == "" {
				fmt.Println("Invalid entropy path")
				os.Exit(1)
			}

			pkgs, err := e.RetrieveRepoPackages(dbPath)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if jsonOut {

				out, err := json.Marshal(pkgs)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				fmt.Println(string(out))

			} else {
				plist := make([]string, 0)

				for _, p := range pkgs {
					if verbose {
						plist = append(plist, fmt.Sprintf("%s:%s", p, p.Slot))
					} else {
						plist = append(plist, fmt.Sprintf("%s/%s", p.Category, p.Name))
					}
				}

				sort.Strings(plist)

				for _, p := range plist {
					fmt.Println(p)
				}
			}
		},
	}

	cmd.Flags().StringP("entropy-db", "d", "", "Path of the entropy database")
	cmd.Flags().BoolP("json", "j", false, "Enable json output on stdout.")

	return cmd
}
