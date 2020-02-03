// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// entryKey and sortEntries have been adapted from
// src/cmd/go/internal/envcmd/env.go in the Go source distribution.
// Copyright 2012 The Go Authors. All rights reserved.

// Package envtest provides support for testing the env package.
package envtest

import (
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/perillo/gocmd/env"
)

// File represents a GOENV file.
type File struct {
	f      *os.File
	t      *testing.T // to avoid returning errors
	Config env.Config
}

// NewFile returns a new GOENV file at a temporary location.
func NewFile(t *testing.T) *File {
	f, err := ioutil.TempFile("", "cmdgo-*env")
	if err != nil {
		t.Fatal(err)

		// not reached.
	}
	c := env.Config{
		Path: f.Name(),
	}

	return &File{f, t, c}
}

// Read reads the entire content of the GOENV file.  It can be called multiple
// times.
func (f *File) Read() string {
	f.f.Seek(0, io.SeekStart)
	data, err := ioutil.ReadAll(f.f)
	if err != nil {
		f.t.Fatal(err)

		// not reached.
	}

	return strings.TrimSpace(string(data))
}

// Remove removes the GOENV file.
func (f *File) Remove() error {
	return os.Remove(f.f.Name())
}

// Utility functions.

// Key returns the sorted sequence of env keys.
func Keys(env map[string]string) []string {
	buf := make([]string, 0, len(env))
	for key, _ := range env {
		buf = append(buf, key)
	}
	sort.Strings(buf)

	return buf
}

// Encode encodes env in the same format used by GOENV.
func Encode(env map[string]string) string {
	buf := entries(env)

	return strings.Join(buf, "\n")
}

// entries returns the sorted sequence of env entries.
func entries(env map[string]string) []string {
	buf := make([]string, 0, len(env))
	for key, val := range env {
		ent := key + "=" + val
		buf = append(buf, ent)
	}
	sortEntries(buf)

	return buf
}

// entryKey returns the KEY part of the entry KEY=VALUE or else an empty
// string.
func entryKey(entry string) string {
	i := strings.Index(entry, "=")
	if i < 0 {
		return ""
	}

	return entry[:i]
}

// sortEntries sorts a sequence of entries by key.
// It differs from sort.Strings in that GO386= sorts after GO=.
//
// sortEntries uses the same sorting algorithm used by go env.
func sortEntries(entries []string) {
	sort.Slice(entries, func(i, j int) bool {
		return entryKey(entries[i]) < entryKey(entries[j])
	})
}
