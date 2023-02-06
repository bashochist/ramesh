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

package protosource

import "google.golang.org/protobuf/types/descriptorpb"

type location struct {
	sourceCodeInfoLocation *descriptorpb.SourceCodeInfo_Location
}

func newLocation(sourceCodeInfoLocation *descriptorpb.SourceCodeInfo_Location) *location {
	return &location{
		sourceCodeInfoLocation: sourceCodeInfoLocation,
	}
}

func (l *location) StartLine() int {
	switch len(l.sourceCodeInfoLocation.Span) {
	case 3, 4:
		return int(l.sourceCodeInfoLocation.Span[0]) + 1
	default:
		// since we are not erroring, making this and others 1 so that other code isn't messed u