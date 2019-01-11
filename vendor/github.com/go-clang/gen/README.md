# go-clang/gen [![GoDoc](https://godoc.org/github.com/go-clang/gen?status.png)](https://godoc.org/github.com/go-clang/gen) [![Build Status](https://travis-ci.org/go-clang/gen.svg?branch=master)](https://travis-ci.org/go-clang/gen) [![Coverage Status](https://coveralls.io/repos/go-clang/gen/badge.png?branch=master)](https://coveralls.io/r/go-clang/gen?branch=master)

Generate native Go bindings for Clang's C API.

## I found a bug/missing a feature in go-clang

Please go through the [open issues](https://github.com/go-clang/gen/issues) in the tracker first. If you cannot find your request just open up a [new issue](https://github.com/go-clang/gen/issues/new).

## Where are the bindings?

The Go bindings are placed in their own repositories to provide the correct bindings for the corresponding Clang version:

- [v3.4](https://github.com/go-clang/v3.4)
- [v3.6](https://github.com/go-clang/v3.6)
- [v3.7](https://github.com/go-clang/v3.7)
- [v3.8](https://github.com/go-clang/v3.8)
- [v3.9](https://github.com/go-clang/v3.9)
- [v4.0](https://github.com/go-clang/v4.0)
- [v5.0](https://github.com/go-clang/v5.0)

## Install `go-clang/gen` to contribute to go-clang

```bash
CGO_LDFLAGS="-L`llvm-config --libdir`" \
  go get github.com/go-clang/bootstrap github.com/go-clang/gen
```

## How to contribute?

You want to contribute to go-clang? GREAT! If you are here because of a bug you want to fix or a feature you want to add, you can just read on. Otherwise we have a list of [open issues in the tracker](https://github.com/go-clang/gen/issues). Just choose something you think you can work on and discuss your plans in the issue by commenting on it.

This repository, [gen](https://github.com/go-clang/gen), holds the code to generate new bindings from headers of Clang's C API. These bindings are then bootstrapped using the [bootstrap](https://github.com/go-clang/bootstrap) repository. The `bootstrap` repository holds all basic files, like the CI configuration and a Makefile, as well as some additional code to make the bindings more complete and powerful.

To ease the development process we have our own development environment based on [Vagrant](https://www.vagrantup.com/). The provided Vagrantfile executed in the root of the repository will setup an Ubuntu VM with our currently used Go version as well as Clang 3.4 and will set up everything that is needed to development and handle new versions of Clang.

> **Please note**, only the major and minor version must be declared if a Clang version is needed in a command.

### Generate bindings for the current Clang version (VM)

Make sure that the `go-clang-gen` command is up to date using `make install` in the repository's root directory. After that execute `go-clang-gen` which will generate bindings in your current directory.

### Switch to a different Clang version (VM)

Replace `3.4` with the Clang version you want to switch to.

```bash
$GOPATH/src/github.com/go-clang/gen/scripts/switch-clang-version.sh 3.4
```

This command will install and configure everything that is needed to develop with the given Clang version. The command will however not generate new bindings for the version.

### Do a PR

Every PR must be prepared using the following commands:

```bash
cd $GOPATH/src/github.com/go-clang/gen
scripts/switch-clang-version.sh 3.4
make install
make test
make lint
```

This will switch to the current Clang version for the `go-clang-gen` command, execute all tests and process the source code with the project's linters. Make sure that you do not introduce new linting problems.

### Regenerate the bindings for the `bootstrap` repository

The `bootstrap` repository holds the base for all version repositories it must therefore be updated if something changes in the `gen` repository to regenerate all version repositories. This tedious task is automated by the following command.

```bash
$GOPATH/src/github.com/go-clang/gen/scripts/update-bootstrap.sh
```

Any changes must be manually verified and pushed via a feature branch.

## Maintainer documentation

The following sections are specific to the maintaining process.

> **Please note**, only the major and minor version must be declared if a Clang version is needed in a command.

### Create a new Clang version (VM)

Every now and then a new Clang version emerges which needs to be generated using the `go-clang-gen` command. The new version has to be available using the VM's and CI's packages. Otherwise, we cannot correctly test and therefore support the version.

If a new version is available create a repository on GitHub named `v<MAJOR>.<MINOR>` and set the repository description to `Go bindings for Clang's C API v<MAJOR>.<MINOR>`. Disable all repository features, e.g. `Issues` and `Wiki`. Enable the repository on TravisCI before you push anything to the repository. Lastly, execute the following command inside the development VM.

```bash
$GOPATH/src/github.com/go-clang/gen/scripts/create-clang-version.sh 3.4
```

This will create a new repository `v3.4` in the `$GOPATH/src/github.com/go-clang` directory and initialize it using the `bootstrap` repository. The command also generates, installs, configures and tests bindings for the given Clang version. The changes must then be manually verified, added, committed and pushed to the already set up remote "origin".

After the repository is online godoc.org needs to know about the new repository. This can be achieved by visiting [https://godoc.org/github.com/go-clang/v&lt;MAJOR&gt;.&lt;MINOR&gt;/clang](https://godoc.org/github.com/go-clang/v<MAJOR>.<MINOR>/clang).

### Update a branch with a new Clang version (VM)

Every now and then a new Clang subminor version is released. The given version can be supported if packages are available inside the VM and CI. The following command can then be executed inside the development VM.

```bash
$GOPATH/src/github.com/go-clang/gen/scripts/update-clang-version.sh 3.4
```

This will rebase the latest commits of the `bootstrap` repository onto the commits of the `v3.4` repository.  The command also generates, installs, configures and tests bindings for the given Clang version. The changes must then be manually verified, added, committed and pushed to the already set up remote "origin".

### Update branches with a new `go-clang/gen` version (VM)

If the `go-clang-gen` command changes its generation output, all branches need to be updated which is basically just updating for a new Clang version.

# License

This project, like all go-clang projects, is licensed under a BSD-3 license which can be found in the [LICENSE file](https://github.com/go-clang/license/blob/master/LICENSE) in [go-clang's license repository](https://github.com/go-clang/license)
