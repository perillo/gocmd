// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/perillo/gocmd/internal/invoke"
)

const (
	prefix = "STDERR: "
	indent = "    "
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("STDERR: ")
	if len(os.Args) == 1 {
		return
	}

	stdout, err := invoke.Go(nil, os.Args[1], os.Args[2:]...)
	if err != nil {
		log.Print(indented(err, prefix, indent))
	}
	if len(stdout) != 0 {
		fmt.Printf("%s\n", stdout)
	}
}

func indented(err error, prefix, indent string) string {
	var b strings.Builder

	for i, line := range strings.Split(err.Error(), "\n") {
		if line == "" {
			continue
		}
		if i > 0 {
			b.WriteString(prefix)
			b.WriteString(indent)
		}

		b.WriteString(line)
		b.WriteByte('\n')
	}

	return b.String()
}
