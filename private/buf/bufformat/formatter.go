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
		f.writeStart(fieldReferenceNode.Name)
	}
	if fieldReferenceNode.Close != nil {
		f.writeInline(fieldReferenceNode.Close)
	}
	if messageFieldNode.Sep != nil {
		f.writeInline(messageFieldNode.Sep)
	}
	f.Space()
}

// writeEnum writes the enum node.
//
// For example,
//
//	enum Foo {
//	  option deprecated = true;
//	  reserved 1 to 5;
//
//	  FOO_UNSPECIFIED = 0;
//	}
func (f *formatter) writeEnum(enumNode *ast.EnumNode) {
	var elementWriterFunc func()
	if len(enumNode.Decls) > 0 {
		elementWriterFunc = func() {
			for _, decl := range enumNode.Decls {
				f.writeNode(decl)
			}
		}
	}
	f.writeStart(enumNode.Keyword)
	f.Space()
	f.writeInline(enumNode.Name)
	f.Space()
	f.writeCompositeTypeBody(
		enumNode.OpenBrace,
		enumNode.CloseBrace,
		elementWriterFunc,
	)
}

// writeEnumValue writes the enum value as a single line. If the enum has
// compact options, it will be written across multiple lines.
//
// For example,
//
//	FOO_UNSPECIFIED = 1 [
//	  deprecated = true
//	];
func (f *formatter) writeEnumValue(enumValueNode *ast.EnumValueNode) {
	f.writeStart(enumValueNode.Name)
	f.Space()
	f.writeInline(enumValueNode.Equals)
	f.Space()
	f.writeInline(enumValueNode.Number)
	if enumValueNode.Options != nil {
		f.Space()
		f.writeNode(enumValueNode.Options)
	}
	f.writeLineEnd(enumValueNode.Semicolon)
}

// writeField writes the field node as a single line. If the field has
// compact options, it will be written across multiple lines.
//
// For example,
//
//	repeated string name = 1 [
//	  deprecated = true,
//	  json_name = "name"
//	];
func (f *formatter) writeField(fieldNode *ast.FieldNode) {
	// We need to handle the comments for the field label specially since
	// a label might not be defined, but it has the leading comments attached
	// to it.
	if fieldNode.Label.KeywordNode != nil {
		f.writeStart(fieldNode.Label)
		f.Space()
		f.writeInline(fieldNode.FldType)
	} else {
		// If a label was not written, the multiline comments will be
		// attached to the type.
		if compoundIdentNode, ok := fieldNode.FldType.(*ast.CompoundIdentNode); ok {
			f.writeCompountIdentForFieldName(compoundIdentNode)
		} else {
			f.writeStart(fieldNode.FldType)
		}
	}
	f.Space()
	f.writeInline(fieldNode.Name)
	f.Space()
	f.writeInline(fieldNode.Equals)
	f.Space()
	f.writeInline(fieldNode.Tag)
	if fieldNode.Options != nil {
		f.Space()
		f.writeNode(fieldNode.Options)
	}
	f.writeLineEnd(fieldNode.Semicolon)
}

// writeMapField writes a map field (e.g. 'map<string, string> pairs = 1;').
func (f *formatter) writeMapField(mapFieldNode *ast.MapFieldNode) {
	f.writeNode(mapFieldNode.MapType)
	f.Space()
	f.writeInline(mapFieldNode.Name)
	f.Space()
	f.writeInline(mapFieldNode.Equals)
	f.Space()
	f.writeInline(mapFieldNode.Tag)
	if mapFieldNode.Options != nil {
		f.Space()
		f.writeNode(mapFieldNode.Options)
	}
	f.writeLineEnd(mapFieldNode.Semicolon)
}

// writeMapType writes a map type (e.g. 'map<string, string>').
func (f *formatter) writeMapType(mapTypeNode *ast.MapTypeNode) {
	f.writeStart(mapTypeNode.Keyword)
	f.writeInline(mapTypeNode.OpenAngle)
	f.writeInline(mapTypeNode.KeyType)
	f.writeInline(mapTypeNode.Comma)
	f.Space()
	f.writeInline(mapTypeNode.ValueType)
	f.writeInline(mapTypeNode.CloseAngle)
}

// writeFieldReference writes a field reference (e.g. '(foo.bar)').
func (f *formatter) writeFieldReference(fieldReferenceNode *ast.FieldReferenceNode) {
	if fieldReferenceNode.Open != nil {
		f.writeInline(fieldReferenceNode.Open)
	}
	f.writeInline(fieldReferenceNode.Name)
	if fieldReferenceNode.Close != nil {
		f.writeInline(fieldReferenceNode.Close)
	}
}

// writeExtend writes the extend node.
//
// For example,
//
//	extend google.protobuf.FieldOptions {
//	  bool redacted = 33333;
//	}
func (f *formatter) writeExtend(extendNode *ast.ExtendNode) {
	var elementWriterFunc func()
	if len(extendNode.Decls) > 0 {
		elementWriterFunc = func() {
			for _, decl := range extendNode.Decls {
				f.writeNode(decl)
			}
		}
	}
	f.writeStart(extendNode.Keyword)
	f.Space()
	f.writeInline(extendNode.Extendee)
	f.Space()
	f.writeCompositeTypeBody(
		extendNode.OpenBrace,
		extendNode.CloseBrace,
		elementWriterFunc,
	)
}

// writeService writes the service node.
//
// For example,
//
//	service FooService {
//	  option deprecated = true;
//
//	  rpc Foo(FooRequest) returns (FooResponse) {};
func (f *formatter) writeService(serviceNode *ast.ServiceNode) {
	var elementWriterFunc func()
	if len(serviceNode.Decls) > 0 {
		elementWriterFunc = func() {
			for _, decl := range serviceNode.Decls {
				f.writeNode(decl)
			}
		}
	}
	f.writeStart(serviceNode.Keyword)
	f.Space()
	f.writeInline(serviceNode.Name)
	f.Space()
	f.writeCompositeTypeBody(
		serviceNode.OpenBrace,
		serviceNode.CloseBrace,
		elementWriterFunc,
	)
}

// writeRPC writes the RPC node. RPCs are formatted in
// the following order:
//
// For example,
//
//	rpc Foo(FooRequest) returns (FooResponse) {
//	  option deprecated = true;
//	};
func (f *formatter) writeRPC(rpcNode *ast.RPCNode) {
	var elementWriterFunc func()
	if len(rpcNode.Decls) > 0 {
		elementWriterFunc = func() {
			for _, decl := range rpcNode.Decls {
				f.writeNode(decl)
			}
		}
	}
	f.writeStart(rpcNode.Keyword)
	f.Space()
	f.writeInline(rpcNode.Name)
	f.writeInline(rpcNode.Input)
	f.Space()
	f.writeInline(rpcNode.Returns)
	f.Space()
	f.writeInline(rpcNode.Output)
	if rpcNode.OpenBrace == nil {
		// This RPC doesn't have any elements, so we prefer the
		// ';' form.
		//
		//  rpc Ping(PingRequest) returns (PingResponse);
		//
		f.writeLineEnd(rpcNode.Semicolon)
		return
	}
	f.Space()
	f.writeCompositeTypeBody(
		rpcNode.OpenBrace,
		rpcNode.CloseBrace,
		elementWriterFunc,
	)
}

// writeRPCType writes the RPC type node (e.g. (stream foo.Bar)).
func (f *formatter) writeRPCType(rpcTypeNode *ast.RPCTypeNode) {
	f.writeInline(rpcTypeNode.OpenParen)
	if rpcTypeNode.Stream != nil {
		f.writeInline(rpcTypeNode.Stream)
		f.Space()
	}
	f.writeInline(rpcTypeNode.MessageType)
	f.writeInline(rpcTypeNode.CloseParen)
}

// writeOneOf writes the oneof node.
//
// For example,
//
//	oneof foo {
//	  option deprecated = true;
//
//	  string name = 1;
//	  int number = 2;
//	}
func (f *formatter) writeOneOf(oneOfNode *ast.OneOfNode) {
	var elementWriterFunc func()
	if len(oneOfNode.Decls) > 0 {
		elementWriterFunc = func() {
			for _, decl := range oneOfNode.Decls {
				f.writeNode(decl)
			}
		}
	}
	f.writeStart(oneOfNode.Keyword)
	f.Space()
	f.writeInline(oneOfNode.Name)
	f.Space()
	f.writeCompositeTypeBody(
		oneOfNode.OpenBrace,
		oneOfNode.CloseBrace,
		elementWriterFunc,
	)
}

// writeGroup writes the group node.
//
// For example,
//
//	optional group Key = 4 [
//	  deprecated = true,
//	  json_name = "key"
//	] {
//	  optional uint64 id = 1;
//	  optional string name = 2;
//	}
func (f *formatter) writeGroup(groupNode *ast.GroupNode) {
	var elementWriterFunc func()
	if len(groupNode.Decls) > 0 {
		elementWriterFunc = func() {
			for _, decl := range groupNode.Decls {
				f.writeNode(decl)
			}
		}
	}
	// We need to handle the comments for the group label specially since
	// a label might not be defined, but it has the leading comments attached
	// to it.
	if groupNode.Label.KeywordNode != nil {
		f.writeStart(groupNode.Label)
		f.Space()
		f.writeInline(groupNode.Keyword)
	} else {
		// If a label was not written, the multiline comments will be
		// attached to the keyword.
		f.writeStart(groupNode.Keyword)
	}
	f.Space()
	f.writeInline(groupNode.Name)
	f.Space()
	f.writeInline(groupNode.Equals)
	f.Space()
	f.writeInline(groupNode.Tag)
	if groupNode.Options != nil {
		f.Space()
		f.writeNode(groupNode.Options)
	}
	f.Space()
	f.writeCompositeTypeBody(
		groupNode.OpenBrace,
		groupNode.CloseBrace,
		elementWriterFunc,
	)
}

// writeExtensionRange writes the extension range node.
//
// For example,
//
//	extensions 5-10, 100 to max [
//	  deprecated = true
//	];
func (f *formatter) writeExtensionRange(extensionRangeNode *ast.ExtensionRangeNode) {
	f.writeStart(extensionRangeNode.Keyword)
	f.Space()
	for i := 0; i < len(extensionRangeNode.Ranges); i++ {
		if i > 0 {
			// The length of this slice must be exactly len(Ranges)-1.
			f.writeInline(extensionRangeNode.Commas[i-1])
			f.Space()
		}
		f.writeNode(extensionRangeNode.Ranges[i])
	}
	if extensionRangeNode.Options != nil {
		f.Space()
		f.writeNode(extensionRangeNode.Options)
	}
	f.writeLineEnd(extensionRangeNode.Semicolon)
}

// writeReserved writes a reserved node.
//
// For example,
//
//	reserved 5-10, 100 to max;
func (f *formatter) writeReserved(reservedNode *ast.ReservedNode) {
	f.writeStart(reservedNode.Keyword)
	// Either names or ranges will be set, but never both.
	elements := make([]ast.Node, 0, len(reservedNode.Names)+len(reservedNode.Ranges))
	switch {
	case reservedNode.Names != nil:
		for _, nameNode := range reservedNode.Names {
			elements = append(elements, nameNode)
		}
	case reservedNode.Ranges != nil:
		for _, rangeNode := range reservedNode.Ranges {
			elements = append(elements, rangeNode)
		}
	}
	f.Space()
	for i := 0; i < len(elements); i++ {
		if i > 0 {
			// The length of this slice must be exactly len({Names,Ranges})-1.
			f.writeInline(reservedNode.Commas[i-1])
			f.Space()
		}
		f.writeInline(elements[i])
	}
	f.writeLineEnd(reservedNode.Semicolon)
}

// writeRange writes the given range node (e.g. '1 to max').
func (f *formatter) writeRange(rangeNode *ast.RangeNode) {
	f.writeInline(rangeNode.StartVal)
	if rangeNode.To != nil {
		f.Space()
		f.writeInline(rangeNode.To)
	}
	// Either EndVal or Max will be set, but never both.
	switch {
	case rangeNode.EndVal != nil:
		f.Space()
		f.writeInline(rangeNode.EndVal)
	case rangeNode.Max != nil:
		f.Space()
		f.writeInline(rangeNode.Max)
	}
}

// writeCompactOptions writes a compact options node.
//
// For example,
//
//	[
//	  deprecated = true,
//	  json_name = "something"
//	]
func (f *formatter) writeCompactOptions(compactOptionsNode *ast.CompactOptionsNode) {
	f.inCompactOptions = true
	defer func() {
		f.inCompactOptions = false
	}()
	if len(compactOptionsNode.Options) == 1 &&
		!f.hasInteriorComments(compactOptionsNode.OpenBracket, compactOptionsNode.Options[0].Name) {
		// If there's only a single compact scalar option without comments, we can write it
		// in-line. For example:
		//
		//  [deprecated = true]
		//
		// However, this does not include the case when the '[' has trailing comments,
		// or the option name has leading comments. In those cases, we write the option
		// across multiple lines. For example:
		//
		//  [
		//    // This type is deprecated.
		//    deprecated = true
		//  ]
		//
		optionNode := compactOptionsNode.Options[0]
		f.writeInline(compactOptionsNode.OpenBracket)
		f.writeInline(optionNode.Name)
		f.Space()
		f.writeInline(optionNode.Equals)
		if node, ok := optionNode.Val.(*ast.CompoundStringLiteralNode); ok {
			// If there's only a single compact option, the value needs to
			// write its comments (if any) in a way that preserves the closing ']'.
			f.writeCompoundStringLiteralNoIndentEndInline(node)
			f.writeInline(compactOptionsNode.CloseBracket)
			return
		}
		f.Space()
		f.writeInline(optionNode.Val)
		f.writeInline(compactOptionsNode.CloseBracket)
		return
	}
	var elementWriterFunc func()
	if len(compactOptionsNode.Options) > 0 {
		elementWriterFunc = func() {
			for i, opt := range compactOptionsNode.Options {
				if i == len(compactOptionsNode.Options)-1 {
					// The last element won't have a trailing comma.
					f.writeLastCompactOption(opt)
					return
				}
				f.writeNode(opt)
				f.writeLineEnd(compactOptionsNode.Commas[i])
			}
		}
	}
	f.writeCompositeValueBody(
		compactOptionsNode.OpenBracket,
		compactOptionsNode.CloseBracket,
		elementWriterFunc,
	)
}

func (f *formatter) hasInteriorComments(nodes ...ast.Node) bool {
	for i, n := range nodes {
		// interior comments mean we ignore leading comments on first
		// token and trailing comments on the last one
		info := f.fileNode.NodeInfo(n)
		if i > 0 && info.LeadingComments().Len() > 0 {
			return true
		}
		if i < len(nodes)-1 && info.TrailingComments().Len() > 0 {
			return true
		}
	}
	return false
}

// writeArrayLiteral writes an array literal across multiple lines.
//
// For example,
//
//	[
//	  "foo",
//	  "bar"
//	]
func (f *formatter) writeArrayLiteral(arrayLiteralNode *ast.ArrayLiteralNode) {
	if len(arrayLiteralNode.Elements) == 1 &&
		!f.hasInteriorComments(arrayLiteralNode.Children()...) &&
		!arrayLiteralHasNestedMessageOrArray(arrayLiteralNode) {
		// arrays with a single scalar value and no comments can be
		// printed all on one line
		valueNode := arrayLiteralNode.Elements[0]
		f.writeInline(arrayLiteralNode.OpenBracket)
		f.writeInline(valueNode)
		f.writeInline(arrayLiteralNode.CloseBracket)
		return
	}

	var elementWriterFunc func()
	if len(arrayLiteralNode.Elements) > 0 {
		elementWriterFunc = func() {
			for i := 0; i < len(arrayLiteralNode.Elements); i++ {
				lastElement := i == len(arrayLiteralNode.Elements)-1
				if compositeNode, ok := arrayLiteralNode.Elements[i].(ast.CompositeNode); ok {
					f.writeCompositeValueForArrayLiteral(compositeNode, lastElement)
					if !lastElement {
						f.writeLineEnd(arrayLiteralNode.Commas[i])
					}
					continue
				}
				if lastElement {
					// The last element won't have a trailing comma.
					f.writeLineElement(arrayLiteralNode.Elements[i])
					return
				}
				f.writeStart(arrayLiteralNode.Elements[i])
				f.writeLineEnd(arrayLiteralNode.Commas[i])
			}
		}
	}
	f.writeCompositeValueBody(
		arrayLiteralNode.OpenBracket,
		arrayLiteralNode.CloseBracket,
		elementWriterFunc,
	)
}

// writeCompositeForArrayLiteral writes the composite node in a way that's suitable
// for array literals. In general, signed integers and compound strings should have their
// comments written in-line because they are one of many components in a single line.
//
// However, each of these composite types occupy a single line in an array literal,
// so they need their comments to be formatted like a standalone node.
//
// For example,
//
//	option (value) = /* In-line comment for '-42' */ -42;
//
//	option (thing) = {
//	  values: [
//	    // Leading comment on -42.
//	    -42, // Trailing comment on -42.
//	  ]
//	}
//
// The lastElement boolean is used to signal whether or not the composite value
// should be written as the last element (i.e. it doesn't have a trailing comma).
func (f *formatter) writeCompositeValueForArrayLiteral(
	compositeNode ast.CompositeNode,
	lastElement bool,
) {
	switch node := compositeNode.(type) {
	case *ast.CompoundStringLiteralNode:
		f.writeCompoundStringLiteralForArray(node, lastElement)
	case *ast.PositiveUintLiteralNode:
		f.writePositiveUintLiteralForArray(node, lastElement)
	case *ast.NegativeIntLiteralNode:
		f.writeNegativeIntLiteralForArray(node, lastElement)
	case *ast.SignedFloatLiteralNode:
		f.writeSignedFloatLiteralForArray(node, lastElement)
	case *ast.MessageLiteralNode:
		f.writeMessageLiteralForArray(node, lastElement)
	default:
		f.err = multierr.Append(f.err, fmt.Errorf("unexpected array value node %T", node))
	}
}

// writeCompositeTypeBody writes the body of a composite type, e.g. message, enum, extend, oneof, etc.
func (f *formatter) writeCompositeTypeBody(
	openBrace *ast.RuneNode,
	closeBrace *ast.RuneNode,
	elementWriterFunc func(),
) {
	f.writeBody(
		openBrace,
		closeBrace,
		elementWriterFunc,
		f.writeOpenBracePrefix,
		f.writeBodyEnd,
	)
}

// writeCompositeValueBody writes the body of a composite value, e.g. compact options,
// array literal, etc. We need to handle the ']' different than composite types because
// there could be more tokens following the final ']'.
func (f *formatter) writeCompositeValueBody(
	openBrace *ast.RuneNode,
	closeBrace *ast.RuneNode,
	elementWriterFunc func(),
) {
	f.writeBody(
		openBrace,
		closeBrace,
		elementWriterFunc,
		f.writeOpenBracePrefix,
		f.writeBodyEndInline,
	)
}

// writeBody writes the body of a type or value, e.g. message, enum, compact options, etc.
// The elementWriterFunc is used to write the declarations within the composite type (e.g.
// fields in a message). The openBraceWriterFunc and closeBraceWriterFunc functions are used
// to customize how the '{' and '} nodes are written, respectively.
func (f *formatter) writeBody(
	openBrace *ast.RuneNode,
	closeBrace *ast.RuneNode,
	elementWriterFunc func(),
	openBraceWriterFunc func(ast.Node),
	closeBraceWriterFunc func(ast.Node, bool),
) {
	if elementWriterFunc == nil && !f.hasInteriorComments(openBrace, closeBrace) {
		// completely empty body
		f.writeInline(openBrace)
		closeBraceWriterFunc(closeBrace, true)
		return
	}

	openBraceWriterFunc(openBrace)
	if elementWriterFunc != nil {
		elementWriterFunc()
	}
	closeBraceWriterFunc(closeBrace, false)
}

// writeOpenBracePrefix writes the open brace with its leading comments in-line.
// This is used for nearly every use case of f.writeBody, excluding the instances
// in array literals.
func (f *formatter) writeOpenBracePrefix(openBrace ast.Node) {
	defer f.SetPreviousNode(openBrace)
	info := f.fileNode.NodeInfo(openBrace)
	if info.LeadingComments().Len() > 0 {
		f.writeInlineComments(info.LeadingComments())
		if info.LeadingWhitespace() != "" {
			f.Space()
		}
	}
	f.writeNode(openBrace)
	if info.TrailingComments().Len() > 0 {
		f.writeTrailingEndComments(info.TrailingComments())
	} else {
		f.P()
	}
}

// writeOpenBracePrefixForArray writes the open brace with its leading comments
// on multiple lines. This is only used for message literals in arrays.
func (f *formatter) writeOpenBracePrefixForArray(openBrace ast.Node) {
	defer f.SetPreviousNode(openBrace)
	info := f.fileNode.NodeInfo(openBrace)
	if info.LeadingComments().Len() > 0 {
		f.writeMultilineComments(info.LeadingComments())
	}
	f.Indent(openBrace)
	f.writeNode(openBrace)
	if info.TrailingComments().Len() > 0 {
		f.writeTrailingEndComments(info.TrailingComments())
	} else {
		f.P()
	}
}

// writeCompoundIdent writes a compound identifier (e.g. '.com.foo.Bar').
func (f *formatter) writeCompoundIdent(compoundIdentNode *ast.CompoundIdentNode) {
	if compoundIdentNode.LeadingDot != nil {
		f.writeInline(compoundIdentNode.LeadingDot)
	}
	for i := 0; i < len(compoundIdentNode.Components); i++ {
		if i > 0 {
			// The length of this slice must be exactly len(Components)-1.
			f.writeInline(compoundIdentNode.Dots[i-1])
		}
		f.writeInline(compoundIdentNode.Components[i])
	}
}

// writeCompountIdentForFieldName writes a compound identifier, but handles comments
// specially for field names.
//
// For example,
//
//	message Foo {
//	  // These are comments attached to bar.
//	  bar.v1.Bar bar = 1;
//	}
func (f *formatter) writeCompountIdentForFieldName(compoundIdentNode *ast.CompoundIdentNode) {
	if compoundIdentNode.LeadingDot != nil {
		f.writeStart(compoundIdentNode.LeadingDot)
	}
	for i := 0; i < len(compoundIdentNode.Components); i++ {
		if i == 0 && compoundIdentNode.LeadingDot == nil {
			f.writeStart(compoundIdentNode.Components[i])
			continue
		}
		if i > 0 {
			// The length of this slice must be exactly len(Components)-1.
			f.writeInline(compoundIdentNode.Dots[i-1])
		}
		f.writeInline(compoundIdentNode.Components[i])
	}
}

// writeFieldLabel writes the field label node.
//
// For example,
//
//	optional
//	repeated
//	required
func (f *formatter) writeFieldLabel(fieldLabel ast.FieldLabel) {
	f.WriteString(fieldLabel.Val)
}

// writeCompoundStringLiteral writes a compound string literal value.
//
// For example,
//
//	"one,"
//	"two,"
//	"three"
func (f *formatter) writeCompoundStringLiteral(
	compoundStringLiteralNode *ast.CompoundStringLiteralNode,
	needsIndent bool,
	hasTrailingPunctuation bool,
) {
	f.P()
	if needsIndent {
		f.In()
	}
	for i, child := range compoundStringLiteralNode.Children() {
		if hasTrailingPunctuation && i == len(compoundStringLiteralNode.Children())-1 {
			// inline because there may be a subsequent comma or punctuation from enclosing element
			f.writeStart(child)
			break
		}
		f.writeLineElement(child)
	}
	if needsIndent {
		f.Out()
	}
}

func (f *formatter) writeCompoundStringLiteralIndent(
	compoundStringLiteralNode *ast.CompoundStringLiteralNode,
) {
	f.writeCompoundStringLiteral(compoundStringLiteralNode, true, false)
}

func (f *formatter) writeCompoundStringLiteralIndentEndInline(
	compoundStringLiteralNode *ast.CompoundStringLiteralNode,
) {
	f.writeCompoundStringLiteral(compoundStringLiteralNode, true, true)
}

func (f *formatter) writeCompoundStringLiteralNoIndentEndInline(
	compoundStringLiteralNode *ast.CompoundStringLiteralNode,
) {
	f.writeCompoundStringLiteral(compoundStringLiteralNode, false, true)
}

// writeCompoundStringLiteralForArray writes a compound string literal value,
// but writes its comments suitable for an element in an array literal.
//
// The lastElement boolean is used to signal whether or not the value should
// be written as the last element (i.e. it doesn't have a trailing comma).
func (f *formatter) writeCompoundStringLiteralForArray(
	compoundStringLiteralNode *ast.CompoundStringLiteralNode,
	lastElement bool,
) {
	for i, child := range compoundStringLiteralNode.Children() {
		if !lastElement && i == len(compoundStringLiteralNode.Children())-1 {
			f.writeStart(child)
			return
		}
		f.writeLineElement(child)
	}
}

// writeFloatLiteral writes a float literal value (e.g. '42.2').
func (f *formatter) writeFloatLiteral(floatLiteralNode *ast.FloatLiteralNode) {
	f.writeRaw(floatLiteralNode)
}

// writeSignedFloatLiteral writes a signed float literal value (e.g. '-42.2').
func (f *formatter) writeSignedFloatLiteral(signedFloatLiteralNode *ast.SignedFloatLiteralNode) {
	f.writeInline(signedFloatLiteralNode.Sign)
	f.writeInline(signedFloatLiteralNode.Float)
}

// writeSignedFloatLiteralForArray writes a signed float literal value, but writes
// its comments suitable for an element in an array literal.
//
// The lastElement boolean is used to signal whether or not the value should
// be written as the last element (i.e. it doesn't have a trailing comma).
func (f *formatter) writeSignedFloatLiteralForArray(
	signedFloatLiteralNode *ast.SignedFloatLiteralNode,
	lastElement bool,
) {
	f.writeStart(signedFloatLiteralNode.Sign)
	if lastElement {
		f.writeLineEnd(signedFloatLiteralNode.Float)
		return
	}
	f.writeInline(signedFloatLiteralNode.Float)
}

// writeSpecialFloatLiteral writes a special float literal value (e.g. "nan" or "inf").
func (f *formatter) writeSpecialFloatLiteral(specialFloatLiteralNode *ast.SpecialFloatLiteralNode) {
	f.WriteString(specialFloatLiteralNode.KeywordNode.Val)
}

// writeStringLiteral writes a string literal value (e.g. "foo").
// Note that the raw string is written as-is so that it preserves
// the quote style used in the original source.
func (f *formatter) writeStringLiteral(stringLiteralNode *ast.StringLiteralNode) {
	f.writeRaw(stringLiteralNode)
}

// writeUintLiteral writes a uint literal (e.g. '42').
func (f *formatter) writeUintLiteral(uintLiteralNode *ast.UintLiteralNode) {
	f.writeRaw(uintLiteralNode)
}

// writeNegativeIntLiteral writes a negative int literal (e.g. '-42').
func (f *formatter) writeNegativeIntLiteral(negativeIntLiteralNode *ast.NegativeIntLiteralNode) {
	f.writeInline(negativeIntLiteralNode.Minus)
	f.writeInline(negativeIntLiteralNode.Uint)
}

func (f *formatter) writeRaw(n ast.Node) {
	info := f.fileNode.NodeInfo(n)
	f.WriteString(info.RawText())
}

// writeNegativeIntLiteralForArray writes a negative int literal value, but writes
// its comments suitable for an element in an array literal.
//
// The lastElement boolean is used to signal whether or not the value should
// be written as the last element (i.e. it doesn't have a trailing comma).
func (f *formatter) writeNegativeIntLiteralForArray(
	negativeIntLiteralNode *ast.NegativeIntLiteralNode,
	lastElement bool,
) {
	f.writeStart(negativeIntLiteralNode.Minus)
	if lastElement {
		f.writeLineEnd(negativeIntLiteralNode.Uint)
		return
	}
	f.writeInline(negativeIntLiteralNode.Uint)
}

// writePositiveUintLiteral writes a positive uint literal (e.g. '+42').
func (f *formatter) writePositiveUintLiteral(positiveIntLiteralNode *ast.PositiveUintLiteralNode) {
	f.writeInline(positiveIntLiteralNode.Plus)
	f.writeInline(positiveIntLiteralNode.Uint)
}

// writePositiveUintLiteralForArray writes a positive uint literal value, but writes
// its comments suitable for an element in an array literal.
//
// The lastElement boolean is used to signal whether or not the value should
// be written as the last element (i.e. it doesn't have a trailing comma).
func (f *formatter) writePositiveUintLiteralForArray(
	positiveIntLiteralNode *ast.PositiveUintLiteralNode,
	lastElement bool,
) {
	f.writeStart(positiveIntLiteralNode.Plus)
	if lastElement {
		f.writeLineEnd(positiveIntLiteralNode.Uint)
		return
	}
	f.writeInline(positiveIntLiteralNode.Uint)
}

// writeIdent writes an identifier (e.g. 'foo').
func (f *formatter) writeIdent(identNode *ast.IdentNode) {
	f.WriteString(identNode.Val)
}

// writeKeyword writes a keyword (e.g. 'syntax').
func (f *formatter) writeKeyword(keywordNode *ast.KeywordNode) {
	f.WriteString(keywordNode.Val)
}

// writeRune writes a rune (e.g. '=').
func (f *formatter) writeRune(runeNode *ast.RuneNode) {
	if strings.ContainsRune("{[(<", runeNode.Rune) {
		f.pendingIndent++
	} else if strings.ContainsRune("}])>", runeNode.Rune) {
		f.pendingIndent--
	}
	f.WriteString(string(runeNode.Rune))
}

// writeNode writes the node by dispatching to a function tailored to its concrete type.
//
// Comments are handled in each respective write function so that it can determine whether
// to write the comments in-line or not.
func (f *formatter) writeNode(node ast.Node) {
	switch element := node.(type) {
	case *ast.ArrayLiteralNode:
		f.writeArrayLiteral(element)
	case *ast.CompactOptionsNode:
		f.writeCompactOptions(element)
	case *ast.CompoundIdentNode:
		f.writeCompoundIdent(element)
	case *ast.CompoundStringLiteralNode:
		f.writeCompoundStringLiteralIndent(element)
	case *ast.EnumNode:
		f.writeEnum(element)
	case *ast.EnumValueNode:
		f.writeEnumValue(element)
	case *ast.ExtendNode:
		f.writeExtend(element)
	case *ast.ExtensionRangeNode:
		f.writeExtensionRange(element)
	case as