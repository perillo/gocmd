// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env_test // in order to avoid import cycle

import (
	"reflect"
	"testing"

	"github.com/perillo/gocmd/env"
	"github.com/perillo/gocmd/internal/envtest"
)

// TestSet tests the Set, Get and Unset functions.
func TestSet(t *testing.T) {
	goenv := envtest.NewFile(t)
	defer goenv.Remove()

	// 1. Set some environment variables.  We assume goenv keeps the key
	// sorted.
	env := map[string]string{
		"AR": "xx",
		"CC": "yy",
	}
	if err := goenv.Config.Set(env); err != nil {
		t.Fatal(err)
	}

	// 2. Check that GOENV file contains the correct environ.
	{
		want := envtest.Encode(env)
		got := goenv.Read()
		if want != got {
			t.Errorf("GOENV: got %q, want %q", got, want)
		}
	}

	// 3. Check that Get returns the same environ.
	{
		want := env
		got, err := goenv.Config.Get("AR", "CC")
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("read: got %q, want %q", got, want)
		}
	}

	// 4. Unset.
	{
		if err := goenv.Config.Unset("AR", "CC"); err != nil {
			t.Fatal(err)
		}
	}

	// 5. Check that GOENV file is empty.
	{
		const want = ""
		got := goenv.Read()
		if got != want {
			t.Errorf("GOENV: got %q, want %q", got, want)
		}
	}

}

// TestSetenv tests the Setenv, Getenv and Unsetenv functions.
func TestSetenv(t *testing.T) {
	goenv := envtest.NewFile(t)
	defer goenv.Remove()

	// 1. Setenv.
	const key, value = "AR", "xx"
	if err := goenv.Config.Setenv(key, value); err != nil {
		t.Fatal(err)
	}

	// 2. Check that GOENV file contains the correct environ.
	{
		want := key + "=" + value
		got := goenv.Read()
		if want != got {
			t.Errorf("GOENV: got %q, want %q", got, want)
		}
	}

	// 3. Check that Getenv returns the same environ.
	{
		want := value
		got, err := goenv.Config.Getenv(key)
		if err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Errorf("getenv: got %q, want %q", got, want)
		}
	}

	// 4. Unsetenv.
	{
		if err := goenv.Config.Unsetenv(key); err != nil {
			t.Fatal(err)
		}
	}

	// 5. Check that GOENV file is empty.
	{
		const want = ""
		got := goenv.Read()
		if got != want {
			t.Errorf("GOENV: got %q, want %q", got, want)
		}
	}

	// 6. Check that Getenv returns the default value.
	{
		// The default value of AR is 'ar', as documented in
		// go help environment.
		const want = "ar"
		got, err := goenv.Config.Getenv(key)
		if err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Errorf("getenv: got %q, want %q", got, want)
		}
	}
}

// TestGetUnknown tests the Get function with an unknown variable.
func TestGetUnknown(t *testing.T) {
	goenv := envtest.NewFile(t)
	defer goenv.Remove()

	got, err := goenv.Config.Get("XX")
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{"XX": ""}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("get: got %q, want %w", got, want)
	}
}

// TestSetUnknown tests the Set function with an unknown variable.
func TestSetUnknown(t *testing.T) {
	goenv := envtest.NewFile(t)
	defer goenv.Remove()

	env := map[string]string{"XX": "yy"}
	if err := goenv.Config.Set(env); err == nil {
		t.Errorf("set: expected error")
	}
}

// TestUnsetUnknown tests the Unset function with an unknown variable.
func TestUnsetUnknown(t *testing.T) {
	goenv := envtest.NewFile(t)
	defer goenv.Remove()

	if err := goenv.Config.Unset("XX"); err == nil {
		t.Errorf("unset: expected error")
	}
}

// TestGetenvUnknown tests the Getenv function with an unknown variable.
func TestGetenvUnknown(t *testing.T) {
	goenv := envtest.NewFile(t)
	defer goenv.Remove()

	value, err := goenv.Config.Getenv("XX")
	if err != nil {
		t.Fatal(err)
	}
	if value != "" {
		t.Errorf("getenv: expected empty, got %q", value)
	}
}

// TestSetenvUnknown tests the Setenv function with an unknown variable.
func TestSetenvUnknown(t *testing.T) {
	goenv := envtest.NewFile(t)
	defer goenv.Remove()

	if err := goenv.Config.Setenv("XX", "yy"); err == nil {
		t.Errorf("setenv: expected error")
	}
}

// TestUnsetenvUnknown tests the Unsetenv function with an unknown variable.
func TestUnsetenvUnknown(t *testing.T) {
	goenv := envtest.NewFile(t)
	defer goenv.Remove()

	if err := goenv.Config.Unsetenv("XX"); err == nil {
		t.Errorf("unsetenv: expected error")
	}
}

// TestNotExisting tests the Set, Unset, Setenv and Unsetenv functions when the
// GOENV file does not exists.
func TestNotExisting(t *testing.T) {
	config := env.Config{
		Path: "/abc/xyz/noenv",
	}

	// We only need to test the Setenv function.  Make sure to use a know
	// variable, to avoid triggering a different error.
	if err := config.Setenv("AR", "ar+"); err == nil {
		t.Errorf("setenv: expected error")
	}
}

// TestInvalidVariable tests the Get and Getenv functions with an invalid
// variable name.
func TestInvalidVariable(t *testing.T) {
	goenv := envtest.NewFile(t)
	defer goenv.Remove()

	// To make go env fail, use a variable name starting with "-".
	const key = "-xxx"

	if _, err := goenv.Config.Get(key); err == nil {
		t.Errorf("get: expected error")
	}
	if _, err := goenv.Config.Getenv(key); err == nil {
		t.Errorf("getenv: expected error")
	}
}
