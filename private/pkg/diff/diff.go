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

// Package diff implements diffing.
//
// Should primarily be used for testing.
package diff

// Largely copied from https://github.com/golang/go/blob/master/src/cmd/gofmt/gofmt.go
//
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// https://github.com/golang/go/blob/master/LICENSE

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/bufbuild/buf/private/pkg/command"
)

// Diff does a diff.
//
// Returns nil if no diff.
func Diff(
	ctx context.Context,
	runner command.Runner,
	b1 []byte,
	b2 []byte,
	filename1 string,
	filename2 string,
	options ...DiffOption,
) ([]byte, error) {
	diffOptions := newDiffOptions()
	for _, option := r