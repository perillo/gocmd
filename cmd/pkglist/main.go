// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command pkglist is used for debugging the pkglist package.
// The output from stdout, stderr and the standard log is redirected to stdout,
// and each line is printed with a prefix indicating the origin.
package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/perillo/gocmd/internal/debug"
	"github.com/perillo/gocmd/pkglist"
)

func main() {
	log.SetFlags(0)

	// Set the GOCMDDEBUG environment variable to debug some corner cases.
	os.Setenv("GOCMDDEBUG", "on")

	// Initialize the debug environment.
	if err := debug.Init(); err != nil {
		log.Fatal(err)
	}

	pkglist, err := pkglist.Load(os.Args[1:]...)
	if err != nil {
		debug.Stderr.WriteString(err.Error())
	}

	// Encode packages.
	enc := json.NewEncoder(debug.Stdout)
	enc.SetIndent("", "\t")
	for _, pkg := range pkglist {
		if err := enc.Encode(pkg); err != nil {
			log.Fatalf("JSON encode: %v", err)
		}
	}
}
