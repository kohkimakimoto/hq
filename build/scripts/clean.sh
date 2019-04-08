#!/usr/bin/env bash
set -eu

# Get the directory path.
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
scripts_dir="$( cd -P "$( dirname "$SOURCE" )/" && pwd )"
build_dir="$(cd $scripts_dir/.. && pwd)"
outputs_dir="$(cd $build_dir/outputs && pwd)"
repo_dir="$(cd $build_dir/.. && pwd)"

echo "Cleaning old files."
rm -rf $outputs_dir/*
