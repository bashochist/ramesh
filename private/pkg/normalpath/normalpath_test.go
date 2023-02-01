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

package normalpath

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChunkByDir(t *testing.T) {
	t.Parallel()
	testChunkByDir(
		t,
		nil,
		0,
	)
	testChunkByDir(
		t,
		nil,
		5,
	)
	testChunkByDir(
		t,
		[]string{},
		0,
	)
	testChunkByDir(
		t,
		[]string{},
		5,
	)
	testChunkByDir(
		t,
		[]string{"a/a.proto"},
		1,
		[]string{"a/a.proto"},
	)
	testChunkByDir(
		t,
		[]string{"a/a.proto"},
		2,
		[]string{"a/a.proto"},
	)
	testChunkByDir(
		t,
		[]stri