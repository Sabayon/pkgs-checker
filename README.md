# Sabayon Packages Checker

Tool used for different tasks on Sabayon build process.

```bash


$# pkgs-checker --help
Sabayon packages checker

Usage:
   [command]

Available Commands:
  filter      Filter bin-host packages/directory.
  hash        Hashing packages
  help        Help about any command

Flags:
  -c, --concurrency       Enable concurrency process.
  -h, --help              help for this command
  -l, --logfile string    Logfile Path. Optional.
  -L, --loglevel string   Set logging level.
                          [DEBUG, INFO, WARN, ERROR] (default "INFO")
  -v, --verbose           Enable verbose logging on stdout.
      --version           version for this command

Use " [command] --help" for more information about a command.

```

### *hash* command

#### Why

Artefacts created with emerge could be compressed with tar+bz2 not in ordered way and some files must be skipped to avoid injection of some packages that are equal.

#### How

*pkgs-checker* processes tarball file and create an MD5 checksum for any file inside package that are not be skipped (by command line options). At EOF create a new MD5 with all MD5 checksum bytes plus list of directories found sorted.

Normally, files to skip are .pyc,.pyo,.mo that contains timestamp data that generate false events for package injection.

#### Usage

Hereinafter, all available options:

```
$# pkgs-checker hash --help
Hashing packages

Usage:
   hash [OPTIONS] [flags]

Examples:
$> pkgs-checker hash -p /usr/portage/packages/sys-app/entropy-9999.tbz2

$> pkgs-checker hash -e .pyc -e .pyo -e .mo -e .bz2 --directory /usr/portage/packages/

Flags:
  -d, --directory string           Artefacts directory with .tbz2 files.
      --hash-empty                 If create a fake hash for empty packages or use 00000000000000000000000000000000.
  -f, --hashfile string            Path of hashfile where write checksum.
                                   Default output on stdout with format: HASH <CHECKSUM> <PACKAGE>
  -h, --help                       help for hash
  -i, --ignore strings             File to ignore.
      --ignore-errors              Ignore errors with broken tarball.
  -e, --ignore-extension strings   Extension to ignore.
  -p, --package strings            Path of package to check.
      --stdin                      Read package data from stdin

Global Flags:
  -c, --concurrency       Enable concurrency process.
  -l, --logfile string    Logfile Path. Optional.
  -L, --loglevel string   Set logging level.
                          [DEBUG, INFO, WARN, ERROR] (default "INFO")
  -v, --verbose           Enable verbose logging on stdout.

```

## Task for Next Release:

  * Add support to stdin processing

  * Add support to different checksum algorithms
