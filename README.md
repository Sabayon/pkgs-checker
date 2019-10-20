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
  pkglist     Manage pkglist files.
  sark        Manage sark process.

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

## *filter* command

*filter* sub-command permits to filter packages that must be excluded to injection phase of a specific repository.

### Usage

```bash
pkgs-checker filter --help
Filter bin-host packages/directory.

Usage:
   filter [OPTIONS] [flags]

Examples:
$> pkgs-checker filter --binhost-dir /usr/portage/packages/ --sark-config ./rules.yaml

Flags:
  -d, --binhost-dir string          bin-hosts directory where filter packages.
      --category strings            Filter specific category.
      --dry-run                     Only check file to remove.
  -t, --filter-type string          Define filter type (whitelist|blacklist)
  -h, --help                        help for filter
  -p, --package strings             Filter specific package.
  -r, --report-prefix-path string   Report file prefix to use for both filtered and unfiltered packages...
  -f, --sark-config string          SARK Configuration file with filter rules or targets.

Global Flags:
  -c, --concurrency       Enable concurrency process.
  -l, --logfile string    Logfile Path. Optional.
  -L, --loglevel string   Set logging level.
                          [DEBUG, INFO, WARN, ERROR] (default "INFO")
  -v, --verbose           Enable verbose logging on stdout.

```

## *pkglist* command

*pkglist* command permits to work with pkglist files.

```bash

$# pkgs-checker pkglist --help
Manage pkglist files.

Usage:
   pkglist [command]

Available Commands:
  create      Create pkglist file.
  intersect   Search duplicate package between multiple pkglist.
  show        Show pkglist from one or multiple resources.

Flags:
  -h, --help   help for pkglist

Global Flags:
  -c, --concurrency       Enable concurrency process.
  -l, --logfile string    Logfile Path. Optional.
  -L, --loglevel string   Set logging level.
                          [DEBUG, INFO, WARN, ERROR] (default "INFO")
  -v, --verbose           Enable verbose logging on stdout.

Use " pkglist [command] --help" for more information about a command.
```

### *pkglist intersect* command

Search packages available in multiple pkglist.

```
$# pkgs-checker pkglist intersect --help
Search duplicate package between multiple pkglist.

Usage:
   pkglist intersect [OPTIONS] [flags]

Examples:
$> pkgs-checker pkglist intersect -r https://server1/sbi/namespace/base-arm/base-arm-binhost/base-arm.pkglist,https://server2/sbi/namespace/core-arm/core-arm-binhost/core-arm.pkglist

Flags:
  -h, --help              help for intersect
  -r, --pkglist strings   Path or URL of pkglist resource.
  -q, --quiet             Quiet output.

Global Flags:
  -c, --concurrency       Enable concurrency process.
  -l, --logfile string    Logfile Path. Optional.
  -L, --loglevel string   Set logging level.
                          [DEBUG, INFO, WARN, ERROR] (default "INFO")
  -v, --verbose           Enable verbose logging on stdout.
```

### *pkglist create* command

Print or write a file with package list available to a binhost directory.

```bash
$# pkgs-checker pkglist create --help
Create pkglist file.

Usage:
   pkglist create [OPTIONS] [flags]

Flags:
  -d, --binhost-dir string    bin-hosts directory where compute pkglist.
  -h, --help                  help for create
  -f, --pkglist-file string   Path of pkglist file.
                              Default output to stdout with format: category/pkgname-pkgversion

Global Flags:
  -c, --concurrency       Enable concurrency process.
  -l, --logfile string    Logfile Path. Optional.
  -L, --loglevel string   Set logging level.
                          [DEBUG, INFO, WARN, ERROR] (default "INFO")
  -v, --verbose           Enable verbose logging on stdout.
```

### *pkglist show* command

Retrieve list of packages from multiple resources (URL or local files).

```
$# pkgs-checker pkglist show --help
Show pkglist from one or multiple resources.

Usage:
   pkglist show [OPTIONS] [flags]

Examples:
$> pkgs-checker pkglist show -r https://server1/sbi/namespace/base-arm/base-arm-binhost/base-arm.pkglist,https://server2/sbi/namespace/core-arm/core-arm-binhost/core-arm.pkglist

Flags:
  -h, --help              help for show
  -p, --parse-pkgname     Parse package version string and hide entropy revision.
  -r, --pkglist strings   Path or URL of pkglist resource.
  -q, --quiet             Quiet output.

Global Flags:
  -c, --concurrency       Enable concurrency process.
  -l, --logfile string    Logfile Path. Optional.
  -L, --loglevel string   Set logging level.
                          [DEBUG, INFO, WARN, ERROR] (default "INFO")
  -v, --verbose           Enable verbose logging on stdout.
```

## *sark* command

Commands for help on SARK processes.

```
$# pkgs-checker sark
Manage sark process.

Usage:
   sark [command]

Available Commands:
  compare     Compare sark targets with pkglist files

Flags:
  -h, --help   help for sark

Global Flags:
  -c, --concurrency       Enable concurrency process.
  -l, --logfile string    Logfile Path. Optional.
  -L, --loglevel string   Set logging level.
                          [DEBUG, INFO, WARN, ERROR] (default "INFO")
  -v, --verbose           Enable verbose logging on stdout.

Use " sark [command] --help" for more information about a command.
```

### *sark compare* command

Check for packages not defined on SARK build files or packages that are defined in SARK build files
but not available on packages list.

```
$ ./pkgs-checker sark compare --help
Compare sark targets with pkglist files

Usage:
   sark compare [OPTIONS] [flags]

Examples:

Show targets not present on package list:
$> pkgs-checker sark compare -s core-staging1-build.yaml -p core-arm.pkglist -v -m

Show packages not present between SARK targets:
$> pkgs-checker sark compare -s core-staging1-build.yaml -p core-arm.pkglist -v -t


Flags:
  -h, --help                    help for compare
  -m, --missing-packages        Show targets not present on pkglist(s).
  -t, --missing-targets         Show packages not present on target.
  -p, --pkglist-files strings   Path or URL of pkglist resources.
  -s, --sark-files strings      Path or URL of sark config resources.

Global Flags:
  -c, --concurrency       Enable concurrency process.
  -l, --logfile string    Logfile Path. Optional.
  -L, --loglevel string   Set logging level.
                          [DEBUG, INFO, WARN, ERROR] (default "INFO")
  -v, --verbose           Enable verbose logging on stdout.
```

## *hash* command

### Why

Artefacts created with emerge could be compressed with tar+bz2 not in ordered way and some files must be skipped to avoid injection of some packages that are equal.

### How

*pkgs-checker* processes tarball file and create an MD5 checksum for any file inside package that are not be skipped (by command line options). At EOF create a new MD5 with all MD5 checksum bytes plus list of directories found sorted.

Normally, files to skip are .pyc,.pyo,.mo that contains timestamp data that generate false events for package injection.

### Usage

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

### Task for Next Release:

  * Add support to stdin processing

  * Add support to different checksum algorithms
