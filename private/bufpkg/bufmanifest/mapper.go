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

package bufmanifest

import (
	"context"
	"fmt"
	"io"

	modulev1alpha1 "github.com/bufbuild/buf/private/gen/proto/go/buf/alpha/module/v1alpha1"
	"github.com/bufbuild/buf/private/pkg/manifest"
	"go.uber.org/multierr"
)

var (
	protoDigestTypeToDigestType = map[modulev1alpha1.DigestType]manifest.DigestType{
		modulev1alpha1.DigestType_DIGEST_TYPE_SHAKE256: manifest.DigestTypeShake256,
	}
	digestTypeToProtoDigestType = map[manifest.DigestType]modulev1alpha1.DigestType{
		manifest.DigestTypeShake256: modulev1alpha1.DigestType_DIGEST_TYPE_SHAKE256,
	}
)

// NewDigestFromProtoDigest maps a modulev1alpha1.Digest to a Digest.
func NewDigestFromProtoDigest(digest *modulev1alpha1.Digest) (*manifest.Digest, error) {
	if digest == nil {
		return nil