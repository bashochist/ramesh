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

// Package ioextended provides io utilities.
package ioextended

import (
	"bytes"
	"io"
	"sync"

	"go.uber.org/multierr"
)

var (
	// DiscardReader is an io.Reader in which all calls return 0 and io.EOF.
	DiscardReader io.Reader = discardReader{}
	// DiscardReadCloser is an io.ReadCloser in which all calls return 0 and io.EOF.
	DiscardReadCloser io.ReadCloser = io.NopCloser(DiscardReader)
	// DiscardWriteCloser is a discard io.WriteCloser.
	DiscardWriteCloser io.WriteCloser = NopWriteCloser(io.Discard)
	// NopCloser is a no-op closer.
	NopCloser = nopCloser{}
)

// NopWriteCloser returns an io.WriteCloser with a no-op Close method wrapping the provided io.Writer.
func NopWriteCloser(writer io.Writer) io.WriteCloser {
	return nopWriteCloser{Writer: writer}
}

// LockedWriter creates a locked Writer.
func LockedWriter(writer io.Writer) io.Writer {
	return &locked