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
	"io"
	"net/http"

	"github.com/bufbuild/buf/private/buf/buffetch/internal"
	"github.com/bufbuild/buf/private/bufpkg/bufmodule"
	"github.com/bufbuild/buf/private/pkg/app"
	"github.com/bufbuild/buf/private/pkg/git"
	"github.com/bufbuild/buf/private/pkg/httpauth"
	"github.com/bufbuild/buf/private/pkg/storage/storageos"
	"github.com/bufbuild/buf/private/pkg/stringutil"
	"go.uber.org/zap"
)

const (
	// ImageEncodingBin is the binary image encoding.
	ImageEncodingBin ImageEncoding = iota + 1
	// ImageEncodingJSON is the JSON image encoding.
	ImageEncodingJSON
)

var (
	// ImageFormatsString is the string representation of all image formats.
	//
	// This does not include deprecated formats.
	ImageFormatsString = stringutil.SliceToString(imageFormatsNotDeprecated)
	// SourceDirFormatsString is the string representation of all source directory formats.
	// This includes all of the formats in SourceFormatsString except the protofile format.
	//
	// This does not include deprecated formats.
	SourceDirFormatsString = stringutil.SliceToString(sourceDirFormatsNotDeprecated)
	// SourceFormatsString is the string representation of all source formats.
	//
	// This does not include deprecated formats.
	SourceFormatsString = stringutil.SliceToString(sourceFormatsNotDeprecated)
	// ModuleFormatsString is the string representation of all module formats.
	//
	// Module formats are also source formats.
	//
	// This does not include deprecated formats.
	ModuleFormatsString = stringutil.SliceToString(moduleFormatsNotDeprecated)
	// SourceOrModuleFormatsString is the string representation of all source or module formats.
	//
	// This does not include deprecated formats.
	SourceOrModuleFormatsString = stringutil.SliceToString(sourceOrModuleFormatsNotDeprecated)
	// AllFormatsString is the string representation of all formats.
	//
	// This does not include deprecated formats.
	AllFormatsString = stringutil.SliceToString(allFormatsNotDeprecated)
)

// ImageEncoding is the encoding of the image.
type ImageEncoding int

// PathResolver resolves external paths to paths.
type PathResolver interface {
	// PathForExternalPath takes a path external to the asset and converts it to
	// a path that is relative to the asset.
	//
	// The returned path will be normalized and validated.
	//
	// Example:
	//   Directory: /foo/bar
	//   ExternalPath: /foo/bar/baz/bat.proto
	//   Path: baz/bat.proto
	//
	// Example:
	//   Directory: .
	//   ExternalPath: baz/bat.proto
	//   Path: baz/bat.proto
	PathForExternalPath(externalPath string) (string, error)
}

// Ref is an image file or source bucket reference.
type Ref interface {
	PathResolver

	internalRef() internal.Ref
}

// ImageRef is an image file reference.
type ImageRef interface {
	Ref
	ImageEncoding() ImageEncoding
	IsNull() bool
	internalFileRef() internal.FileRef
}

// SourceOrModuleRef is a source bucket or module reference.
type SourceOrModuleRef interface {
	Ref
	isSourceOrModuleRef()
}

// SourceRef is a source bucket reference.
type SourceRef interface {
	SourceOrModuleRef
	internalBucketRef() internal.BucketRef
}

// ModuleRef is a module reference.
type ModuleRef interface {
	SourceOrModuleRef
	internalModuleRef() internal.ModuleRef
}

// ProtoFileRef is a proto file reference.
type ProtoFileRef interface {
	SourceRef
	IncludePackageFiles() bool
	internalProtoFileRef() internal.ProtoFileRef
}

// ImageRefParser is an image ref parser for Buf.
type ImageRefParser interface {
	// GetImageRef gets the reference for the image file.
	GetImageRef(ctx context.Context, value string) (ImageRef, error)
}

// SourceRefParser is a source ref parser for Buf.
type SourceRefParser interface {
	// GetSourceRef gets the reference for the source file.
	GetSourceRef(ctx context.Context, value string) (SourceRef, error)
}

// ModuleRefParser is a source ref parser for Buf.
type ModuleRefParser interface {
	// GetModuleRef gets the reference for the source file.
	//
	// A module is a special type of source with additional properties.
	GetModuleRef(ctx context.Context, value string) (ModuleRef, error)
}

// SourceOrModuleRefParser is a source or module ref parser for Buf.
type SourceOrModuleRefParser interface {
	SourceRefParser
	ModuleRefParser

	// GetSourceOrModuleRef gets the reference for the image file or source bucket.
	GetSourceOrModuleRef(ctx context.Context, value string) (SourceOrModuleRef, error)
}

// RefParser is a ref parser for Buf.
type RefParser interface {
	ImageRefParser
	SourceOrModuleRefParser

	// GetRef gets the re