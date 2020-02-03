// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package modfetch

import (
	"os"
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
