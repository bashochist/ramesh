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

package bufpluginconfig

import (
	"context"
	"math"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/bufbuild/buf/private/bufpkg/bufplugin/bufpluginref"
	"github.com/bufbuild/buf/private/pkg/storage/storageos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestGetConfigForBucket(t *testing.T) {
	t.Parallel()
	storageosProvider := storageos.NewProvider()
	readWriteBucket, err := storageosProvider.NewReadWriteBucket(filepath.Join("testdata", "success", "go"))
	require.NoError(t, err)
	pluginConfig, err := GetConfigForBucket(context.Background(), readWriteBucket)
	require.NoError(t, err)
	pluginIdentity, err := bufpluginref.PluginIdentityForString("buf.build/library/go-grpc")
	require.NoError(t, err)
	pluginDependency, err := bufpluginref.PluginReferenceForString("buf.build/library/go:v1.28.0", 1)
	require.NoError(t, err)
	require.Equal(
		t,
		&Config{
			Name:          pluginIdentity,
			PluginVersion: "v1.2.0",
			SourceURL:     "https://github.com/grpc/grpc-go",
			Description:   "Generates Go language bindings of services in protobuf definition files for gRPC.",
			Dependencies: []bufpluginref.PluginReference{
				pluginDependency,
			},
			OutputLanguages: []string{"go"},
			Registry: &RegistryConfig{
				Go: &GoRegistryConfig{
					MinVersion: "1.18",
					Deps: []*GoRegistryDependencyConfig{
						{
							Module:  "google.golang.org/grpc",
							Version: "v1.32.0",
						},
					},
				},
				Options: map[string]string{
					"separate_package": "true",
				},
			},
			SPDXLicenseID: "Apache-2.0",
			LicenseURL:    "https://github.com/grpc/grpc-go/blob/master/LICENSE",
		},
		pluginConfig,
	)
}

func TestParsePluginConfigGoYAML(t *testing.T) {
	t.Parallel()
	pluginConfig, err := ParseConfig(filepath.Join("testdata", "success", "go", "buf.plugin.yaml"))
	require.NoError(t, err)
	pluginIdentity, err := bufpluginref.PluginIdentityForString("buf.build/library/go-grpc")
	require.NoError(t, err)
	pluginDependency, err := bufpluginref.PluginReferenceForString("buf.build/library/go:v1.28.0", 1)
	require.NoError(t, err)
	require.Equal(
		t,
		&Config{
			Name:          pluginIdentity,
			PluginVersion: "v1.2.0",
			SourceURL:     "https://github.com/grpc/grpc-go",
			Description:   "Generates Go language bindings of services in protobuf definition files for gRPC.",
			Dependencies: []bufpluginref.PluginReference{
				pluginDependency,
			},
			OutputLanguages: []string{"go"},
			Registry: &RegistryConfig{
				Go: &GoRegistryConfig{
					MinVersion: "1.18",
					Deps: []*GoRegistryDependencyConfig{
						{
							Module:  "google.golang.org/grpc",
							Version: "v1.32.0",
						},
					},
				},
				Options: map[string]string{
					"separate_package": "true",
				},
			},
			SPDXLicenseID: "Apache-2.0",
			LicenseURL:    "https://github.com/grpc/grpc-go/blob/master/LICENSE",
		},
		pluginConfig,
	)
}

func TestParsePluginConfigGoYAMLOverrideRemote(t *testing.T) {
	t.Parallel()
	pluginConfig, err := ParseConfig(filepath.Join("testdata", "success", "go", "buf.plugin.yaml"), WithOverrideRemote("buf.mydomain.com"))
	require.NoError(t, err)
	pluginIdentity, err := bufpluginref.PluginIdentityForString("buf.mydomain.com/library/go-grpc")
	require.NoError(t, err)
	pluginDependency, err := bufpluginref.PluginReferenceForString("buf.mydomain.com/library/go:v1.28.0", 1)
	require.NoError(t, err)
	assert.Equal(t, pluginIdentity, pluginConfig.Name)
	require.Len(t, pluginConfig.Dependencies, 1)
	assert.Equal(t, pluginDependency, pluginConfig.Dependencies[0])
}

func TestParsePluginConfigNPMYAML(t *testing.T) {
	t.Parallel()
	pluginConfig, err := ParseConfig(filepath.Join("testdata", "success", "npm", "buf.plugin.yaml"))
	require.NoError(t, err)
	pluginIdentity, err := bufpluginref.PluginIdentityForString("buf.build/protocolbuffers/js")
	require.NoError(t, err)
	require.Equal(
		t,
		&Config{
			Name:            pluginIdentity,
			PluginVer