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
		// We only want to write an indent if we're
		// writing elements (not just a newline).
		f.Indent(nil)
		for _, elem := range elements {
			f.WriteString(elem)
		}
	}
	f.WriteString("\n")

	if f.pendingIndent > 0 {
		f.In()
	} else if f.pendingIndent < 0 {
		f.Out()
	}
	f.pendingIndent = 0
}

// Space adds a space to the generated output.
func (f *formatter) Space() {
	f.pendingSpace = true
}

// In increases the current level of indentation.
func (f *formatter) In() {
	f.indent++
}

// Out reduces the current level of indentation.
func (f *formatter) Out() {
	if f.indent <= 0 {
		// Unreachable.
		f.err = multierr.Append(
			f.err,
			errors.New("internal error: attempted to decrement indentation at zero"),
		)
		return
	}
	f.indent--
}

// Indent writes the number of spaces associated
// with the current level of indentation.
func (f *formatter) Indent(nextNode ast.Node) {
	// only indent at beginning of line
	if f.lastWritten != '\n' {
		return
	}
	indent := f.indent
	if rn, ok := nextNode.(*ast.RuneNode); ok && indent > 0 {
		if strings.ContainsRune("}])>", rn.Rune) {
			indent--
		}
	}
	f.WriteString(strings.Repeat("  ", indent))
}

// WriteString writes the given element to the generated output.
func (f *formatter) WriteString(elem string) {
	if f.pendingSpace {
		f.pendingSpace = false
		first, _ := utf8.DecodeRuneInString(elem)

		// We don't want "dangling spaces" before certain characters:
		// newlines, commas, and semicolons. Also, when writing
		// elements inline, we don't want spaces before close parens
		// and braces. Similarly, we don't want extra/doubled spaces
		// or dangling spaces after certain characters when printing
		// inline, like open parens/braces. So only print the space
		// if the previous and next character don't match above
		// conditions.

		prevBlockList := "\x00 \t\n"
		nextBlockList := "\n;,"
		if f.inline {
			prevBlockList = "\x00 \t\n<[{("
			nextBlockList = "\n;,)]}>"
		}

		if !strings.ContainsRune(prevBlockList, f.lastWritten) &&
			!strings.ContainsRune(nextBlockList, first) {
			if _, err := f.writer.Write([]byte{' '}); err != nil {
				f.err = multierr.Append(f.err, err)
				return
			}
		}
	}
	if len(elem) == 0 {
		return
	}
	f.lastWritten, _ = utf8.DecodeLastRuneInString(elem)
	if _, err := f.writer.Write([]byte(elem)); err != nil {
		f.err = multierr.Append(f.err, err)
	}
}

// SetPreviousNode sets the previously written node. This should
// be called in all of the comment writing functions.
func (f *formatter) SetPreviousNode(node ast.Node) {
	f.previousNode = node
}

// writeFile writes the file node.
func (f *formatter) writeFile() {
	f.writeFileHeader()
	f.writeFileTypes()
	if f.fileNode.EOF != nil {
		info := f.fileNode.NodeInfo(f.fileNode.EOF)
		f.writeMultilineComments(info.LeadingComments())
	}
	if f.lastWritten != 0 && f.lastWritten != '\n' {
		// If anything was written, we always conclude with
		// a newline.
		f.P()
	}
}

// writeFileHeader writes the header of a .proto file. This includes the syntax,
// package, imports, and options (in that order). The imports and options are
// sorted. All other file elements are handled by f.writeFileTypes.
//
// For example,
//
//	syntax = "proto3";
//
//	package acme.v1.weather;
//
//	import "acme/payment/v1/payment.proto";
//	import "google/type/datetime.proto";
//
//	option cc_enable_arenas = true;
//	option optimize_for = SPEED;
func (f *formatter) writeFileHeader() {
	var (
		packageNode *ast.PackageNode
		importNodes []*ast.ImportNode
		optionNodes []*ast.OptionNode
	)
	for _, fileElement := range f.fileNode.Decls {
		switch node := fileElement.(type) {
		case *ast.PackageNode:
			packageNode = node
		case *ast.ImportNode:
			importNodes = append(importNodes, node)
		case *ast.OptionNode:
			optionNodes = append(optionNodes, node)
		default:
			continue
		}
	}
	if f.fileNode.Syntax == nil && packageNode == nil && importNodes == nil && optionNodes == nil {
		// There aren't any header values, so we can return early.
		return
	}
	if syntaxNode := f.fileNode.Syntax; syntaxNode != nil {
		f.writeSyntax(syntaxNode)
	}
	if packageNode != nil {
		f.writePackage(packageNode)
	}
	sort.Slice(importNodes, func(i, j int) bool {
		return importNodes[i].Name.AsString() < importNodes[j].Name.AsString()
	})
	for i, importNode := range importNodes {
		if i == 0 && f.previousNode != nil && !f.leadingCommentsContainBlankLine(importNode) {
			f.P()
		}
		f.writeImport(importNode, i > 0)
	}
	sort.Slice(optionNodes, func(i, j int) bool {
		// The default options (e.g. cc_enable_arenas) should always
		// be sorted above custom options (which are identified by a
		// leading '(').
		left := stringForOptionName(optionNodes[i].Name)
		right := stringForOptionName(optionNodes[j].Name)
		if strings.HasPrefix(left, "(") && !strings.HasPrefix(right, "(") {
			// Prefer the default option on the right.
			return false
		}
		if !strings.HasPrefix(left, "(") && strings.HasPrefix(right, "(") {
			// Prefer the default option on the left.
			return true
		}
		// Both options are custom, so we defer to the standard sorting.
		return left < right
	})
	for i, optionNode := range optionNodes {
		if i == 0 && f.previousNode != nil && !f.leadingCommentsContainBlankLine(optionNode) {
			f.P()
		}
		f.writeFileOption(optionNode, i > 0)
	}
}

// writeFileTypes writes the types defined in a .proto file. This includes the messages, enums,
// services, etc. All other elements are ignored since they are handled by f.writeFileHeader.
func (f *formatter) writeFileTypes() {
	for i, fileElement := range f.fileNode.Decls {
		switch node := fileElement.(type) {
		case *ast.PackageNode, *ast.OptionNode, *ast.ImportNode, *ast.EmptyDeclNode:
			// These elements have already been written by f.writeFileHeader.
			continue
		default:
			info := f.fileNode.NodeInfo(node)
			wantNewline := f.previousNode != nil && (i == 0 || info.LeadingComments().Len() > 0)
			if wantNewline && !f.leadingCommentsContainBlankLine(node) {
				f.P()
			}
			f.writeNode(node)
		}
	}
}

// writeSyntax writes the syntax.
//
// For example,
//
//	syntax = "proto3";
func (f *formatter) writeSyntax(syntaxNode *ast.SyntaxNode) {
	f.writeStart(syntaxNode.Keyword)
	f.Space()
	f.writeInline(syntaxNode.Equals)
	f.Space()
	f.writeInline(syntaxNode.Syntax)
	f.writeLineEnd(syntaxNode.Semicolon)
}

// writePackage writes the package.
//
// For example,
//
//	package acme.weather.v1;
func (f *formatter) writePackage(packageNode *ast.PackageNode) {
	f.writeStart(packageNode.Keyword)
	f.Space()
	f.writeInline(packageNode.Name)
	f.writeLineEnd(packageNode.Semicolon)
}

// writeImport writes an import statement.
//
// For example,
//
//	import "google/protobuf/descriptor.proto";
func (f *formatter) writeImport(importNode *ast.ImportNode, forceCompact bool) {
	f.writeStartMaybeCompact(importNode.Keyword, forceCompact)
	f.Space()
	// We don't want to write the "public" and "weak" nodes
	// if they aren't defined. One could be set, but never both.
	switch {
	case importNode.Public != nil:
		f.writeInline(importNode.Public)
		f.Space()
	case importNode.Weak != nil:
		f.writeInline(importNode.Weak)
		f.Space()
	}
	f.writeInline(importNode.Name)
	f.writeLineEnd(importNode.Semicolon)
}

// writeFileOption writes a file option. This function is slightly
// different than f.writeOption because file options are sorted at
// the top of the file, and leading comments are adjusted accordingly.
func (f *formatter) writeFileOption(optionNode *ast.OptionNode, forceCompact bool) {
	f.writeStartMaybeCompact(optionNode.Keyword, forceCompact)
	f.Space()
	f.writeNode(optionNode.Name)
	f.Space()
	f.writeInline(optionNode.Equals)
	if node, ok := optionNode.Val.(*ast.CompoundStringLiteralNode); ok {
		// Compound string literals are written across multiple lines
		// immediately after the '=', so we don't need a trailing
		// space in the option prefix.
		f.writeCompoundStringLiteralIndentEndInline(node)
		f.writeLineEnd(optionNode.Semicolon)
		return
	}
	f.Space()
	f.writeInline(optionNode.Val)
	f.writeLineEnd(optionNode.Semicolon)
}

// writeOption writes an option.
//
// For example,
//
//	option go_package = "github.com/foo/bar";
func (f *formatter) writeOption(optionNode *ast.OptionNode) {
	f.writeOptionPrefix(optionNode)
	if optionNode.Semicolon != nil {
		if node, ok := optionNode.Val.(*ast.CompoundStringLiteralNode); ok {
			// Compound string literals are written across multiple lines
			// immediately after the '=', so we don't need a trailing
			// space in the option prefix.
			f.writeCompoundStringLiteralIndentEndInline(node)
			f.writeLineEnd(optionNode.Semicolon)
			return
		}
		f.writeInline(optionNode.Val)
		f.writeLineEnd(optionNode.Semicolon)
		return
	}

	if node, ok := optionNode.Val.(*ast.CompoundStringLiteralNode); ok {
		f.writeCompoundStringLiteralIndent(node)
		return
	}
	f.writeInline(optionNode.Val)
}

// writeLastCompactOption writes a compact option but preserves its the
// trailing end comments. This is only used for the last compact option
// since it's the only time a trailing ',' will be omitted.
//
// For example,
//
//	[
//	  deprecated = true,
//	  json_name = "something" // Trailing comment on the last element.
//	]
func (f *formatter) writeLastCompactOption(optionNode *ast.OptionNode) {
	f.writeOptionPrefix(optionNode)
	f.writeLineEnd(optionNode.Val)
}

// writeOptionValue writes the option prefix, which makes up all of the
// option's definition, excluding the final token(s).
//
// For example,
//
//	deprecated =
func (f *formatter) writeOptionPrefix(optionNode *ast.OptionNode) {
	if optionNode.Keyword != nil {
		// Compact options don't have the keyword.
		f.writeStart(optionNode.Keyword)
		f.Space()
		f.writeNode(optionNode.Name)
	} else {
		f.writeStart(optionNode.Name)
	}
	f.Space()
	f.writeInline(optionNode.Equals)
	f.Space()
}

// writeOptionName writes an option name.
//
// For example,
//
//	go_package
//	(custom.thing)
//	(custom.thing).bridge.(another.thing)
func (f *formatter) writeOptionName(optionNameNode *ast.OptionNameNode) {
	for i := 0; i < len(optionNameNode.Parts); i++ {
		if f.inCompactOptions && i == 0 {
			// The leading comments of the first token (either open rune or the
			// name) will have already been written, so we need to handle this
			// case specially.
			fieldReferenceNode := optionNameNode.Parts[0]
			if fieldReferenceNode.Open != nil {
				f.writeNode(fieldReferenceNode.Open)
				if info := f.fileNode.NodeInfo(fieldReferenceNode.Open); info.TrailingComments().Len() > 0 {
					f.writeInlineComments(info.TrailingComments())
				}
				f.writeInline(fieldReferenceNode.Name)
			} else {
				f.writeNode(fieldReferenceNode.Name)
				if info := f.fileNode.NodeInfo(fieldReferenceNode.Name); info.TrailingComments().Len() > 0 {
					f.writeInlineComments(info.TrailingComments())
				}
			}
			if fieldReferenceNode.Close != nil {
				f.writeInline(fieldReferenceNode.Close)
			}
			continue
		}
		if i > 0 {
			// The length of this slice must be exactly len(Parts)-1.
			f.writeInline(optionNameNode.Dots[i-1])
		}
		f.writeNode(optionNameNode.Parts[i])
	}
}

// writeMessage writes the message node.
//
// For example,
//
//	message Foo {
//	  option deprecated = true;
//	  reserved 50 to 100;
//	  extensions 150 to 200;
//
//	  message Bar {
//	    string name = 1;
//	  }
//	  enum Baz {
//	    BAZ_UNSPECIFIED = 0;
//	  }
//	  extend Bar {
//	    string value = 2;
//	  }
//
//	  Bar bar = 1;
//	  Baz baz = 2;
//	}
func (f *formatter) writeMessage(messageNode *ast.MessageNode) {
	var elementWriterFunc func()
	if len(messageNode.Decls) != 0 {
		elementWriterFunc = func() {
			for _, decl := range messageNode.Decls {
				f.writeNode(decl)
			}
		}
	}
	f.writeStart(messageNode.Keyword)
	f.Space()
	f.writeInline(messageNode.Name)
	f.Space()
	f.writeCompositeTypeBody(
		messageNode.OpenBrace,
		messageNode.CloseBrace,
		elementWriterFunc,
	)
}

// writeMessageLiteral writes a message literal.
//
// For example,
//
//	{
//	  foo: 1
//	  foo: 2
//	  foo: 3
//	  bar: <
//	    name:"abc"
//	    id:123
//	  >
//	}
func (f *formatter) writeMessageLiteral(messageLiteralNode *ast.MessageLiteralNode) {
	if f.maybeWriteCompactMessageLiteral(messageLiteralNode, false) {
		return
	}
	var elementWriterFunc func()
	if len(messageLiteralNode.Elements) > 0 {
		elementWriterFunc = func() {
			f.writeMessageLiteralElements(messageLiteralNode)
		}
	}
	f.writeCompositeValueBody(
		messageLiteralNode.Open,
		messageLiteralNode.Close,
		elementWriterFunc,
	)
}

// writeMessageLiteral writes a message literal suitable for
// an element in an array literal.
func (f *formatter) writeMessageLiteralForArray(
	messageLiteralNode *ast.MessageLiteralNode,
	lastElement bool,
) {
	if f.maybeWriteCompactMessageLiteral(messageLiteralNode, true) {
		return
	}
	var elementWriterFunc func()
	if len(messageLiteralNode.Elements) > 0 {
		elementWriterFunc = func() {
			f.writeMessageLiteralElements(messageLiteralNode)
		}
	}
	closeWriter := f.writeBodyEndInline
	if lastElement {
		closeWriter = f.writeBodyEnd
	}
	f.writeBody(
		messageLiteralNode.Open,
		messageLiteralNode.Close,
		elementWriterFunc,
		f.writeOpenBracePrefixForArray,
		closeWriter,
	)
}

func (f *formatter) maybeWriteCompactMessageLiteral(
	messageLiteralNode *ast.MessageLiteralNode,
	inArrayLiteral bool,
) bool {
	if len(messageLiteralNode.Elements) == 0 || len(messageLiteralNode.Elements) > 1 ||
		f.hasInteriorComments(messageLiteralNode.Children()...) ||
		messageLiteralHasNestedMessageOrArray(messageLiteralNode) {
		return false
	}
	// messages with a single scalar field and no comments can be
	// printed all on one line
	if inArrayLiteral {
		f.Indent(messageLiteralNode.Open)
	}
	f.writeInline(messageLiteralNode.Open)
	fieldNode := messageLiteralNode.Elements[0]
	f.writeInline(fieldNode.Name)
	if fieldNode.Sep != nil {
		f.writeInline(fieldNode.Sep)
	}
	f.Space()
	f.writeInline(fieldNode.Val)
	f.writeInline(messageLiteralNode.Close)
	return true
}

func messageLiteralHasNestedMessageOrArray(messageLiteralNode *ast.MessageLiteralNode) bool {
	for _, elem := range messageLiteralNode.Elements {
		switch elem.Val.(type) {
		case *ast.ArrayLiteralNode, *ast.MessageLiteralNode:
			return true
		}
	}
	return false
}

func arrayLiteralHasNestedMessageOrArray(arrayLiteralNode *ast.ArrayLiteralNode) bool {
	for _, elem := range arrayLiteralNode.Elements {
		switch elem.(type) {
		case *ast.ArrayLiteralNode, *ast.MessageLiteralNode:
			return true
		}
	}
	return false
}

// writeMessageLiteralElements writes the message literal's elements.
//
// For example,
//
//	foo: 1
//	foo: 2
func (f *formatter) writeMessageLiteralElements(messageLiteralNode *ast.MessageLiteralNode) {
	for i := 0; i < len(messageLiteralNode.Elements); i++ {
		if sep := messageLiteralNode.Seps[i]; sep != nil {
			f.writeMessageFieldWithSeparator(messageLiteralNode.Elements[i])
			f.writeLineEnd(messageLiteralNode.Seps[i])
			continue
		}
		f.writeNode(messageLiteralNode.Elements[i])
	}
}

// writeMessageField writes the message field node, and concludes the
// line without leaving room for a trailing separator in the parent
// message literal.
func (f *formatter) writeMessageField(messageFieldNode *ast.MessageFieldNode) {
	f.writeMessageFieldPrefix(messageFieldNode)
	if compoundStringLiteral, ok := messageFieldNode.Val.(*ast.CompoundStringLiteralNode); ok {
		f.writeCompoundStringLiteralIndent(compoundStringLiteral)
		return
	}
	f.writeLineEnd(messageFieldNode.Val)
}

// writeMessageFieldWithSeparator writes the message field node,
// but leaves room for a trailing separator in the parent message
// literal.
func (f *formatter) writeMessageFieldWithSeparator(messageFieldNode *ast.MessageFieldNode) {
	f.writeMessageFieldPrefix(messageFieldNode)
	if compoundStringLiteral, ok := messageFieldNode.Val.(*ast.CompoundStringLiteralNode); ok {
		f.writeCompoundStringLiteralIndentEndInline(compoundStringLiteral)
		return
	}
	f.writeInline(messageFieldNode.Val)
}

// writeMessageFieldPrefix writes the message field node as a single line.
//
// For example,
//
//	foo:"bar"
func (f *formatter) writeMessageFieldPrefix(messageFieldNode *ast.MessageFieldNode) {
	// The comments need to be written as a multiline comment above
	// the message field name.
	//
	// Note that this is different than how field reference nodes are
	// normally formatted in-line (i.e. as option name components).
	fieldReferenceNode := messageFieldNode.Name
	if fieldReferenceNode.Open != nil {
		f.writeStart(fieldReferenceNode.Open)
		f.writeInline(fieldReferenceNode.Name)
	} else {
		f.writeSt