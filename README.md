# Sabayon Packages Checker

[![Coverage Status](https://coveralls.io/repos/github/Sabayon/pkgs-checker/badge.svg?branch=master)](https://coveralls.io/github/Sabayon/pkgs-checker?branch=master)

Tool used for different tasks on Sabayon build process.

```bash
$# pkgs-checker --help
Sabayon packages checker

Usage:
   [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  entropy     Entropy tool.
  filter      Filter bin-host packages/directory.
  hash        Hashing packages
  help        Help about any command
  pkg         Package parsing tool.
  pkglist     Manage pkglist files.
  portage     Manage portage files/metafiles.
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

## *portage gen-pkgs-uses* command

Read the content of Portage `/var/db/pkg` and then it permits to generate the existing `package.use` file of
all packages or all the packages selected on filter file.

The filter file is in this format:

```yaml
# Define the regex for the use flags to exclude
use_filters:
  - "^userland_"
  - "^kernel_"
  - "^x86"
  - "^x64"
  - "^ppc"
  - "^arm"
  - "^amd64"
  - "^prefix"
  - "^m68k"
  - "^ia64"
  - "^riscv"
  - "^s390"
  - "^hppa"
  - "^s390"
  - "^mips"
  - "^alpha"
  - "^sparc"
  - "^elibc_"
  - "^abi_"

# Define the categories to elaborate.
# If the list is empty all categories are elaborated.
categories:
#  - media-libs

# Define the regex string of the packages to elaborate in format cat/pkg:slot.
# If the list is empty all packages are elaborated.
pkgs_filters:
  #- cat/foo
  #- ^media
```

With the options `--luet-portage-converter-format` and `--treePath` is possible to generate
the YAML content used by the `luet-portage-converter` tool.

## *portage metadata*

Show metadata of a Portage package.

```bash
$> pkgs-checker portage metadata --help
Show metadata of a package.

Usage:
   portage metadata cat/pkg[:slot] [OPTIONS] [flags]

Flags:
  -p, --db-pkgs-dir-path string         Path of the portage metadata. (default "/var/db/pkg")
      --filter-opts string              Using filter rules through YAML file.
  -h, --help                            help for metadata
  -j, --json                            Output in JSON format
      --luet-portage-converter-format   Generate luet-portage-converter YAML output.
      --treePath string                 Define the tree path to use on luet-portage-converter artefacts. (default "packages/atoms")

Global Flags:
  -c, --concurrency       Enable concurrency process.
  -l, --logfile string    Logfile Path. Optional.
  -L, --loglevel string   Set logging level.
                          [DEBUG, INFO, WARN, ERROR] (default "INFO")
  -v, --verbose           Enable verbose logging on stdout.
```

Using `-j|--json` option to see data in JSON format.

```bash
$> pkgs-checker portage metadata sys-devel/bc:0 -j | jq
[
  {
    "package": {
      "name": "bc",
      "category": "sys-devel",
      "version": "1.07.1",
      "version_suffix": "-r3",
      "slot": "0",
      "Condition": 5,
      "repository": "gentoo",
      "use_flags": [
        "abi_x86_64",
        "-libedit",
        "readline",
        "-sh",
        "-static"
      ],
      "license": "GPL-2 LGPL-2.1"
    },
    "iuse": [
      "libedit",
      "readline",
      "static"
    ],
    "iuse_effective": [
      "abi_x86_64",
      "alpha",
      "amd64",
      "amd64-fbsd",
      "amd64-linux",
      "arm",
      "arm64",
      "elibc_AIX",
      "elibc_Cygwin",
      "elibc_Darwin",
      "elibc_DragonFly",
      "elibc_FreeBSD",
      "elibc_HPUX",
      "elibc_Interix",
      "elibc_NetBSD",
      "elibc_OpenBSD",
      "elibc_SunOS",
      "elibc_Winnt",
      "elibc_bionic",
      "elibc_glibc",
      "elibc_mingw",
      "elibc_mintlib",
      "elibc_musl",
      "elibc_uclibc",
      "hppa",
      "ia64",
      "kernel_AIX",
      "kernel_Darwin",
      "kernel_FreeBSD",
      "kernel_HPUX",
      "kernel_NetBSD",
      "kernel_OpenBSD",
      "kernel_SunOS",
      "kernel_Winnt",
      "kernel_freemint",
      "kernel_linux",
      "libedit",
      "m68k",
      "m68k-mint",
      "mips",
      "ppc",
      "ppc-aix",
      "ppc-macos",
      "ppc64",
      "ppc64-linux",
      "prefix",
      "prefix-guest",
      "prefix-stack",
      "readline",
      "riscv",
      "s390",
      "sh",
      "sparc",
      "sparc-solaris",
      "sparc64-solaris",
      "static",
      "userland_BSD",
      "userland_GNU",
      "x64-cygwin",
      "x64-macos",
      "x64-solaris",
      "x86",
      "x86-cygwin",
      "x86-fbsd",
      "x86-linux",
      "x86-macos",
      "x86-solaris",
      "x86-winnt"
    ],
    "use": [
      "abi_x86_64",
      "amd64",
      "elibc_glibc",
      "kernel_linux",
      "readline",
      "userland_GNU"
    ],
    "eapi": "6",
    "cxxflags": "-O2 -march=x86-64 -pipe",
    "cflags": "-O2 -march=x86-64 -pipe",
    "ldflags": "-Wl,-O1 -Wl,--as-needed",
    "chost": "x86_64-pc-linux-gnu",
    "rdepend": ">=sys-libs/readline-4.1:0/7= >=sys-libs/ncurses-5.2:0/6=",
    "depend": ">=sys-libs/readline-4.1:0/7= >=sys-libs/ncurses-5.2:0/6= sys-devel/flex virtual/yacc",
    "requires": "x86_64: libc.so.6 libreadline.so.7",
    "keywords": "~alpha ~amd64 ~arm ~arm64 ~hppa ~ia64 ~m68k ~mips ~ppc ~ppc64 ~riscv ~s390 ~sh ~sparc ~x86 ~ppc-aix ~x64-cygwin ~amd64-linux ~x86-linux ~ppc-macos ~x64-macos ~x86-macos ~m68k-mint ~sparc-solaris ~sparc64-solaris ~x64-solaris ~x86-solaris",
    "size": "337574",
    "build_time": "1582248755",
    "cbuild": "x86_64-pc-linux-gnu",
    "counter": "46679",
    "defined_phases": "compile configure prepare",
    "description": "Handy console-based calculator utility",
    "features": "assume-digests binpkg-docompress binpkg-dostrip binpkg-logs compressdebug config-protect-if-modified distlocks ebuild-locks fixlafiles ipc-sandbox merge-sync multilib-strict network-sandbox news parallel-fetch pid-sandbox preserve-libs protect-owned sandbox sfperms splitdebug strict unknown-features-warn unmerge-logs unmerge-orphans userfetch userpriv usersandbox usersync xattr",
    "homepage": "https://www.gnu.org/software/bc/bc.html",
    "inherited": "desktop estack epatch toolchain-funcs multilib ltprune preserve-libs vcs-clean eutils flag-o-matic",
    "needed": "/usr/bin/dc libc.so.6\n/usr/bin/bc libreadline.so.7,libc.so.6",
    "ebuild": "# Copyright 1999-2020 Gentoo Authors\n# Distributed under the terms of the GNU General Public License v2\n\nEAPI=\"6\"\n\ninherit flag-o-matic toolchain-funcs\n\nDESCRIPTION=\"Handy console-based calculator utility\"\nHOMEPAGE=\"https://www.gnu.org/software/bc/bc.html\"\nSRC_URI=\"mirror://gnu/bc/${P}.tar.gz\"\n\nLICENSE=\"GPL-2 LGPL-2.1\"\nSLOT=\"0\"\nKEYWORDS=\"~alpha ~amd64 ~arm ~arm64 ~hppa ~ia64 ~m68k ~mips ~ppc ~ppc64 ~riscv ~s390 ~sh ~sparc ~x86 ~ppc-aix ~x64-cygwin ~amd64-linux ~x86-linux ~ppc-macos ~x64-macos ~x86-macos ~m68k-mint ~sparc-solaris ~sparc64-solaris ~x64-solaris ~x86-solaris\"\nIUSE=\"libedit readline static\"\n\nRDEPEND=\"\n\t!readline? ( libedit? ( dev-libs/libedit:= ) )\n\treadline? (\n\t\t>=sys-libs/readline-4.1:0=\n\t\t>=sys-libs/ncurses-5.2:=\n\t)\n\"\nDEPEND=\"\n\t${RDEPEND}\n\tsys-devel/flex\n\tvirtual/yacc\n\"\n\nPATCHES=(\n\t\"${FILESDIR}/${PN}-1.07.1-no-ed-its-sed.patch\"\n)\n\nsrc_prepare() {\n\tdefault\n\n\t# Avoid bad build tool usage when cross-compiling.  #627126\n\ttc-is-cross-compiler && eapply \"${FILESDIR}/${PN}-1.07.1-use-system-bc.patch\"\n}\n\nsrc_configure() {\n\tlocal myconf=(\n\t\t$(use_with readline)\n\t)\n\tif use readline ; then\n\t\tmyconf+=( --without-libedit )\n\telse\n\t\tmyconf+=( $(use_with libedit) )\n\tfi\n\tuse static && append-ldflags -static\n\n\teconf \"${myconf[@]}\"\n\n\t# Do not regen docs -- configure produces a small fragment that includes\n\t# the version info which causes all pages to regen (newer file). #554774\n\ttouch -r doc doc/*\n}\n\nsrc_compile() {\n\temake AR=\"$(tc-getAR)\"\n}",
    "content": [
      {
        "type": "dir",
        "file": "/usr"
      },
      {
        "type": "dir",
        "file": "/usr/bin"
      },
      {
        "type": "obj",
        "file": "/usr/bin/bc",
        "hash": "6fcde9e94b835cc7c297d065ed3fa929",
        "timestamp": "1582248756"
      },
      {
        "type": "obj",
        "file": "/usr/bin/dc",
        "hash": "62eb780a0d49ea03b72de1ed0d1c2874",
        "timestamp": "1582248756"
      },
      {
        "type": "dir",
        "file": "/usr/share"
      },
      {
        "type": "dir",
        "file": "/usr/share/man"
      },
      {
        "type": "dir",
        "file": "/usr/share/man/man1"
      },
      {
        "type": "obj",
        "file": "/usr/share/man/man1/bc.1.bz2",
        "hash": "cc4ce6ce3f974799beb6bdce52f6e75f",
        "timestamp": "1582248752"
      },
      {
        "type": "obj",
        "file": "/usr/share/man/man1/dc.1.bz2",
        "hash": "f22e5292e72c127d73eb16d5ea5cb169",
        "timestamp": "1582248752"
      },
      {
        "type": "dir",
        "file": "/usr/share/info"
      },
      {
        "type": "obj",
        "file": "/usr/share/info/dc.info.bz2",
        "hash": "b63625653f5a7a8f1b67475e6337f556",
        "timestamp": "1582248752"
      },
      {
        "type": "obj",
        "file": "/usr/share/info/bc.info.bz2",
        "hash": "1b1de10f3db04d8d9d1265382ad10417",
        "timestamp": "1582248752"
      },
      {
        "type": "dir",
        "file": "/usr/share/doc"
      },
      {
        "type": "dir",
        "file": "/usr/share/doc/bc-1.07.1-r3"
      },
      {
        "type": "obj",
        "file": "/usr/share/doc/bc-1.07.1-r3/ChangeLog.bz2",
        "hash": "895b8be7441d208ab172aa7ad0ad6a70",
        "timestamp": "1582248752"
      },
      {
        "type": "obj",
        "file": "/usr/share/doc/bc-1.07.1-r3/FAQ.bz2",
        "hash": "3427caa93df52778061f689e06f59b5a",
        "timestamp": "1582248753"
      },
      {
        "type": "obj",
        "file": "/usr/share/doc/bc-1.07.1-r3/AUTHORS.bz2",
        "hash": "bd85627e0a94fa543bb00f37a63f67a8",
        "timestamp": "1582248753"
      },
      {
        "type": "obj",
        "file": "/usr/share/doc/bc-1.07.1-r3/README.bz2",
        "hash": "441f92ee728835fb8705d58e9b18f388",
        "timestamp": "1582248752"
      },
      {
        "type": "obj",
        "file": "/usr/share/doc/bc-1.07.1-r3/NEWS.bz2",
        "hash": "63def554e7ec3255a329277bcd30c8df",
        "timestamp": "1582248753"
      },
      {
        "type": "dir",
        "file": "/usr/lib"
      }
    ]
  }
]
```

## *portage metadata* command

Permit to read and generate the metadata files used by Portage to a defined path or one or more
packages.

```bash
$> pkgs-checker portage gen-metadata --help
Generate metadata of a package to a specific path.

Usage:
   portage gen-metadata cat/pkg[:slot] ... catN/pkgN[:slot] [OPTIONS] [flags]

Flags:
  -p, --db-pkgs-dir-path string   Path of the portage metadata. (default "/var/db/pkg")
  -h, --help                      help for gen-metadata
      --to string                 Generate medata tree to the specified path. (default "/gen-metadata")

Global Flags:
  -c, --concurrency       Enable concurrency process.
  -l, --logfile string    Logfile Path. Optional.
  -L, --loglevel string   Set logging level.
                          [DEBUG, INFO, WARN, ERROR] (default "INFO")
  -v, --verbose           Enable verbose logging on stdout.
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
