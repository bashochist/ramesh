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

package protoc

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/bufbuild/buf/private/bufpkg/buftesting"
	"github.com/bufbuild/buf/private/pkg/app"
	"github.com/bufbuild/buf/private/pkg/app/appcmd"
	"github.com/bufbuild/buf/private/pkg/app/appcmd/appcmdtesting"
	"github.com/bufbuild/buf/private/pkg/app/appflag"
	"github.com/bufbuild/buf/private/pkg/command"
	"github.com/bufbuild/buf/private/pkg/protoencoding"
	"github.com/bufbuild/buf/private/pkg/prototesting"
	"github.com/bufbuild/buf/private/pkg/storage"
	"github.com/bufbuild/buf/private/pkg/storage/storagearchive"
	"github.com/bufbuild/buf/private/pkg/storage/storagemem"
	"github.com/bufbuild/buf/private/pkg/storage/storageos"
	"github.com/bufbuild/buf/private/pkg/testingextended"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/descriptorpb"
)

var buftestingDirPath = filepath.Join(
	"..",
	"..",
	"..",
	"..",
	"..",
	"..",
	"..",
	"private",
	"bufpkg",
	"buftesting",
)

type testPluginInfo struct {
	name string
	opt  string
}

func TestOverlap(t *testing.T) {
	t.Parallel()
	// https://github.com/bufbuild/buf/issues/113
	appcmdtesting.RunCommandSuccess(
		t,
		func(name string) *appcmd.Command {
			return NewCommand(
				name,
				appflag.NewBuilder(name),
			)
		},
		nil,
		nil,
		nil,
		"-I",
		filepath.Join("testdata", "overlap", "a"),
		"-I",
		filepath.Join("testdata", "overlap", "b"),
		"-o",
		app.DevNullFilePath,
		filepath.Join("testdata", "overlap", "a", "1.proto"),
		filepath.Join("testdata", "overlap", "b", "2.proto"),
	)
}

func TestComparePrintFreeFieldNumbersGoogleapis(t *testing.T) {
	t.Parallel()
	googleapisDirPath := buftesting.GetGoogleapisDirPath(t, buftestingDirPath)
	filePaths := buftesting.GetProtocFilePaths(t, googleapisDirPath, 100)
	actualProtocStdout := bytes.NewBuffer(nil)
	runner := command.NewRunner()
	buftesting.RunActualProtoc(
		t,
		runner,
		false,
		false,
		googleapisDirPath,
		filePaths,
		nil,
		actualProtocStdout,
		fmt.Sprintf("--%s", printFreeFieldNumbersFlagName),
	)
	appcmdtesting.RunCommandSuccessStdout(
		t,
		func(name string) *appcmd.Command {
			return NewCommand(
				name,
				appflag.NewBuilder(name),
			)
		},
		actualProtocStdout.String(),
		nil,
		nil,
		append(
			[]string{
				"-I",
				googleapisDirPath,
				fmt.Sprintf("--%s", printFreeFieldNumbersFlagName),
			},
			filePaths...,
		)...,
	)
}

func TestCompareOutputGoogleapis(t *testing.T) {
	testingextended.SkipIfShort(t)
	t.Parallel()
	googleapisDirPath := buftesting.GetGoogleapisDirPath(t, buftestingDirPath)
	filePaths := buftesting.GetProtocFilePaths(t, googleapisDirPath, 100)
	runner := command.NewRunner()
	actualProtocFileDescriptorSet := buftesting.GetActualProtocFileDescriptorSet(
		t,
		runner,
		false,
		false,
		googleapisDirPath,
		filePaths,
	)
	bufProtocFileDescriptorSet := testGetBufProtocFileDescriptorSet(t, googleapisDirPath)
	prototesting.AssertFileDescriptorSetsEqual(t, command.NewRunner(), bufProtocFileDescriptorSet, actualProtocFileDescriptorSet)
}

func TestCompareGeneratedStubsGoogleapisGo(t *testing.T) {
	testingextended.SkipIfShort(t)
	t.Parallel()
	googleapisDirPath := buftesting.GetGoogleapisDirPath(t, buftestingDirPath)
	testCompareGeneratedStubs(
		t,
		command.NewRunner(),
		googleapisDirPath,
		[]testPluginInfo{
			{name: "go", opt: "Mgoogle/api/auth.proto=foo"},
		},
	)
}

fu