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

package bufformat

import (
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/bufbuild/protocompile/ast"
	"go.uber.org/multierr"
)

// formatter writes an *ast.FileNode as a .proto file.
type formatter struct {
	writer   io.Writer
	fileNode *ast.FileNode

	// Current level of indentation.
	indent int
	// The last character written to writer.
	lastWritten rune

	// The last node written. This must be updated from all functions
	// that write comments with a node. This flag informs how the next
	// node's leading comments and whitespace should be written.
	previousNode ast.Node

	// If true, a space will be written to the output unless the next character
	// written is a newline (don't wait errant trailing spaces).
	pendingSpace bool
	// If true, the formatter is in the middle of printing compact options.
	inCompactOptions bool

	// Track runes that open blocks/scopes and are expected to increase indention
	// level. For example, when runes "{" "[" "(" ")" are written, the pending
	// value is 2 (increment three times for "{" "[" "("; decrement once for ")").
	// If it's greater than zero at the end of a line, we call In() so that
	// subsequent lines are indented. If it's less than zero at the end of a line,
	// we call Out(). This minimizes the amount of explicit indent/unindent code
	// that is needed and makes it less error-prone.
	pendingIndent int
	// If true, an inline node/sequence is being written. We treat whitespace a
	// little differently for when blocks are printed inline vs. across multiple
	// lines. So this flag informs the logic that makes those whitespace decisions.
	inline bool

	// Records all errors that occur during the formatting process. Nearly any
	// non-nil error represents a bug in the implementation.
	err error
}

// newFormatter returns a new formatter for the given file.
func newFormatter(
	writer io.Writer,
	fileNode *ast.FileNode,
) *formatter {
	return &formatter{
		writer:   writer,
		fileNode: fileNode,
	}
}

// Run runs the formatter and writes the file's content to the formatter's writer.
func (f *formatter) Run() error {
	f.writeFile()
	return f.err
}

// P prints a line to the generated output.
func (f *formatter) P(elements ...string) {
	if len(elements) > 0 {
		// We only want to write an inde