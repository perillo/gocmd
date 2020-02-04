// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package modfetch

// For the actual definition of moduleJSON, see
// src/cmd/go/internal/modcmd/download.go.

// moduleJSON is the internal representation used by go mod download.
// It has the problem that the Error field is defined as a string, and not as a
// ModuleError, thus preventing the error to be easily wrapped.
type moduleJSON struct {
	Path     string `json:",omitempty"` // module path
	Version  string `json:",omitempty"` // module version
	Error    string `json:",omitempty"` // error loading module
	Info     string `json:",omitempty"` // absolute path to cached .info file
	GoMod    string `json:",omitempty"` // absolute path to cached .mod file
	Zip      string `json:",omitempty"` // absolute path to cached .zip file
	Dir      string `json:",omitempty"` // absolute path to cached source root directory
	Sum      string `json:",omitempty"` // checksum for path, version (as in go.sum)
	GoModSum string `json:",omitempty"` // checksum for go.mod (as in go.sum)
}

// Module represents a cached module.
type Module struct {
	Path     string       `json:",omitempty"` // module path
	Version  string       `json:",omitempty"` // module version
	Info     string       `json:",omitempty"` // absolute path to cached .info file
	GoMod    string       `json:",omitempty"` // absolute path to cached .mod file
	Zip      string       `json:",omitempty"` // absolute path to cached .zip file
	Dir      string       `json:",omitempty"` // absolute path to cached source root directory
	Sum      string       `json:",omitempty"` // checksum for path, version (as in go.sum)
	GoModSum string       `json:",omitempty"` // checksum for go.mod (as in go.sum)
	Error    *ModuleError `json:",omitempty"` // error loading module
}

// String implements the Stringer interface.
func (m *Module) String() string {
	s := m.Path
	if m.Version != "" {
		s += "@" + m.Version
	}

	return s
}

// ModuleError represents a module error.
type ModuleError struct {
	Err string // the error itself
}

// Error implements the error interface.
func (me *ModuleError) Error() string {
	return me.Err
}
