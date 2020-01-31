// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command modlist is used for debugging the modlist package.
//
// When the -debug flag is set, the output from stdout, stderr and the standard
// log is redirected to stdout, and each line is printed with a prefix
// indicating the origin.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/perillo/gocmd/internal/debug"
	"github.com/perillo/gocmd/modlist"
)

var (
	debugging = flag.Bool("debug", false, "enable debugging")
)

var (
	stdout io.Writer = os.Stdout
	stderr io.Writer = os.Stderr
)

func main() {
	log.SetFlags(0)
	flag.Parse()

	if *debugging {
		// Set the GOCMDDEBUG environment variable to debug some corner cases.
		os.Setenv("GOCMDDEBUG", "on")

		// Initialize the debug environment.
		if err := debug.Init(); err != nil {
			log.Fatal(err)
		}
		stdout = debug.Stdout
		stderr = debug.Stderr
	}

	modlist, err := modlist.Load(flag.Args()...)
	if err != nil {
		fmt.Fprint(stderr, err)
	}

	// Encode packages.
	enc := json.NewEncoder(stdout)
	enc.SetIndent("", "\t")
	for _, mod := range modlist {
		if err := enc.Encode(mod); err != nil {
			log.Fatalf("JSON encode: %v", err)
		}
	}
}
