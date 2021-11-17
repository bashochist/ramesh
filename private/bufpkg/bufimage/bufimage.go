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

package bufimage

import (
	"fmt"
	"sort"

	"github.com/bufbuild/buf/private/bufpkg/bufmodule/bufmoduleref"
	imagev1 "github.com/bufbuild/buf/private/gen/proto/go/buf/alpha/image/v1"
	"github.com/bufbuild/buf/private/pkg/normalpath"
	"github.com/bufbuild/buf/private/pkg/protodescriptor"
	"github.com/bufbuild/buf/private/pkg/protoencoding"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

// ImageFile is a Protobuf file within an image.
type ImageFile interface {
	bufmoduleref.FileInfo
	// Proto is the backing *descriptorpb.FileDescriptorProto for this File.
	//
	// FileDescriptor should be preferred to Proto. We keep this method around
	// because we have code that does modification to the ImageFile via this.
	//
	// This will never be nil.
	// The value Path() is equal to Proto.GetName() .
	Proto() *descriptorpb.FileDescriptorProto
	// FileDescriptor is the backing FileDescriptor for this File.
	//
	// This will never be nil.
	// The value Path() is equal to FileDescriptor.GetName() .
	FileDescriptor() protodescriptor.FileDescriptor
	// IsSyntaxUnspecified will be true if the syntax was not explicitly specified.
	IsSyntaxUnspecified() bool
	// UnusedDependencyIndexes returns the indexes of the unused dependencies within
	// FileDescriptor.GetDependency().
	//
	// All indexes will be valid.
	// Will return nil if empty.
	UnusedDependencyIndexes() []int32

	withIsImport(isImport bool) ImageFile
	isImageFile()
}

// NewImageFile returns a new ImageFile.
//
// If externalPath is empty, path is used.
//
// TODO: moduleIdentity and commit should be options since they are optional.
func NewImageFile(
	fileDescriptor protodescriptor.FileDescriptor,
	moduleIdentity bufmoduleref.ModuleIdentity,
	commit string,
	externalPath string,
	isImport bool,
	isSyntaxUnspecified bool,
	unusedDependencyIndexes []int32,
) (ImageFile, error) {
	return newImageFile(
		fileDescriptor,
		moduleIdentity,
		commit,
		externalPath,
		isImport,
		isSyntaxUnspecified,
		unusedDependencyIndexes,
	)
}

// Image is a buf image.
type Image interface {
	// Files are the files that comprise the image.
	//
	// This contains all files, including imports if available.
	// The returned files are in correct DAG order.
	Files() []ImageFile
	// GetFile gets the file for the root relative file path.
	//
	// If the file does not exist, nil is returned.
	// The path is expected to be normalized and validated.
	// Note that all values of GetDependency() can be used here.
	GetFile(path string) ImageFile
	isImage()
}

// NewImage returns a new Image for the given ImageFiles.
//
// The input ImageFiles are expected to be in correct DAG order!
// TODO: Consider checking the above, and if not, reordering the Files.
// If imageFiles is empty, returns error
func NewImage(imageFiles []ImageFile) (Image, error) {
	return newImage(imageFiles, false)
}

// NewMultiImage returns a new Image for the given Images.
//
// Reorders the ImageFiles to be in DAG order.
// Duplicates cannot exist across the Images.
func NewMultiImage(images ...Image) (Image, error) {
	switch len(images) {
	case 0:
		return nil, nil
	case 1:
		return images[0], nil
	default:
		var imageFiles []ImageFile
		for _, image := range images {
			imageFiles = append(imageFiles, image.Files()...)
		}
		return newImage(imageFiles, true)
	}
}

// MergeImages returns a new Image for the given Images. ImageFiles
// treated as non-imports in at least one of the given Images will
// be treated as non-imports in the returned Image. The first non-import
// version of a file will be used in the result.
//
// Reorders the ImageFiles to be in DAG order.
// Duplicates can exist across the Images, but only if duplicates are non-imports.
func MergeImages(images ...Image) (Image, error) {
	switch len(images) {
	case 0:
		return nil, nil
	case 1:
		return images[0], nil
	default:
		var paths []string
		imageFileSet := make(map[string]ImageFile)
		for _, image := range images {
			for _, currentImageFile := range image.Files() {
				storedImageFile, ok := imageFileSet[currentImageFile.Path()]
				if !ok {
					imageFileSet[currentImageFile.Path()] = currentImageFile
					paths = append(paths, currentImageFile.Path())
					continue
				}
				if !storedImageFile.IsImport() && !currentImageFile.IsImport() {
					return nil, fmt.Errorf("%s is a non-import in multiple images", currentImageFile.Path())
				}
				if storedImageFile.IsImport() && !currentImageFile.IsImport() {
					imageFileSet[currentImageFile.Path()] = currentImageFile
				}
			}
		}
		// We need to preserve order for deterministic results, so we add
		// the files in the order they're given, but base our selection
		// on the imageFileSet.
		imageFiles := make([]ImageFile, 0, len(imageFileSet))
		for _, path := range paths {
			imageFiles = append(imageFiles, imageFileSet[path] /* Guaranteed to exist */)
		}
		return newImage(imageFiles, true)
	}
}

// NewImageForProto returns a new Image for the given proto Image.
//
// The input Files are expected to be in correct DAG order!
// TODO: Consider checking the above, and if not, reordering the Files.
//
// TODO: do we want to add the ability to do external path resolution here?
func NewImageForProto(protoImage *imagev1.Image, options ...NewImageForProtoOption) (Image, error) {
	var newImageOptions newImageForProtoOptions
	for _, option := range options {
		option(&newImageOptions)
	}
	if !newImageOptions.noReparse {
		if err := reparseImageProto(protoImage); err != nil {
			return nil, err
		}
	}
	if err := validateProtoImage(protoImage); err != nil {
		return nil, err
	}
	imageFiles := make([]ImageFile, len(protoImage.File))
	for i, protoImageFile := range protoImage.File {
		var isImport bool
		var isSyntaxUnspecified bool
		var unusedDependencyIndexes []int32
		var moduleIdentity bufmoduleref.ModuleIdentity
		var commit string
		var err error
		if protoImageFileExtension := protoImageFile.GetBufExtension(); protoImageFileExtension != nil {
			isImport = protoImageFileExtension.GetIsImport()
			isSyntaxUnspecified = protoImageFileExtension.GetIsSyntaxUnspecified()
			unusedDependencyIndexes = protoImageFileExtension.GetUnusedDependency()
			if protoModuleInfo := protoImageFileExtension.GetModuleInfo(); protoModuleInfo != nil {
				if protoModuleName := protoModuleInfo.GetName(); protoModuleName != nil {
					moduleIdentity, err = bufmoduleref.NewModuleIdentity(
						protoModuleName.GetRemote(),
						protoModuleName.GetOwner(),
						protoModuleName.GetRepository(),
					)
					if err != nil {
						return nil, err
					}
					// we only want to set this if there is a module name
					commit = protoModuleInfo.GetCommit()
				}
			}
		}
		imageFile, err := NewImageFile(
			protoImageFile,
			moduleIdentity,
			commit,
			protoImageFile.GetName(),
			isImport,
			isSyntaxUnspecified,
			unusedDependencyIndexes,
		)
		if err != nil {
			return nil, err
		}
		imageFiles[i] = imageFile
	}
	return NewImage(imageFiles)
}

// NewImageForCodeGeneratorRequest returns a new Image from a given CodeGeneratorRequest.
//
// The input Files are expected to be in correct DAG order!
// TODO: Consider checking the above, and if not, reordering the Files.
func NewImageForCodeGeneratorRequest(request *pluginpb.CodeGeneratorRequest, options ...NewImageForProtoOption) (Image, error) {
	if err := protodescriptor.ValidateCodeGeneratorRequestExceptFileDescriptorProtos(request); err != nil {
		return nil, err
	}
	protoImageFiles := make([]*imagev1.ImageFile, len(request.GetProtoFile()))
	for i, fileDescriptorProto := range request.GetProtoFile() {
		// we filter whether something is an import or not in ImageWithOnlyPaths
		// we cannot determine if the syntax was unset
		protoImageFiles[i] = fileDescriptorProtoToProtoImageFile(fileDescriptorProto, false, false, nil, nil, "")
	}
	image, err := NewImageForProto(
		&imagev1.Image{
			File: protoImageFiles,
		},
		options...,
	)
	if err != nil {
		return nil, err
	}
	return ImageWithOnlyPaths(