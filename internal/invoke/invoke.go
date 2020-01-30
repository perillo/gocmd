// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package invoke implements a simple api to invoke a cmd/go command and return
// its stdout content or an error.
package invoke

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Attr holds the attributes that will be applied to the Go command.
type Attr struct {
	// Env specifies the environment of the Go command.
	// Each entry is of the form "key=value".
	// If Env is nil, the Go command uses the current process's environment.
	Env []string

	// Dir specifies the working directory of the Go command.
	// If Dir is the empty string, the Go command runs in the calling process's
	// current directory.
	Dir string
}

// Go invokes a cmd/go command and returns its stdout content or an error.  It
// implicitly assumes that the cmd/go command is invoked with the -json flag
// set.
//
// If the go command returns a non 0 exit status, Go will return the stdout
// content, or nil if empty, and the error as returned by the exec package,
// with the stderr content as additional context.
//
// If the go command returns a 0 exit status, Go will return the possibly empty
// stdout content and a nil error.  The stderr content will be ignored, unless
// the GOCMDDEBUG environment variable is not empty, in which case it will be
// logged using the log package.
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
		// Just return the error, including the stderr output as is.
		// Make sure to also return the stdout content, since it may be
		// important.  But only if it is not empty.
		argv := strings.Trim(fmt.Sprint(argv), "[]")
		var buf []byte
		if stdout.Len() > 0 {
			buf = stdout.Bytes()
		}

		return buf, fmt.Errorf("go %v: %w: %s", argv, err, stderr)
	}
	if stderr.Len() != 0 && os.Getenv("GOCMDDEBUG") != "" {
		argv := strings.Trim(fmt.Sprint(argv), "[]")
		log.Printf("go %v: %s", argv, stderr)
	}

	return stdout.Bytes(), nil
}
