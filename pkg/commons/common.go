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

package commons

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"

	logger "github.com/sirupsen/logrus"
	settings "github.com/spf13/viper"
)

const PKGS_CHECKER_VERSION = "0.6.2"
const PKGS_CHECKER_EMPTY_PKGHASH = "00000000000000000000000000000000"

func InitConcurrency() {
	if settings.GetInt("maxconcurrency") > runtime.NumCPU() {
		logger.Warnf("maxconcurrency value %d is bigger of number of host CPU. I force %d.",
			settings.GetInt("maxconcurrency"), runtime.NumCPU())
		settings.Set("maxconcurrency", runtime.NumCPU())
	}

	runtime.GOMAXPROCS(settings.GetInt("maxconcurrency"))
}

func InitLogging() *os.File {
	var logFile *os.File

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

	return logFile
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func SanitizeDuplicate(i []string) (o []string) {
	m := make(map[string]bool)
	for _, s := range i {
		m[s] = true
	}
	for k, _ := range m {
		o = append(o, k)
	}
	return
}
