// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkglist

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

	const want = "flag"
	pkgs, err := l.Load(want)
	if err != nil {
		t.Fatal(err)
	}
	if len(pkgs) != 1 {
		t.Errorf("load: expected 1, got %d packages", len(pkgs))
	}
	got := pkgs[0].Name
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

	pkgs, err := l.Load("xxx")
	if err == nil {
		t.Error("expected an error")
	}
	if pkgs != nil {
		t.Errorf("expected the data to be nil, got %v", pkgs)
	}

	// TODO(mperillo): Ensure that the test is not brittle.
	err = errors.Unwrap(err)
	stderr := string(err.(*Error).Stderr)
	pattern := "package xxx: malformed module path"

	if !strings.Contains(stderr, pattern) {
		t.Errorf("stderr does not contain pattern %q, got %q", pattern, stderr)
	}
}
