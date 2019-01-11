#!/bin/bash

set -exuo pipefail

if [ -z "$1" ]; then
	exit
fi

export LLVM_VERSION=$1

# Switch Clang version
$(dirname "$0")/switch-clang-version.sh $LLVM_VERSION

# Update the repository
cd $GOPATH/src/github.com/go-clang/v${LLVM_VERSION}

git fetch --prune bootstrap

git rebase master bootstrap/master

# Generate the new Clang version
$(dirname "$0")/generate-and-test.sh $LLVM_VERSION
