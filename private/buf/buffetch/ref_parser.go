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

package buffetch

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bufbuild/buf/private/buf/buffetch/internal"
	"github.com/bufbuild/buf/private/bufpkg/bufmodule/bufmoduleref"
	"github.com/bufbuild/buf/private/pkg/app"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const (
	loggerName = "buffetch"
	tracerName = "bufbuild/buf"
)

type refParser struct {
	allowProtoFileRef bool
	logger            *zap.Logger
	fetchRefParser    internal.RefParser
	tracer            trace.Tracer
}

func newRefParser(logger *zap.Logger, options ...RefParserOption) *refParser {
	refParser := &refParser{}
	for _, option := range options {
		option(refParser)
	}
	fetchRefParserOptions := []internal.RefParserOption{
		internal.WithRawRefProcessor(newRawRefProcessor(refParser.allowProtoFileRef)),
		internal.WithSingleFormat(formatBin),
		internal.WithSingleFormat(formatJSON),
		internal.WithSingleFormat(
			formatBingz,
			internal.WithSingleDefaultCompressionType(
				internal.CompressionTypeGzip,
			),
		),
		internal.WithSingleFormat(
			formatJSONGZ,
			internal.WithSingleDefaultCompressionType(
				internal.CompressionTypeGzip,
			),
		),
		internal.WithArchiveFormat(
			formatTar,
			internal.ArchiveTypeTar,
		),
		internal.WithArchiveFormat(
			formatTargz,
			internal.ArchiveTypeTar,
			internal.WithArchiveDefaultCompressionType(
				internal.CompressionTypeGzip,
			),
		),
		internal.WithArchiveFormat(
			formatZip,
			internal.ArchiveTypeZip,
		),
		internal.WithGitFormat(formatGit),
		internal.WithDirFormat(formatDir),
		internal.WithModuleFormat(formatMod),
	}
	if refParser.allowProtoFileRef {
		fetchRefParserOptions = append(fetchRefParserOptions, internal.WithProtoFileFormat(formatProtoFile))
	}
	refParser.logger = logger.Named(loggerName)
	refParser.tracer = otel.GetTracerProvider().Tracer(tracerName)
	refParser.fetchRefParser = internal.NewRefParser(
		logger,
		fetchRefParserOptions...,
	)
	return refParser
}

func newImageRefParser(logger *zap.Logger) *refParser {
	return &refParser{
		logger: logger.Named(loggerName),
		fetchRefParser: internal.NewRefParser(
			logger,
			internal.WithRawRefProcessor(processRawRefImage),
			internal.WithSingleFormat(formatBin),
			internal.WithSingleFormat(formatJSON),
			internal.WithSingleFormat(
				formatBingz,
				internal.WithSingleDefaultCompressionType(
					internal.CompressionTypeGzip,
				),
			),
			internal.WithSingleFormat(
				formatJSONGZ,
				internal.WithSingleDefaultCompressionType(
					internal.CompressionTypeGzip,
				),
			),
		),
		tracer: otel.GetTracerProvider().Tracer(tracerName),
	}
}

func newSourceRefParser(logger *zap.Logger) *refParser {
	return &refParser{
		logger: logger.Named(loggerName),
		fetchRefParser: internal.NewRefParser(
			logger,
			internal.WithRawRefProcessor(processRawRefSource),
			internal.WithArchiveFormat(
				formatTar,
				internal.ArchiveTypeTar,
			),
			internal.WithArchiveFormat(
				formatTargz,
				internal.ArchiveTypeTar,
				internal.WithArchiveDefaultCompressionType(
					internal.CompressionTypeGzip,
				),
			),
			internal.WithArchiveFormat(
				formatZip,
				internal.ArchiveTypeZip,
			),
			internal.WithGitFormat(formatGit),
			internal.WithDirFormat(formatDir),
		),
		tracer: otel.GetTracerProvider().Tracer(tracerName),
	}
}

func newModuleRefParser(logger *zap.Logger) *refParser {
	return &refParser{
		logger: logger.Named(loggerName),
		fetchRefParser: internal.NewRefParser(
			logger,
			internal.WithRawRefProcessor(processRawRefModule),
			internal.WithModuleFormat(formatMod),
		),
		tracer: otel.GetTracerProvider().Tracer(tracerName),
	}
}

func newSourceOrModuleRefParser(logger *zap.Logger) *refParser {
	return &refParser{
		logger: logger.Name