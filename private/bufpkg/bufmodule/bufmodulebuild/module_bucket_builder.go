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

package bufmodulebuild

import (
	"context"

	"github.com/bufbuild/buf/private/bufpkg/bufconfig"
	"github.com/bufbuild/buf/private/bufpkg/buflock"
	"github.com/bufbuild/buf/private/bufpkg/bufmodule"
	"github.com/bufbuild/buf/private/bufpkg/bufmodule/bufmoduleconfig"
	"github.com/bufbuild/buf/private/pkg/normalpath"
	"github.com/bufbuild/buf/private/pkg/storage"
	"github.com/bufbuild/buf/private/pkg/storage/storagemem"
)

// BuiltModule ties a bufmodule.Module with the configuration and a bucket
// containing just the files required to build it.
type BuiltModule struct {
	bufmodule.Module
	Bucket storage.ReadBucket
}

type moduleBucketBuilder struct {
	opt buildOptions
}

func newModuleBucketBuilder(
	options ...BuildOption,
) *moduleBucketBuilder {
	opt := buildOptions{}
	for _, option := range options {
		option(&opt)
	}
	return &moduleBucketBuilder{opt: opt}
}

// BuildForBucket is an alternative constructo