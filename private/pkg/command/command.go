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

package command

import (
	"bytes"
	"context"
	"io"

	"github.com/bufbuild/buf/private/pkg/app"
)

// Runner runs external commands.
//
// A Runner will limit the number of concurrent commands, as well as explicitly
// set stdin, stdout, stderr, and env to nil/empty values if not set with options.
//
// All external commands in buf MUST use command.Runner instead of
// exec.Command, exec.CommandContext.
type Runner interface {
	// Run runs the external command.
	//
	// This should be used instead of exec.CommandContext(...).Run().
	Run(ctx context.Context, name string, options ...RunOption) error
}

// RunOption is an option for Run.
type RunOption func(*runOptions)

// RunWithArgs returns a new RunOption that sets the arguments other
// than the name.
//
// The default is no arguments.
func RunWithArgs(args ...string) RunOption {
	return func(runOptions *runOptions) {
		runOptions.args = args
	}
}

// RunWithEnv returns a new RunOption that sets the environment variables.
//
// The default is to use the single environment variable __EMPTY_ENV__=1 as we
// cannot explicitly set an empty environment with the exec package.
func RunWithEnv(env map[string]string) 