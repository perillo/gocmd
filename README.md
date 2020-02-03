# gocmd [![GoDoc](https://godoc.org/github.com/perillo/gocmd?status.svg)](http://godoc.org/github.com/perillo/gocmd)

The `github.com/perillo/gocmd` module provides the `env`, `pkglist`, `modlist`
and `modfetch` packages that implement a simple *API* for, respectively, the
`go env`, `go list`, `go list -m` and `go mod download` commands.

## env

The `github.com/perillo/gocmd/env` package provides support for accessing the
*Go* environment.  With the default `Config` any changes to the *Go*
environment are recorded in the default `GOENV` file.

`env` is a wrapper for the `go env` command.

## pkglist

The `github.com/perillo/gocmd/pkglist` package provides support for loading
packages.  The `Load` function accepts the same patters accepted by the
`go list` command.

In case of errors, no packages are returned.  The error details are available
in `Error.Stderr`.

`pkglist` is a wrapper for the `go list -json` command.

## modlist

The `github.com/perillo/gocmd/modlist` package provides support for loading
modules.  The `Load` function accepts the same patterns accepted by the
`go list -m` command.

In case of errors, no modules are returned.  The error details are available in
`Error.Stderr`.

`modlist` is a wrapper for the `go list -m -json` command.

## modfetch

The `github.com/perillo/gocmd/modfetch` package provides support for fetching
modules and pre-filling the module cache.  The `Load` function accepts the
same patterns accepted by the `go mod download` command.

In case of errors, no cached modules are returned.  The error details are
available in `Error.Stderr`.

`modfetch` is a wrapper for the `go mod download -json` command,


## API design

For each of the `Load` function, if one or more `packages`/`modules` can not be
loaded, the function returns an error, and a nil slice.  The error will be of
type `Error`, and the error messages are available in the `Stderr` field.

If the `Load` function returns with a `nil` error, it will returns only
correctly loaded `packages`/`modules`.

The *API* is designed to be easy to use and to implement.


## TODO

  - [ ] Add support for the `-compiled`, `-deps`, `-export`, `-find` and
    `-test` options for the `pkglist` package.

  - [ ] Add support for the `-u` and `-versions` options for the `modlist`
    package.
