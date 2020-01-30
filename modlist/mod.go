// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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

// ModuleError represents a module error.
type ModuleError struct {
	Err string // the error itself
}

func (me *ModuleError) Error() string {
	return me.Err
}
