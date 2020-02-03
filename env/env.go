// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package env is a wrapper for the go env command.
package env

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/perillo/gocmd/internal/invoke"
)

// Config is used to provide custom options for accessing the Go environment.
type Config struct {
	// Path is the path in which the Go environment configuration file is
	// stored.
	Path string
}

// Get returns the entire Go environment as a map.
//
// It one or more variable names is given as arguments, Get returns each named
// variable.  If a variable does not exists, the associated value will be the
// empty string.
func (c *Config) Get(vars ...string) (map[string]string, error) {
	argv := []string{"-json"}
	argv = append(argv, vars...)

	stdout, err := c.invokeGo(argv)
	if err != nil {
		return nil, fmt.Errorf("env: read: %w", err)
	}
	env, err := decode(stdout)
	if err != nil {
		return nil, fmt.Errorf("env: read: %w", err)
	}

	return env, nil
}

// Set changes the default settings of the named environment variables
// specified in env.
//
// If one or more variables does not exist, Set returns an error.
func (c *Config) Set(env map[string]string) error {
	argv := []string{"-w"}
	argv = append(argv, flatenv(env)...)

	_, err := c.invokeGo(argv)
	if err != nil {
		return fmt.Errorf("env: write: %w", err)
	}

	return nil
}

// Unset unsets the default settings for the named Go environment variables, if
// one has been set with Set or Setenv.
//
// If one or more variables does not exist, Unset returns an error.
func (c *Config) Unset(vars ...string) error {
	argv := []string{"-u"}
	argv = append(argv, vars...)

	_, err := c.invokeGo(argv)
	if err != nil {
		return fmt.Errorf("env: unset %q: %w", vars, err)
	}

	return nil
}

// Getenv returns the named Go environment variable.
//
// If key does not exist, Getenv returns an empty string.
func (c *Config) Getenv(key string) (string, error) {
	argv := []string{key}

	stdout, err := c.invokeGo(argv)
	if err != nil {
		return "", fmt.Errorf("env: getenv %q: %w", key, err)
	}

	return strings.TrimSpace(string(stdout)), nil
}

// Setenv changes the default setting of the named Go environment variable to
// the given value.
//
// If key does not exist, Setenv returns an error.
func (c *Config) Setenv(key, value string) error {
	argv := []string{"-w", key + "=" + value}

	_, err := c.invokeGo(argv)
	if err != nil {
		return fmt.Errorf("env: setenv \"%s=%s\": %w", key, value, err)
	}

	return nil
}

// Unsetenv unsets the default setting for the named Go environment variable,
// if one has been set with Set or Setenv.
//
// If key does not exist, Unsetenv returns an error.
func (c *Config) Unsetenv(key string) error {
	argv := []string{"-u", key}

	_, err := c.invokeGo(argv)
	if err != nil {
		return fmt.Errorf("env: unsetenv %q: %w", key, err)
	}

	return nil
}

func (c *Config) invokeGo(argv []string) ([]byte, error) {
	if c.Path == "" {
		return invoke.Go("env", argv, nil)
	}
	attr := invoke.Attr{
		Env: append(os.Environ(), "GOENV="+c.Path),
	}

	return invoke.Go("env", argv, &attr)
}

// Get returns the entire Go environment as a map, using the default
// configuration.
//
// It one or more variable names is given as arguments, Get returns each named
// variable.  If a variable does not exists, the associated value will be the
// empty string.
func Get(vars ...string) (map[string]string, error) {
	var c Config

	return c.Get(vars...)
}

// Set changes the default settings of the named environment variables
// specified in env, using the default configuration.
//
// If one or more variables does not exist, Set returns an error.
func Set(env map[string]string) error {
	var c Config

	return c.Set(env)
}

// Unset unsets the default settings for the named Go environment variables, if
// one has been set with Set or Setenv, using the default configuration.
//
// If one or more variables does not exist, Unset returns an error.
func Unset(vars ...string) error {
	var c Config

	return c.Unset(vars...)
}

// Getenv returns the named Go environment variable, using the default
// configuration.
//
// If key does not exist, Getenv returns an empty string.
func Getenv(key string) (string, error) {
	var c Config

	return c.Getenv(key)
}

// Setenv changes the default setting of the named Go environment variable to
// the given value, using the default configuration.
//
// If key does not exist, Setenv returns an error.
func Setenv(key, value string) error {
	var c Config

	return c.Setenv(key, value)
}

// Unsetenv unsets the default setting for the named Go environment variable,
// if one has been set with Set or Setenv, using the default configuration.
//
// If key does not exist, Unsetenv returns an error.
func Unsetenv(key string) error {
	var c Config

	return c.Unsetenv(key)
}

func decode(data []byte) (env map[string]string, err error) {
	if err := json.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("JSON decode: %w", err)
	}

	return env, nil
}

// flatenv flattens env.
// TODO(mperillo): flatenv does not sort the data since it assumes goenv -w
// will keep the environment variables sorted.
func flatenv(env map[string]string) []string {
	buf := make([]string, 0, len(env))
	for k, v := range env {
		ent := k + "=" + v
		buf = append(buf, ent)
	}

	return buf
}
