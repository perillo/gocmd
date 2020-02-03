// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package invoke implements a simple api to invoke a cmd/go command and return
// its stdout content and an error.
package invoke

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Attr holds the attributes that will be applied to the cmd/go command.
type Attr struct {
	// Env specifies the environment of the cmd/go command.
	// Each entry is of the form "key=value".
	// If Env is nil, the cmd/go command uses the current process's
	// environment.
	Env []string

	// Dir specifies the working directory of the cmd/go command.
	// If Dir is the empty string, the cmd/go command runs in the calling
	// process's current directory.
	Dir string
}

// Error is returned by Go in case the go command returns an error.
type Error struct {
	Argv   []string // arguments to the go command
	Stderr []byte   // the entire content of the go command stderr
	Err    error    // the original error from exec.Command.Run
}

// Error implements the error interface.
func (e *Error) Error() string {
	argv := strings.Trim(fmt.Sprint(e.Argv), "[]")
	stderr := string(e.Stderr)
	msg := "go " + argv + ": " + e.Err.Error()

	if stderr == "" {
		return msg
	}

	return msg + ": " + stderr
}

// Unwrap implements the error unwrapping interface.
func (e *Error) Unwrap() error {
	return e.Err
}

// Go invokes a cmd/go command and returns its stdout content and an error.  It
// implicitly assumes that the cmd/go command is invoked with the -json flag
// set.
//
// If the cmd/go command returns a non 0 exit status, Go will return the stdout
// content, or nil if empty, and the error will be of type *Error.
//
// If the cmd/go command returns a 0 exit status, Go will return the possibly
// empty stdout content and a nil error.  The stderr content will be ignored,
// unless the GOCMDDEBUG environment variable is not empty, in which case it
// will be logged using the log package.
func Go(verb string, argv []string, attr *Attr) ([]byte, error) {
	argv = append([]string{verb}, argv...)
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	cmd := exec.Command("go", argv...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if attr != nil {
		cmd.Dir = attr.Dir
		cmd.Env = attr.Env
	}

	if err := cmd.Run(); err != nil {
		// Just return the error, including the non empty stderr output as is.
		// Make sure to also return the stdout content, since it may be
		// important.  But only if it is not empty.
		var buf []byte
		if stdout.Len() > 0 {
			buf = stdout.Bytes()
		}
		err := &Error{
			Argv:   argv,
			Stderr: stderr.Bytes(),
			Err:    err,
		}

		return buf, err
	}
	if stderr.Len() != 0 && os.Getenv("GOCMDDEBUG") != "" {
		argv := strings.Trim(fmt.Sprint(argv), "[]")
		log.Printf("go %v: %s", argv, stderr)
	}

	return stdout.Bytes(), nil
}
