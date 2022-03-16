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

package bufstudioagent

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/textproto"
	"net/url"

	studiov1alpha1 "github.com/bufbuild/buf/private/gen/proto/go/buf/alpha/studio/v1alpha1"
	"github.com/bufbuild/buf/private/pkg/protoencoding"
	"github.com/bufbuild/connect-go"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"google.golang.org/protobuf/proto"
)

// MaxMessageSizeBytesDefault determines the maximum number of bytes to read
// from the request body.
const MaxMessageSizeBytesDefault = 1024 * 1024 * 5

// plainPostHandler implements a POST handler for forwarding requests that can
// be called with simple CORS requests.
//
// Simple CORS requests are limited [1] to certain headers and content types, so
// this handler expects base64 encoded protobuf messages in the body and writes
// out base64 encoded protobuf messages to be able to use Content-Type: text/plain.
//
// Because of the content-type restriction we do not define a protobuf service
// that gets served by