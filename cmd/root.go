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
	"runtime"

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
	Example: `$> pkgs-checker -p /usr/portage/packages/sys-app/entropy-9999.tbz2

$> pkgs-checker -e .pyc -e .pyo -e .mo -e .bz2 --directory /usr/portage/packages/`,

	PreRun: func(cmd *cobra.Command, args []string) {
		if settings.GetBool("stdin") == false &&
			len(settings.GetStringSlice("package")) == 0 &&
			settings.GetString("directory") == "" {
			fmt.Println("Both package and directory not present or stdin option is not present.")
			os.Exit(1)
		}
	},

	Run: func(cmd *cobra.Command, args []string) {
		initLogging()

		var err error
		var checker commons.CheckerExecutor

		logger.WithFields(logger.Fields{
			"package": settings.GetStringSlice("package"),
			"dir":     settings.GetString("directory"),
			"stdin":   settings.GetBool("stdin"),
		}).Debugf("[*] Starting Calculate hashing...")

		if settings.GetBool("concurrency") == true {
			initConcurrency()
			checker, err = commons.NewCheckerConcurrent(settings.GetViper(), logger.StandardLogger())
		} else {
			checker, err = commons.NewChecker(settings.GetViper(), logger.StandardLogger())
		}
		if err != nil {
			panic("Error on create Checker object")
		}

		err = checker.Run()
		checkErr(err)

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

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func writeHashfile(checker commons.CheckerExecutor) {
	var err error
	var hashfile *os.File

	hashfile, err = os.OpenFile(
		settings.GetString("hashfile"),
		os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0660)
	if err != nil {
		panic(fmt.Sprintf("Error on open hashfile ", settings.GetString("hashfile")))
	}
	defer hashfile.Close()

	for _, p := range checker.GetPackages() {
		// Skip package in errors from file
		if p.CheckSum() != "" {
			_, err = fmt.Fprintf(hashfile, "%s %s\n", p.CheckSum(), p.Name())
		}
		checkErr(err)
	}

	hashfile.Sync()
}

func initConcurrency() {
	if settings.GetInt("maxconcurrency") > runtime.NumCPU() {
		logger.Warnf("maxconcurrency value %d is bigger of number of host CPU. I force %d.",
			settings.GetInt("maxconcurrency"), runtime.NumCPU())
		settings.Set("maxconcurrency", runtime.NumCPU())
	}

	runtime.GOMAXPROCS(settings.GetInt("maxconcurrency"))
}

func initLogging() {

	// Initialize logging
	logger.SetFormatter(&logger.TextFormatter{
		FullTimestamp: true,
	})

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

	// Configure logging
	if settings.GetString("loglevel") == "" || settings.GetString("loglevel") == "INFO" {
		logger.SetLevel(logger.InfoLevel)
	} else if settings.GetString("loglevel") == "ERROR" {
		logger.SetLevel(logger.ErrorLevel)
	} else if settings.GetString("loglevel") == "WARN" {
		logger.SetLevel(logger.WarnLevel)
	} else if settings.GetString("loglevel") == "DEBUG" {
		logger.SetLevel(logger.DebugLevel)
	} else {
		// For invalid loglevel force INFO
		logger.SetLevel(logger.InfoLevel)
	}

}

func init() {
	// Initialize command flags and settings binding
	rootCmd.PersistentFlags().Bool("stdin", false, "Read package data from stdin")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose logging on stdout.")
	rootCmd.PersistentFlags().BoolP("concurrency", "c", false, "Enable concurrency process.")
	rootCmd.PersistentFlags().Bool("hash-empty", false,
		fmt.Sprintf("If create a fake hash for empty packages or use %s.", commons.PKGS_CHECKER_EMPTY_PKGHASH))
	rootCmd.PersistentFlags().Bool("ignore-errors", false, "Ignore errors with broken tarball.")
	rootCmd.PersistentFlags().StringSliceP("package", "p", []string{}, "Path of package to check.")
	rootCmd.PersistentFlags().StringSliceP("ignore", "i", []string{}, "File to ignore.")
	rootCmd.PersistentFlags().StringSliceP("ignore-extension", "e", []string{}, "Extension to ignore.")
	rootCmd.PersistentFlags().StringP("logfile", "l", "", "Logfile Path. Optional.")
	rootCmd.PersistentFlags().StringP("loglevel", "L", "INFO", `Set logging level.
[DEBUG, INFO, WARN, ERROR]`)
	rootCmd.PersistentFlags().StringP("directory", "d", "", "Artefacts directory with .tbz2 files.")
	rootCmd.PersistentFlags().StringP("hashfile", "f", "", `Path of hashfile where write checksum.
Default output on stdout with format: HASH <CHECKSUM> <PACKAGE>`)

	settings.BindPFlag("stdin", rootCmd.PersistentFlags().Lookup("stdin"))
	settings.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	settings.BindPFlag("package", rootCmd.PersistentFlags().Lookup("package"))
	settings.BindPFlag("directory", rootCmd.PersistentFlags().Lookup("directory"))
	settings.BindPFlag("hashfile", rootCmd.PersistentFlags().Lookup("hashfile"))
	settings.BindPFlag("concurrency", rootCmd.PersistentFlags().Lookup("concurrency"))
	settings.BindPFlag("hash-empty", rootCmd.PersistentFlags().Lookup("hash-empty"))
	settings.BindPFlag("logfile", rootCmd.PersistentFlags().Lookup("logfile"))
	settings.BindPFlag("loglevel", rootCmd.PersistentFlags().Lookup("loglevel"))
	settings.BindPFlag("ignoreFiles", rootCmd.PersistentFlags().Lookup("ignore"))
	settings.BindPFlag("ignoreExt", rootCmd.PersistentFlags().Lookup("ignore-extension"))
	settings.BindPFlag("ignoreErrors", rootCmd.PersistentFlags().Lookup("ignore-errors"))
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
