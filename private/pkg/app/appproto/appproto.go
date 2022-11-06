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
	// toResponse returns the resulting CodeGeneratorResponse. This must
	// only be called after all writing has been completed.
	toResponse() *pluginpb.CodeGeneratorResponse
}

// Handler is a protoc plugin handler.
type Handler interface {
	// Handle handles the plugin.
	//
	// This function can assume the request is valid.
	// This should only return error on system error.
	// Plugin generation errors should be added with AddError.
	// See https://github.com/protocolbuffers/protobuf/blob/95e6c5b4746dd7474d540ce4fb375e3f79a086f8/src/google/protobuf/compiler/plugin.proto#L100
	Handle(
		ctx context.Context,
		container app.EnvStderrContainer,
		responseWriter ResponseBuilder,
		request *pluginpb.CodeGeneratorRequest,
	) error
}

// HandlerFunc is a handler function.
type HandlerFunc func(
	context.Context,
	app.EnvStderrContainer,
	ResponseBuilder,
	*pluginpb.CodeGeneratorRequest,
) error

// Handle implements Handler.
func (h HandlerFunc) Handle(
	ctx context.Context,
	container app.EnvStderrContainer,
	responseWriter ResponseBuilder,
	request *pluginpb.CodeGeneratorRequest,
) error {
	return h(ctx, container, responseWriter, request)
}

// Main runs the plugin using app.Main and the Handler.
func Main(ctx context.Context, handler Handler) {
	app.Main(ctx, newRunFunc(handler))
}

// Run runs the plugin using app.Main and the Handler.
//
// The exit code can be determined using app.GetExitCode.
func Run(ctx context.Context, container app.Container, handler Handler) error {
	return app.Run(ctx, container, newRunFunc(handler))
}

// Generator executes the Handler using protoc's plugin execution logic.
//
// If multiple requests are specified, these are executed in parallel and the
// result is combined into one response that is written.
type Generator interface {
	// Generate generates a CodeGeneratorResponse for the given CodeGeneratorRequests.
	//
	// A new ResponseBuilder is constructed for every invocation of Generate and is
	// used to consolidate all of the CodeGeneratorResponse_Files returned from a single
	// plugin into a single CodeGeneratorResponse.
	Generate(
		ctx context.Context,
		container app.EnvStderrContainer,
		requests []*pluginpb.CodeGeneratorRequest,
	) (*pluginpb.CodeGeneratorResponse, error)
}

// NewGenerator returns a new Generator.
func NewGenerator(
	logger *zap.Logger,
	handler Handler,
) Generator {
	return newGenerator(logger, handler)
}

// ResponseWriter handles the response and writes it to the given storage.WriteBucket
// without executing any plugins and handles insertion points as needed.
type ResponseWriter interface {
	// WriteResponse writes to the bucket with the given response. In practice, the
	// WriteBucket is most often an in-memory bucket.
	//
	// CodeGeneratorResponses are consolidated into the bucket, and insertion points
	// are applied in-place so that they can only access the files created in a single
	// generation invocation (just like protoc).
	WriteResponse(
		ctx context.Context,
		writeBucket storage.WriteBucket,
		response *pluginpb.CodeGeneratorResponse,
		options ...WriteResponseOption,
	) error
}

// NewResponseWriter returns a new ResponseWriter.
func NewResponseWriter(logger *zap.Logger) ResponseWriter {
	return newResponseWriter(logger)
}

// WriteResponseOption is an option for WriteResponse.
type WriteResponseOption func(*writeResponseOptions)

// WriteResponseWithInsertionPointReadBucket returns a new WriteResponseOption that uses the given
// Read