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

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	settings "github.com/spf13/viper"

	"github.com/Sabayon/pkgs-checker/commons"
)

// Logfile file descriptor pointer
var logFile *os.File

// Program command declaration
var rootCmd = &cobra.Command{
	Short:   "Sabayon packages checker",
	Version: commons.PKGS_CHECKER_VERSION,
	Args:    cobra.OnlyValidArgs,

	PreRun: func(cmd *cobra.Command, args []string) {
		logFile = commons.InitLogging()
	},
}

func init() {
	// Initialize command flags and settings binding
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose logging on stdout.")
	rootCmd.PersistentFlags().BoolP("concurrency", "c", false, "Enable concurrency process.")
	rootCmd.PersistentFlags().StringP("logfile", "l", "", "Logfile Path. Optional.")
	rootCmd.PersistentFlags().StringP("loglevel", "L", "INFO", `Set logging level.
[DEBUG, INFO, WARN, ERROR]`)

	settings.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	settings.BindPFlag("concurrency", rootCmd.PersistentFlags().Lookup("concurrency"))
	settings.BindPFlag("logfile", rootCmd.PersistentFlags().Lookup("logfile"))
	settings.BindPFlag("loglevel", rootCmd.PersistentFlags().Lookup("loglevel"))

	rootCmd.AddCommand(
		newHashCommand(),
	)
}

func Execute() {
	// Start command execution
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if logFile != nil {
		defer logFile.Close()
	}
}
