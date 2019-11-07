#!/usr/bin/env bash
set -eu

echo "Running packaging script in '$DOCKER_IMAGE' container..."
echo "PRODUCT_NAME: $PRODUCT_NAME"
echo "PRODUCT_VERSION: $PRODUCT_VERSION"
echo "COMMIT_HASH: $COMMIT_HASH"

echo "Copying files..."

repo_dir=$(pwd)
platform=el${RHEL_VERSION}

cp -pr build/scripts/rpm/SPECS $HOME/rpmbuild/
cp -pr build/scripts/rpm/SOURCES $HOME/rpmbuild/
cp -pr build/outputs/dist/${PRODUCT_NAME}_linux_amd64.zip $HOME/rpmbuild/SOURCES/${PRODUCT_NAME}_linux_amd64.zip

echo "Building RPM..."
cd $HOME
rpmbuild \
    --define "_product_name ${PRODUCT_NAME}" \
    --define "_product_version ${PRODUCT_VERSION}" \
    --define "_rhel_version ${RHEL_VERSION}" \
    -ba rpmbuild/SPECS/${PRODUCT_NAME}.spec

echo "Copying generated files to shared folder..."
cd $repo_dir

mkdir -p build/outputs/packaging/${platform}
cp -pr $HOME/rpmbuild/RPMS build/outputs/packaging/${platform}
cp -pr $HOME/rpmbuild/SRPMS build/outputs/packaging/${platform}
