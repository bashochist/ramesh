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

package bufimagebuildtesting

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bufbuild/buf/private/bufpkg/bufanalysis"
	"github.com/bufbuild/buf/private/bufpkg/bufimage"
	"github.com/bufbuild/buf/private/bufpkg/bufimage/bufimagebuild"
	"github.com/bufbuild/buf/private/bufpkg/bufmodule"
	"github.com/bufbuild/buf/private/bufpkg/bufmodule/bufmodulebuild"
	"github.com/bufbuild/buf/private/bufpkg/bufmodule/bufmoduleconfig"
	"github.com/bufbuild/buf/private/bufpkg/buftesting"
	"github.com/bufbuild/buf/private/pkg/command"
	"github.com/bufbuild/buf/private/pkg/prototesting"
	"github.com/bufbuild/buf/private/pkg/storage/storageos"
	"github.com/bufbuild/buf/private/pkg/tmp"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"golang.org/x/tools/txtar"
	"google.golang.org/protobuf/types/descriptorpb"
)

// Fuzz is the entrypoint for the fuzzer.
// We use https://github.com/dvyukov/go-fuzz for fuzzing.
// Please follow the instructions
// in their README for help with running the fuzz targets.
func Fuzz(data []byte) int {
	ctx := context.Background()
	runner := command.NewRunner()
	result, err := fuzz(ctx, runner, data)
	if err != nil {
		// data was invalid in some way
		return -1
	}
	return result.panicOrN(ctx)
}

func fuzz(ctx context.Context, runner command.Runner, data []byte) (_ *fuzzResult, retErr error) {
	dir, err := tmp.NewDir()
	if err != nil {
		return nil, err
	}
	defer func() {
		retErr = multierr.Append(retErr, dir.Close())
	}()
	if err := untxtar(data, dir.AbsPath()); err != nil {
		return nil, err
	}

	filePaths, err := buftesting.GetProtocFilePathsErr(ctx, dir.AbsPath(), 0)
	if err != nil {
		return nil, err
	}

	actualProtocFileDescriptorSet, protocErr := prototesting.GetProtocFileDescriptorSet(
		ctx,
		runner,
		[]string{dir.AbsPath()},
		filePaths,
		false,
		false,
	)

	image, bufAnnotations, bufErr := fuzzBuild(ctx, dir.AbsPath())
	return newFuzzResult(
		runner,
		bufAnnotations,
		bufErr,
		protocErr,
		actualProtocFileDescriptorSet,
		image,
	), nil
}

// fuzzBuild does a builder.Build for a fuzz test.
func fuzzBuild(ctx context.Context, dirPath string) (bufimage.Image, []bufanalysis.FileAnnotation, error) {
	moduleFileSet, err := fuzzGetModuleFileSet(ctx, dirPath)
	if err != nil {
		return nil, nil, err
	}
	builder := bufimagebuild.NewBuilder(zap.NewNop())
	opt := bufimagebuild.WithExcludeSourceCodeInfo()
	return builder.Build(ctx, moduleFileSet, opt)
}

// fuzzGetModuleFileSet gets the bufmodule.ModuleFileSet for a fuzz test.
func fuzzGetModuleFileSet(ctx context.Context, dirPath string) (bufmodule.ModuleFileSet, error) {
	storageosProvider 