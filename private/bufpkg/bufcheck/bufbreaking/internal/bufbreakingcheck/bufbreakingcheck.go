
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

// Package bufbreakingcheck impelements the check functions.
//
// These are used by bufbreakingbuild to create RuleBuilders.
package bufbreakingcheck

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bufbuild/buf/private/pkg/protodescriptor"
	"github.com/bufbuild/buf/private/pkg/protosource"
	"github.com/bufbuild/buf/private/pkg/stringutil"
	"google.golang.org/protobuf/types/descriptorpb"
)

// CheckEnumNoDelete is a check function.
var CheckEnumNoDelete = newFilePairCheckFunc(checkEnumNoDelete)

func checkEnumNoDelete(add addFunc, corpus *corpus, previousFile protosource.File, file protosource.File) error {
	previousNestedNameToEnum, err := protosource.NestedNameToEnum(previousFile)
	if err != nil {
		return err
	}
	nestedNameToEnum, err := protosource.NestedNameToEnum(file)
	if err != nil {
		return err
	}
	for previousNestedName := range previousNestedNameToEnum {
		if _, ok := nestedNameToEnum[previousNestedName]; !ok {
			// TODO: search for enum in other files and return that the enum was moved?
			descriptor, location, err := getDescriptorAndLocationForDeletedEnum(file, previousNestedName)
			if err != nil {
				return err
			}
			add(descriptor, nil, location, `Previously present enum %q was deleted from file.`, previousNestedName)
		}
	}
	return nil
}

// CheckEnumValueNoDelete is a check function.
var CheckEnumValueNoDelete = newEnumPairCheckFunc(checkEnumValueNoDelete)

func checkEnumValueNoDelete(add addFunc, corpus *corpus, previousEnum protosource.Enum, enum protosource.Enum) error {
	return checkEnumValueNoDeleteWithRules(add, previousEnum, enum, false, false)
}

// CheckEnumValueNoDeleteUnlessNumberReserved is a check function.
var CheckEnumValueNoDeleteUnlessNumberReserved = newEnumPairCheckFunc(checkEnumValueNoDeleteUnlessNumberReserved)

func checkEnumValueNoDeleteUnlessNumberReserved(add addFunc, corpus *corpus, previousEnum protosource.Enum, enum protosource.Enum) error {
	return checkEnumValueNoDeleteWithRules(add, previousEnum, enum, true, false)
}

// CheckEnumValueNoDeleteUnlessNameReserved is a check function.
var CheckEnumValueNoDeleteUnlessNameReserved = newEnumPairCheckFunc(checkEnumValueNoDeleteUnlessNameReserved)

func checkEnumValueNoDeleteUnlessNameReserved(add addFunc, corpus *corpus, previousEnum protosource.Enum, enum protosource.Enum) error {
	return checkEnumValueNoDeleteWithRules(add, previousEnum, enum, false, true)
}

func checkEnumValueNoDeleteWithRules(add addFunc, previousEnum protosource.Enum, enum protosource.Enum, allowIfNumberReserved bool, allowIfNameReserved bool) error {
	previousNumberToNameToEnumValue, err := protosource.NumberToNameToEnumValue(previousEnum)
	if err != nil {
		return err
	}
	numberToNameToEnumValue, err := protosource.NumberToNameToEnumValue(enum)
	if err != nil {
		return err
	}
	for previousNumber, previousNameToEnumValue := range previousNumberToNameToEnumValue {
		if _, ok := numberToNameToEnumValue[previousNumber]; !ok {
			if !isDeletedEnumValueAllowedWithRules(previousNumber, previousNameToEnumValue, enum, allowIfNumberReserved, allowIfNameReserved) {
				suffix := ""
				if allowIfNumberReserved && allowIfNameReserved {
					return errors.New("both allowIfNumberReserved and allowIfNameReserved set")
				}
				if allowIfNumberReserved {
					suffix = fmt.Sprintf(` without reserving the number "%d"`, previousNumber)
				}
				if allowIfNameReserved {
					nameSuffix := ""
					if len(previousNameToEnumValue) > 1 {
						nameSuffix = "s"
					}
					suffix = fmt.Sprintf(` without reserving the name%s %s`, nameSuffix, stringutil.JoinSliceQuoted(getSortedEnumValueNames(previousNameToEnumValue), ", "))
				}
				add(enum, nil, enum.Location(), `Previously present enum value "%d" on enum %q was deleted%s.`, previousNumber, enum.Name(), suffix)
			}
		}
	}
	return nil
}

func isDeletedEnumValueAllowedWithRules(previousNumber int, previousNameToEnumValue map[string]protosource.EnumValue, enum protosource.Enum, allowIfNumberReserved bool, allowIfNameReserved bool) bool {
	if allowIfNumberReserved {
		return protosource.NumberInReservedRanges(previousNumber, enum.ReservedTagRanges()...)
	}
	if allowIfNameReserved {
		// if true for all names, then ok
		for previousName := range previousNameToEnumValue {
			if !protosource.NameInReservedNames(previousName, enum.ReservedNames()...) {
				return false
			}
		}
		return true
	}
	return false
}

// CheckEnumValueSameName is a check function.
var CheckEnumValueSameName = newEnumValuePairCheckFunc(checkEnumValueSameName)

func checkEnumValueSameName(add addFunc, corpus *corpus, previousNameToEnumValue map[string]protosource.EnumValue, nameToEnumValue map[string]protosource.EnumValue) error {
	previousNames := getSortedEnumValueNames(previousNameToEnumValue)
	names := getSortedEnumValueNames(nameToEnumValue)
	// all current names for this number need to be in the previous set
	// ie if you now have FOO=2, BAR=2, you need to have had FOO=2, BAR=2 previously
	// FOO=2, BAR=2, BAZ=2 now would pass
	// FOO=2, BAR=2, BAZ=2 previously would fail
	if !stringutil.SliceElementsContained(names, previousNames) {
		previousNamesString := stringutil.JoinSliceQuoted(previousNames, ", ")
		namesString := stringutil.JoinSliceQuoted(names, ", ")
		nameSuffix := ""
		if len(previousNames) > 1 && len(names) > 1 {
			nameSuffix = "s"
		}
		for _, enumValue := range nameToEnumValue {
			add(enumValue, nil, enumValue.NumberLocation(), `Enum value "%d" on enum %q changed name%s from %s to %s.`, enumValue.Number(), enumValue.Enum().Name(), nameSuffix, previousNamesString, namesString)
		}
	}
	return nil
}

// CheckExtensionMessageNoDelete is a check function.
var CheckExtensionMessageNoDelete = newMessagePairCheckFunc(checkExtensionMessageNoDelete)

func checkExtensionMessageNoDelete(add addFunc, corpus *corpus, previousMessage protosource.Message, message protosource.Message) error {
	previousStringToExtensionRange := protosource.StringToExtensionMessageRange(previousMessage)
	stringToExtensionRange := protosource.StringToExtensionMessageRange(message)
	for previousString := range previousStringToExtensionRange {
		if _, ok := stringToExtensionRange[previousString]; !ok {
			add(message, nil, message.Location(), `Previously present extension range %q on message %q was deleted.`, previousString, message.Name())
		}
	}
	return nil
}

// CheckFieldNoDelete is a check function.
var CheckFieldNoDelete = newMessagePairCheckFunc(checkFieldNoDelete)

func checkFieldNoDelete(add addFunc, corpus *corpus, previousMessage protosource.Message, message protosource.Message) error {
	return checkFieldNoDeleteWithRules(add, previousMessage, message, false, false)
}

// CheckFieldNoDeleteUnlessNumberReserved is a check function.
var CheckFieldNoDeleteUnlessNumberReserved = newMessagePairCheckFunc(checkFieldNoDeleteUnlessNumberReserved)

func checkFieldNoDeleteUnlessNumberReserved(add addFunc, corpus *corpus, previousMessage protosource.Message, message protosource.Message) error {
	return checkFieldNoDeleteWithRules(add, previousMessage, message, true, false)
}

// CheckFieldNoDeleteUnlessNameReserved is a check function.
var CheckFieldNoDeleteUnlessNameReserved = newMessagePairCheckFunc(checkFieldNoDeleteUnlessNameReserved)

func checkFieldNoDeleteUnlessNameReserved(add addFunc, corpus *corpus, previousMessage protosource.Message, message protosource.Message) error {
	return checkFieldNoDeleteWithRules(add, previousMessage, message, false, true)
}

func checkFieldNoDeleteWithRules(add addFunc, previousMessage protosource.Message, message protosource.Message, allowIfNumberReserved bool, allowIfNameReserved bool) error {
	previousNumberToField, err := protosource.NumberToMessageField(previousMessage)
	if err != nil {
		return err
	}
	numberToField, err := protosource.NumberToMessageField(message)
	if err != nil {
		return err
	}
	for previousNumber, previousField := range previousNumberToField {
		if _, ok := numberToField[previousNumber]; !ok {
			if !isDeletedFieldAllowedWithRules(previousField, message, allowIfNumberReserved, allowIfNameReserved) {
				// otherwise prints as hex
				previousNumberString := strconv.FormatInt(int64(previousNumber), 10)
				suffix := ""
				if allowIfNumberReserved && allowIfNameReserved {
					return errors.New("both allowIfNumberReserved and allowIfNameReserved set")
				}
				if allowIfNumberReserved {
					suffix = fmt.Sprintf(` without reserving the number "%d"`, previousField.Number())
				}
				if allowIfNameReserved {
					suffix = fmt.Sprintf(` without reserving the name %q`, previousField.Name())
				}
				add(message, nil, message.Location(), `Previously present field %q with name %q on message %q was deleted%s.`, previousNumberString, previousField.Name(), message.Name(), suffix)
			}
		}
	}
	return nil
}

func isDeletedFieldAllowedWithRules(previousField protosource.Field, message protosource.Message, allowIfNumberReserved bool, allowIfNameReserved bool) bool {
	return (allowIfNumberReserved && protosource.NumberInReservedRanges(previousField.Number(), message.ReservedTagRanges()...)) ||
		(allowIfNameReserved && protosource.NameInReservedNames(previousField.Name(), message.ReservedNames()...))
}

// CheckFieldSameCType is a check function.
var CheckFieldSameCType = newFieldPairCheckFunc(checkFieldSameCType)

func checkFieldSameCType(add addFunc, corpus *corpus, previousField protosource.Field, field protosource.Field) error {
	if previousField.CType() != field.CType() {
		// otherwise prints as hex
		numberString := strconv.FormatInt(int64(field.Number()), 10)
		add(field, nil, withBackupLocation(field.CTypeLocation(), field.Location()), `Field %q with name %q on message %q changed option "ctype" from %q to %q.`, numberString, field.Name(), field.Message().Name(), previousField.CType().String(), field.CType().String())
	}
	return nil
}

// CheckFieldSameJSONName is a check function.
var CheckFieldSameJSONName = newFieldPairCheckFunc(checkFieldSameJSONName)

func checkFieldSameJSONName(add addFunc, corpus *corpus, previousField protosource.Field, field protosource.Field) error {
	if previousField.JSONName() != field.JSONName() {
		// otherwise prints as hex
		numberString := strconv.FormatInt(int64(field.Number()), 10)
		add(field, nil, withBackupLocation(field.JSONNameLocation(), field.Location()), `Field %q with name %q on message %q changed option "json_name" from %q to %q.`, numberString, field.Name(), field.Message().Name(), previousField.JSONName(), field.JSONName())
	}
	return nil
}

// CheckFieldSameJSType is a check function.
var CheckFieldSameJSType = newFieldPairCheckFunc(checkFieldSameJSType)

func checkFieldSameJSType(add addFunc, corpus *corpus, previousField protosource.Field, field protosource.Field) error {
	if previousField.JSType() != field.JSType() {
		// otherwise prints as hex
		numberString := strconv.FormatInt(int64(field.Number()), 10)
		add(field, nil, withBackupLocation(field.JSTypeLocation(), field.Location()), `Field %q with name %q on message %q changed option "jstype" from %q to %q.`, numberString, field.Name(), field.Message().Name(), previousField.JSType().String(), field.JSType().String())
	}
	return nil
}

// CheckFieldSameLabel is a check function.
var CheckFieldSameLabel = newFieldPairCheckFunc(checkFieldSameLabel)

func checkFieldSameLabel(add addFunc, corpus *corpus, previousField protosource.Field, field protosource.Field) error {
	if previousField.Label() != field.Label() {
		// otherwise prints as hex
		numberString := strconv.FormatInt(int64(field.Number()), 10)
		// TODO: specific label location
		add(field, nil, field.Location(), `Field %q on message %q changed label from %q to %q.`, numberString, field.Message().Name(), protodescriptor.FieldDescriptorProtoLabelPrettyString(previousField.Label()), protodescriptor.FieldDescriptorProtoLabelPrettyString(field.Label()))
	}
	return nil
}

// CheckFieldSameName is a check function.
var CheckFieldSameName = newFieldPairCheckFunc(checkFieldSameName)

func checkFieldSameName(add addFunc, corpus *corpus, previousField protosource.Field, field protosource.Field) error {
	if previousField.Name() != field.Name() {
		// otherwise prints as hex
		numberString := strconv.FormatInt(int64(field.Number()), 10)
		add(field, nil, field.NameLocation(), `Field %q on message %q changed name from %q to %q.`, numberString, field.Message().Name(), previousField.Name(), field.Name())
	}
	return nil
}

// CheckFieldSameOneof is a check function.
var CheckFieldSameOneof = newFieldPairCheckFunc(checkFieldSameOneof)

func checkFieldSameOneof(add addFunc, corpus *corpus, previousField protosource.Field, field protosource.Field) error {
	previousOneof := previousField.Oneof()
	oneof := field.Oneof()
	previousInsideOneof := previousOneof != nil
	insideOneof := oneof != nil
	if !previousInsideOneof && !insideOneof {
		return nil
	}
	if previousInsideOneof && insideOneof {
		if previousOneof.Name() != oneof.Name() {
			// otherwise prints as hex
			numberString := strconv.FormatInt(int64(field.Number()), 10)
			add(field, nil, field.Location(), `Field %q on message %q moved from oneof %q to oneof %q.`, numberString, field.Message().Name(), previousOneof.Name(), oneof.Name())
		}
		return nil
	}

	previous := "inside"
	current := "outside"
	if insideOneof {
		previous = "outside"
		current = "inside"
	}
	// otherwise prints as hex
	numberString := strconv.FormatInt(int64(field.Number()), 10)
	add(field, nil, field.Location(), `Field %q on message %q moved from %s to %s a oneof.`, numberString, field.Message().Name(), previous, current)
	return nil
}

// TODO: locations not working for map entries
// TODO: weird output for map entries:
//
// breaking_field_same_type/1.proto:1:1:Field "2" on message "SixEntry" changed type from ".a.One" to ".a.Two".
// breaking_field_same_type/1.proto:1:1:Field "2" on message "SixEntry" changed type from ".a.One" to ".a.Two".
// breaking_field_same_type/1.proto:1:1:Field "2" on message "SixEntry" changed type from ".a.One" to ".a.Two".
// breaking_field_same_type/1.proto:8:3:Field "1" on message "Two" changed type from "int32" to "int64".
// breaking_field_same_type/1.proto:9:3:Field "2" on message "Two" changed type from ".a.One" to ".a.Two".
// breaking_field_same_type/1.proto:11:3:Field "4" on message "Two" changed type from ".a.One" to ".a.Two".
// breaking_field_same_type/1.proto:12:3:Field "5" on message "Two" changed type from ".a.Two.FiveEntry" to ".a.Two".
// breaking_field_same_type/1.proto:19:7:Field "1" on message "Five" changed type from "int32" to "int64".
// breaking_field_same_type/1.proto:20:7:Field "2" on message "Five" changed type from ".a.One" to ".a.Two".
// breaking_field_same_type/1.proto:22:7:Field "4" on message "Five" changed type from ".a.One" to ".a.Two".
// breaking_field_same_type/1.proto:23:7:Field "5" on message "Five" changed type from ".a.Three.Four.Five.FiveEntry" to ".a.Two".
// breaking_field_same_type/1.proto:36:5:Field "1" on message "Seven" changed type from "int32" to "int64".
// breaking_field_same_type/1.proto:37:5:Field "2" on message "Seven" changed type from ".a.One" to ".a.Two".
// breaking_field_same_type/1.proto:39:5:Field "4" on message "Seven" changed type from ".a.One" to ".a.Two".
// breaking_field_same_type/1.proto:40:5:Field "5" on message "Seven" changed type from ".a.Three.Seven.FiveEntry" to ".a.Two".
// breaking_field_same_type/2.proto:64:5:Field "1" on message "Nine" changed type from "int32" to "int64".
// breaking_field_same_type/2.proto:65:5:Field "2" on message "Nine" changed type from ".a.One" to ".a.Nine".

// CheckFieldSameType is a check function.
var CheckFieldSameType = newFieldPairCheckFunc(checkFieldSameType)

func checkFieldSameType(add addFunc, corpus *corpus, previousField protosource.Field, field protosource.Field) error {
	if previousField.Type() != field.Type() {
		addFieldChangedType(add, previousField, field)
		return nil
	}

	switch field.Type() {
	case descriptorpb.FieldDescriptorProto_TYPE_ENUM,
		descriptorpb.FieldDescriptorProto_TYPE_GROUP,
		descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
		if previousField.TypeName() != field.TypeName() {
			addEnumGroupMessageFieldChangedTypeName(add, previousField, field)
		}
	}
	return nil
}

// CheckFieldWireCompatibleType is a check function.
var CheckFieldWireCompatibleType = newFieldPairCheckFunc(checkFieldWireCompatibleType)

func checkFieldWireCompatibleType(add addFunc, corpus *corpus, previousField protosource.Field, field protosource.Field) error {
	previousWireCompatibilityGroup, ok := fieldDescriptorProtoTypeToWireCompatiblityGroup[previousField.Type()]
	if !ok {
		return fmt.Errorf("unknown FieldDescriptorProtoType: %v", previousField.Type())
	}
	wireCompatibilityGroup, ok := fieldDescriptorProtoTypeToWireCompatiblityGroup[field.Type()]
	if !ok {
		return fmt.Errorf("unknown FieldDescriptorProtoType: %v", field.Type())
	}
	if previousWireCompatibilityGroup != wireCompatibilityGroup {
		extraMessages := []string{
			"See https://developers.google.com/protocol-buffers/docs/proto3#updating for wire compatibility rules.",
		}
		switch {
		case previousField.Type() == descriptorpb.FieldDescriptorProto_TYPE_STRING && field.Type() == descriptorpb.FieldDescriptorProto_TYPE_BYTES:
			// It is OK to evolve from string to bytes
			return nil
		case previousField.Type() == descriptorpb.FieldDescriptorProto_TYPE_BYTES && field.Type() == descriptorpb.FieldDescriptorProto_TYPE_STRING:
			extraMessages = append(
				extraMessages,
				"Note that while string and bytes are compatible if the data is valid UTF-8, there is no way to enforce that a field is UTF-8, so these fields may be incompatible.",
			)
		}
		addFieldChangedType(add, previousField, field, extraMessages...)
		return nil
	}
	switch field.Type() {
	case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
		if previousField.TypeName() != field.TypeName() {
			return checkEnumWireCompatibleForField(add, corpus, previousField, field)
		}
	case descriptorpb.FieldDescriptorProto_TYPE_GROUP,
		descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
		if previousField.TypeName() != field.TypeName() {
			addEnumGroupMessageFieldChangedTypeName(add, previousField, field)
			return nil
		}
	}
	return nil
}

// CheckFieldWireJSONCompatibleType is a check function.
var CheckFieldWireJSONCompatibleType = newFieldPairCheckFunc(checkFieldWireJSONCompatibleType)

func checkFieldWireJSONCompatibleType(add addFunc, corpus *corpus, previousField protosource.Field, field protosource.Field) error {
	previousWireJSONCompatibilityGroup, ok := fieldDescriptorProtoTypeToWireJSONCompatiblityGroup[previousField.Type()]
	if !ok {
		return fmt.Errorf("unknown FieldDescriptorProtoType: %v", previousField.Type())
	}
	wireJSONCompatibilityGroup, ok := fieldDescriptorProtoTypeToWireJSONCompatiblityGroup[field.Type()]
	if !ok {
		return fmt.Errorf("unknown FieldDescriptorProtoType: %v", field.Type())
	}
	if previousWireJSONCompatibilityGroup != wireJSONCompatibilityGroup {
		addFieldChangedType(
			add,
			previousField,
			field,
			"See https://developers.google.com/protocol-buffers/docs/proto3#updating for wire compatibility rules and https://developers.google.com/protocol-buffers/docs/proto3#json for JSON compatibility rules.",
		)
		return nil
	}
	switch field.Type() {
	case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
		if previousField.TypeName() != field.TypeName() {
			return checkEnumWireCompatibleForField(add, corpus, previousField, field)
		}
	case descriptorpb.FieldDescriptorProto_TYPE_GROUP,
		descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
		if previousField.TypeName() != field.TypeName() {
			addEnumGroupMessageFieldChangedTypeName(add, previousField, field)
			return nil
		}
	}
	return nil
}

func checkEnumWireCompatibleForField(add addFunc, corpus *corpus, previousField protosource.Field, field protosource.Field) error {
	previousEnum, err := getEnumByFullName(
		corpus.previousFiles,
		strings.TrimPrefix(previousField.TypeName(), "."),
	)
	if err != nil {
		return err
	}
	enum, err := getEnumByFullName(
		corpus.files,
		strings.TrimPrefix(field.TypeName(), "."),
	)
	if err != nil {
		return err
	}
	if previousEnum.Name() != enum.Name() {
		// If the short names are not equal, we say that this is a different enum.
		addEnumGroupMessageFieldChangedTypeName(add, previousField, field)
		return nil
	}
	isSubset, err := protosource.EnumIsSubset(enum, previousEnum)
	if err != nil {
		return err
	}
	if !isSubset {
		// If the previous enum is not a subset of the new enum, we say that
		// this is a different enum.
		// We allow subsets so that enum values can be added within the
		// same change.
		addEnumGroupMessageFieldChangedTypeName(add, previousField, field)
		return nil
	}
	return nil
}

func addFieldChangedType(add addFunc, previousField protosource.Field, field protosource.Field, extraMessages ...string) {
	combinedExtraMessage := ""
	if len(extraMessages) > 0 {
		// protect against mistakenly added empty extra messages
		if joined := strings.TrimSpace(strings.Join(extraMessages, " ")); joined != "" {
			combinedExtraMessage = " " + joined
		}
	}
	var fieldLocation protosource.Location
	switch field.Type() {
	case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE,
		descriptorpb.FieldDescriptorProto_TYPE_ENUM,
		descriptorpb.FieldDescriptorProto_TYPE_GROUP:
		fieldLocation = field.TypeNameLocation()
	default:
		fieldLocation = field.TypeLocation()
	}
	// otherwise prints as hex
	previousNumberString := strconv.FormatInt(int64(previousField.Number()), 10)
	add(
		field,
		nil,
		fieldLocation,
		`Field %q on message %q changed type from %q to %q.%s`,
		previousNumberString,
		field.Message().Name(),
		protodescriptor.FieldDescriptorProtoTypePrettyString(previousField.Type()),
		protodescriptor.FieldDescriptorProtoTypePrettyString(field.Type()),
		combinedExtraMessage,
	)
}

func addEnumGroupMessageFieldChangedTypeName(add addFunc, previousField protosource.Field, field protosource.Field) {
	// otherwise prints as hex
	numberString := strconv.FormatInt(int64(previousField.Number()), 10)
	add(
		field,
		nil,
		field.TypeNameLocation(),
		`Field %q on message %q changed type from %q to %q.`,
		numberString,
		field.Message().Name(),
		strings.TrimPrefix(previousField.TypeName(), "."),
		strings.TrimPrefix(field.TypeName(), "."),
	)
}

// CheckFileNoDelete is a check function.
var CheckFileNoDelete = newFilesCheckFunc(checkFileNoDelete)

func checkFileNoDelete(add addFunc, corpus *corpus) error {
	previousFilePathToFile, err := protosource.FilePathToFile(corpus.previousFiles...)
	if err != nil {
		return err
	}
	filePathToFile, err := protosource.FilePathToFile(corpus.files...)
	if err != nil {
		return err
	}
	for previousFilePath, previousFile := range previousFilePathToFile {
		if _, ok := filePathToFile[previousFilePath]; !ok {
			// Add previous descriptor to check for ignores. This will mean that if
			// we have ignore_unstable_packages set, this file will cause the ignore
			// to happen.
			add(nil, []protosource.Descriptor{previousFile}, nil, `Previously present file %q was deleted.`, previousFilePath)
		}
	}
	return nil
}

// CheckFileSameCsharpNamespace is a check function.
var CheckFileSameCsharpNamespace = newFilePairCheckFunc(checkFileSameCsharpNamespace)

func checkFileSameCsharpNamespace(add addFunc, corpus *corpus, previousFile protosource.File, file protosource.File) error {
	return checkFileSameValue(add, previousFile.CsharpNamespace(), file.CsharpNamespace(), file, file.CsharpNamespaceLocation(), `option "csharp_namespace"`)
}

// CheckFileSameGoPackage is a check function.
var CheckFileSameGoPackage = newFilePairCheckFunc(checkFileSameGoPackage)

func checkFileSameGoPackage(add addFunc, corpus *corpus, previousFile protosource.File, file protosource.File) error {
	return checkFileSameValue(add, previousFile.GoPackage(), file.GoPackage(), file, file.GoPackageLocation(), `option "go_package"`)
}

// CheckFileSameJavaMultipleFiles is a check function.
var CheckFileSameJavaMultipleFiles = newFilePairCheckFunc(checkFileSameJavaMultipleFiles)

func checkFileSameJavaMultipleFiles(add addFunc, corpus *corpus, previousFile protosource.File, file protosource.File) error {
	return checkFileSameValue(add, strconv.FormatBool(previousFile.JavaMultipleFiles()), strconv.FormatBool(file.JavaMultipleFiles()), file, file.JavaMultipleFilesLocation(), `option "java_multiple_files"`)
}

// CheckFileSameJavaOuterClassname is a check function.
var CheckFileSameJavaOuterClassname = newFilePairCheckFunc(checkFileSameJavaOuterClassname)

func checkFileSameJavaOuterClassname(add addFunc, corpus *corpus, previousFile protosource.File, file protosource.File) error {
	return checkFileSameValue(add, previousFile.JavaOuterClassname(), file.JavaOuterClassname(), file, file.JavaOuterClassnameLocation(), `option "java_outer_classname"`)
}

// CheckFileSameJavaPackage is a check function.
var CheckFileSameJavaPackage = newFilePairCheckFunc(checkFileSameJavaPackage)

func checkFileSameJavaPackage(add addFunc, corpus *corpus, previousFile protosource.File, file protosource.File) error {
	return checkFileSameValue(add, previousFile.JavaPackage(), file.JavaPackage(), file, file.JavaPackageLocation(), `option "java_package"`)
}

// CheckFileSameJavaStringCheckUtf8 is a check function.
var CheckFileSameJavaStringCheckUtf8 = newFilePairCheckFunc(checkFileSameJavaStringCheckUtf8)

func checkFileSameJavaStringCheckUtf8(add addFunc, corpus *corpus, previousFile protosource.File, file protosource.File) error {
	return checkFileSameValue(add, strconv.FormatBool(previousFile.JavaStringCheckUtf8()), strconv.FormatBool(file.JavaStringCheckUtf8()), file, file.JavaStringCheckUtf8Location(), `option "java_string_check_utf8"`)
}

// CheckFileSameObjcClassPrefix is a check function.
var CheckFileSameObjcClassPrefix = newFilePairCheckFunc(checkFileSameObjcClassPrefix)

func checkFileSameObjcClassPrefix(add addFunc, corpus *corpus, previousFile protosource.File, file protosource.File) error {
	return checkFileSameValue(add, previousFile.ObjcClassPrefix(), file.ObjcClassPrefix(), file, file.ObjcClassPrefixLocation(), `option "objc_class_prefix"`)
}

// CheckFileSamePackage is a check function.
var CheckFileSamePackage = newFilePairCheckFunc(checkFileSamePackage)

func checkFileSamePackage(add addFunc, corpus *corpus, previousFile protosource.File, file protosource.File) error {
	return checkFileSameValue(add, previousFile.Package(), file.Package(), file, file.PackageLocation(), `package`)
}

// CheckFileSamePhpClassPrefix is a check function.
var CheckFileSamePhpClassPrefix = newFilePairCheckFunc(checkFileSamePhpClassPrefix)

func checkFileSamePhpClassPrefix(add addFunc, corpus *corpus, previousFile protosource.File, file protosource.File) error {
	return checkFileSameValue(add, previousFile.PhpClassPrefix(), file.PhpClassPrefix(), file, file.PhpClassPrefixLocation(), `option "php_class_prefix"`)
}

// CheckFileSamePhpNamespace is a check function.
var CheckFileSamePhpNamespace = newFilePairCheckFunc(checkFileSamePhpNamespace)

func checkFileSamePhpNamespace(add addFunc, corpus *corpus, previousFile protosource.File, file protosource.File) error {
	return checkFileSameValue(add, previousFile.PhpNamespace(), file.PhpNamespace(), file, file.PhpNamespaceLocation(), `option "php_namespace"`)
}

// CheckFileSamePhpMetadataNamespace is a check function.
var CheckFileSamePhpMetadataNamespace = newFilePairCheckFunc(checkFileSamePhpMetadataNamespace)

func checkFileSamePhpMetadataNamespace(add addFunc, corpus *corpus, previousFile protosource.File, file protosource.File) error {
	return checkFileSameValue(add, previousFile.PhpMetadataNamespace(), file.PhpMetadataNamespace(), file, file.PhpMetadataNamespaceLocation(), `option "php_metadata_namespace"`)
}

// CheckFileSameRubyPackage is a check function.
var CheckFileSameRubyPackage = newFilePairCheckFunc(checkFileSameRubyPackage)

func checkFileSameRubyPackage(add addFunc, corpus *corpus, previousFile protosource.File, file protosource.File) error {
	return checkFileSameValue(add, previousFile.RubyPackage(), file.RubyPackage(), file, file.RubyPackageLocation(), `option "ruby_package"`)
}

// CheckFileSameSwiftPrefix is a check function.
var CheckFileSameSwiftPrefix = newFilePairCheckFunc(checkFileSameSwiftPrefix)

func checkFileSameSwiftPrefix(add addFunc, corpus *corpus, previousFile protosource.File, file protosource.File) error {
	return checkFileSameValue(add, previousFile.SwiftPrefix(), file.SwiftPrefix(), file, file.SwiftPrefixLocation(), `option "swift_prefix"`)
}

// CheckFileSameOptimizeFor is a check function.
var CheckFileSameOptimizeFor = newFilePairCheckFunc(checkFileSameOptimizeFor)

func checkFileSameOptimizeFor(add addFunc, corpus *corpus, previousFile protosource.File, file protosource.File) error {
	return checkFileSameValue(add, previousFile.OptimizeFor().String(), file.OptimizeFor().String(), file, file.OptimizeForLocation(), `option "optimize_for"`)
}

// CheckFileSameCcGenericServices is a check function.
var CheckFileSameCcGenericServices = newFilePairCheckFunc(checkFileSameCcGenericServices)

func checkFileSameCcGenericServices(add addFunc, corpus *corpus, previousFile protosource.File, file protosource.File) error {
	return checkFileSameValue(add, strconv.FormatBool(previousFile.CcGenericServices()), strconv.FormatBool(file.CcGenericServices()), file, file.CcGenericServicesLocation(), `option "cc_generic_services"`)
}
