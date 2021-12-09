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

package bufimagemodify

import (
	"context"
	"strings"

	"github.com/bufbuild/buf/private/bufpkg/bufimage"
	"github.com/bufbuild/buf/private/pkg/stringutil"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// PhpNamespaceID is the ID of the php_namespace modifier.
const PhpNamespaceID = "PHP_NAMESPACE"

var (
	// phpNamespacePath is the SourceCodeInfo path for the php_namespace option.
	// Ref: https://github.com/protocolbuffers/protobuf/blob/61689226c0e3ec88287eaed66164614d9c4f2bf7/src/google/protobuf/descriptor.proto#L443
	phpNamespacePath = []int32{8, 41}

	// Keywords and classes that could be produced by our heuristic.
	// They must not be used in a php_namespace.
	// Ref: https://www.php.net/manual/en/reserved.php
	phpReservedKeywords = map[string]struct{}{
		// Reserved classes as per above.
		"directory":           {},
		"exception":           {},
		"errorexception":      {},
		"closure":   