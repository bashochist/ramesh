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

package bufimageutil

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sort"
	"testing"

	"github.com/bufbuild/buf/private/bufpkg/bufimage"
	"github.com/bufbuild/buf/private/bufpkg/bufimage/bufimagebuild"
	"github.com/bufbuild/buf/private/bufpkg/bufmodule"
	"github.com/bufbuild/buf/private/bufpkg/bufmodule/bufmoduleref"
	"github.com/bufbuild/buf/private/pkg/protoencoding"
	"github.com/bufbuild/buf/private/pkg/storage"
	"github.com/bufbuild/buf/private/pkg/storage/storagemem"
	"github.com/bufbuild/buf/private/pkg/storage/storageos"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoprint"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"golang.org/x/tools/txtar"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// IF YOU HAVE ANY FAILING TESTS IN HERE, ESPECIALLY AFTER A PROTOC UPGRADE,
// SWITCH THIS TO TRUE, RE-RUN THE TESTS ***TWICE***, AND THEN SWITCH BACK TO FALSE
const shouldUpdateExpectations = false

func TestOptions(t *testing.T) {
	t.Parallel()
	t.Run("message", func(t *testing.T) {
		runDiffTest(t, "testdata/options", []string{"pkg.Foo"}, "pkg.Foo.txtar")
	})
	t.Run("enum", func(t *testing.T) {
		runDiffTest(t, "testdata/options", []string{"pkg.FooEnum"}, "pkg.FooEnum.txtar")
	})
	t.Run("service", func(t *testing.T) {
		runDiffTest(t, "testdata/options", []string{"pkg.FooService"}, "pkg.FooService.txtar")
	})
	t.Run("method", func(t *testing.T) {
		runDiffTest(t, "testdata/options", []string{"pkg.FooService.Do"}, "pkg.FooService.Do.txtar")
	})
	t.Run("all", func(t *testing.T) {
		runDiffTest(t, "testdata/opt