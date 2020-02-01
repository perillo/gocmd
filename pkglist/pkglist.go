// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The pkglist package api has been partially influenced by
// golang.org/x/tools/go/packages.

// Package pkglist is a wrapper for the go list command.
package pkglist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/perillo/gocmd/internal/invoke"
)

// Loader is used to provide custom options for loading packages.
type Loader struct {
	// Dir is the directory in which to run the go list command.
	// If Dir is empty, go list is run in the current directory.
	Dir string

	// Env is the environment to use when invoking go list.
	// If Env is nil, the current environment is used.
	Env []string
}

// Load loads and returns the Go packages named by the given patterns.
// The patterns are the same as the ones used by go list.
func (l *Loader) Load(patterns ...string) ([]*Package, error) {
	attr := invoke.Attr{
		Dir: l.Dir,
		Env: l.Env,
	}
	argv := []string{"-json"} // note no -e flag for now
	argv = append(argv, patterns...)

	stdout, err := invoke.Go("list", argv, &attr)
	if err != nil {
		return nil, fmt.Errorf("pkglist: load: %w", err)
	}
	pkglist, err := decode(stdout)
	if err != nil {
		return nil, fmt.Errorf("pkglist: load: %w", err)
	}

	return pkglist, nil
}

// Load loads and returns the Go packages named by the given patterns, using
// the default loader configuration.
// The patterns are the same as the ones used by go list.
func Load(patterns ...string) ([]*Package, error) {
	var l Loader

	return l.Load(patterns...)
}

func decode(data []byte) ([]*Package, error) {
	pkglist := make([]*Package, 0, 10)
	buf := bytes.NewBuffer(data)
	for dec := json.NewDecoder(buf); dec.More(); {
		pkg := new(Package)
		if err := dec.Decode(pkg); err != nil {
			return nil, fmt.Errorf("JSON decode: %w", err)
		}

		pkg = normalize(pkg)
		pkglist = append(pkglist, pkg)
	}

	return pkglist, nil
}

// normalize ensures all the source file paths are absolute, for consistency.
func normalize(pkg *Package) *Package {
	abspaths(pkg.Dir, pkg.GoFiles)
	abspaths(pkg.Dir, pkg.CgoFiles)
	//abspaths(pkg.Dir, pkg.CompiledGoFiles)
	abspaths(pkg.Dir, pkg.IgnoredGoFiles)
	abspaths(pkg.Dir, pkg.CFiles)
	abspaths(pkg.Dir, pkg.CXXFiles)
	abspaths(pkg.Dir, pkg.MFiles)
	abspaths(pkg.Dir, pkg.HFiles)
	abspaths(pkg.Dir, pkg.FFiles)
	abspaths(pkg.Dir, pkg.SFiles)
	abspaths(pkg.Dir, pkg.SwigFiles)
	abspaths(pkg.Dir, pkg.SwigCXXFiles)
	abspaths(pkg.Dir, pkg.SysoFiles)
	abspaths(pkg.Dir, pkg.TestGoFiles)
	abspaths(pkg.Dir, pkg.XTestGoFiles)

	return pkg
}

func abspaths(dir string, names []string) []string {
	for i, name := range names {
		path := filepath.Join(dir, name)
		names[i] = path
	}

	return names
}
