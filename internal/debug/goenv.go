// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package debug

import (
	"bytes"
	"path/filepath"

	"github.com/perillo/gocmd/env"
)

func initenv() (func([]byte) []byte, error) {
	env, err := goenv()
	if err != nil {
		// Should never happen.
		return nil, err
	}

	mkrel := func(b []byte) []byte {
		// Make all absolute paths in b relative to the paths in env.
		for _, ent := range env {
			old := []byte(ent.value)
			new := []byte("$" + ent.key)
			b = bytes.ReplaceAll(b, old, new)
		}

		return b
	}

	return mkrel, nil
}

type entry struct {
	key   string
	value string
}

func goenv() ([]entry, error) {
	// We need GOBIN, GOCACHE and GOROOT in addition to GOPATH because, as an
	// example:
	//
	// go list -json -compiled
	// returns path relative to $GOCACHE in CompiledGoFiles.
	//
	// go list -json -export
	// returns path relative to $GOBIN in Target and relative to $GOCACHE in
	// Export.
	//
	// go list -json flag
	// returns a path relative to $GOROOT in Target.
	environ, err := env.Get("GOBIN", "GOCACHE", "GOPATH", "GOROOT")
	if err != nil {
		return nil, err
	}

	return flatenv(environ), nil
}

// flatenv flattens env.  It splits $GOPATH into duplicate entries.
func flatenv(env map[string]string) []entry {
	buf := make([]entry, 0, len(env))
	for k, v := range env {
		// Handle GOPATH, since it may contain multiple entries.
		if k == "GOPATH" {
			for _, path := range filepath.SplitList(v) {
				ent := entry{
					key:   k,
					value: path,
				}
				buf = append(buf, ent)
			}

			continue
		}
		if v == "" {
			// Should not happen.
			continue
		}

		ent := entry{
			key:   k,
			value: v,
		}
		buf = append(buf, ent)
	}

	return buf
}
