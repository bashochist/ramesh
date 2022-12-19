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

// Matching the unix-like build tags in the Golang source i.e. https://github.com/golang/go/blob/912f0750472dd4f674b69ca1616bfaf377af1805/src/os/file_unix.go#L6

//go:build aix || darwin || dragonfly || freebsd || (js && wasm) || linux || netbsd || openbsd || solaris
// +build aix darwin dragonfly freebsd js,wasm linux netbsd openbsd solaris

package filepathextended

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRealClean(t *testing.T) {
	t.Parallel()
	path, err := RealClean("../filepathextended")
	assert.NoError(t, err)
	assert.Equal(t, ".", path)
	path, err = RealClean("../filepathextended/foo")
	assert.NoError(t, err)
	assert.Equal(t, "foo", path)
	path, err = RealClean("/foo")
	assert.NoError(t, err)
	assert.Equal(t, "/foo", path)
	path, err = RealClean("/foo/../bar")
	assert.NoError(t, err)
	assert.Equal(t, "/bar", path)
}

func TestWalkSymlinkSuccessNoSymlinks(t *testing.T) {
	t.Parallel()
	filePaths, err := testWalkGetRegularFilePaths(
		filepath.Join("testdata", "symlin