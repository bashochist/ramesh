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
	pyGenericServicesPathKey    = getP