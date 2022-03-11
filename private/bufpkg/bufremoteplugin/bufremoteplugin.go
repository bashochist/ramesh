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

package bufremoteplugin

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	registryv1alpha1 "github.com/bufbuild/buf/private/gen/proto/go/buf/alpha/registry/v1alpha1"
	"github.com/bufbuild/buf/private/pkg/app/appcmd"
	"github.com/bufbuild/buf/private/pkg/encoding"
)

const (
	// PluginsPathName is the path prefix used to signify that
	// a name belongs to a plugin.
	PluginsPathName = "plugins"

	// TemplatesPathName is the path prefix used to signify that
	// a name belongs to a template.
	TemplatesPathName = "templates"

	v1Version = "v1"
)

// ParsePluginPath parses a string in the format <buf.build/owner/plugins/name>
// into remote, owner and name.
func ParsePluginPath(pluginPath string) (remote string, owner string, name string, _ error) {
	if pluginPath == "" {
		return "", "", "", appcmd.NewInvalidArgumentError("you must specify a plugin path")
	}
	components := strings.Split(pluginPath, "/")
	if len(components) != 4 || components[2] != PluginsPath