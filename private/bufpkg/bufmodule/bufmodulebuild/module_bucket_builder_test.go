
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
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bufbuild/buf/private/bufpkg/bufcheck/bufbreaking/bufbreakingconfig"
	"github.com/bufbuild/buf/private/bufpkg/bufcheck/buflint/buflintconfig"
	"github.com/bufbuild/buf/private/bufpkg/bufmodule/bufmoduleconfig"
	"github.com/bufbuild/buf/private/bufpkg/bufmodule/bufmoduleref"
	"github.com/bufbuild/buf/private/bufpkg/bufmodule/bufmoduletesting"
	"github.com/bufbuild/buf/private/pkg/normalpath"
	"github.com/bufbuild/buf/private/pkg/storage"
	"github.com/bufbuild/buf/private/pkg/storage/storagemem"
	"github.com/bufbuild/buf/private/pkg/storage/storageos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBucketGetFileInfos1(t *testing.T) {
	config, err := bufmoduleconfig.NewConfigV1(
		bufmoduleconfig.ExternalConfigV1{
			Excludes: []string{"proto/b"},
		},
	)
	require.NoError(t, err)
	testBucketGetFileInfos(
		t,
		"testdata/1",
		config,
		bufmoduletesting.NewFileInfo(t, "proto/a/1.proto", "testdata/1/proto/a/1.proto", false, nil, ""),
		bufmoduletesting.NewFileInfo(t, "proto/a/2.proto", "testdata/1/proto/a/2.proto", false, nil, ""),