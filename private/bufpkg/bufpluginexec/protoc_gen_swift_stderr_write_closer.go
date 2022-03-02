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

package bufpluginexec

import (
	"bytes"
	"errors"
	"io"
)

type protocGenSwiftStderrWriteCloser struct {
	delegate io.Writer
	buffer   *bytes.Buffer
}

func newProtocGenSwiftStderrWriteCloser(delegate io.Writer) io.WriteCloser {
	return &protocGenSwiftStderrWriteCloser{
		delegate: delegate,
		buffer:   bytes.NewBuffer(nil),
	}
}

func (p *protocGenSwiftStderrWriteCloser) Write(data []byte) (int, error) {
	// If protoc-gen-swift, we want to capture all the stderr so we can process it.
	return p.buffer.Write(data)
}

func (p *protocGenSwiftStderrWriteCloser) Close() error {
	data := p.buffer.Bytes()
	if len(data) == 0 {
		return nil
	}
	newData := bytes.ReplaceAll(
		data,
		// If swift-protobuf changes their error message, this may not longer filter properly
		//