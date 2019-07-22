#!/bin/sh

set -e

if [ ! -f "build/env.sh" ]; then
    echo "$0 must be run from the root of the repository."
    exit 2
fi

# Create fake Go workspace if it doesn't exist yet.
workspace="$PWD/build/_workspace"
root="$PWD"
vntdir="$workspace/src/github.com/vntchain"
if [ ! -L "$vntdir/go-vnt" ]; then
    mkdir -p "$vntdir"
    cd "$vntdir"
    ln -s ../../../../../. go-vnt
    cd "$root"
fi

# Set up the environment to use the workspace.
GOPATH="$workspace"
export GOPATH

# Run the command inside the workspace.
cd "$vntdir/go-vnt"
PWD="$vntdir/go-vnt"

# Launch the arguments with the configured environment.
exec "$@"
