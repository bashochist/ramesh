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

// Package bufwire wires everything together.
//
// TODO: This package should be split up into individual functionality.
package bufwire

import (
	"context"

	"github.com/bufbuild/buf/private/buf/bufconvert"
	"github.com/bufbuild/buf/private/buf/buffetch"
	"github.com/bufbuild/buf/private/bufpkg/bufanalysis"
	"github.com/bufbuild/buf/private/bufpkg/bufconfig"
	"github.com/bufbuild/buf/private/bufpkg/bufimage"
	"github.com/bufbuild/buf/private/bufpkg/bufimage/bufimagebuild"
	"github.com/bufbuild/buf/private/bufpkg/bufmodule"
	"github.com/bufbuild/buf/private/bufpkg/bufmodule/bufmodulebuild"
	"github.com/bufbuild/buf/private/bufpkg/bufmodule/bufmoduleref"
	"github.com/bufbuild/buf/private/pkg/app"
	"github.com/bufbuild/buf/private/pkg/storage/storageos"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// ImageConfig is an image and configuration.
type ImageConfig interface {
	Image() bufimage.Image
	Config() *bufconfig.Config
}

// ImageConfigReader is an ImageConfig reader.
type ImageConfigReader interface {
	// GetImageConfigs gets the ImageConfig for the fetch value.
	//
	// If externalDirOrFilePaths is empty, this builds all files under Buf control.
	GetImageConfigs(
		ctx context.Context,
		container app.EnvStdinContainer,
		ref buffetch.Ref,
		configOverride string,
		externalDirOrFilePaths []string,
		externalExcludeDirOrFilePaths []string,
		externalDirOrFilePathsAllowNotExist bool,
		excludeSourceCodeInfo bool,
	) ([]ImageConfig, []bufanalysis.FileAnnotation, error)
}

// NewImageC