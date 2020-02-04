// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The Module.String method has been adapted from
// src/cmd/go/internal/modinfo/info.go in the Go source distribution.
// Copyright 2018 The Go Authors. All rights reserved.

package modlist

import (
	"time"
)

// For the actual definition of Module, see
// src/cmd/go/internal/modinfo/info.go.

// Module represents a module.
type Module struct {
	Path      string       `json:",omitempty"` // module path
	Version   string       `json:",omitempty"` // module version
	Versions  []string     `json:",omitempty"` // available module versions (with -versions)
	Replace   *Module      `json:",omitempty"` // replaced by this module
	Time      *time.Time   `json:",omitempty"` // time version was created
	Update    *Module      `json:",omitempty"` // available update, if any (with -u)
	Main      bool         `json:",omitempty"` // is this the main module?
	Indirect  bool         `json:",omitempty"` // is this module only an indirect dependency of main module?
	Dir       string       `json:",omitempty"` // directory holding files for this module, if any
	GoMod     string       `json:",omitempty"` // path to go.mod file for this module, if any
	GoVersion string       `json:",omitempty"` // go version used in module
	Error     *ModuleError `json:",omitempty"` // error loading module
}

// String implements the Stringer interface.
func (m *Module) String() string {
	s := m.Path
	if m.Version != "" {
		s += " " + m.Version
		if m.Update != nil {
			s += " [" + m.Update.Version + "]"
		}
	}
	if m.Replace != nil {
		s += " => " + m.Replace.Path
		if m.Replace.Version != "" {
			s += " " + m.Replace.Version
			if m.Replace.Update != nil {
				s += " [" + m.Replace.Update.Version + "]"
			}
		}
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
