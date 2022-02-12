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

package bufmodule

import (
	"context"
	"fmt"

	"github.com/bufbuild/buf/private/bufpkg/bufcheck/bufbreaking/bufbreakingconfig"
	"github.com/bufbuild/buf/private/bufpkg/bufcheck/buflint/buflintconfig"
	"github.com/bufbuild/buf/private/bufpkg/bufconfig"
	"github.com/bufbuild/buf/private/bufpkg/bufmodule/bufmoduleref"
	breakingv1 "github.com/bufbuild/buf/private/gen/proto/go/buf/alpha/breaking/v1"
	lintv1 "github.com/bufbuild/buf/private/gen/proto/go/buf/alpha/lint/v1"
	modulev1alpha1 "github.com/bufbuild/buf/private/gen/proto/go/buf/alpha/module/v1alpha1"
	"github.com/bufbuild/buf/private/pkg/manifest"
	"github.com/bufbuild/buf/private/pkg/normalpath"
	"github.com/bufbuild/buf/private/pkg/storage"
	"github.com/bufbuild/buf/private/pkg/storage/storagemem"
)

type module struct {
	sourceReadBucket     storage.ReadBucket
	dependencyModulePins []bufmoduleref.ModulePin
	moduleIdentity       bufmoduleref.ModuleIdentity
	commit               string
	documentation        string
	license              string
	breakingConfig       *bufbreakingconfig.Config
	lintConfig           *buflintconfig.Config
	manifest             *manifest.Manifest
	blobSet              *manifest.BlobSet
}

func newModuleForProto(
	ctx context.Context,
	protoModule *modulev1alpha1.Module,
	options ...ModuleOption,
) (*module, error) {
	if err := ValidateProtoModule(protoModule); err != nil {
		return nil, err
	}
	// We store this as a ReadBucket as this should never be modified outside of this function.
	readWriteBucket := storagemem.NewReadWriteBucket()
	for _, moduleFile := range protoModule.Files {
		if normalpath.Ext(moduleFile.Path) != ".proto" {
			return nil, fmt.Errorf("expected .proto file but got %q", moduleFile)
		}
		// we already know that paths are unique from validation
		if err := storage.PutPath(ctx, readWriteBucket, moduleFile.Path, moduleFile.Content); err != nil {
			return nil, err
		}
	}
	dependencyModulePins, err := bufmoduleref.NewModulePinsForProtos(protoModule.Dependencies...)
	if err != nil {
		return nil, err
	}
	breakingConfig, lintConfig, err := configsForProto(protoModule.GetBreakingConfig(), protoModule.GetLintConfig())
	if err != nil {
		return nil, err
	}
	return newModule(
		ctx,
		readWriteBucket,
		dependencyModulePins,
		nil, // The module identity is not stored on the proto. We rely on the layer above, (e.g. `ModuleReader`) to set this as needed.
		protoModule.GetDocumentation(),
		protoModule.GetLicense(),
		breakingConfig,
		lintConfig,
		options...,
	)
}

func configsForProto(
	protoBreakingConfig *breakingv1.Config,
	protoLintConfig *lintv1.Config,
) (*bufbreakingconfig.Config, *buflintconfig.Config, error) {
	var breakingConfig *bufbreakingconfig.Config
	var breakingConfigVersion string
	if protoBreakingConfig != nil {
		breakingConfig = bufbreakingconfig.ConfigForProto(protoBreakingConfig)
		breakingConfigVersion = breakingConfig.Version
	}
	var lintConfig *buflintconfig.Config
	var lintConfigVersion string
	if protoLintConfig != nil {
		lintConfig = buflintconfig.ConfigForProto(protoLintConfig)
		lintConfigVersi