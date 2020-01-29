// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command gocmd is used for debugging the invoke.Go function.
// The output from stdout, stderr and the standard log is redirected to stdout,
// and each line is printed with a prefix indicating the origin.
package main

import (
	"bytes"
	"log"
	"os"
	"strings"

	"github.com/perillo/gocmd/internal/invoke"
)

var (
	tstdout *writer = &writer{prefix: "STDOUT "}
	tstderr *writer = &writer{prefix: "STDERR "}
	tlog    *writer = &writer{prefix: "LOG    "}
)

func main() {
	// Set the GOCMDDEBUG environment variable to debug some corner cases.
	os.Setenv("GOCMDDEBUG", "on")

	log.SetFlags(0)
	log.SetOutput(tlog)
	if len(os.Args) == 1 {
		return
	}

	stdout, err := invoke.Go(os.Args[1], os.Args[2:], nil)
	if err != nil {
		tstderr.WriteString(err.Error())
	}
	if len(stdout) != 0 {
		tstdout.Write(stdout)
	}
}

// writer writes on stdout each line with the specified prefix.
type writer struct {
	prefix string
}

func (w *writer) Write(buf []byte) (int, error) {
	var b bytes.Buffer

	for _, line := range bytes.Split(buf, []byte("\n")) {
		b.WriteString(w.prefix)
		b.Write(line)
		b.WriteByte('\n')
	}

	return os.Stdout.Write(b.Bytes())
}

func (w *writer) WriteString(buf string) (int, error) {
	var b bytes.Buffer

	for _, line := range strings.Split(buf, "\n") {
		b.WriteString(w.prefix)
		b.WriteString(line)
		b.WriteByte('\n')
	}

	return os.Stdout.Write(b.Bytes())
}
