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
	"io"
	"io/ioutil"
	"os"

	logger "github.com/sirupsen/logrus"
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
	Example: "pkgs-checker -p /usr/portage/packages/sys-app/entropy-9999.tbz2",

	PreRun: func(cmd *cobra.Command, args []string) {
		if settings.GetBool("stdin") == false &&
			len(settings.GetStringSlice("package")) == 0 {
			fmt.Println("No package supply or stdin option is not present.")
			os.Exit(1)
		}
	},

	Run: func(cmd *cobra.Command, args []string) {
		initLogging()

		logger.WithFields(logger.Fields{
			"package": settings.GetStringSlice("package"),
			"stdin":   settings.GetBool("stdin"),
		}).Infof("[*] Calculate hashing of %s.",
			settings.GetStringSlice("package"))

		var checker, err = commons.New(settings.GetViper(), logger.StandardLogger())
		if err != nil {
			panic("Error on create Checker object")
		}

		err = checker.Run()
		if err != nil {
			panic(err)
		}

		for _, p := range checker.GetPackages() {
			fmt.Printf("HASH %s %s\n", p.CheckSum(), p.Name())
		}

	},
}

func initLogging() {

	// Initialize logging
	logger.SetFormatter(&logger.TextFormatter{
		FullTimestamp: true,
	})

	// Configure logging
	if settings.GetString("loglevel") == "" || settings.GetString("loglevel") == "INFO" {
		logger.SetLevel(logger.InfoLevel)
	} else if settings.GetString("loglevel") == "ERROR" {
		logger.SetLevel(logger.ErrorLevel)
	} else if settings.GetString("loglevel") == "WARN" {
		logger.SetLevel(logger.WarnLevel)
	} else if settings.GetString("loglevel") == "DEBUG" {
		logger.SetLevel(logger.DebugLevel)
	}

	if settings.GetString("logfile") != "" {
		var err error
		logFile, err = os.OpenFile(
			settings.GetString("logfile"),
			os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0660)
		if err != nil {
			fmt.Println("Error on openfile ", settings.GetString("logfile"))
		}
	}

	if logFile != nil && settings.GetBool("verbose") {
		logger.SetOutput(io.MultiWriter(os.Stdout, logFile))
	} else if logFile != nil {
		logger.SetOutput(logFile)
	} else if settings.GetBool("verbose") {
		// Default is to stderr
		logger.SetOutput(os.Stdout)
	} else {
		logger.SetOutput(ioutil.Discard)
	}

}

func init() {
	// Initialize command flags and settings binding
	rootCmd.PersistentFlags().Bool("stdin", false, "Read package data from stdin")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose logging")
	rootCmd.PersistentFlags().StringSliceP("package", "p", []string{}, "Path of package to check.")
	rootCmd.PersistentFlags().StringSliceP("ignore", "i", []string{}, "File to ignore.")
	rootCmd.PersistentFlags().StringSliceP("ignore-extension", "e", []string{}, "Extension to ignore.")
	rootCmd.PersistentFlags().StringP("logfile", "l", "", "Logfile Path. Optional.")
	rootCmd.PersistentFlags().StringP("loglevel", "L", "INFO", `Set logging level.
[DEBUG, INFO, WARN, ERROR]`)
	settings.BindPFlag("stdin", rootCmd.PersistentFlags().Lookup("stdin"))
	settings.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	settings.BindPFlag("package", rootCmd.PersistentFlags().Lookup("package"))
	settings.BindPFlag("logfile", rootCmd.PersistentFlags().Lookup("logfile"))
	settings.BindPFlag("loglevel", rootCmd.PersistentFlags().Lookup("loglevel"))
	settings.BindPFlag("ignoreFiles", rootCmd.PersistentFlags().Lookup("ignore"))
	settings.BindPFlag("ignoreExt", rootCmd.PersistentFlags().Lookup("ignore-extension"))
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
