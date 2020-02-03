// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkglist

import (
	"os"
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
