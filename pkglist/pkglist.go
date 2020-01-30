// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package pkglist is a wrapper for the go list command.
package pkglist

import (
	"bytes"
	"encoding/json"
	"fmt"

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
		return nil, err
	}

	// Decode the packages.
	pkglist := make([]*Package, 0, 10)
	buf := bytes.NewBuffer(stdout)
	for dec := json.NewDecoder(buf); dec.More(); {
		pkg := new(Package)
		if err := dec.Decode(pkg); err != nil {
			return nil, fmt.Errorf("JSON decode: %w", err)
		}

		// TODO(mperillo): Make the source file paths absolute.
		pkglist = append(pkglist, pkg)
	}

	return pkglist, err
}

// Load loads and returns the Go packages named by the given patterns, using
// the default loader configuration.
// The patterns are the same as the ones used by go list.
func Load(patterns ...string) ([]*Package, error) {
	var l Loader

	return l.Load(patterns...)
}
