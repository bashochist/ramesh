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

	"go.uber.org/atomic"
)

// LimitWriteBucket returns a [WriteBucket] that writes to [writeBucket]
// but stops with an error after [limit] bytes are written.
//
// The error can be checked using [IsWriteLimitReached].
//
// A negative [limit] is same as 0 limit.
func LimitWriteBucket(writeBucket WriteBucket, limit int) WriteBucket {
	if limit < 0 {
		limit = 0
	}
	return newLimitedWriteBucket(writeBucket, int64(limit))
}

type limitedWriteBucket struct {
	WriteBucket
	currentSize *atomic.Int64
	limit       int64
}

func newLimitedWriteBucket(buc