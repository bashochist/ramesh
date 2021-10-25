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

package bufanalysis

import (
	"crypto/sha256"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

const (
	// FormatText is the text format for FileAnnotations.
	FormatText Format = iota + 1
	// FormatJSON is the JSON format for FileAnnotations.
	FormatJSON
	// FormatMSVS is the MSVS format for FileAnnotations.
	FormatMSVS
	// FormatJUnit is the JUnit format for FileAnnotations.
	FormatJUnit
)

var (
	// AllFormatStrings is all format strings without aliases.
	//
	// Sorted in the order we want to display them.
	AllFormatStrings = []string{
		"text",
		"json",
		"msvs",
		"junit",
	}
	// AllFormatStringsWithAliases is all format strings with aliases.
	//
	// Sorted in the order we want to display them.
	AllFormatStringsWithAliases = []string{
		"text",
		"gcc",
		"json",
		"msvs",
		"junit",
	}

	stringToFormat = map[string]Format{
		"text": FormatText,
		// alias for text
		"gcc":   FormatText,
		"json":  FormatJSON,
		"msvs":  FormatMSVS,
		"junit": FormatJUnit,
	}
	formatToString = map[Format]string{
		FormatText:  "text",
		FormatJSON:  "json",
		FormatMSVS:  "msvs",
		FormatJUnit: "junit",
	}
)

// Format is a FileAnnotation format.
type Format int

// String implements fmt.Stringer.
func (f Format) String() string {
	s, ok := formatToString[f]
	if !ok {
		return strconv.Itoa(int(f))
	}
	return s
}

// ParseFormat parses the Format.
//
// The empty strings defaults to FormatText.
func ParseFormat(s string) (Format, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	if s == "" {
		return FormatText, nil
	}
	f, ok := stringToFormat[s]
	if ok {
		return f, nil
	}
	return 0, fmt.Errorf("unknown format: %q", s)
}

// FileInfo is a minimal FileInfo interfac