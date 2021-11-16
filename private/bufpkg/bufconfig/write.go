
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

package bufconfig

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"text/template"

	"github.com/bufbuild/buf/private/bufpkg/bufcheck/bufbreaking/bufbreakingconfig"
	"github.com/bufbuild/buf/private/bufpkg/bufcheck/buflint/buflintconfig"
	"github.com/bufbuild/buf/private/bufpkg/bufmodule/bufmoduleref"
	"github.com/bufbuild/buf/private/pkg/storage"
	"gopkg.in/yaml.v3"
)

// If this is updated, make sure to update docs.buf.build TODO automate this

const (
	exampleName = "buf.build/acme/weather"
	// This is only used for `buf mod init`.
	tmplDocumentationCommentsData = `{{$top := .}}# This specifies the configuration file version.
#
# This controls the configuration file layout, defaults, and lint/breaking
# rules and rule categories. Buf takes breaking changes seriously in
# all aspects, and none of these will ever change for a given version.
#
# The only valid versions are "v1beta1", "v1".
# This key is required.
version: {{.Version}}

# name is the module name.
{{if .NameUnset}}#{{end}}name: {{.Name}}

# deps are the module dependencies
{{if .DepsUnset}}#{{end}}deps:
{{range $dep := .Deps}}{{if $top.DepsUnset}}#{{end}}  - {{$dep}}
{{end}}
# build contains the options for builds.
#
# This affects the behavior of buf build, as well as the build behavior
# for source lint and breaking change rules.
#
# If you want to build all files in your repository, this section can be
# omitted.
build:

  # excludes is the list of directories to exclude.
  #
  # These directories will not be built or checked. If a directory is excluded,
  # buf treats the directory as if it does not exist.
  #
  # All directory paths in exclude must be relative to the directory of
  # your buf.yaml. Only directories can be specified, and all specified
  # directories must within the root directory.
  {{if not .Uncomment}}#{{end}}excludes:
  {{if not .Uncomment}}#{{end}}  - foo
  {{if not .Uncomment}}#{{end}}  - bar/baz

# lint contains the options for lint rules.
lint:

  # use is the list of rule categories and ids to use for buf lint.
  #
  # Categories are sets of rule ids.
  # Run buf config ls-lint-rules --all to get a list of all rules.
  #
  # The union of the categories and ids will be used.
  #
  # The default is [DEFAULT].
  use:
{{range $lint_id := .LintConfig.Use}}    - {{$lint_id}}
{{end}}
  # except is the list of rule ids or categories to remove from the use
  # list.
  #