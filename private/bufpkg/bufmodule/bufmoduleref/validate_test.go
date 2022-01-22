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

package bufmoduleref

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateRemoteNotEmpty(t *testing.T) {
	err := ValidateRemoteNotEmpty("")
	require.Equal(t, "you must specify a remote module", err.Error())
	require.NoError(t, ValidateRemoteNotEmpty("buf.build"))
}

func TestValidateRemoteHasNoPaths(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		Name        string
		Input       string
		InvalidPath string
	}{
		{
			Name:        "Remote with two trailing slashes",
			Input:       "buf.build//",
			InvalidPath: "//",
		},
		{
			Name:        "Remote with a single path",
			Input:       "buf.build/path1",
			InvalidPath: "/path1",
		},
		{
			Name:        "Remote with a single path and trailing slash",
			Input:       "buf.b