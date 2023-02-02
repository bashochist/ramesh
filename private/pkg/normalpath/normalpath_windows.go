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

//go:build windows
// +build windows

package normalpath

import (
	"os"
	"path/filepath"
	"strings"
)

// NormalizeAndValidate normalizes and validates the given path.
//
// This calls Normalize on the path.
// Returns Error if the path is not relative or jumps context.
// This can be used to validate that paths are valid to use with Buckets.
// The error message is safe to pass to users.
func NormalizeAndValidate(path string) (string, error) {
	normalizedPath := Normalize(path)
	if filepath.IsAbs(normalizedPath) || (len(normalizedPath) > 0 && normalizedPath[0] == '/') {
		// the stdlib implementation of `IsAbs` assumes that a volume name is required for a path to
		// be absolute, however Windows treats a `/` (normalized) rooted path as absolute _within_ the current volume.
		// In the context of validating that a path is _not_ relative, we need to reject a path that begins
		// with `/`.
		return "", NewError(path, errNotRelative)
	}
	// https://github.com/bufbuild/buf/issues/51
	if strings.HasPrefix(normalizedPath, normalizedRelPathJumpContextPrefix) {
		return "", NewError(path, errOutsideContextDir)
	}
	return normalizedPath, nil
}

// EqualsOrContainsPath returns true if the value is equal to or contains the path.
// path is compared at each directory level to value for equivalency under simple unicode
// codepoint folding. This means it is context and locale independent. This matching
// will not support the few rare cases, primarily in Turkish and Lithuanian, noted
// in the caseless matching section of Unicode 13.0 https://www.unicode.org/versions/Unicode13.0.0/ch05.pdf#page=47.
//
// The path and value are expected to be normalized and validated if Relative is used.
// The path and value are expected to be normalized and absolute if Absolute is used.
func EqualsOrContainsPath(value string, path string, pathType PathType) bool {
	curPath := path
	var lastSeen string
	for {
		if strings.EqualFold(value, curPath) {
			return true
		}
		curPath = Dir(curPath)
		if lastSeen == curPath {
			break
		}
		lastSeen = curPath
	}
	return false
}

// MapHasEqualOrContainingPath returns true if the path matches any file or directory in the map.
//
// The path and keys in m are expected to be normalized and validated if Relative is used.
// The path and keys in m are expected to be normalized and absolute if Absolute is used.
//
// If the map is empty, returns false.
func MapHasEqualOrContainingPath(m map[string]struct{}, path string, pathType PathType) bool {
	if len(m) == 0 {
		return false
	}

	for value := range m {
		if EqualsOrContainsPath(value, path, pathType) {
			return true
		}
	}

	return false
}

// MapAllEqualOrContainingPathMap returns the paths in m that are equal to, or contain
// path, in a new map.
//
// The path and keys in m are expected to be normalized and validated if Relative is used.
// The path and keys in m are expected to be normalized and absolute if Absolute is used.
//
// If the map is empty, returns nil.
func MapAllEqualOrContainingPathMap(m map[string]struct{}, path string, pathType PathType) map[string]struct{} {
	if len(m) == 0 {
		return nil
	}

	n := make(map[string]struct{})

	for potentialMatch 