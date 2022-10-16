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

// Package appcmd contains helper functionality for applications using commands.
package appcmd

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/bufbuild/buf/private/pkg/app"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/pflag"
)

// Command is a command.
type Command struct {
	// Use is the one-line usage message.
	// Required.
	Use string
	// Aliases are aliases that can be used instead of the first word in Use.
	Aliases []string
	// Short is the short message shown in the 'help' output.
	// Required if Long is set.
	Short string
	// Long is the long message shown in the 'help <this-command>' output.
	// The Short field will be prepended to the Long field with two newlines.
	// Must be unset if short is unset.
	Long string
	// Args are the expected arguments.
	//
	// TODO: make specific types for appcmd to limit what can be done.
	Args cobra.PositionalArgs
	// Deprecated says to print this deprecation string.
	Deprecated string
	// Hidden says to hide this command.
	Hidden bool
	// BindFlags allows binding of flags on build.
	BindFlags func(*pflag.FlagSet)
	// BindPersistentFlags allows binding of flags on build.
	BindPersistentFlags func(*pflag.FlagSet)
	// NormalizeFlag allows for normalization of flag names.
	NormalizeFlag func(*pflag.FlagSet, string) string
	// NormalizePersistentFlag allows for normalization of flag names.
	NormalizePersistentFlag func(*pflag.FlagSet, string) string
	// Run is the command to run.
	// Required if there are no sub-commands.
	// Must be unset if there are sub-commands.
	Run func(context.Context, app.Container) error
	// SubCommands are the sub-commands. Optional.
	// Must be unset if there is a run function.
	SubCommands []*Command
	// Version the version of the command.
	//
	// If this is specified, a flag --version will be added to the command
	// that precedes all other functionality, and which prints the version
	// to stdout.
	Version string
}

// NewInvalidArgumentError creates a new invalidArgumentError, indicating that
// the error was caused by argument validation. This causes us to print the usage
// help text for the command that it is returned from.
func NewInvalidArgumentError(message string) error {
	return newInvalidArgumentError(message)
}

// NewInvalidArgumentErrorf creates a new InvalidArgumentError, indicating that
// the error was caused by argument validation. This causes us to print the usage
// help text for the command that it is returned from.
func NewInvalidArgumentErrorf(format string, args ...interface{}) error {
	return NewInvalidArgumentError(fmt.Sprintf(format, args...))
}

// Main runs the application using the OS container and calling os.Exit on the return value of Run.
func Main(ctx context.Context, command *Command) {
	app.Main(ctx, newRunFunc(command))
}

// Run runs the application using the container.
func Run(ctx context.Context, container app.Container, command *Command) error {
	return app.Run(ctx, container, newRunFunc(command))
}

// BindMultiple is a convenience function for binding multiple flag functions.
func BindMultiple(bindFuncs ...func(*pflag.FlagSet)) func(*pflag.FlagSet) {
	return func(flagSet *pflag.FlagSet) {
		for _, bindFunc := range bindFuncs {
			bindFunc(flagSet)
		}
	}
}

func newRunFunc(command *Command) func(context.Context, app.Container) error {
	return func(ctx context.Context, container app.Container) error {
		return run(ctx, container, command)
	}
}

func run(
	ctx context.Context,
	container app.Container,
	command *Command,
) error {
	var runErr error

	cobraCommand, err := commandToCobra(ctx, container, command, &runErr)
	if err != nil {
		return err
	}

	// Cobra 1.2.0 introduced default completion commands under
	// "<binary> completion <bash/zsh/fish/powershell>"". Since we have
	// our own completion commands, disable the generation of the default
	// commands.
	cobraCommand.CompletionOptions.DisableDefaultCmd = true

	// If the root command is not the only command, add a hidden manpages command
	// and a visible completion command.
	if len(command.SubCommands) > 0 {
		shellCobraCommand, err := commandToCobra(
			ctx,
			container,
			&Command{
				Use:   "completion",
				Short: "Generate auto-completion scripts for commonly used shells",
				SubCommands: []*Command{
					{
						Use:   "bash",
						Short: "Generate auto-completion scripts for bash",
						Args:  cobra.NoArgs,
						Run: func(ctx context.Context, container app.Container) error {
							return cobraCommand.GenBashCompletion(container.Stdout())
						},
					},
					{
						Use:   "fish",
						Short: "Generate auto-completion scripts for fish",
						Args:  cobra.NoArgs,
						Run: func(ctx context.Context, container app.Container) error {
							return cobraCommand.GenFishCompletion(container.Stdout(), true)
						},
					},
					{
						Use:   "powershell",
						Short: "Generate auto-completion scripts for powershell",
						Args:  cobra.NoArgs,
						Run: func(ctx context.Context, container app.Container) error {
							return cobraCommand.GenPowerShellCompletion(cont