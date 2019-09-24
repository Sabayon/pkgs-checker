/*

Copyright (C) 2017-2018  Daniele Rondina <geaaru@sabayonlinux.org>

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

	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	settings "github.com/spf13/viper"

	"github.com/Sabayon/pkgs-checker/pkg/commons"
	"github.com/Sabayon/pkgs-checker/pkg/pkglist"
)

func newPkglistCreateCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "create [OPTIONS]",
		Short: "Create pkglist file.",
		Args:  cobra.OnlyValidArgs,

		PreRun: func(cmd *cobra.Command, args []string) {
			if settings.GetString("pkglist-binhost-dir") == "" {
				fmt.Errorf("Missing mandatory binhost-dir option\n")
				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {

			binHostDir := settings.GetString("pkglist-binhost-dir")

			pkgs, err := pkglist.PkgListCreate(binHostDir, logger.StandardLogger())
			if err != nil {
				fmt.Errorf("Error: %s\n", string(err.Error()))
				os.Exit(1)
			}

			if settings.GetString("pkglist-file") != "" {
				err = pkglist.PkgListWriteFile(pkgs, settings.GetString("pkglist-file"))
			} else {
				err = pkglist.PkgListWrite(pkgs, os.Stdout)
			}
			commons.CheckErr(err)
		},
	}

	var flags = cmd.Flags()

	flags.StringP("binhost-dir", "d", "", "bin-hosts directory where compute pkglist.")
	flags.StringP("pkglist-file", "f", "", `Path of pkglist file.
Default output to stdout with format: category/pkgname-pkgversion`)

	settings.BindPFlag("pkglist-binhost-dir", flags.Lookup("binhost-dir"))
	settings.BindPFlag("pkglist-file", flags.Lookup("pkglist-file"))

	return cmd
}
