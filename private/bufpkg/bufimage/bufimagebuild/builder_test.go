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

package bufimagebuild

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/bufbuild/buf/private/bufpkg/bufanalysis"
	"github.com/bufbuild/buf/private/bufpkg/bufimage"
	"github.com/bufbuild/buf/private/bufpkg/bufimage/bufimageutil"
	"github.com/bufbuild/buf/private/bufpkg/bufmodule"
	"github.com/bufbuild/buf/private/bufpkg/bufmodule/bufmodulebuild"
	"github.com/bufbuild/buf/private/bufpkg/bufmodule/bufmoduleconfig"
	"github.com/bufbuild/buf/private/bufpkg/buftesting"
	"github.com/bufbuild/buf/private/pkg/command"
	"github.com/bufbuild/buf/private/pkg/normalpath"
	"github.com/bufbuild/buf/private/pkg/protosource"
	"github.com/bufbuild/buf/private/pkg/prototesting"
	"github.com/bufbuild/buf/private/pkg/storage/storageos"
	"github.com/bufbuild/buf/private/pkg/testingextended"
	"github.com/bufbuild/buf/private/pkg/thread"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var buftestingDirPath = filepath.Join(
	"..",
	"..",
	"buftesting",
)

func TestGoogleapis(t *testing.T) {
	testingextended.SkipIfShort(t)
	t.Parallel()
	image := testBuildGoogleapis(t, true)
	assert.Equal(t, buftesting.NumGoogleapisFilesWithImports, len(image.Files()))
	assert.Equal(
		t,
		[]string{
			"google/protobuf/any.proto",
			"google/protobuf/api.proto",
			"google/protobuf/descriptor.proto",
			"google/protobuf/duration.proto",
			"google/protobuf/empty.proto",
			"google/protobuf/field_mask.proto",
			"google/protobuf/source_context.proto",
			"google/protobuf/struct.proto",
			"google/protobuf/timestamp.proto",
			"google/protobuf/type.proto",
			"google/protobuf/wrappers.proto",
		},
		testGetImageImportPaths(image),
	)

	imageWithoutImports := bufimage.ImageWithoutImports(image)
	assert.Equal(t, buftesting.NumGoogleapisFiles, len(imageWithoutImports.Files()))
	imageWithoutImports = bufimage.ImageWithoutImports(imageWithoutImports)
	assert.Equal(t, buftesting.NumGoogleapisFiles, len(imageWithoutImports.Files()))

	imageWithSpecificNames, err := bufimage.ImageWithOnlyPathsAllowNotExist(
		image,
		[]string{
			"google/protobuf/descriptor.proto",
			"google/protobuf/api.proto",
			"google/type/date.proto",
			"google/foo/nonsense.proto",
		},
		nil,
	)
	assert.NoError(t, err)
	assert.Equal(
		t,
		[]string{
			"google/protobuf/any.proto",
			"google/protobuf/api.proto",
			"google/protobuf/descriptor.proto",
			"google/protobuf/source_context.proto",
			"google/protobuf/type.proto",
			"google/type/date.proto",
		},
		testGetImageFilePaths(imageWithSpecificNames),
	)
	imageWithSpecificNames, err = bufimage.ImageWithOnlyPathsAllowNotExist(
		image,
		[]string{
			"google/protobuf/descriptor.proto",
			"google/protobuf/api.proto",
			"google/type",
			"google/foo",
		},
		nil,
	)
	assert.NoError(t, err)
	assert.Equal(
		t,
		[]string{
			"google/protobuf/any.proto",
			"google/protobuf/api.proto",
			"google/protobuf/descriptor.proto",
			"google/protobuf/source_context.proto",
			"google/protobuf/type.proto",
			"google/protobuf/wrappers.proto",
			"google/type/calendar_period.proto",
			"google/type/color.proto",
			"google/type/date.proto",
			"google/type/dayofweek.proto",
			"google/type/expr.proto",
			"google/type/fraction.proto",
			"google/type/latlng.proto",
			"google/type/money.proto",
			"google/type/postal_address.proto",
			"google/type/quaternion.