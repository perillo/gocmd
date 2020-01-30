// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package debug provides support for debugging the invoke.Go function.
// It provides a custom Stdout, Stderr and Stdlog that will prefix each line
// with the string "STDOUT", "STDERR" and "STDLOG" respectively, in order to
// multiplex the output from stdout, stderr and the standard log package to
// stdout.
// It will also make all the absolute paths read from cmd/go relative to
// $GOPATH.
package debug

import (
	"bytes"
	"log"
	"os"
	"strings"
)

// mkrel is a function that will make all absolute paths in b relative to
// $GOBIN, $GOCACHE and $GOPATH.
var mkrel func(b []byte) []byte

// Stdout, Stderr and Stdlog are io.Writer that allow to multiplex output
// directed to stdout, stderr or logged by the standard log package to stdout.
var (
	Stdout = &Writer{prefix: "STDOUT "}
	Stderr = &Writer{prefix: "STDERR "}
	Stdlog = &Writer{prefix: "STDLOG "}
)

// Writer writes on stdout each line with the specified prefix.
type Writer struct {
	prefix string
}

// Write implements the Read interface.
func (w *Writer) Write(buf []byte) (int, error) {
	var b bytes.Buffer

	for _, line := range bytes.Split(buf, []byte("\n")) {
		b.WriteString(w.prefix)
		b.Write(line)
		b.WriteByte('\n')
	}

	return w.emit(b)
}

func (w *Writer) WriteString(buf string) (int, error) {
	var b bytes.Buffer

	for _, line := range strings.Split(buf, "\n") {
		b.WriteString(w.prefix)
		b.WriteString(line)
		b.WriteByte('\n')
	}

	return w.emit(b)
}

func (w *Writer) emit(b bytes.Buffer) (int, error) {
	// Make any path inside $GOPATH relative to $GOPATH.
	buf := b.Bytes()
	if mkrel != nil {
		buf = mkrel(buf)
	}

	return os.Stdout.Write(buf)
}

// Init initializes the debug environment.
// It also sets the standard log output to Stdlog.
func Init() error {
	// Initialize the environment and set the global mkrel variable.
	f, err := initenv()
	if err != nil {
		return err
	}
	mkrel = f

	// Configure the standard logger to use debug.Stdlog.
	log.SetOutput(Stdlog)

	// Unfortunately we can not change os.Stdout and os.Stderr.
	return nil
}
