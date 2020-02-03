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

// Error is returned by Load in case the go command returns an error.
type Error = invoke.Error

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
//
// If one or more modules cannot be loaded, Load returns a nil slice and an
// error of type *Error.  If Load returns successfully, the returned modules
// have all been correctly loaded.
func (l *Loader) Load(patterns ...string) ([]*Module, error) {
	attr := invoke.Attr{
		Dir: l.Dir,
		Env: l.Env,
	}
	argv := []string{"download", "-json"}
	argv = append(argv, patterns...)

	stdout, err := invoke.Go("mod", argv, &attr)
	if err != nil {
		return nil, fmt.Errorf("modfetch: load: %w", err)
	}
	modlist, err := decode(stdout)
	if err != nil {
		return nil, fmt.Errorf("modfetch: load: %w", err)
	}

	return modlist, nil
}

// Load loads and returns the Go modules named by the given patterns, using
// the default loader configuration.
// The patterns are the same as the ones used by go mod download.
//
// If one or more modules cannot be loaded, Load returns a nil slice and an
// error of type *Error.  If Load returns successfully, the returned modules
// have all been correctly loaded.
func Load(patterns ...string) ([]*Module, error) {
	var l Loader

	return l.Load(patterns...)
}

func decode(data []byte) ([]*Module, error) {
	modlist := make([]*Module, 0, 10)
	buf := bytes.NewBuffer(data)
	for dec := json.NewDecoder(buf); dec.More(); {
		tmp := new(moduleJSON)
		if err := dec.Decode(tmp); err != nil {
			return nil, fmt.Errorf("JSON decode: %w", err)
		}

		mod := fromInternal(tmp)
		modlist = append(modlist, mod)
	}

	return modlist, nil
}

// fromInternal converts mod from the moduleJSON type used by go mod download,
// to the Module type used by the modfetch package.
func fromInternal(mod *moduleJSON) *Module {
	r := new(Module)
	r.Path = mod.Path
	r.Version = mod.Version
	r.Info = mod.Info
	r.GoMod = mod.GoMod
	r.Zip = mod.Zip
	r.Dir = mod.Zip
	r.Sum = mod.Sum
	r.GoModSum = mod.GoModSum
	if mod.Error != "" {
		r.Error = &ModuleError{Err: mod.Error}
	}

	return r
}
