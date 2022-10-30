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

// Package appproto contains helper functionality for protoc plugins.
//
// Note this is currently implicitly tested through buf's protoc command.
// If this were split out into a separate package, testing would need to be
// moved to this package.
package appproto

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"unicode"
	"unicode/utf8"

	"github.com/bufbuild/buf/private/pkg/app"
	"github.com/bufbuild/buf/private/pkg/protodescriptor"
	"github.com/bufbuild/buf/private/pkg/protoencoding"
	"github.com/bufbuild/buf/private/pkg/storage"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/pluginpb"
)

const (
	// Our generated files in `private/gen/proto` are on average 15KB which isn't
	// an unreasonable amount of memory to reserve each time we process an insertion
	// point and will save a significant number of allocations.
	averageGeneratedFileSize = 15 * 1024
	// We don't use insertion points internally, but assume they are smaller than
	// entire generated files.
	averageInsertionPointSize = 1024
)

// ResponseBuilder builds CodeGeneratorResponses.
type ResponseBuilder interface {
	// AddFile adds the file to the response.
	//
	// Returns error if nil or the name is empty.
	// Warns to stderr if the name is already added or the name is not normalized.
	AddFile(*pluginpb.CodeGeneratorResponse_File) error
	// AddError adds the error message to the response.
	//
	// If there is an existing error message, this will be concatenated with a newline.
	// If message is empty, a message "error" will be added.
	AddError(message string)
	// SetFeatureProto3Optional sets the proto3 optional feature.
	SetFeatureProto3Optional()
	// toResponse returns the resulting CodeGe