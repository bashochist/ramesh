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
	if len(components) != 4 || components[2] != PluginsPathName {
		return "", "", "", appcmd.NewInvalidArgumentErrorf("%s is not a valid plugin path", pluginPath)
	}
	return components[0], components[1], components[3], nil
}

// ParsePluginVersionPath parses a string in the format <buf.build/owner/plugins/name[:version]>
// into remote, owner, name and version. The version is empty if not specified.
func ParsePluginVersionPath(pluginVersionPath string) (remote string, owner string, name string, version string, _ error) {
	remote, owner, name, err := ParsePluginPath(pluginVersionPath)
	if err != nil {
		return "", "", "", "", err
	}
	components := strings.Split(name, ":")
	switch len(components) {
	case 2:
		return remote, owner, components[0], components[1], nil
	case 1:
		return remote, owner, name, "", nil
	default:
		return "", "", "", "", fmt.Errorf("invalid version: %q", name)
	}
}

// ParseTemplatePath parses a string in the format <buf.build/owner/templates/name>
// into remote, owner and name.
func ParseTemplatePath(templatePath string) (remote string, owner string, name string, _ error) {
	if templatePath == "" {
		return "", "", "", appcmd.NewInvalidArgumentError("you must specify a template path")
	}
	components := strings.Split(templatePath, "/")
	if len(components) != 4 || components[2] != TemplatesPathName {
		return "", "", "", appcmd.NewInvalidArgumentErrorf("%s is not a valid template path", templatePath)
	}
	return components[0], components[1], components[3], nil
}

// ValidateTemplateName validates the format of the template name.
// This is only used for client side validation and attempts to avoid
// validation constraints that we may want to change.
func ValidateTemplateName(templateName string) error {
	if templateName == "" {
		return errors.New("template name is required")
	}
	return nil
}

// TemplateConfig is the config used to describe the plugins
// of a new template.
type TemplateConfig struct {
	Plugins []PluginConfig
}

// TemplateConfigToProtoPluginConfigs converts the template config to a slice of proto plugin configs,
// suitable for use with the Plugin Service CreateTemplate RPC.
func TemplateConfigToProtoPluginConfigs(templateConfig *TemplateConfig) []*registryv1alpha1.PluginConfig {
	pluginConfigs := make([]*registryv1alpha1.PluginConfig, 0, len(templateConfig.Plugins))
	for _, plugin := range templateConfig.Plugins {
		pluginConfigs = append(
			pluginConfigs,
			&registryv1alpha1.PluginConfig{
				PluginOwner: plugin.Owner,
				PluginName:  plugin.Name,
				Parameters:  plugin.Parameters,
			},
		)
	}
	return pluginConfigs
}

// PluginConfig is the config used to describe a plugin in
// a new template.
type PluginConfig struct {
	Owner      string
	Name       string
	Parameters []string
}

// ParseTemplateConfig parses the input template config as a path or JSON/YAML literal.
func ParseTemplateConfig(config string) (*TemplateConfig, error) {
	var data []byte
	var err error
	switch filepath.Ext(config) {
	case ".json", ".yaml", ".yml":
		data, err = os.ReadFile(config)
		if err != nil {
			return nil, fmt.Errorf("could not read file: %v", err)
		}
	default:
		data = []byte(config)
	}
	var version externalTemplateConfigVersion
	if err := encoding.UnmarshalJSONOrYAMLNonStrict(data, &version); err != nil {
		return nil, fmt.Errorf("failed to determine version of template config: %w", err)
	}
	switch version.Version {
	case "":
		return nil