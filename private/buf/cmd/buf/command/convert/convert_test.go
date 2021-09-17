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

package convert

import (
	"strings"
	"testing"

	"github.com/bufbuild/buf/private/pkg/app/appcmd"
	"github.com/bufbuild/buf/private/pkg/app/appcmd/appcmdtesting"
	"github.com/bufbuild/buf/private/pkg/app/appflag"
)

// This test is in its own file as opposed to buf_test because it needs to test a single module in testdata.
func TestConvertDir(t *testing.T) {
	cmd := func(use string) *appcmd.Command { return NewCommand("convert", appflag.NewBuilder("convert")) }
	t.Run("default-input", func(t *testing.T) {
		appcmdtesting.RunCommandExitCodeStdout(
			t,
			cmd,
			0,
			`{"one":"55"}`,
			nil,
			nil,
			"--type",
			"buf.Foo",
			"--from",
			"testdata/convert/bin_json/payload.bin",
		)
	})
	t.Run("from-stdin", func(t *testing.T) {
		appcmdtesting.RunCommandExitCodeStdoutStdinFile(
			t,
			cmd,
			0,
			`{"one":"55"}`,
			nil,
			"testdata/convert/bin_json/payload.bin",

			"--type",