# Sabayon Packages Checker

## Why

Artefacts created with emerge could be compressed with tar+bz2 not in ordered way and some files must be skipped to avoid injection of some packages that are equal.

## How

*pkgs-checker* processes tarball file and create an MD5 checksum for any file inside package that are not be skipped (by command line options). At EOF create a new MD5 with all MD5 checksum bytes plus list of directories found sorted.

Normally, files to skip are .pyc,.pyo,.mo that contains timestamp data that generate false events for package injection.

## Usage

Hereinafter, all available options:

```
$# pkgs-checker --help
Sabayon packages checker

Usage:
   [flags]

Examples:
$> pkgs-checker -p /usr/portage/packages/sys-app/entropy-9999.tbz2

$> pkgs-checker -e .pyc -e .pyo -e .mo -e .bz2 --directory /usr/portage/packages/

Flags:
  -c, --concurrency                Enable concurrency process.
  -d, --directory string           Artefacts directory with .tbz2 files.
  -f, --hashfile string            Path of hashfile where write checksum.
                                   Default output on stdout with format: HASH <CHECKSUM> <PACKAGE>
  -h, --help                       help for this command
  -i, --ignore strings             File to ignore.
  -e, --ignore-extension strings   Extension to ignore.
  -l, --logfile string             Logfile Path. Optional.
  -L, --loglevel string            Set logging level.
                                   [DEBUG, INFO, WARN, ERROR] (default "INFO")
  -p, --package strings            Path of package to check.
      --stdin                      Read package data from stdin
  -v, --verbose                    Enable verbose logging on stdout.
      --version                    version for this command

```

## Features to Complete for release 0.1.0

  * Add support to stdin processing

  * Add support to different checksum algorithms

  * Add support for directly create hashing file
