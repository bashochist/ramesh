
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

// Package bufwork defines the primitives used to enable workspaces.
//
// If a buf.work.yaml file exists in a parent directory (up to the root of
// the filesystem), the directory containing the file is used as the root of
// one or more modules. With this, modules can import from one another, and a
// variety of commands work on multiple modules rather than one. For example, if
// `buf lint` is run for an input that contains a buf.work.yaml, each of
// the modules contained within the workspace will be linted. Other commands, such
// as `buf build`, will merge workspace modules into one (i.e. a "supermodule")
// so that all of the files contained are consolidated into a single image.
//
// In the following example, the workspace consists of two modules: the module
// defined in the petapis directory can import definitions from the paymentapis
// module without vendoring the definitions under a common root. To be clear,
// `import "acme/payment/v2/payment.proto";` from the acme/pet/v1/pet.proto file
// will suffice as long as the buf.work.yaml file exists.
//
//	// buf.work.yaml
//	version: v1
//	directories:
//	  - paymentapis
//	  - petapis
//
//	$ tree
//	.
//	├── buf.work.yaml
//	├── paymentapis
//	│   ├── acme
//	│   │   └── payment
//	│   │       └── v2
//	│   │           └── payment.proto
//	│   └── buf.yaml
//	└── petapis
//	    ├── acme
//	    │   └── pet
//	    │       └── v1
//	    │           └── pet.proto
//	    └── buf.yaml
//
// Note that inputs MUST NOT overlap with any of the directories defined in the buf.work.yaml
// file. For example, it's not possible to build input "paymentapis/acme" since the image
// would otherwise include the content defined in paymentapis/acme/payment/v2/payment.proto as
// acme/payment/v2/payment.proto and payment/v2/payment.proto.
//
// EVERYTHING IN THIS PACKAGE SHOULD ONLY BE CALLED BY THE CLI AND CANNOT BE USED IN SERVICES.
package bufwork

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/bufbuild/buf/private/bufpkg/bufconfig"
	"github.com/bufbuild/buf/private/bufpkg/bufmodule"
	"github.com/bufbuild/buf/private/bufpkg/bufmodule/bufmodulebuild"
	"github.com/bufbuild/buf/private/pkg/normalpath"
	"github.com/bufbuild/buf/private/pkg/storage"
)

const (
	// ExternalConfigV1FilePath is the default configuration file path for v1.
	ExternalConfigV1FilePath = "buf.work.yaml"
	// V1Version is the version string used to indicate the v1 version of the buf.work.yaml file.
	V1Version = "v1"

	// BackupExternalConfigV1FilePath is another acceptable configuration file path for v1.
	//
	// Originally we thought we were going to use buf.work, and had this around for
	// a while, but then moved to buf.work.yaml. We still need to support buf.work as
	// we released with it, however.
	BackupExternalConfigV1FilePath = "buf.work"
)

var (
	// AllConfigFilePaths are all acceptable config file paths without overrides.
	//
	// These are in the order we should check.
	AllConfigFilePaths = []string{
		ExternalConfigV1FilePath,
		BackupExternalConfigV1FilePath,
	}
)

// WorkspaceBuilder builds workspaces. A single WorkspaceBuilder should NOT be persisted
// acorss calls because the WorkspaceBuilder caches the modules used in each workspace.
type WorkspaceBuilder interface {
	// BuildWorkspace builds a bufmodule.Workspace for the given targetSubDirPath.
	BuildWorkspace(
		ctx context.Context,
		workspaceConfig *Config,
		readBucket storage.ReadBucket,
		relativeRootPath string,
		targetSubDirPath string,
		configOverride string,
		externalDirOrFilePaths []string,
		externalExcludeDirOrFilePaths []string,
		externalDirOrFilePathsAllowNotExist bool,
	) (bufmodule.Workspace, error)

	// GetModuleConfig returns the bufmodule.Module and *bufconfig.Config, associated with the given
	// targetSubDirPath, if it exists.
	GetModuleConfig(targetSubDirPath string) (bufmodule.Module, *bufconfig.Config, bool)
}

// NewWorkspaceBuilder returns a new WorkspaceBuilder.
func NewWorkspaceBuilder(
	moduleBucketBuilder bufmodulebuild.ModuleBucketBuilder,
) WorkspaceBuilder {
	return newWorkspaceBuilder(moduleBucketBuilder)
}

// BuildOptionsForWorkspaceDirectory returns the bufmodulebuild.BuildOptions required for
// the given subDirPath based on the workspace configuration.
//
// The subDirRelPaths are the relative paths of the externalDirOrFilePaths that map to the
// provided subDirPath.
// The subDirRelExcludePaths are the relative paths of the externalExcludeDirOrFilePaths that map to the
// provided subDirPath.
func BuildOptionsForWorkspaceDirectory(
	ctx context.Context,
	workspaceConfig *Config,
	moduleConfig *bufconfig.Config,
	externalDirOrFilePaths []string,
	externalExcludeDirOrFilePaths []string,
	subDirRelPaths []string,
	subDirRelExcludePaths []string,
	externalDirOrFilePathsAllowNotExist bool,
) ([]bufmodulebuild.BuildOption, error) {
	buildOptions := []bufmodulebuild.BuildOption{
		// We can't determine the module's commit from the local file system.
		// This also may be nil.
		//
		// This is particularly useful for the GoPackage modifier used in
		// managed mode, which supports module-specific overrides.
		bufmodulebuild.WithModuleIdentity(moduleConfig.ModuleIdentity),
	}
	if len(externalDirOrFilePaths) == 0 && len(externalExcludeDirOrFilePaths) == 0 {
		return buildOptions, nil