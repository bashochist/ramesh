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
		runDiffTest(t, "testdata/options", []string{"pkg.Foo", "pkg.FooEnum", "pkg.FooService"}, "all.txtar")
	})
	t.Run("exclude-options", func(t *testing.T) {
		runDiffTest(t, "testdata/options", []string{"pkg.Foo", "pkg.FooEnum", "pkg.FooService"}, "all-exclude-options.txtar", WithExcludeCustomOptions())
	})
	t.Run("files", func(t *testing.T) {
		runDiffTest(t, "testdata/options", []string{"Files"}, "Files.txtar")
	})
	t.Run("all-with-files", func(t *testing.T) {
		runDiffTest(t, "testdata/options", []string{"pkg.Foo", "pkg.FooEnum", "pkg.FooService", "Files"}, "all-with-Files.txtar")
	})
}

func TestNesting(t *testing.T) {
	t.Parallel()
	t.Run("message", func(t *testing.T) {
		runDiffTest(t, "testdata/nesting", []string{"pkg.Foo"}, "message.txtar")
	})
	t.Run("recursenested", func(t *testing.T) {
		runDiffTest(t, "testdata/nesting", []string{"pkg.Foo.NestedFoo.NestedNestedFoo"}, "recursenested.txtar")
	})
	t.Run("enum", func(t *testing.T) {
		runDiffTest(t, "testdata/nesting", []string{"pkg.FooEnum"}, "enum.txtar")
	})
	t.Run("usingother", func(t *testing.T) {
		runDiffTest(t, "testdata/nesting", []string{"pkg.Baz"}, "usingother.txtar")
	})
}

func TestImportModifiers(t *testing.T) {
	t.Parallel()
	t.Run("regular_weak", func(t *testing.T) {
		runDiffTest(t, "testdata/importmods", []string{"ImportRegular", "ImportWeak"}, "regular_weak.txtar")
	})
	t.Run("weak_public", func(t *testing.T) {
		runDiffTest(t, "testdata/importmods", []string{"ImportWeak", "ImportPublic"}, "weak_public.txtar")
	})
	t.Run("regular_public", func(t *testing.T) {
		runDiffTest(t, "testdata/importmods", []string{"ImportRegular", "ImportPublic"}, "regular_public.txtar")
	})
	t.Run("noimports", func(t *testing.T) {
		runDiffTest(t, "testdata/importmods", []string{"NoImports"}, "noimports.txtar")
	})
}

func TestExtensions(t *testing.T) {
	t.Parallel()
	runDiffTest(t, "testdata/extensions", []string{"pkg.Foo"}, "extensions.txtar")
	runDiffTest(t, "testdata/extensions", []string{"pkg.Foo"}, "extensions-excluded.txtar", WithExcludeKnownExtensions())
}

func TestPackages(t *testing.T) {
	t.Parallel()
	runDiffTest(t, "testdata/packages", []string{""}, "root.txtar")
	runDiffTest(t, "testdata/packages", []string{"foo"}, "foo.txtar")
	runDiffTest(t, "testdata/packages", []string{"foo.bar"}, "foo.bar.txtar")
	runDiffTest(t, "testdata/packages", []string{"foo.bar.baz"}, "foo.bar.baz.txtar")
}

func TestAny(t *testing.T) {
	t.Parallel()
	runDiffTest(t, "testdata/any", []string{"ExtendedAnySyntax"}, "c1.txtar")
	runDiffTest(t, "testdata/any", []string{"ExtendedAnySyntax_InField"}, "c2.txtar")
	runDiffTest(t, "testdata/any", []string{"ExtendedAnySyntax_InList"}, "c3.txtar")
	runDiffTest(t, "testdata/any", []string{"ExtendedAnySyntax_InMap"}, "c4.txtar")
	runDiffTest(t, "testdata/any", []string{"NormalMessageSyntaxValidType"}, "d.txtar")
	runDiffTest(t, "testdata/any", []string{"NormalMessageSyntaxInvalidType"}, "e.txtar")
}

func TestTransitivePublic(t *testing.T) {
	ctx := context.Background()
	bucket, err := storagemem.NewReadBucket(map[string][]byte{
		"a.proto": []byte(`syntax = "proto3";package a;message Foo{}`),
		"b.proto": []byte(`syntax = "proto3";package b;import public "a.proto";message Bar {}`),
		"c.proto": []byte(`syntax = "proto3";package c;import "b.proto";message Baz{ a.Foo foo = 1; }`),
	})
	require.NoError(t, err)
	module, er