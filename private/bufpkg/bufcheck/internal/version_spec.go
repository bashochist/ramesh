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

package internal

import (
	"sort"

	"github.com/bufbuild/buf/private/pkg/stringutil"
)

// VersionSpec specifies the rules, ids, and categories for a given version.
type VersionSpec struct {
	RuleBuilders      []*RuleBuilder
	DefaultCategories []string
	// May include IDs without any categories.
	// To get all categories, use AllCategoriesForVersionSpec.
	IDToCategories map[string][]string
}

// AllCategoriesForVersionSpec returns all categories for the VersionSpec.
//
// Sorted by category priority.
func AllCategoriesForVersionSpec(versionSpec *VersionSpec) []string {
	categoriesMap := make(map[string]struct{})
	for _, categories := range versionSpec.IDToCategories {
		for _, category := range categories {
			categoriesMap[category] = struct{}{}
		}
	}
	categories := stringutil.MapToSlice(categoriesMap)
	sort.Slice(
		categories,
		func(i int, j int) bool {
			return categoryLess