#!/usr/bin/env bash
set -eu

# Get the directory path.
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
scripts_dir="$( cd -P "$( dirname "$SOURCE" )/" && pwd )"
build_dir="$(cd $scripts_dir/.. && pwd)"
outputs_dir="$(cd $build_dir/outputs && pwd)"
repo_dir="$(cd $build_dir/.. && pwd)"

# Move the parent (repository) directory
cd "$repo_dir"

# Load config
source $scripts_dir/config

echo "Removing old files."
rm -rf $outputs_dir/dev/*

COMMIT_HASH=`git log --pretty=format:%H -n 1`

echo "Building dev binary..."
echo "PRODUCT_NAME: $PRODUCT_NAME"
echo "PRODUCT_VERSION: $PRODUCT_VERSION"
echo "COMMIT_HASH: $COMMIT_HASH"

go build \
    -ldflags=" -w \
        -X github.com/kohkimakimoto/$PRODUCT_NAME/$PRODUCT_NAME.CommitHash=$COMMIT_HASH \
        -X github.com/kohkimakimoto/$PRODUCT_NAME/$PRODUCT_NAME.Version=$PRODUCT_VERSION" \
    -o="$outputs_dir/dev/$PRODUCT_NAME" \
    ./cmd/${PRODUCT_NAME}
echo "Results:"
ls -hl "$outputs_dir/dev"