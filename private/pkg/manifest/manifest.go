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

// A manifest is a file containing a list of paths and their hash digests,
// canonically ordered by path in increasing lexicographical order. Manifests
// are encoded as:
//
//	<digest type>:<digest>[SP][SP]<path>[LF]
//
// "shake256" is the only supported digest type. The digest is 64 bytes of hex
// encoded output of SHAKE256. See golang.org/x/crypto/sha3 and FIPS 202 for
// details on the SHAKE hash.
//
// [Manifest] can read and write manifest files. Canonical form is produced
// when serialized ([Manifest.MarshalText]). Non-canonical form is a valid
// manifest and will not produce errors when deserializing.
//
// Interacting with a manifest is typically by path ([Manifest.Paths],
// [Manifest.DigestFor]) or by a [Digest] ([Manifest.PathsFor]).
//
// [Blob] represents file content and its digest. [BlobSet] collects related
// blobs together into a set. [NewMemoryBlob] provides an in-memory
// implementation. A manifest, being a file, is also a blob ([Manifest.Blob]).
//
// Blobs are anonymous files and a manifest gives names to anonymous files.
// It's possible to view a manifest and its associated blobs as a file system.
// [NewBucket] creates a storage bucket from a manifest and blob set.
// [NewFromBucket] does the inverse: the creation of a manifest and blob set
// from a storage bucket.
package manifest

import (
	"bufio"
	"bytes"
	"encoding"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
)

var errNoFinalNewline = errors.New("partial record: missing newline")

func newError(lineno int, msg string) error {
	return fmt.Errorf("invalid manifest: %d: %s", lineno, msg)
}

func newErrorWrapped(lineno int, err error) error {
	return fmt.Errorf("invalid manifest: %d: %w", lineno, err)
}

// Manifest represents a list of paths and their digests.
type Manifest struct {
	pathToDigest  map[string]Digest
	digestToPaths map[string][]string
}

var _ encoding.TextMarshaler = (*Manifest)(nil)
