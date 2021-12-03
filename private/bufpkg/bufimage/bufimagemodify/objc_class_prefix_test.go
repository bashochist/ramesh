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

package bufimagemodify

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/bufbuild/buf/private/bufpkg/bufmodule/bufmoduleref"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestObjcClassPrefixEmptyOptions(t *testing.T) {
	t.Parallel()
	dirPath := filepath.Join("testdata", "emptyoptions")
	t.Run("with SourceCodeInfo", func(t *testing.T) {
		t.Parallel()
		image := testGetImage(t, dirPath, true)
		assertFileOptionSourceCodeInfoEmpty(t, image, objcClassPrefixPath, true)

		sweeper := NewFileOptionSweeper()
		objcClassPrefixModifier := ObjcClassPrefix(zap.NewNop(), sweeper, "", nil, nil, nil)

		modifier := NewMultiModifier(objcClassPrefixModifier, ModifierFunc(sweeper.Sweep))
		err := modifier.Modify(
			context.Background(),
			image,
		)
		require.NoError(t, err)
		assert.Equal(t, testGetImage(t, dirPath, true), image)
	})

	t.Run("without SourceCodeInfo", func(t *testing.T) {
		t.Parallel()
		image := testGetImage(t, dirPath, false)
		assertFileOptionSourceCodeInfoEmpty(t, image, objcClassPrefixPath, false)

		sweeper := NewFileOptionSweeper()
		modifier := ObjcClassPrefix(zap.NewNop(), sweeper, "", nil, nil, nil)
		err := modifier.Modify(
			context.Background(),
			image,
		)
		require.NoError(t, err)
		assert.Equal(t, testGetImage(t, dirPath, false), image)
	})

	t.Run("with SourceCodeInfo and per-file overrides", func(t *testing.T) {
		t.Parallel()
		image := testGetImage(t, dirPath, true)
		assertFileOptionSourceCodeInfoEmpty(t, image, objcClassPrefixPath, true)

		sweeper := NewFileOptionSweeper()
		objcClassPrefixModifier := ObjcClassPrefix(zap.NewNop(), sweeper, "", nil, nil, map[string]string{"a.proto": "override"})

		modifier := NewMultiModifier(objcClassPrefixModifier, ModifierFunc(sweeper.Sweep))
		err := modifier.Modify(
			context.Background(),
			image,
		)
		require.NoError(t, err)
		require.Equal(t, 1, len(image.Files()))
		descriptor := image.Files()[0].Proto()
		require.Equal(t, "override", descriptor.GetOptions().GetObjcClassPrefix())
	})

	t.Run("without SourceCodeInfo and with per-file overrides", func(t *testing.T) {
		t.Parallel()
		image := testGetImage(t, dirPath, false)
		assertFileOptionSourceCodeInfoEmpty(t, image, objcClassPrefixPath, false)

		sweeper := NewFileOptionSweeper()
		modifier := ObjcClassPrefix(zap.NewNop(), sweeper, "", nil, nil, map[string]string{"a.proto": "override"})
		err := modifier.Modify(
			context.Background(),
			image,
		)
		require.NoError(t, err)
		require.Equal(t, 1, len(image.Files()))
		descriptor := image.Files()[0].Proto()
		require.Equal(t, "override", descriptor.GetOptions().GetObjcClassPrefix())
	})
}

func TestObjcClassPrefixAllOptions(t *testing.T) {
	t.Parallel()
	dirPath := filepath.Join("testdata", "alloptions")
	t.Run("with SourceCodeInfo", func(t *testing.T) {
		t.Parallel()
		image := testGetImage(t, dirPath, true)
		assertFileOptionSourceCodeInfoNotEmpty(t, image, objcClassPrefixPath)

		sweeper := NewFileOptionSweeper()
		objcClassPrefixModifier := ObjcClassPrefix(zap.NewNop(), sweeper, "", nil, nil, nil)

		modifier := NewMultiModifier(objcClassPrefixModifier, ModifierFunc(sweeper.Sweep))
		err := modifier.Modify(
			context.Background(),
			image,
		)
		require.NoError(t, err)

		for _, imageFile := range image.Files() {
			descriptor := imageFile.Proto()
			assert.Equal(t, "foo", descriptor.GetOptions().GetObjcClassPrefix())
		}
		assertFileOptionSourceCodeInfoNotEmpty(t, image, objcClassPrefixPath)
	})

	t.Run("without SourceCodeInfo", func(t *testing.T) {
		t.Parallel()
		image := testGetImage(t, dirPath, false)
		assertFileOptionSourceCodeInfoEmpty(t, image, objcClassPrefixPath, false)

		sweeper := NewFileOptionSweeper()
		modifier := ObjcClassPrefix(zap.NewNop(), sweeper, "", nil, nil, nil)
		err := modifier.Modify(
			context.Background(),
			image,
		)
		require.NoError(t, err)

		for _, imageFile := range image.Files() {
			descriptor := imageFile.Proto()
			assert.Equal(t, "foo", descriptor.GetOptions().GetObjcClassPrefix())
		}
		assertFileOptionSourceCodeInfoEmpty(t, image, objcClassPrefixPath, false)
	})

	t.Run("with SourceCodeInfo and per-file overrides", func(t *testing.T) {
		t.Parallel()
		image := testGetImage(t, dirPath, true)
		assertFileOptionSourceCodeInfoNotEmpty(t, image, objcClassPrefixPath)

		sweeper := NewFileOptionSweeper()
		objcClassPrefixModifier := ObjcClassPrefix(zap.NewNop(), sweeper, "", nil, nil, map[string]string{"a.proto": "override"})

		modifier := NewMultiModifier(objcClassPrefixModifier, ModifierFunc(sweeper.Sweep))
		err := modifier.Modify(
			context.Background(),
			image,
		)
		require.NoError(t, err)

		for _, imageFile := range image.Files() {
			descriptor := imageFile.Proto()
			if imageFile.Pa