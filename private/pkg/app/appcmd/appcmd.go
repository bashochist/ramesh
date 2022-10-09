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
	// If this is specified, a flag --ver