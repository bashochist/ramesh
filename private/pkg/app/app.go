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

// NewStdoutContainer returns a new StdoutContainer.
func NewStdoutContainer(writer io.Writer) StdoutContainer {
	return newStdoutContainer(writer)
}

// NewStdoutContainerForOS returns a new StdoutContainer for the operatoutg system.
func NewStdoutContainerForOS() StdoutContainer {
	return newStdoutContainer(os.Stdout)
}

// StderrContainer provides stderr.
type StderrContainer interface {
	// Stderr provides stderr.
	//
	// If no value was passed when Stdio was created, this will return io.EOF on any call.
	Stderr() io.Writer
}

// NewStderrContainer returns a new StderrContainer.
func NewStderrContainer(writer io.Writer) StderrContainer {
	return newStderrContainer(writer)
}

// NewStderrContainerForOS returns a new StderrContainer for the operaterrg system.
func NewStderrContainerForOS() StderrContainer {
	return newStderrContainer(os.Stderr)
}

// ArgContainer provides the arguments.
type ArgContainer interface {
	// NumArgs gets the number of arguments.
	NumArgs() int
	// Arg gets the ith argument.
	//
	// Panics if i < 0 || i >= Len().
	Arg(i int) string
}

// NewArgContainer returns a new ArgContainer.
func NewArgContainer(args ...string) ArgContainer {
	return newArgContainer(args)
}

// NewArgContainerForOS returns a new ArgContainer for the operating system.
func NewArgContainerForOS() ArgContainer {
	return newArgContainer(os.Args)
}

// Container contains environment variables, args, and stdio.
type Container interface {
	EnvContainer
	StdinContainer
	StdoutContainer
	StderrContainer
	ArgContainer
}

// NewContainer returns a new Container.
func NewContainer(
	env map[string]string,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
	args ...string,
) Container {
	return newContainer(
		NewEnvContainer(env),
		NewStdinContainer(stdin),
		NewStdoutContainer(stdout),
		NewStderrContainer(stderr),
		NewArgContainer(args...),
	)
}

// NewContainerForOS returns a new Container for the operating system.
func NewContainerForOS() (Container, error) {
	envContainer, err := NewEnvContainerForOS()
	if err != nil {
		return nil, err
	}
	return newContainer(
		envContainer,
		NewStdinContainerForOS(),
		NewStdoutContainerForOS(),
		NewStderrContainerForOS(),
		NewArgContainerForOS(),
	), nil
}

// NewContainerForArgs returns a new Container with the replacement args.
func NewContainerForArgs(container Container, newArgs ...string) Container {
	return newContainer(
		container,
		container,
		container,
		container,
		NewArgContainer(newArgs...),
	)
}

// StdioContainer is a stdio container.
type StdioContainer interface {
	StdinContainer
	StdoutContainer
	StderrContainer
}

// EnvStdinContainer is an environment and stdin container.
type EnvStdinContainer interface {
	EnvContainer
	StdinContainer
}

// EnvStdoutContainer is an environment and stdout container.
type EnvStdoutContainer interface {
	EnvContainer
	StdoutContainer
}

// EnvStderrContainer is an environment and stderr container.
type EnvStderrContainer interface {
	EnvContainer
	StderrContainer
}

// EnvStdioContainer is an environment and stdio container.
type EnvStdioContainer interface {
	EnvContainer
	StdioContainer
}

// Environ returns all environment variables in the form "KEY=VALUE".
//
// Equivalent to os.Environ.
//
// Sorted.
func Environ(envContainer EnvContainer) []string 