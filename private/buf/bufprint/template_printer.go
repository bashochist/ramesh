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

package bufprint

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	registryv1alpha1 "github.com/bufbuild/buf/private/gen/proto/go/buf/alpha/registry/v1alpha1"
)

type templatePrinter struct {
	writer io.Writer
}

func newTemplatePrinter(writer io.Writer) *templatePrinter {
	return &templatePrinter{
		writer: writer,
	}
}

func (t *templatePrinter) PrintTemplate(ctx context.Context, format Format, template *registryv1alpha1.Template) error {
	switch format {
	case FormatText:
		return t.printTe