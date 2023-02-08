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

var (
	csharpNamespacePathKey      = getPathKey([]int32{8, 37})
	goPackagePathKey            = getPathKey([]int32{8, 11})
	javaMultipleFilesPathKey    = getPathKey([]int32{8, 10})
	javaOuterClassnamePathKey   = getPathKey([]int32{8, 8})
	javaPackagePathKey          = getPathKey([]int32{8, 1})
	javaStringCheckUtf8PathKey  = getPathKey([]int32{8, 27})
	objcClassPrefixPathKey      = getPathKey([]int32{8, 36})
	packagePathKey              = getPathKey([]int32{2})
	phpClassPrefixPathKey       = getPathKey([]int32{8, 40})
	phpNamespacePathKey         = getPathKey([]int32{8, 41})
	phpMetadataNamespacePathKey = getPathKey([]int32{8, 44})
	rubyPackagePathKey          = getPathKey([]int32{8, 45})
	swiftPrefixPathKey          = getPathKey([]int32{8, 39})
	optimizeForPathKey          = getPathKey([]int32{8, 9})
	ccGenericServicesPathKey    = getPathKey([]int32{8, 16})
	javaGenericServicesPathKey  = getPathKey([]int32{8, 17})
	pyGenericServicesPathKey    = getPathKey([]int32{8, 18})
	phpGenericServicesPathKey   = getPathKey([]int32{8, 42})
	ccEnableArenasPathKey       = getPathKey([]int32{8, 31})
	syntaxPathKey               = getPathKey([]int32{12})
)

func getDependencyPath(dependencyIndex int) []int32 {
	return []int32{3, int32(dependencyIndex)}
}

func getMessagePath(topLevelMessageIndex int, nestedMessageIndexes ...int) []int32 {
	path := []int32{4, int32(topLevelMessageIndex)}
	for _, nestedMessageIndex := range nestedMessageIndexes {
		path = append(path, 3, int32(nestedMessageIndex))
	}
	return path
}

func getMessageNamePath(messageIndex int, nestedMessageIndexes ...int) []int32 {
	return append(getMessagePath(messageIndex, nestedMessageIndexes...), 1)
}

func getMessageMessageSetWireFormatPath(messageIndex int, nestedMessageIndexes ...int) []int32 {
	return append(getMessagePath(messageIndex, nestedMessageIndexes...), 7, 1)
}

func getMessageNoStandardDescriptorAccessorPath(messageIndex int, nestedMessageIndexes ...int) []int32 {
	return append(getMessagePath(messageIndex, nestedMessageIndexes...), 7, 2)
}

func getMessageFieldPath(fieldIndex int, topLevelMessageIndex int, nestedMessageIndexes ...int) []int32 {
	return append(getMessagePath(topLevelMessageIndex, nestedMessageIndexes...), 2, int32(fieldIndex))
}

func getMessageFieldNamePath(fieldIndex int, topLevelMessageIndex int, nestedMessageIndexes ...int) []int32 {
	return append(getMessageFieldPath(fieldIndex, topLevelMessageIndex, nestedMessageIndexes...), 1)
}

func getMessageFieldNumberPath(fieldIndex int, topLevelMessageIndex int, nestedMessageIndexes ...int) []int32 {
	return append(getMessageFieldPath(fieldIndex, topLevelMessageIndex, nestedMessageIndexes...), 3)
}

func getMessageFieldTypePath(fieldIndex int, topLevelMessageIndex int, nestedMessageIndexes ...int) []int32 {
	return append(getMessageFieldPath(fieldIndex, topLevelMessageIndex, nestedMessageIndexes...), 5)
}

func getMessageFieldTypeNamePath(fieldIndex int, topLevelMessageIndex int, nestedMessageIndexes ...int) []int32 {
	return append(getMessageFieldPath(fieldIndex, topLevelMessageIndex, nestedMessageIndexes...), 6)
}

func getMessageFieldJSONNamePath(fieldIndex int, topLevelMessageIndex int, nestedMessageIndexes ...int) []int32 {
	return append(getMessageFieldPath(fieldIndex, topLevelMessageIndex, nestedMessageIndexes...), 10)
}

func getMessageFieldJSTypePath(fieldIndex int, topLevelMessageIndex int, nestedMessageIndexes ...int) []int32 {
	return append(getMessageFieldPath(fieldIndex, topLevelMessageIndex, nestedMessageIndexes...), 8, 6)
}

func getMessageFieldCTypePath(fieldIndex int, topLevelMessageIndex int, nestedMessageIndexes ...int) []int32 {
	return append(getMessageFieldPath(fieldIndex, topLevelMessageIndex, nestedMessageIndexes...), 8, 1)
}

func getMessageFieldPackedPath(fieldIndex int, topLevelMessageIndex int, nestedMessageIndexes ...int) []int32 {
	return append(getMessageFieldPath(fieldIndex, topLevelMessageIndex, nestedMessageIndexes...), 8, 2)
}

func getMessageFieldExtendeePath(fieldIndex int, topLevelMessageIndex int, nestedMessageIndexes ...int) []int32 {
	return append(getMessageFieldPath(fieldIndex, topLevelMessageIndex, nestedMessageIndexes...), 2)
}

func getMessageExtensionPath(ext