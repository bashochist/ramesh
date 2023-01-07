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

package git

import (
	"context"
	"errors"
	"net/http/cgi"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bufbuild/buf/private/pkg/app"
	"github.com/bufbuild/buf/private/pkg/command"
	"github.com/bufbuild/buf/private/pkg/storage"
	"github.com/bufbuild/buf/private/pkg/storage/storagemem"
	"github.com/bufbuild/buf/private/pkg/storage/storageos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestGitCloner(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	container, err := app.NewContainerForOS()
	require.NoError(t, err)
	runner := command.NewRunner()
	originDir, workDir := createGitDirs(ctx, t, container, runner)

	t.Run("default", func(t *testing.T) {
		t.Parallel()
		readBucket := readBucketForName(ctx, t, runner, workDir, 1, nil, false)

		content, err := storage.ReadPath(ctx, readBucket, "test.proto")
		require.NoError(t, err)
		assert.Equal(t, "// commit 2", string(content), "expected the commit on local-branch to be checked out")
		_, err = readBucket.Stat(ctx, "nonexistent")
		assert.True(t, storage.IsNotExist(err))
		_, err = storage.ReadPath(ctx, readBucket, "submodule/test.proto")
		assert.True(t, storage.IsNotExist(err))
	})

	t.Run("default_submodule", func(t *testing.T) {
		t.Parallel()
		readBucket := readBucketForName(ctx, t, runner, workDir, 1, nil, true)

		content, err := storage.ReadPath(ctx, readBucket, "test.proto")
		require.NoError(t, err)
		assert.Equal(t, "// commit 2", string(content), "expected the commit on local-branch to be checked out")
		_, err = readBucket.Stat(ctx, "nonexistent")
		assert.True(t, storage.IsNotExist(err))
		content, err = storage.ReadPath(ctx, readBucket, "submodule/test.proto")
		require.NoError(t, err)
		assert.Equal(t, "// submodule", string(content))
	})

	t.Run("main", func(t *testing.T) {
		t.Parallel()
		readBucket := readBucketForName(ctx, t, runner, workDir, 1, NewBranchName("main"), false)

		content, err := storage.ReadPath(ctx, readBucket, "test.proto")
		require.NoError(t, err)
		assert.Equal(t, "// commit 1", string(content))
		_, err = readBucket.Stat(ctx, "nonexistent")
		assert.True(t, storage.IsNotExist(err))
	})

	t.Run("origin/main", func(t *testing.T) {
		t.Parallel()
		readBucket := readBucketForName(ctx, t, runner, workDir, 1, NewBranchName("origin/main"), false)

		content, err := storage.ReadPath(ctx, readBucket, "test.proto")
		require.NoError(t, err)
		assert.Equal(t, "// commit 3", string(content))
		_, err = readBucket.Stat(ctx, "nonexistent")
		assert.True(t, storage.IsNotExist(err))
	})

	t.Run("origin/remote-branch", func(t *testing.T) {
		t.Parallel()
		readBucket := readBucketForName(ctx, t, runner, workDir, 1, NewBranchName("origin/remote-branch"), false)

		content, err := storage.ReadPath(ctx, readBucket, "test.proto")
		require.NoError(t, err)
		assert.Equal(t, "// commit 4", string(content))
		_, err = readBucket.Stat(ctx, "nonexistent")
		assert.True(t, storage.IsNotExist(err))
	})

	t.Run("remote-tag", func(t *testing.T) {
		t.Parallel()
		readBucket := readBucketForName(ctx, t, runner, workDir, 1, NewTagName("remote-tag"), false)

		content, err := storage.ReadPath(ctx, readBucket, "test.proto")
		require.NoError(t, err)
		assert.Equal(t, "// commit 4", string(content))
		_, err = readBucket.Stat(ctx, "nonexistent")
		assert.True(t, storage.IsNotExist(err))
	})

	t.Run("branch_and_main_ref", func(t *testing.T) {
		t.Parallel()
		readBucket := readBucketForName(ctx, t, runner, workDir, 2, NewRefNameWithBranch("HEAD~", "main"), false)

		content, err := storage.ReadPath(ctx, readBucket, "test.proto")
		require.NoError(t, err)
		assert.Equal(t, "// commit 0", string(content))
		_, err = readBucket.Stat(ctx, "nonexistent")
		assert.True(t, storage.IsNotExist(err))
	})

	t.Run("branch_and_ref", func(t *testing.T) {
		t.Parallel()
		readBucket := readBucketForName(ctx, t, runner, workDir, 2, NewRefNameWithBranch("local-branch~", "local-branch"), false)

		content, err := storage.ReadPath(ctx, readBucket, "test.proto")
		require.NoError(t, err)
		assert.Equal(t, "// commit 1", string(content))
		_, err = readBucket.Stat(ctx, "nonexistent")
		assert.True(t, storage.IsNotExist(err))
	})

	t.Run("HEAD", func(t *testing.T) {
		t.Parallel()
		readBucket := readBucketForName(ctx, t, runner, workDir, 1, NewRefName("HEAD"), false)

		content, err := storage.ReadPath(ctx, readBucket, "test.proto")
		require.NoError(t, err)
		assert.Equal(t, "// commit 2", string(content))
		_, err = readBucket.Stat(ctx, "nonexistent")
		assert.True(t, storage.IsNotExist(err))
	})

	t.Run("commit-local", func(t *testing.T) {
		t.Parallel()
		revParseBytes, err := command.RunStdout(ctx, container, runner, "git", "-C", workDir, "rev-parse", "HEAD~")
		require.NoError(t, err)
		readBucket := readBucketForName(ctx, t, runner, workDir, 2, NewRefName(strings.TrimSpace(string(revParseBytes))), false)

		content, err := storage.ReadPath(ctx, readBucket, "test.proto")
		require.NoError(t, err)
		assert.Equal(t, "// commit 1", string(content))
		_, err = readBucket.Stat(ctx, "nonexistent")
		assert.True(t, storage.IsNotExist(err))
	})

	t.Run("commit-remote", func(t *testing.T) {
		t.Parallel()
		revParseBytes, err := command.RunStdout(ctx, container, runner, "git", "-C", originDir, "rev-parse", "remote-branch~")
		require.NoError(t, err)
		readBucket := readBucketForName(ctx, t, runner, workDir, 2, NewRefNameWithBranch(strings.TrimSpace(string(revParseBytes)), "origin/remote-branch"), false)

		content, err := storage.ReadPath(ctx, readBucket, "test.proto")
		require.NoError(t, err)
		assert.Equal(t, "// commit 3", string(content))
		_, err = readBucket.Stat(ctx, "nonexistent")
		assert.True(t, storage.IsNotExist(err))
	})
}

func readBucketForName(ctx context.Context, t *testing.T, runner command.Runner, path string, depth uint32, name Name, recurseSubmodules bool) storage.ReadBucket {
	t.Helper()
	storageosProvider := storageos.NewProvider(storageos.ProviderWithSymlinks())
	cloner := NewCloner(zap.NewNop(), storageosProvider, runner, ClonerOptions{})
	envContainer, err := app.NewEnvContainerForOS()
	require.NoError(t, err)

	readWriteBucket := storagemem.NewReadWriteBucket()
	err = cloner.CloneToBucket(
		ctx,
		envContainer,
		"file://"+filepath.Join(path, ".git"),
		depth,
		readWriteBucket,
		CloneToBucketOptions{
			Mapper:            storage.MatchPathExt(".proto"),
			Name:              name,
			RecurseSubmodules: recurseSubmodules,
		},
	)
	require.NoError(t, err)
	return readWriteBucket
}

func createGitDirs(
	ctx context.Context,
	t *testing.T,
	container app.EnvStdioContainer,
	runner command.Runner,
) (string, string) {
	tmpDir := t.TempDir()

	submodulePath := filepath.Join(tmpDir, "submodule")
	require.NoError(t, os.MkdirAll(submodulePath, os.ModePerm))
	runCommand(ctx, t, container, runner, "git", "-C", submodulePath, "init")
	runCommand(ctx, t, container, runner, "git", "-C", submodulePath, "config", "user.email", "tests@buf.build")
	runComma