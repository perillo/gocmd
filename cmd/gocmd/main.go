// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command gocmd is used for debugging the invoke.Go function.
// The output from stdout, stderr and the standard log is redirected to stdout,
// and each line is printed with a prefix indicating the origin.
package main

import (
	"log"
	"os"

	"github.com/perillo/gocmd/internal/debug"
	"github.com/perillo/gocmd/internal/invoke"
)

func main() {
	// Set the GOCMDDEBUG environment variable to debug some corner cases.
	os.Setenv("GOCMDDEBUG", "on")

	log.SetFlags(0)
	log.SetOutput(debug.Stdlog)
	if len(os.Args) == 1 {
		return
	}

	stdout, err := invoke.Go(os.Args[1], os.Args[2:], nil)
	if err != nil {
		debug.Stderr.WriteString(err.Error())
	}
	if stdout != nil {
		debug.Stdout.Write(stdout)
	}
}
