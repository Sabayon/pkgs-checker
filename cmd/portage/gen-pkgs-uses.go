/*

Copyright (C) 2017-2021  Daniele Rondina <geaaru@sabayonlinux.org>

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
package portage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Sabayon/pkgs-checker/pkg/gentoo"
	"github.com/Sabayon/pkgs-checker/pkg/luet"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func parseFilterFile(file string) (gentoo.PortageUseParseOpts, error) {
	var ans gentoo.PortageUseParseOpts

	_, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return ans, err
		}
		return ans, err
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return ans, err
	}

	if err := yaml.Unmarshal(data, &ans); err != nil {
		return ans, err
	}

	return ans, nil
}

func newGenPkgsUsesCommand() *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "gen-pkgs-uses [OPTIONS]",
		Short: "Generate packages.use of the installed pkgs.",
		Args:  cobra.OnlyValidArgs,
		Run: func(cmd *cobra.Command, args []string) {

			dbPkgsDir, _ := cmd.Flags().GetString("db-pkgs-dir-path")
			jsonOutput, _ := cmd.Flags().GetBool("json")
			filterFile, _ := cmd.Flags().GetString("filter-opts")
			treePath, _ := cmd.Flags().GetString("treePath")
			lpcFormat, _ := cmd.Flags().GetBool("luet-portage-converter-format")
			verbose, _ := cmd.Flags().GetBool("verbose")

			if dbPkgsDir == "" {
				fmt.Println("Invalid Path of the portage metadata.")
				os.Exit(1)
			}

			var opts gentoo.PortageUseParseOpts
			var err error

			if filterFile != "" {
				opts, err = parseFilterFile(filterFile)
				if err != nil {
					fmt.Println("Error on read filter file: " + err.Error())
					os.Exit(1)
				}
			} else {

				opts = gentoo.PortageUseParseOpts{
					UseFilters: []string{
						"^userland_",
						"^kernel_",
						"^x86",
						"^x64",
						"^ppc",
						"^arm",
						"^amd64",
						"^prefix",
						"^m68k",
						"^ia64",
						"^riscv",
						"^s390",
						"^hppa",
						"^mips",
						"^alpha",
						"^sparc",
						"^elibc_",
					},
				}
			}

			opts.Verbose = verbose

			pkgs, err := gentoo.ParseMetadataDir(dbPkgsDir, opts)
			if err != nil {
				fmt.Println("ERROR: " + err.Error())
				os.Exit(1)
			}

			if lpcFormat {
				artefacts := luet.ConvertPortageMeta2PortageConverter(pkgs, treePath)

				data, err := yaml.Marshal(&artefacts)
				if err != nil {
					fmt.Println(fmt.Sprintf("Error on convert data to YAML: %s", err.Error()))
					os.Exit(1)
				}

				fmt.Println(string(data))
			} else if jsonOutput {
				data, err := json.Marshal(pkgs)
				if err != nil {
					fmt.Println(fmt.Sprintf("Error on convert data to json: %s", err.Error()))
					os.Exit(1)
				}
				fmt.Println(string(data))
			} else {
				for _, p := range pkgs {
					fmt.Println(
						fmt.Sprintf("%s %s", p.GetPackageNameWithSlot(), strings.Join(p.UseFlags, " ")),
					)
				}
			}
		},
	}

	var flags = cmd.Flags()

	flags.StringP("db-pkgs-dir-path", "p", "/var/db/pkg",
		"Path of the portage metadata.")
	flags.BoolP("json", "j", false, "Output in JSON format")
	flags.String("filter-opts", "", "Using filter rules through YAML file.")
	flags.Bool("luet-portage-converter-format", false,
		"Generate luet-portage-converter YAML output.")
	flags.String("treePath", "packages/atoms",
		"Define the tree path to use on luet-portage-converter artefacts.")

	return cmd
}
