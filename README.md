# Sabayon Packages Checker

## Why

Artefacts created with emerge could be compressed with tar+bz2 not in ordered way and some files must be skipped to avoid injection of some packages that are equal.

## How

*pkgs-checker* processes tarball file and create an MD5 checksum for any file inside package that are not be skipped (by command line options). At EOF create a new MD5 with all MD5 checksum bytes plus list of directories found sorted.

Normally, files to skip are .pyc,.pyo,.mo that contains timestamp data that generate false events for package injection.

## Features to Complete for release 0.1.0

  * Add support to process a directory

  * Add support to stdin processing

  * Add support to different checksum algorithms
