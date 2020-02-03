// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package modlist

import (
	"errors"
	"os"
	"strings"
	"testing"
)

// TestLoad tests that the Load function works correctly.
func TestLoad(t *testing.T) {
	l := Loader{
		Dir: os.TempDir(),
	}

	const want = "golang.org/x/text@v0.1.0" // released on 2017-09-14T09:07:07Z
	mods, err := l.Load(want)
	if err != nil {
		t.Fatal(err)
	}
	if len(mods) != 1 {
		t.Errorf("load: expected 1, got %d modules", len(mods))
	}
	got := mods[0].Path + "@" + mods[0].Version
	if got != want {
		t.Errorf("load: got %q, want %q", got, want)
	}
}

// TestLoadFail tests that the Load function in case of failure reports the
// error details in Error.Stderr.
func TestLoadFail(t *testing.T) {
	l := Loader{
		Dir: os.TempDir(),
	}

	mods, err := l.Load("xxx@latest")
	if err == nil {
		t.Error("expected an error")
	}
	if mods != nil {
		t.Errorf("expected the data to be nil, got %v", mods)
	}

	// TODO(mperillo): Ensure that the test is not brittle.
	err = errors.Unwrap(err)
	stderr := string(err.(*Error).Stderr)
	pattern := "xxx@latest: malformed module path"

	if !strings.Contains(stderr, pattern) {
		t.Errorf("stderr does not contain pattern %q, got %q", pattern, stderr)
	}
}
