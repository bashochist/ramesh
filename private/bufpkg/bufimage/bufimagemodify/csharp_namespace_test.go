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
	"strings"
	"testing"

	"github.com/bufbuild/buf/private/bufpkg/bufmodule/bufmoduleref"
	"github.com/bufbuild/buf/private/pkg/normalpath"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestCsharpNamespaceEmptyOptions(t *testing.T) {
	t.Parallel()
	dirPath := filepath.Join("testdata", "emptyoptions")
	t.Run("with SourceCodeInfo", func(t *testing.T) {
		t.Parallel()
		image := testGetImage(t, dirPath, true)
		assertFileOptionSourceCodeInfoEmpty(t, image, csharpNamespacePath, true)

		sweeper := NewFileOptionSweeper()
		csharpNamespaceModifier := CsharpNamespace(zap.NewNop(), sweeper, nil, nil, nil)

		modifier := NewMultiModifier(csharpNamespaceModifier, ModifierFunc(sweeper.Sweep))
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
		assertFileOptionSourceCodeInfoEmpty(t, image, csharpNamespacePath, false)

		sweeper := NewFileOptionSweeper()
		modifier := CsharpNamespace(zap.NewNop(), sweeper, nil, nil, nil)
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
		assertFileOptionSourceCodeInfoEmpty(t, image, csharpNamespacePath, true)

		sweeper := NewFileOptionSweeper()
		csharpNamespaceModifier := CsharpNamespace(zap.NewNop(), sweeper, nil, nil, map[string]string{"a.proto": "foo"})

		modifier := NewMultiModifier(csharpNamespaceModifier, ModifierFunc(sweeper.Sweep))
		err := modifier.Modify(
			context.Background(),
			image,
		)
		require.NoError(t, err)
		// Overwritten with "foo" in the namespace
		assertFileOptionSourceCodeInfoEmpty(t, image, csharpNamespacePath, true)
		require.Equal(t, 1, len(image.Files()))
		descriptor := image.Files()[0].Proto()
		assert.Equal(t, "foo", descriptor.GetOptions().GetCsharpNamespace())
	})

	t.Run("without SourceCodeInfo and with per-file overrides", func(t *testing.T) {
		t.Parallel()
		image := testGetImage(t, dirPath, false)
		assertFileOptionSourceCodeInfoEmpty(t, image, csharpNamespacePath, false)

		sweeper := NewFileOptionSweeper()
		modifier := CsharpNamespace(zap.NewNop(), sweeper, nil, nil, map[string]string{"a.proto": "foo"})
		err := modifier.Modify(
			context.Background(),
			image,
		)
		require.NoError(t, err)
		// Overwritten with "foo" in the namespace
		require.Equal(t, 1, len(image.Files()))
		descriptor := image.Files()[0].Proto()
		assert.