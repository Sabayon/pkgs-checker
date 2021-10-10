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
	"fmt"
	"os"

	"github.com/Sabayon/pkgs-checker/pkg/gentoo"

	"github.com/spf13/cobra"
)

func newGenMetadataCommand() *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "gen-metadata cat/pkg[:slot] ... catN/pkgN[:slot] [OPTIONS]",
		Short: "Generate metadata of a package to a specific path.",
		Args:  cobra.OnlyValidArgs,
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Fprintf(os.Stderr, "No packages defined.\n")
				os.Exit(1)
			}

			dbPkgsDir, _ := cmd.Flags().GetString("db-pkgs-dir-path")
			if dbPkgsDir == "" {
				fmt.Println("Invalid Path of the portage metadata.")
				os.Exit(1)
			}
			to, _ := cmd.Flags().GetString("to")
			if to == "" {
				fmt.Println("Invalid path where generate metadata.")
				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {

			dbPkgsDir, _ := cmd.Flags().GetString("db-pkgs-dir-path")
			verbose, _ := cmd.Flags().GetBool("verbose")
			to, _ := cmd.Flags().GetString("to")

			var err error
			var opts *gentoo.PortageUseParseOpts = &gentoo.PortageUseParseOpts{
				UseFilters: []string{},
				Categories: []string{},
				Packages:   []string{},
			}

			for _, pkg := range args {
				gp, err := gentoo.ParsePackageStr(pkg)
				if err != nil {
					fmt.Println(fmt.Sprintf("Invalid pkg %s: %s",
						pkg, err.Error()))
					os.Exit(1)
				}

				opts.Packages = append(opts.Packages, gp.GetPackageNameWithSlot())
				opts.AddCategory(gp.Category)
			}

			opts.Verbose = verbose

			pkgs, err := gentoo.ParseMetadataDir(dbPkgsDir, opts)
			if err != nil {
				fmt.Println("ERROR: " + err.Error())
				os.Exit(1)
			}

			for _, p := range pkgs {

				fmt.Println(
					fmt.Sprintf("Package: %s-%s:%s", p.GetPackageName(), p.GetPVR(), p.Slot),
				)

				if len(p.CONTENTS) > 0 {
					fmt.Println("CONTENTS:")
					for _, e := range p.CONTENTS {
						fmt.Println(e)
					}
				}

				err = p.WriteMetadata2Dir(to, opts)
				if err != nil {
					fmt.Println(
						fmt.Sprintf("Error on generate metadata for %s: %s",
							p.GetPackageNameWithSlot(), err.Error()))
					os.Exit(1)
				}

			}
		},
	}

	var flags = cmd.Flags()

	flags.String("to", "/gen-metadata",
		"Generate medata tree to the specified path.")
	flags.StringP("db-pkgs-dir-path", "p", "/var/db/pkg",
		"Path of the portage metadata.")

	return cmd
}
