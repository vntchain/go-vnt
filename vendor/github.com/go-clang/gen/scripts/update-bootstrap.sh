#!/bin/bash

set -exuo pipefail

# Switch to the correct Clang version for the bootstrap and the gen repositories
$GOPATH/src/github.com/go-clang/gen/scripts/switch-clang-version.sh 3.4

# Install the current go-clang-gen command
make -C $GOPATH/src/github.com/go-clang/gen install

# Generate the Clang bindings
cd $GOPATH/src/github.com/go-clang/bootstrap/clang/

rm -rf clang-c/
rm -f *_gen.go

go-clang-gen

cd ..

# Install and test the bindings
make install
make test

# Show the current state of the repository
git status
