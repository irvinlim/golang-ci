# golang-ci

Provides CI tools for Go programs.

## Installation

```sh
GO111MODULE=off go get github.com/irvinlim/golang-ci
```

## Usage

```sh
$ golang-ci
CI tool for Golang

Usage:
  golang-ci [command]

Available Commands:
  help             Help about any command
  lint             Runs linters using golangci-lint
  lint-install     Installs golangci-lint
  lint-install-all Installs all versions of golangci-lint

Flags:
  -h, --help   help for golang-ci

Use "golang-ci [command] --help" for more information about a command.
```

### `lint`

```sh
$ golang-ci lint --help
Executes a specific version of golangci-lint, downloading to $GOPATH/bin as necessary.

Example usage:
  # Run default version. Note that the `--` is necessary to pass arguments to golangci-lint.
  golang-ci lint -- run -v

  # Run a specific version of golangci-lint.
  golang-ci lint --version v1.31.0 -- run -v

Usage:
  golang-ci lint [flags]

Flags:
  -h, --help             help for lint
      --version string   Version of golangci-lint to use. If not specified, will use latest version available.
```

### `lint-install`

```sh
$ golang-ci lint-install --help
Installs a specific version of golangci-lint, downloading to $GOPATH/bin as necessary.

Example usage:
  # Installs the latest version at $GOPATH/bin/golangci-lint.
  golang-ci lint-install latest

  # Installs a specific version of golangci-lint at $GOPATH/bin/golangci-lint@v1.31.0.
  golang-ci lint-install v1.31.0

Usage:
  golang-ci lint-install VERSION [flags]

Flags:
  -h, --help   help for lint-install
```

### `lint-install-all`

```sh
$ golang-ci lint-install-all --help
Installs all versions of golangci-lint, downloading to $GOPATH/bin.

Usage:
  golang-ci lint-install-all [flags]

Flags:
  -h, --help                 help for lint-install-all
      --min-version string   Minimum version of golangci-lint that should be installed.
```

## Features

### Virtual environment for `golangci-lint`

Runs Go linters using a specific version of [`golangci-lint`](https://golangci-lint.run/). 

When working on multiple Go projects which may have different versions of golangci-lint, having mismatched versions of golangci-lint in CI versus the copy of the `golangci-lint` binary on your local computer may cause your CI jobs to fail inadvertently.

This tool aims to provide a "virtual environment" for golangci-lint, by lazily downloading a specific version of `golangci-lint` if it is not present.

Example usage:

```sh
# Previously you might be calling golangci-lint as follows:
golangci-lint run -v --timeout=5m

# Instead, to use a specific version of golangci-lint using irvinlim/golang-ci, try this:
golang-ci lint --version v1.31.0 -- run -v --timeout=5m
```

## License

MIT
