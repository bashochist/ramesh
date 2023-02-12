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

package storage

import (
	"context"

	"github.com/bufbuild/buf/private/pkg/storage/storageutil"
)

// MultiReadBucket takes the union of the ReadBuckets.
//
// If no readBuckets are given, this returns a no-op ReadBucket.
// If one readBucket is given, this returns the original ReadBucket.
// Otherwise, this returns a ReadBucket that will get from all buckets.
//
// This expects and validates that no paths overlap between the ReadBuckets.
// This assumes that buckets are logically unique.
func MultiReadBucket(readBuckets ...ReadBucket) ReadBucket {
	switch len(readBuckets) {
	case 0:
		return nopReadBucket{}
	case 1:
		return readBuckets[0]
	default:
		return newMultiReadBucket(readBuckets, false)
	}
}

// MultiRead