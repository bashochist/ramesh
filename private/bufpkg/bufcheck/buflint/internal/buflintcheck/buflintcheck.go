
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

// Package buflintcheck impelements the check functions.
//
// These are used by buflintbuild to create RuleBuilders.
package buflintcheck

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bufbuild/buf/private/bufpkg/bufanalysis"
	"github.com/bufbuild/buf/private/bufpkg/bufcheck/internal"
	"github.com/bufbuild/buf/private/pkg/normalpath"
	"github.com/bufbuild/buf/private/pkg/protosource"
	"github.com/bufbuild/buf/private/pkg/protoversion"
	"github.com/bufbuild/buf/private/pkg/stringutil"
)

const (
	// CommentIgnorePrefix is the comment ignore prefix.
	//
	// Comments with this prefix do not count towards valid comments in the comment checkers.
	// This is also used in buflint when constructing a new Runner, and is passed to the
	// RunnerWithIgnorePrefix option.
	CommentIgnorePrefix = "buf:lint:ignore"
)

var (
	// CheckCommentEnum is a check function.
	CheckCommentEnum = newEnumCheckFunc(checkCommentEnum)
	// CheckCommentEnumValue is a check function.
	CheckCommentEnumValue = newEnumValueCheckFunc(checkCommentEnumValue)
	// CheckCommentField is a check function.
	CheckCommentField = newFieldCheckFunc(checkCommentField)
	// CheckCommentMessage is a check function.
	CheckCommentMessage = newMessageCheckFunc(checkCommentMessage)
	// CheckCommentOneof is a check function.
	CheckCommentOneof = newOneofCheckFunc(checkCommentOneof)
	// CheckCommentService is a check function.
	CheckCommentService = newServiceCheckFunc(checkCommentService)
	// CheckCommentRPC is a check function.
	CheckCommentRPC = newMethodCheckFunc(checkCommentRPC)
)

func checkCommentEnum(add addFunc, value protosource.Enum) error {
	return checkCommentNamedDescriptor(add, value, "Enum")
}

func checkCommentEnumValue(add addFunc, value protosource.EnumValue) error {
	return checkCommentNamedDescriptor(add, value, "Enum value")
}

func checkCommentField(add addFunc, value protosource.Field) error {
	return checkCommentNamedDescriptor(add, value, "Field")
}

func checkCommentMessage(add addFunc, value protosource.Message) error {
	return checkCommentNamedDescriptor(add, value, "Message")
}

func checkCommentOneof(add addFunc, value protosource.Oneof) error {
	return checkCommentNamedDescriptor(add, value, "Oneof")
}

func checkCommentRPC(add addFunc, value protosource.Method) error {
	return checkCommentNamedDescriptor(add, value, "RPC")
}

func checkCommentService(add addFunc, value protosource.Service) error {
	return checkCommentNamedDescriptor(add, value, "Service")
}

func checkCommentNamedDescriptor(
	add addFunc,
	namedDescriptor protosource.NamedDescriptor,
	typeName string,
) error {
	location := namedDescriptor.Location()
	if location == nil {
		// this will magically skip map entry fields as well as a side-effect, although originally unintended
		return nil
	}
	if !validLeadingComment(location.LeadingComments()) {
		add(namedDescriptor, location, nil, "%s %q should have a non-empty comment for documentation.", typeName, namedDescriptor.Name())
	}
	return nil
}

// CheckDirectorySamePackage is a check function.
var CheckDirectorySamePackage = newDirToFilesCheckFunc(checkDirectorySamePackage)

func checkDirectorySamePackage(add addFunc, dirPath string, files []protosource.File) error {
	pkgMap := make(map[string]struct{})
	for _, file := range files {
		// works for no package set as this will result in "" which is a valid map key
		pkgMap[file.Package()] = struct{}{}
	}
	if len(pkgMap) > 1 {
		var messagePrefix string
		if _, ok := pkgMap[""]; ok {
			delete(pkgMap, "")
			if len(pkgMap) > 1 {
				messagePrefix = fmt.Sprintf("Multiple packages %q and file with no package", strings.Join(stringutil.MapToSortedSlice(pkgMap), ","))
			} else {
				// Join works with only one element as well by adding no comma
				messagePrefix = fmt.Sprintf("Package %q and file with no package", strings.Join(stringutil.MapToSortedSlice(pkgMap), ","))
			}
		} else {
			messagePrefix = fmt.Sprintf("Multiple packages %q", strings.Join(stringutil.MapToSortedSlice(pkgMap), ","))
		}
		for _, file := range files {
			add(file, file.PackageLocation(), nil, "%s detected within directory %q.", messagePrefix, dirPath)
		}
	}
	return nil
}

// CheckEnumNoAllowAlias is a check function.
var CheckEnumNoAllowAlias = newEnumCheckFunc(checkEnumNoAllowAlias)

func checkEnumNoAllowAlias(add addFunc, enum protosource.Enum) error {
	if enum.AllowAlias() {
		add(enum, enum.AllowAliasLocation(), nil, `Enum option "allow_alias" on enum %q must be false.`, enum.Name())
	}
	return nil
}

// CheckEnumPascalCase is a check function.
var CheckEnumPascalCase = newEnumCheckFunc(checkEnumPascalCase)

func checkEnumPascalCase(add addFunc, enum protosource.Enum) error {
	name := enum.Name()
	expectedName := stringutil.ToPascalCase(name)
	if name != expectedName {
		add(enum, enum.NameLocation(), nil, "Enum name %q should be PascalCase, such as %q.", name, expectedName)
	}
	return nil
}

// CheckEnumFirstValueZero is a check function.
var CheckEnumFirstValueZero = newEnumCheckFunc(checkEnumFirstValueZero)

func checkEnumFirstValueZero(add addFunc, enum protosource.Enum) error {
	if values := enum.Values(); len(values) > 0 {
		if firstEnumValue := values[0]; firstEnumValue.Number() != 0 {
			// proto3 compilation references the number
			add(
				firstEnumValue,
				firstEnumValue.NumberLocation(),
				// also check the name location for this comment ignore, as the number location might not have the comment
				// see https://github.com/bufbuild/buf/issues/1186
				// also check the enum for this comment ignore
				// this allows users to set this "globally" for an enum
				// see https://github.com/bufbuild/buf/issues/161
				[]protosource.Location{
					firstEnumValue.NameLocation(),
					firstEnumValue.Enum().Location(),
				},
				"First enum value %q should have a numeric value of 0",
				firstEnumValue.Name(),
			)
		}
	}
	return nil
}

// CheckEnumValuePrefix is a check function.
var CheckEnumValuePrefix = newEnumValueCheckFunc(checkEnumValuePrefix)

func checkEnumValuePrefix(add addFunc, enumValue protosource.EnumValue) error {
	name := enumValue.Name()
	expectedPrefix := fieldToUpperSnakeCase(enumValue.Enum().Name()) + "_"
	if !strings.HasPrefix(name, expectedPrefix) {
		add(
			enumValue,
			enumValue.NameLocation(),
			// also check the enum for this comment ignore
			// this allows users to set this "globally" for an enum
			// this came up in https://github.com/bufbuild/buf/issues/161
			[]protosource.Location{
				enumValue.Enum().Location(),
			},
			"Enum value name %q should be prefixed with %q.",
			name,
			expectedPrefix,
		)
	}
	return nil
}

// CheckEnumValueUpperSnakeCase is a check function.
var CheckEnumValueUpperSnakeCase = newEnumValueCheckFunc(checkEnumValueUpperSnakeCase)

func checkEnumValueUpperSnakeCase(add addFunc, enumValue protosource.EnumValue) error {
	name := enumValue.Name()
	expectedName := fieldToUpperSnakeCase(name)
	if name != expectedName {
		add(
			enumValue,
			enumValue.NameLocation(),
			// also check the enum for this comment ignore
			// this allows users to set this "globally" for an enum
			[]protosource.Location{
				enumValue.Enum().Location(),
			},
			"Enum value name %q should be UPPER_SNAKE_CASE, such as %q.",
			name,
			expectedName,
		)
	}
	return nil
}

// CheckEnumZeroValueSuffix is a check function.
var CheckEnumZeroValueSuffix = func(
	id string,
	ignoreFunc internal.IgnoreFunc,
	files []protosource.File,
	suffix string,
) ([]bufanalysis.FileAnnotation, error) {
	return newEnumValueCheckFunc(
		func(add addFunc, enumValue protosource.EnumValue) error {
			return checkEnumZeroValueSuffix(add, enumValue, suffix)
		},
	)(id, ignoreFunc, files)
}

func checkEnumZeroValueSuffix(add addFunc, enumValue protosource.EnumValue, suffix string) error {
	if enumValue.Number() != 0 {
		return nil
	}
	name := enumValue.Name()
	if !strings.HasSuffix(name, suffix) {
		add(
			enumValue,
			enumValue.NameLocation(),
			// also check the enum for this comment ignore
			// this allows users to set this "globally" for an enum
			[]protosource.Location{
				enumValue.Enum().Location(),
			},
			"Enum zero value name %q should be suffixed with %q.",
			name,
			suffix,
		)
	}
	return nil
}

// CheckFieldLowerSnakeCase is a check function.
var CheckFieldLowerSnakeCase = newFieldCheckFunc(checkFieldLowerSnakeCase)

func checkFieldLowerSnakeCase(add addFunc, field protosource.Field) error {
	message := field.Message()
	if message == nil {
		// just a sanity check
		return errors.New("field.Message() was nil")
	}
	if message.IsMapEntry() {
		// this check should always pass anyways but just in case
		return nil
	}
	name := field.Name()
	expectedName := fieldToLowerSnakeCase(name)
	if name != expectedName {
		add(
			field,
			field.NameLocation(),
			// also check the message for this comment ignore
			// this allows users to set this "globally" for a message
			[]protosource.Location{
				field.Message().Location(),
			},
			"Field name %q should be lower_snake_case, such as %q.",
			name,
			expectedName,
		)
	}
	return nil
}

// CheckFieldNoDescriptor is a check function.
var CheckFieldNoDescriptor = newFieldCheckFunc(checkFieldNoDescriptor)

func checkFieldNoDescriptor(add addFunc, field protosource.Field) error {
	name := field.Name()
	if strings.ToLower(strings.Trim(name, "_")) == "descriptor" {
		add(
			field,
			field.NameLocation(),
			// also check the message for this comment ignore
			// this allows users to set this "globally" for a message
			[]protosource.Location{
				field.Message().Location(),
			},
			`Field name %q cannot be any capitalization of "descriptor" with any number of prefix or suffix underscores.`,
			name,
		)
	}
	return nil
}

// CheckFileLowerSnakeCase is a check function.
var CheckFileLowerSnakeCase = newFileCheckFunc(checkFileLowerSnakeCase)

func checkFileLowerSnakeCase(add addFunc, file protosource.File) error {
	filename := file.Path()
	base := normalpath.Base(filename)
	ext := normalpath.Ext(filename)
	baseWithoutExt := strings.TrimSuffix(base, ext)
	expectedBaseWithoutExt := stringutil.ToLowerSnakeCase(baseWithoutExt)
	if baseWithoutExt != expectedBaseWithoutExt {
		add(file, nil, nil, `Filename %q should be lower_snake_case%s, such as "%s%s".`, base, ext, expectedBaseWithoutExt, ext)
	}
	return nil
}

var (
	// CheckImportNoPublic is a check function.
	CheckImportNoPublic = newFileImportCheckFunc(checkImportNoPublic)
	// CheckImportNoWeak is a check function.
	CheckImportNoWeak = newFileImportCheckFunc(checkImportNoWeak)
	// CheckImportUsed is a check function.
	CheckImportUsed = newFileImportCheckFunc(checkImportUsed)
)

func checkImportNoPublic(add addFunc, fileImport protosource.FileImport) error {
	return checkImportNoPublicWeak(add, fileImport, fileImport.IsPublic(), "public")
}

func checkImportNoWeak(add addFunc, fileImport protosource.FileImport) error {
	return checkImportNoPublicWeak(add, fileImport, fileImport.IsWeak(), "weak")
}

func checkImportNoPublicWeak(add addFunc, fileImport protosource.FileImport, value bool, name string) error {
	if value {
		add(fileImport, fileImport.Location(), nil, `Import %q must not be %s.`, fileImport.Import(), name)
	}
	return nil
}

func checkImportUsed(add addFunc, fileImport protosource.FileImport) error {
	if fileImport.IsUnused() {
		add(fileImport, fileImport.Location(), nil, `Import %q is unused.`, fileImport.Import())
	}
	return nil
}

// CheckMessagePascalCase is a check function.
var CheckMessagePascalCase = newMessageCheckFunc(checkMessagePascalCase)

func checkMessagePascalCase(add addFunc, message protosource.Message) error {
	if message.IsMapEntry() {
		// map entries should always be pascal case but we don't want to check them anyways
		return nil
	}
	name := message.Name()
	expectedName := stringutil.ToPascalCase(name)
	if name != expectedName {
		add(message, message.NameLocation(), nil, "Message name %q should be PascalCase, such as %q.", name, expectedName)
	}
	return nil
}

// CheckOneofLowerSnakeCase is a check function.
var CheckOneofLowerSnakeCase = newOneofCheckFunc(checkOneofLowerSnakeCase)

func checkOneofLowerSnakeCase(add addFunc, oneof protosource.Oneof) error {
	name := oneof.Name()
	expectedName := fieldToLowerSnakeCase(name)
	if name != expectedName {
		// if this is an implicit oneof for a proto3 optional field, do not error
		// https://github.com/protocolbuffers/protobuf/blob/master/docs/implementing_proto3_presence.md
		if fields := oneof.Fields(); len(fields) == 1 {
			if fields[0].Proto3Optional() {
				return nil
			}
		}
		add(
			oneof,
			oneof.NameLocation(),
			// also check the message for this comment ignore
			// this allows users to set this "globally" for a message
			[]protosource.Location{
				oneof.Message().Location(),
			},
			"Oneof name %q should be lower_snake_case, such as %q.",
			name,
			expectedName,
		)
	}
	return nil
}

// CheckPackageDefined is a check function.
var CheckPackageDefined = newFileCheckFunc(checkPackageDefined)

func checkPackageDefined(add addFunc, file protosource.File) error {
	if file.Package() == "" {
		add(file, nil, nil, "Files must have a package defined.")
	}
	return nil
}

// CheckPackageDirectoryMatch is a check function.
var CheckPackageDirectoryMatch = newFileCheckFunc(checkPackageDirectoryMatch)

func checkPackageDirectoryMatch(add addFunc, file protosource.File) error {
	pkg := file.Package()
	if pkg == "" {
		return nil
	}
	expectedDirPath := strings.ReplaceAll(pkg, ".", "/")
	dirPath := normalpath.Dir(file.Path())
	// need to check case where in root relative directory and no package defined
	// this should be valid although if SENSIBLE is turned on this will be invalid
	if dirPath != expectedDirPath {
		add(file, file.PackageLocation(), nil, `Files with package %q must be within a directory "%s" relative to root but were in directory "%s".`, pkg, normalpath.Unnormalize(expectedDirPath), normalpath.Unnormalize(dirPath))
	}
	return nil
}

// CheckPackageLowerSnakeCase is a check function.
var CheckPackageLowerSnakeCase = newFileCheckFunc(checkPackageLowerSnakeCase)

func checkPackageLowerSnakeCase(add addFunc, file protosource.File) error {
	pkg := file.Package()
	if pkg == "" {
		return nil
	}
	split := strings.Split(pkg, ".")
	for i, elem := range split {
		split[i] = stringutil.ToLowerSnakeCase(elem)
	}
	expectedPkg := strings.Join(split, ".")
	if pkg != expectedPkg {
		add(file, file.PackageLocation(), nil, "Package name %q should be lower_snake.case, such as %q.", pkg, expectedPkg)
	}
	return nil
}

// CheckPackageNoImportCycle is a check function.
var CheckPackageNoImportCycle = newFilesCheckFunc(checkPackageNoImportCycle)

func checkPackageNoImportCycle(add addFunc, files []protosource.File) error {
	packageToDirectlyImportedPackageToFileImports, err := protosource.PackageToDirectlyImportedPackageToFileImports(files...)
	if err != nil {
		return err
	}
	// This is way more algorithmically complex than it needs to be.
	//
	// We're doing a DFS starting at each package. What we should do is start from any package,
	// do the DFS and keep track of the packages hit, and then don't ever do DFS from a given
	// package twice. The problem is is that with the current janky package -> direct -> file imports
	// setup, we would then end up with error messages like "import cycle: a -> b -> c -> b", and
	// attach the error message to a file with package a, and we want to just print "b -> c -> b".
	// So to get this to market, we just do a DFS from each package.
	//
	// This may prove to be too expensive but early testing say it is not so far.
	for pkg := range packageToDirectlyImportedPackageToFileImports {
		// Can equal "" per the function signature of PackageToDirectlyImportedPackageToFileImports
		if pkg == "" {
			continue
		}
		// Go one deep in the potential import cycle so that we can get the file imports
		// we want to potentially attach errors to.
		//
		// We know that pkg is never equal to directlyImportedPackage due to the signature
		// of PackageToDirectlyImportedPackageToFileImports.
		for directlyImportedPackage, fileImports := range packageToDirectlyImportedPackageToFileImports[pkg] {
			// Can equal "" per the function signature of PackageToDirectlyImportedPackageToFileImports
			if directlyImportedPackage == "" {
				continue
			}
			if importCycle := getImportCycleIfExists(
				directlyImportedPackage,
				packageToDirectlyImportedPackageToFileImports,
				map[string]struct{}{
					pkg: {},
				},
				[]string{
					pkg,
				},
			); len(importCycle) > 0 {
				for _, fileImport := range fileImports {
					add(fileImport, fileImport.Location(), nil, `Package import cycle: %s`, strings.Join(importCycle, ` -> `))
				}
			}
		}
	}
	return nil
}

// CheckPackageSameDirectory is a check function.
var CheckPackageSameDirectory = newPackageToFilesCheckFunc(checkPackageSameDirectory)

func checkPackageSameDirectory(add addFunc, pkg string, files []protosource.File) error {
	dirMap := make(map[string]struct{})
	for _, file := range files {
		dirMap[normalpath.Dir(file.Path())] = struct{}{}
	}
	if len(dirMap) > 1 {
		dirs := stringutil.MapToSortedSlice(dirMap)
		for _, file := range files {
			add(file, file.PackageLocation(), nil, "Multiple directories %q contain files with package %q.", strings.Join(dirs, ","), pkg)
		}
	}
	return nil
}

var (
	// CheckPackageSameCsharpNamespace is a check function.
	CheckPackageSameCsharpNamespace = newPackageToFilesCheckFunc(checkPackageSameCsharpNamespace)
	// CheckPackageSameGoPackage is a check function.
	CheckPackageSameGoPackage = newPackageToFilesCheckFunc(checkPackageSameGoPackage)
	// CheckPackageSameJavaMultipleFiles is a check function.
	CheckPackageSameJavaMultipleFiles = newPackageToFilesCheckFunc(checkPackageSameJavaMultipleFiles)
	// CheckPackageSameJavaPackage is a check function.
	CheckPackageSameJavaPackage = newPackageToFilesCheckFunc(checkPackageSameJavaPackage)
	// CheckPackageSamePhpNamespace is a check function.
	CheckPackageSamePhpNamespace = newPackageToFilesCheckFunc(checkPackageSamePhpNamespace)
	// CheckPackageSameRubyPackage is a check function.
	CheckPackageSameRubyPackage = newPackageToFilesCheckFunc(checkPackageSameRubyPackage)
	// CheckPackageSameSwiftPrefix is a check function.
	CheckPackageSameSwiftPrefix = newPackageToFilesCheckFunc(checkPackageSameSwiftPrefix)
)

func checkPackageSameCsharpNamespace(add addFunc, pkg string, files []protosource.File) error {
	return checkPackageSameOptionValue(add, pkg, files, protosource.File.CsharpNamespace, protosource.File.CsharpNamespaceLocation, "csharp_namespace")
}

func checkPackageSameGoPackage(add addFunc, pkg string, files []protosource.File) error {
	return checkPackageSameOptionValue(add, pkg, files, protosource.File.GoPackage, protosource.File.GoPackageLocation, "go_package")
}

func checkPackageSameJavaMultipleFiles(add addFunc, pkg string, files []protosource.File) error {
	return checkPackageSameOptionValue(
		add,
		pkg,
		files,
		func(file protosource.File) string {
			return strconv.FormatBool(file.JavaMultipleFiles())
		},
		protosource.File.JavaMultipleFilesLocation,
		"java_multiple_files",
	)
}

func checkPackageSameJavaPackage(add addFunc, pkg string, files []protosource.File) error {
	return checkPackageSameOptionValue(add, pkg, files, protosource.File.JavaPackage, protosource.File.JavaPackageLocation, "java_package")
}

func checkPackageSamePhpNamespace(add addFunc, pkg string, files []protosource.File) error {
	return checkPackageSameOptionValue(add, pkg, files, protosource.File.PhpNamespace, protosource.File.PhpNamespaceLocation, "php_namespace")
}

func checkPackageSameRubyPackage(add addFunc, pkg string, files []protosource.File) error {
	return checkPackageSameOptionValue(add, pkg, files, protosource.File.RubyPackage, protosource.File.RubyPackageLocation, "ruby_package")
}

func checkPackageSameSwiftPrefix(add addFunc, pkg string, files []protosource.File) error {
	return checkPackageSameOptionValue(add, pkg, files, protosource.File.SwiftPrefix, protosource.File.SwiftPrefixLocation, "swift_prefix")
}

func checkPackageSameOptionValue(
	add addFunc,
	pkg string,
	files []protosource.File,
	getOptionValue func(protosource.File) string,
	getOptionLocation func(protosource.File) protosource.Location,
	name string,
) error {
	optionValueMap := make(map[string]struct{})
	for _, file := range files {
		optionValueMap[getOptionValue(file)] = struct{}{}
	}
	if len(optionValueMap) > 1 {
		_, noOptionValue := optionValueMap[""]
		delete(optionValueMap, "")
		optionValues := stringutil.MapToSortedSlice(optionValueMap)
		for _, file := range files {
			if noOptionValue {
				add(file, getOptionLocation(file), nil, "Files in package %q have both values %q and no value for option %q and all values must be equal.", pkg, strings.Join(optionValues, ","), name)
			} else {
				add(file, getOptionLocation(file), nil, "Files in package %q have multiple values %q for option %q and all values must be equal.", pkg, strings.Join(optionValues, ","), name)
			}
		}
	}
	return nil
}

// CheckPackageVersionSuffix is a check function.
var CheckPackageVersionSuffix = newFileCheckFunc(checkPackageVersionSuffix)

func checkPackageVersionSuffix(add addFunc, file protosource.File) error {
	pkg := file.Package()
	if pkg == "" {
		return nil
	}
	if _, ok := protoversion.NewPackageVersionForPackage(pkg); !ok {
		add(file, file.PackageLocation(), nil, `Package name %q should be suffixed with a correctly formed version, such as %q.`, pkg, pkg+".v1")
	}
	return nil
}

// CheckRPCNoClientStreaming is a check function.
var CheckRPCNoClientStreaming = newMethodCheckFunc(checkRPCNoClientStreaming)

func checkRPCNoClientStreaming(add addFunc, method protosource.Method) error {
	if method.ClientStreaming() {
		add(
			method,
			method.Location(),
			// also check the service for this comment ignore
			// this allows users to set this "globally" for a service
			[]protosource.Location{
				method.Service().Location(),
			},
			"RPC %q is client streaming.",
			method.Name(),
		)
	}
	return nil
}

// CheckRPCNoServerStreaming is a check function.
var CheckRPCNoServerStreaming = newMethodCheckFunc(checkRPCNoServerStreaming)

func checkRPCNoServerStreaming(add addFunc, method protosource.Method) error {
	if method.ServerStreaming() {
		add(
			method,
			method.Location(),
			// also check the service for this comment ignore
			// this allows users to set this "globally" for a service
			[]protosource.Location{
				method.Service().Location(),
			},
			"RPC %q is server streaming.",
			method.Name(),
		)
	}
	return nil
}

// CheckRPCPascalCase is a check function.
var CheckRPCPascalCase = newMethodCheckFunc(checkRPCPascalCase)

func checkRPCPascalCase(add addFunc, method protosource.Method) error {
	name := method.Name()
	expectedName := stringutil.ToPascalCase(name)
	if name != expectedName {
		add(
			method,
			method.NameLocation(),
			// also check the service for this comment ignore
			// this allows users to set this "globally" for a service
			[]protosource.Location{
				method.Service().Location(),
			},
			"RPC name %q should be PascalCase, such as %q.",
			name,
			expectedName,
		)
	}
	return nil
}

// CheckRPCRequestResponseUnique is a check function.
var CheckRPCRequestResponseUnique = func(
	id string,
	ignoreFunc internal.IgnoreFunc,
	files []protosource.File,
	allowSameRequestResponse bool,
	allowGoogleProtobufEmptyRequests bool,
	allowGoogleProtobufEmptyResponses bool,
) ([]bufanalysis.FileAnnotation, error) {
	return newFilesCheckFunc(
		func(add addFunc, files []protosource.File) error {
			return checkRPCRequestResponseUnique(
				add,
				files,
				allowSameRequestResponse,
				allowGoogleProtobufEmptyRequests,
				allowGoogleProtobufEmptyResponses,
			)
		},
	)(id, ignoreFunc, files)
}

func checkRPCRequestResponseUnique(
	add addFunc,
	files []protosource.File,
	allowSameRequestResponse bool,