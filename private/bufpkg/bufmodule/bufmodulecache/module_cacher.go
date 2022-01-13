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

package bufmodulecache

import (
	"context"

	"github.com/bufbuild/buf/private/bufpkg/buflock"
	"github.com/bufbuild/buf/private/bufpkg/bufmodule"
	"github.com/bufbuild/buf/private/bufpkg/bufmodule/bufmoduleref"
	"github.com/bufbuild/buf/private/pkg/storage"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type moduleCache interface {
	bufmodule.ModuleReader
	// PutModule stores the module in the cache.
	PutModule(
		ctx context.Context,
		modulePin bufmoduleref.ModulePin,
		module bufmodule.Module,
	) error
}

type moduleCacher struct {
	logger                  *zap.Logger
	dataReadWriteBucket     storage.ReadWriteBucket
	sumReadWriteBucket      storage.ReadWriteBucket
	allowCacheExternalPaths bool
}

var _ moduleCache = (*moduleCacher)(nil)

func newModuleCacher(
	logger *zap.Logger,
	dataReadWriteBucket storage.ReadWriteBucket,
	sumReadWriteBucket storage.ReadWriteBucket,
	allowCacheExternalPaths bool,
) *moduleCacher {
	return &moduleCacher{
		logger:                  logger,
		dataReadWriteBucket:     dataReadWriteBucket,
		sumReadWriteBucket:      sumReadWriteBucket,
		allowCacheExternalPaths: allowCacheExternalPaths,
	}
}

func (m *moduleCacher) GetModule(
	ctx context.Context,
	modulePin bufmoduleref.ModulePin,
) (b