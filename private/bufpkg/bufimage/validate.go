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
	"errors"
	"fmt"

	imagev1 "github.com/bufbuild/buf/private/gen/proto/go/buf/alpha/image/v1"
)

// we validate the actual fields of the FileDescriptorProtos as part of newImageFile
func validateProtoImage(protoImage *imagev1.Image) error {
	if protoImage == nil {
		return errors.New("nil Image")
	}
	if len(protoImage.File) == 0 {
		return errors.New("image contains no files")
	}
	for _, protoImageFile := range protoImage.File {
		if err := validateProtoImageFile(protoImageFile); err != nil {
			return err
		}
	}
	return nil
}

func validateProtoImageFile(protoImageFile *imagev1.ImageFile) error {
	if protoImageFileExtension := protoImageFile.GetBufExtension(); protoImageFileExtension != nil {
		lenDependencies := len(protoImageFile.GetDependency())
		fo