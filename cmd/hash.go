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
package cmd

import (
	"fmt"
	"os"

	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	settings "github.com/spf13/viper"

	"github.com/Sabayon/pkgs-checker/pkg/commons"
)

func newHashCommand() *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "hash [OPTIONS]",
		Short: "Hashing packages",
		Args:  cobra.OnlyValidArgs,

		Example: `$> pkgs-checker hash -p /usr/portage/packages/sys-app/entropy-9999.tbz2

$> pkgs-checker hash -e .pyc -e .pyo -e .mo -e .bz2 --directory /usr/portage/packages/`,

		PreRun: func(cmd *cobra.Command, args []string) {
			if settings.GetBool("stdin") == false &&
				len(settings.GetStringSlice("package")) == 0 &&
				settings.GetString("directory") == "" {
				fmt.Println("Both package and directory not present or stdin option is not present.")
				os.Exit(1)
			}
		},

		Run: func(cmd *cobra.Command, args []string) {

			var err error
			var checker commons.CheckerExecutor

			logger.WithFields(logger.Fields{
				"package": settings.GetStringSlice("package"),
				"dir":     settings.GetString("directory"),
				"stdin":   settings.GetBool("stdin"),
			}).Debugf("[*] Starting Calculate hashing...")

			if settings.GetBool("concurrency") == true {
				commons.InitConcurrency()
				checker, err = commons.NewCheckerConcurrent(settings.GetViper(), logger.StandardLogger())
			} else {
				checker, err = commons.NewChecker(settings.GetViper(), logger.StandardLogger())
			}
			if err != nil {
				panic("Error on create Checker object")
			}

			err = checker.Run()
			commons.CheckErr(err)

			if settings.GetString("hashfile") != "" {
				writeHashfile(checker)
			} else {
				for _, p := range checker.GetPackages() {
					// Skip package in errors from file
					if p.CheckSum() != "" {
						fmt.Printf("HASH %s %s\n", p.CheckSum(), p.Name())
					}
				}
			}

		},
	}

	var flags = cmd.Flags()

	flags.Bool("stdin", false, "Read package data from stdin")
	flags.Bool("hash-empty", false,
		fmt.Sprintf("If create a fake hash for empty packages or use %s.",
			commons.PKGS_CHECKER_EMPTY_PKGHASH))
	flags.Bool("ignore-errors", false, "Ignore errors with broken tarball.")
	flags.StringSliceP("package", "p", []string{}, "Path of package to check.")
	flags.StringSliceP("ignore", "i", []string{}, "File to ignore.")
	flags.StringSliceP("ignore-extension", "e", []string{}, "Extension to ignore.")

	flags.StringP("directory", "d", "", "Artefacts directory with .tbz2 files.")
	flags.StringP("hashfile", "f", "", `Path of hashfile where write checksum.
Default output on stdout with format: HASH <CHECKSUM> <PACKAGE>`)

	settings.BindPFlag("stdin", flags.Lookup("stdin"))
	settings.BindPFlag("package", flags.Lookup("package"))
	settings.BindPFlag("directory", flags.Lookup("directory"))
	settings.BindPFlag("hashfile", flags.Lookup("hashfile"))
	settings.BindPFlag("hash-empty", flags.Lookup("hash-empty"))
	settings.BindPFlag("ignoreFiles", flags.Lookup("ignore"))
	settings.BindPFlag("ignoreExt", flags.Lookup("ignore-extension"))
	settings.BindPFlag("ignoreErrors", flags.Lookup("ignore-errors"))

	return cmd
}

func writeHashfile(checker commons.CheckerExecutor) {
	var err error
	var hashfile *os.File

	hashfile, err = os.OpenFile(
		settings.GetString("hashfile"),
		os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0660)
	if err != nil {
		panic(fmt.Sprintf("Error on open hashfile %s.", settings.GetString("hashfile")))
	}
	defer hashfile.Close()

	for _, p := range checker.GetPackages() {
		// Skip package in errors from file
		if p.CheckSum() != "" {
			_, err = fmt.Fprintf(hashfile, "%s %s\n", p.CheckSum(), p.Name())
		}
		commons.CheckErr(err)
	}

	hashfile.Sync()
}
