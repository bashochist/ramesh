// Copyright 2020-2023 Buf Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package app provides application primitives.
package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/bufbuild/buf/private/pkg/interrupt"
)

// EnvContainer provides environment variables.
type EnvContainer interface {
	// Env gets the environment variable value for the key.
	//
	// Returns empty string if the key is not set or the value is empty.
	Env(key string) string
	// ForEachEnv iterates over all non-empty environment variables and calls the function.
	//
	// The value will never be empty.
	ForEachEnv(func(string, string))
}

// NewEnvContainer returns a new EnvContainer.
//
// Empty values are effectively ignored.
func NewEnvContainer(m map[string]string) EnvContainer {
	return newEnvContainer(m)
}

// NewEnvContainerForOS returns a new EnvContainer for the operating system.
func NewEnvContainerForOS() (EnvContainer, error) {
	return newEnvContainerForEnviron(os.Environ())
}

// NewEnvContainerWithOverrides returns a new EnvContainer with the values of the input
// EnvContainer, overridden by the values in overrides.
//
// Empty values are effectively ignored. To unset a key, set the value to "" in overrides.
func NewEnvContainerWithOverrides(envContainer EnvContainer, overrides map[string]string) EnvContainer {
	m := EnvironMap(envContainer)
	for key, value := range overrides {
		m[key] = value
	}
	return newEnvContainer(m)
}

// StdinContainer provides stdin.
type StdinContainer interface {
	// Stdin provides stdin.
	//
	// If no value was passed when Stdio was created, this will return io.EOF on any call.
	Stdin() io.Reader
}

// NewStdinContainer returns a new StdinContainer.
func NewStdinContainer(reader io.Reader) StdinContainer {
	return newStdinContainer(reader)
}

// NewStdinContainerForOS returns a new StdinContainer for the operating system.
func NewStdinContainerForOS() StdinContainer {
	return newStdinContainer(os.Stdin)
}

// StdoutContainer provides stdout.
type StdoutContainer interface {
	// Stdout provides stdout.
	//
	// If no value was passed when Stdio was created, this will return io.EOF on any call.
	Stdout() io.Writer
}

// NewStdoutContainer returns a new StdoutC