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
rm -rf $outputs_dir/packaging/*

echo "Building RPM packages..."
cd $scripts_dir/rpm
for image in 'kohkimakimoto/rpmbuild:el7'; do
    docker run \
        --env DOCKER_IMAGE=${image}  \
        --env PRODUCT_NAME=${PRODUCT_NAME}  \
        --env PRODUCT_VERSION=${PRODUCT_VERSION}  \
        --env COMMIT_HASH=${COMMIT_HASH}  \
        -v $repo_dir:/tmp/repo \
        -w /tmp/repo \
        --rm \
        ${image} \
        bash ./build/scripts/rpm/run.sh
done

cd "$repo_dir"

echo "Results:"
ls -hl "$outputs_dir/packaging"
