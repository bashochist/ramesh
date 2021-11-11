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

package buflintv1

import (
	"github.com/bufbuild/buf/private/bufpkg/bufcheck/buflint/internal/buflintbuild"
	"github.com/bufbuild/buf/private/bufpkg/bufcheck/internal"
)

var (
	// v1RuleBuilders are the rule builders.
	v1RuleBuilders = []*internal.RuleBuilder{
		buflintbuild.CommentEnumRuleBuilder,
		buflintbuild.CommentEnumValueRuleBuilder,
		buflintbuild.CommentFieldRuleBuilder,
		buflintbuild.CommentMessageRuleBuilder,
		buflintbuild.CommentOneofRuleBuilder,
		buflintbuild.CommentRPCRuleBuilder,
		buflintbuild.CommentServiceRuleBuilder,
		buflintbuild.DirectorySamePackageRuleBuilder,
		buflintbuild.EnumFirstValueZeroRuleBuilder,
		buflintbuild.EnumNoAllowAliasRuleBuilder,
		buflintbuild.EnumPascalCaseRuleBuilder,
		buflintbuild.EnumValuePrefixRuleBuilder,
		buflintbuild.EnumValueUpperSnakeCaseRuleBuilder,
		buflintbuild.EnumZeroValueSuffixRuleBuilder,
		buflintbuild.FieldLowerSnakeCaseRuleBuilder,
		buflintbuild.FileLowerSnakeCaseRuleBuilder,
		buflintbuild.ImportNoPublicRuleBuilder,
		buflintbuild.ImportNoWeakRuleBuilder,
		buflintbuild.ImportUsedRuleBuilder,
		buflintbuild.MessagePascalCaseRuleBuilder,
		buflintbuild.OneofLowerSnakeCaseRuleBuilder,
		buflintbuild.PackageDefinedRuleBuilder,
		buflintbuild.PackageDirectoryMatchRuleBuilder,
		buflintbuild.PackageLowerSnakeCaseRuleBuilder,
		buflintbuild.PackageNoImportCycleRuleBuilder,
		buflintbuild.PackageSameCsharpNamespaceRuleBuilder,
		buflintbuild.PackageSameDirectoryRuleBuilder,
		buflintbuild.PackageSameGoPackageRuleBuilder,
		buflintbuild.PackageSameJavaMultipleFilesRuleBuilder,
		buflintbuild.PackageSameJavaPackageRuleBuilder,
		buflintbuild.PackageSamePhpNamespaceRuleBuilder,
		buflintbuild.PackageSameRubyPackageRuleBuilder,
		buflintbuild.PackageSameSwiftPrefixRuleBuilder,
		buflintbuild.PackageVersionSuffixRuleBuilder,
		buflintbuild.RPCNoClientStreamingRuleBuilder,
		buflintbuild.RPCNoServerStreamingRuleBuilder,
		buflintbuild.RPCPascalCaseRuleBuilder,
		buflintbuild.RPCRequestResponseUniqueRuleBuilder,
		buflintbuild.RPCRequestStandardNameRuleBuilder,
		buflintbuild.RPCResponseStandardNameRuleBuilder,
		buflintbuild.ServicePascalCaseRuleBuilder,
		buflintbuild.ServiceSuffixRuleBuilder,
		buflintbuild.SyntaxSpecifiedRuleBuilder,
	}

	// v1DefaultCategories are the default categories.
	v1DefaultCategories = []string{
		"DEFAULT",
	}
	// v1IDToCategories associates IDs to categories.
	v1IDToCategories = map[string][]string{
		"COMMENT_ENUM": {
			"COMMENTS",
		},
		"COMMENT_ENUM_VALUE": {
			"COMMENTS",
		},
		"COMMENT_FIELD": {
			"COMMENTS",
		},
		"COMMENT_MESSAGE": {
			"COMMENTS",
		},
		"COMMENT_ONEOF": {
			"COMMENTS",
		},
		"COMMENT_RPC": {
			"COMMENTS",
		},
		"COMMENT_SERVICE": {
			"COMMENTS",
		},
		"DIRECTORY_SAME_PACKAGE": {
			"MINIMAL",
			"BASIC",
			"DEFAULT",
		},
		"ENUM_FIRST_VALUE_ZERO": {
			"BASIC",
			"DEFAULT",
		},
		"ENUM_NO_ALLOW_ALIAS": {
			"BASIC