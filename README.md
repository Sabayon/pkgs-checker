# Sabayon Packages Checker

## Why

Artefacts created with emerge could be compressed with tar+bz2 not in ordered way and some files must be skipped to avoid injection of some packages that are equal.

## How

*pkgs-checker* processes tarball file and create an MD5 checksum for any file inside package that are not be skipped (by command line options). At EOF create a new MD5 with all MD5 checksum bytes plus list of directories found sorted.

Normally, files to skip are .pyc,.pyo,.mo that countains timestamp data that generate false events for package injection.
