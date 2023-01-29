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

package manifest

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/sha3"
)

// DigestType is the type for digests in this package.
type DigestType string

const (
	DigestTypeShake256 DigestType = "shake256"

	shake256Length = 64
)

// Digest represents a hash function's value.
type Digest struct {
	dtype  DigestType
	digest []byte
	hexstr string
}

// NewDigestFromBytes builds a digest from a type and the digest bytes.
func NewDigestFromBytes(dtype DigestType, digest []byte) (*Digest, error) {
	if dtype == "" {
		return nil, errors.New("digest type cannot be empty")
	}
	if dtype != DigestTypeShake256 {
		return nil, fmt.Errorf("unsupported digest type: %q", dtype)
	}
	if len(digest) != shake256