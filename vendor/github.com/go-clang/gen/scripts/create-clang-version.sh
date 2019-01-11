#!/bin/bash

set -exuo pipefail

if [ -z "$1" ]; then
	exit
fi

export LLVM_VERSION=$1

# Switch Clang version
$(dirname "$0")/switch-clang-version.sh $LLVM_VERSION

# Initialize the repository
cd $GOPATH/src/github.com/go-clang/
mkdir v${LLVM_VERSION}
cd v${LLVM_VERSION}

git clone https://github.com/go-clang/bootstrap.git .
git remote rename origin bootstrap
git remote add origin git@github.com:go-clang/v${LLVM_VERSION}.git

# Generate the new Clang version
$(dirname "$0")/generate-and-test.sh $LLVM_VERSION
