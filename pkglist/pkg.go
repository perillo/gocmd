// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkglist

import (
	"github.com/perillo/gocmd/modlist"
)

// For the actual definition of Package, see
// src/cmd/go/internal/load/pkg.go.

// Package represents a package.
//
// The Dir, Target, Shlib, Root, ConflictDir, and Export file paths are all
// absolute paths.
//
// The lists GoFiles, CgoFiles, and so on hold absolute paths.
// The generated files added when using the -compiled and -test flags are
// absolute paths referring to cached copies of generated Go source files.
// Although they are Go source files, the paths may not end in ".go".
type Package struct {
	Dir           string   `json:",omitempty"` // directory containing package sources
	ImportPath    string   `json:",omitempty"` // import path of package in dir
	ImportComment string   `json:",omitempty"` // path in import comment on package statement
	Name          string   `json:",omitempty"` // package name
	Doc           string   `json:",omitempty"` // package documentation string
	Target        string   `json:",omitempty"` // install path
	Shlib         string   `json:",omitempty"` // the shared library that contains this package (only set when -linkshared)
	Goroot        bool     `json:",omitempty"` // is this package in the Go root?
	Standard      bool     `json:",omitempty"` // is this package part of the standard Go library?
	Stale         bool     `json:",omitempty"` // would 'go install' do anything for this package?
	StaleReason   string   `json:",omitempty"` // explanation for Stale==true
	Root          string   `json:",omitempty"` // Go root or Go path dir containing this package
	ConflictDir   string   `json:",omitempty"` // this directory shadows Dir in $GOPATH
	BinaryOnly    bool     `json:",omitempty"` // binary-only package (no longer supported)
	ForTest       string   `json:",omitempty"` // package is only for use in named test
	Export        string   `json:",omitempty"` // file containing export data (when using -export)
	Module        *Module  `json:",omitempty"` // info about package's containing module, if any (can be nil)
	Match         []string `json:",omitempty"` // command-line patterns matching this package
	DepOnly       bool     `json:",omitempty"` // package is only a dependency, not explicitly listed

	// Source files
	GoFiles         []string `json:",omitempty"` // .go source files (excluding CgoFiles, TestGoFiles, XTestGoFiles)
	CgoFiles        []string `json:",omitempty"` // .go source files that import "C"
	CompiledGoFiles []string `json:",omitempty"` // .go files presented to compiler (when using -compiled)
	IgnoredGoFiles  []string `json:",omitempty"` // .go source files ignored due to build constraints
	CFiles          []string `json:",omitempty"` // .c source files
	CXXFiles        []string `json:",omitempty"` // .cc, .cxx and .cpp source files
	MFiles          []string `json:",omitempty"` // .m source files
	HFiles          []string `json:",omitempty"` // .h, .hh, .hpp and .hxx source files
	FFiles          []string `json:",omitempty"` // .f, .F, .for and .f90 Fortran source files
	SFiles          []string `json:",omitempty"` // .s source files
	SwigFiles       []string `json:",omitempty"` // .swig files
	SwigCXXFiles    []string `json:",omitempty"` // .swigcxx files
	SysoFiles       []string `json:",omitempty"` // .syso object files to add to archive
	TestGoFiles     []string `json:",omitempty"` // _test.go files in package
	XTestGoFiles    []string `json:",omitempty"` // _test.go files outside package

	// Cgo directives
	CgoCFLAGS    []string `json:",omitempty"` // cgo: flags for C compiler
	CgoCPPFLAGS  []string `json:",omitempty"` // cgo: flags for C preprocessor
	CgoCXXFLAGS  []string `json:",omitempty"` // cgo: flags for C++ compiler
	CgoFFLAGS    []string `json:",omitempty"` // cgo: flags for Fortran compiler
	CgoLDFLAGS   []string `json:",omitempty"` // cgo: flags for linker
	CgoPkgConfig []string `json:",omitempty"` // cgo: pkg-config names

	// Dependency information
	Imports      []string          `json:",omitempty"` // import paths used by this package
	ImportMap    map[string]string `json:",omitempty"` // map from source import to ImportPath (identity entries omitted)
	Deps         []string          `json:",omitempty"` // all (recursively) imported dependencies
	TestImports  []string          `json:",omitempty"` // imports from TestGoFiles
	XTestImports []string          `json:",omitempty"` // imports from XTestGoFiles

	// Error information
	Incomplete bool            `json:",omitempty"` // this package or a dependency has an error
	Error      *PackageError   `json:",omitempty"` // error loading package
	DepsErrors []*PackageError `json:",omitempty"` // errors loading dependencies
}

// PackageError represents a package error.
type PackageError struct {
	ImportStack []string // shortest path from package named on command line to this one
	Pos         string   // position of error (if present, file:line:col)
	Err         string   // the error itself
}

func (pe *PackageError) Error() string {
	return pe.Err
}

// Module represents a module.
type Module = modlist.Module
