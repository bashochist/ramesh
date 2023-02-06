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

import (
	"sync"

	"google.golang.org/protobuf/types/descriptorpb"
)

type locationStore struct {
	sourceCodeInfoLocations []*descriptorpb.SourceCodeInfo_Location

	pathToLocation               map[string]Location
	pathToSourceCodeInfoLocation map[string]*descriptorpb.SourceCodeInfo_Location
	locationLock                 sync.RWMutex
	sourceCodeInfoLocationLock   sync.RWMutex
}

func newLocationStore(sourceCodeInfoLocations []*descriptorpb.SourceCodeInfo_Location) *locationStore {
	return &locationStore{
		sourceCodeInfoLocations: sourceCodeInfoLocations,
		pathToLocation:          make(map[string]Location),
	}
}

func (l *locationStore) getLocation(path []int32) Location {
	return l.getLocationByPathKey(getPathKey(path))
}

// optimization for keys we know ahead of time such as package location, certain file options
func (l *locationStore) getLocationByPathKey(pathKey string) Location {
	// check cache first
	l.locationLock.RLock()
	location, ok := l.pathToLocation[pathKey]
	l.locationLock.RUnlock()
	if ok {
		return location
	}

	// build index and get sourceCodeInfoLocation
	l.sourceCodeInfoLocationLo