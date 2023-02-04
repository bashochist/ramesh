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

package protogenutil

import (
	"fmt"
	"path"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

const namedHelperGoPackageOptionKey = "named_go_package"

type namedHelper struct {
	pluginNameToGoPackage map[string]string
}

func newNamedHelper() *namedHelper {
	return &namedHelper{
		pluginNameToGoPackage: make(map[string]string),
	}
}

func (h *namedHelper) NewGoPackageName(
	baseGoPackageName protogen.GoPackageName,
	pluginName string,
) protogen.GoPackageName {
	return protogen.GoPackageName(string(baseGoPackageName) + pluginName)
}

func (h *namedHelper) NewGoImportPath(
	file *protogen.File,
	pluginName string,
) (protogen.GoImportPath, error) {
	return h.newGoImportPath(
		path.Dir(file.GeneratedFilenamePrefix),
		file.GoPackageName,
		pluginName,
	)
}

func (h *namedHelper) NewPackageGoImportPath(
	goPackageFileSet *GoPackageFileSet,
	pluginName string,
) (protogen.GoImportPath, error) {
	return h.newGoImportPath(
		goPackageFileSet.GeneratedDir,
		goPackageFileSet.GoPackageName,
		pluginName,
	)
}

func (h *namedHelper) NewGlobalGoImportPath(
	pluginName string,
) (protogen.GoImportPath, error) {
	goPackage, ok := h.pluginNameToGoPackage[pluginName]
	if !ok {
		return "", fmt.Errorf("no %s specified for plugin %s", namedHelperGoPackageOptionKey, pluginName)
	}
	return protogen.GoImportPath(goPackage), nil
}

func (h *namedHelpe