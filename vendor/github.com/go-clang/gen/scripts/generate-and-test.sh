#!/bin/bash

set -exuo pipefail

if [ -z "$1" ]; then
	exit
fi

export LLVM_VERSION=$1

# Generate the new Clang version
cd $GOPATH/src/github.com/go-clang/v${LLVM_VERSION}/clang/

rm -rf clang-c/
rm -f *_gen.go

go-clang-gen

cd ..

# Change versions in files
sed -i -e "s/3.4/${LLVM_VERSION}/g" .travis.yml
find . -type f -not -path '*/\.*' -exec sed -i -e "s/bootstrap/v${LLVM_VERSION}/g" {} +

# Install and test the version
make install
make test

# Show the current state of the repository
git status
