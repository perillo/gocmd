// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command gocmd is used for debugging the invoke.Go function.
// The output from stdout, stderr and the standard log is redirected to stdout,
// and each line is printed with a prefix indicating the origin.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/perillo/gocmd/internal/debug"
	"github.com/perillo/gocmd/internal/invoke"
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

	// check command line arguments.
	if flag.NArg() == 0 {
		return
	}

	data, err := invoke.Go(flag.Arg(0), flag.Args()[1:], nil)
	if err != nil {
		fmt.Fprint(stderr, err)
	}
	if stdout != nil {
		stdout.Write(data)
	}
}
