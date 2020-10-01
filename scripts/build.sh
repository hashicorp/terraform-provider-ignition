#!/usr/bin/env bash
#
# This script is used to build releases for multiple platforms.
# * Ensure you have added $GOPATH/bin to $PATH
# 1. Checkout the release tag: # git checkout <tag>
# 2. Run this script: ./scripts/build.sh

if [[ ! -d "./scripts" ]]; then
  echo "ERROR: scripts not found in pwd."
  echo "Please run this from the root of the terraform provider repository"
  exit 1
fi

GIT_TAG=$(git describe --tags)

# The arch/os
XC_ARCH=${XC_ARCH:-"amd64"}
XC_OS=${XC_OS:-linux darwin windows}
XC_EXCLUDE_OSARCH="!darwin/arm !darwin/386"

# Delete the old dir
echo "==> Cleaning up..."
rm -f bin/* && rm -rf pkg/* && mkdir -p bin/

# If required for only current platform set LOCAL=1
if [ "${LOCAL}x" != "x" ]; then
    XC_OS=$(go env GOOS)
    XC_ARCH=$(go env GOARCH)
fi

if ! which gox > /dev/null; then
    echo "==> Installing gox..."
    go install github.com/mitchellh/gox
fi

# instruct gox to build statically linked binaries
export CGO_ENABLED=0

echo "==> Building..."
LD_FLAGS="-s -w"
gox \
    -os="${XC_OS}" \
    -arch="${XC_ARCH}" \
    -osarch="${XC_EXCLUDE_OSARCH}" \
    -ldflags "${LD_FLAGS}" \
    -output "pkg/{{.OS}}_{{.Arch}}/terraform-provider-ignition_${GIT_TAG}" \
    .


echo "==> Packaging..."
for PLATFORM in $(find ./pkg -mindepth 1 -maxdepth 1 -type d); do
    OSARCH=$(basename ${PLATFORM})
    echo "--> ${OSARCH}"

    pushd $PLATFORM >/dev/null 2>&1
    zip ../terraform-provider-ignition_${GIT_TAG}_${OSARCH}.zip ./*
    popd >/dev/null 2>&1
done

echo
echo "==> Completed!"
