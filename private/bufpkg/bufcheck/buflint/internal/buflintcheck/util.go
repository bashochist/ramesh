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

package buflintcheck

import (
	"strings"

	"github.com/bufbuild/buf/private/bufpkg/bufanalysis"
	"github.com/bufbuild/buf/private/bufpkg/bufcheck/internal"
	"github.com/bufbuild/buf/private/pkg/protosource"
	"github.com/bufbuild/buf/private/pkg/stringutil"
)

// addFunc adds a FileAnnotation.
//
// Both the Descriptor and Locations can be nil.
type addFunc func(protosource.Descriptor, protosource.Location, []protosource.Location, string, ...interface{})

func fieldToLowerSnakeCase(s string) string {
	// Try running this on googleapis and watch
	// We allow both effectively by not passing the option
	//return stringutil.ToLowerSnakeCase(s, stringutil.SnakeCaseWithNewWordOnDigits())
	return stringutil.ToLowerSnakeCase(s)
}

func fieldToUpperSnakeCase(s string) string {
	// Try running this on googleapis and watch
	// We allow both effectively by not passing the option
	//return stringutil.ToUpperSnakeCase(s, stringutil.SnakeCaseWithNewWordOnDigits())
	return stringutil.ToUpperSnakeCase(s)
}

// validLeadingComment returns true if comment has at least one line that isn't empty
// and doesn't start with CommentIgnorePrefix.
func validLeadingComment(comment string) bool {
	for _, line := range strings.Split(comment, "\n") {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, CommentIgnorePrefix) {
			return true
		}
	}
	return false
}

// Returns the usedPackageList if there is an import cycle.
//
// Note this stops on the first import cycle detected, it doesn't attempt to get all of them - not perfect.
func getImportCycleIfExists(
	// Should never be ""
	pkg string,
	packageToDirectlyImportedPackageToFileImports map[string]map[string][]protosource.FileImport,
	usedPackageMap map[string]struct{},
	usedPackageList []string,
) []string {
	// Append before checking so that the returned import cycle is actually a cycle
	usedPackageList = append(usedPackageList, pkg)
	if _, ok := usedPackageMap[pkg]; ok {
		// We have an import cycle, but if the first package in the list does not
		// equal the last, do not return as an import cycle unless the first
		// element equals the last - we do DFS from each package so this will
		// be picked up separately
		if usedPackageList[0] == usedPackageList[len(usedPackageList)-1] {
			return usedPackageList
		}
		return nil
	}
	usedPackageMap[pkg] = struct{}{}
	//