// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package modlist is a wrapper for the go list -m command.
package modlist

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/perillo/gocmd/internal/invoke"
)

// Loader is used to provide custom options for loading modules.
type Loader struct {
	// Dir is the directory in which to run the go list -m command.
	// If Dir is empty, go list is run in the current directory.
	Dir string

	// Env is the environment to use when invoking go list -m.
	// If Env is nil, the current environment is used.
	Env []string
}

// Load loads and returns the Go modules named by the given patterns.
// The patterns are the same as the ones used by go list -m.
func (l *Loader) Load(patterns ...string) ([]*Module, error) {
	attr := invoke.Attr{
		Dir: l.Dir,
		Env: l.Env,
	}
	argv := []string{"-json", "-m"} // note no -e flag for now
	argv = append(argv, patterns...)

	stdout, err := invoke.Go("list", argv, &attr)
	if err != nil {
		return nil, err
	}

	// Decode the modules.
	modlist := make([]*Module, 0, 10)
	buf := bytes.NewBuffer(stdout)
	for dec := json.NewDecoder(buf); dec.More(); {
		mod := new(Module)
		if err := dec.Decode(mod); err != nil {
			return nil, fmt.Errorf("JSON decode: %w", err)
		}

		modlist = append(modlist, mod)
	}

	return modlist, err
}

// Load loads and returns the Go modules named by the given patterns, using
// the default loader configuration.
// The patterns are the same as the ones used by go list -m.
func Load(patterns ...string) ([]*Module, error) {
	var l Loader

	return l.Load(patterns...)
}
