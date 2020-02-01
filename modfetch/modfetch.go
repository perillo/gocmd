// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package modfetch is a wrapper for the go mod download command.
package modfetch

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/perillo/gocmd/internal/invoke"
)

// Loader is used to provide custom options for fetching modules.
type Loader struct {
	// Dir is the directory in which to run the go mod download command.
	// If Dir is empty, go mod download is run in the current directory.
	Dir string

	// Env is the environment to use when invoking go mod download.
	// If Env is nil, the current environment is used.
	Env []string
}

// Load downloads and returns the Go modules named by the given patterns.
// The patterns are the same as the ones used by go mod download.
func (l *Loader) Load(patterns ...string) ([]*Module, error) {
	attr := invoke.Attr{
		Dir: l.Dir,
		Env: l.Env,
	}
	argv := []string{"download", "-json"}
	argv = append(argv, patterns...)

	stdout, err := invoke.Go("mod", argv, &attr)
	if err != nil {
		return nil, err
	}

	// Decode the modules.
	modlist := make([]*Module, 0, 10)
	buf := bytes.NewBuffer(stdout)
	for dec := json.NewDecoder(buf); dec.More(); {
		tmp := new(moduleJSON)
		if err := dec.Decode(tmp); err != nil {
			return nil, fmt.Errorf("JSON decode: %w", err)
		}

		mod := fromInternal(tmp)
		modlist = append(modlist, mod)
	}

	return modlist, err
}

// Load loads and returns the Go modules named by the given patterns, using
// the default loader configuration.
// The patterns are the same as the ones used by go mod download.
func Load(patterns ...string) ([]*Module, error) {
	var l Loader

	return l.Load(patterns...)
}
