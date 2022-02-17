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
	"errors"
	"fmt"

	"github.com/bufbuild/buf/private/bufpkg/bufmodule/bufmoduleref"
	"github.com/bufbuild/buf/private/pkg/normalpath"
	"github.com/bufbuild/buf/private/pkg/storage"
	"github.com/bufbuild/buf/private/pkg/stringutil"
)

type targetingModule struct {
	Module
	targetPaths              []string
	pathsAllowNotExistOnWalk bool
	excludePaths             []string
}

func newTargetingModule(
	delegate Module,
	targetPaths []string,
	excludePaths []string,
	pathsAllowNotExistOnWalk bool,
) (*targetingModule, error) {
	if err := normalpath.ValidatePathsNormalizedValidatedUnique(targetPaths); err != nil {
		return nil, err
	}
	return &targetingModule{
		Module:                   delegate,
		targetPaths:              targetPaths,
		pathsAllowNotExistOnWalk: pathsAllowNotExistOnWalk,
		excludePaths:             excludePaths,
	}, nil
}

func (m *targetingModule) TargetFileInfos(ctx context.Context) (fileInfos []bufmoduleref.FileInfo, retErr error) {
	defer func() {
		if retErr == nil {
			bufmoduleref.SortFileInfos(fileInfos)
		}
	}()
	excludePathMap := stringutil.SliceToMap(m.excludePaths)
	// We start by ensuring that no paths have been duplicated between target and exclude pathes.
	for _, targetPath := range m.targetPaths {
		if _, ok := excludePathMap[targetPath]; ok {
			return nil, fmt.Errorf(
				"cannot set the same path for both --path and --exclude-path flags: %s",
				normalpath.Unnormalize(targetPath),
			)
		}
	}
	sourceReadBucket := m.getSourceReadBucket()
	// potentialDirPaths are paths that we need to check if they are directories.
	// These are any files that do not end in .proto, as well as files that end in .proto, but
	// do not have a corresponding file in the source ReadBucket.
	// If there is not an f