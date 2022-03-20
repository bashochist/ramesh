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
// that gets served by connect but instead use a plain post handler.
//
// [1] https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS#simple_requests).
type plainPostHandler struct {
	Logger              *zap.Logger
	MaxMessageSizeBytes int64
	B64Encoding         *base64.Encoding
	TLSClient           *http.Client
	H2CClient           *http.Client
	DisallowedHeaders   map[string]struct{}
	ForwardHeaders      map[string]string
}

func newPlainPostHandler(
	logger *zap.Logger,
	disallowedHeaders map[string]struct{},
	forwardHeaders map[string]string,
	tlsClientConfig *tls.Config,
) *plainPostHandler {
	canonicalDisallowedHeaders := make(map[string]struct{}, len(disallowedHeaders))
	for k := range disallowedHeaders {
		canonicalDisallowedHeaders[textproto.CanonicalMIMEHeaderKey(k)] = struct{}{}
	}
	canonicalForwardHeaders := make(map[string]string, len(forwardHeaders))
	for k, v := range forwardHeaders {
		canonicalForwardHeaders[textproto.CanonicalMIMEHeaderKey(k)] = v
	}
	return &plainPostHandler{
		B64Encoding:       base64.StdEncoding,
		DisallowedHeaders: canonicalDisallowedHeaders,
		ForwardHeaders:    canonicalForwardHeaders,
		H2CClient: &http.Client{
			Transport: &http2.Transport{
				AllowHTTP: true,
				DialTLS: func(netw, addr string, config *tls.Config) (net.Conn, error) {
					return net.Dial(netw, addr)
				},
			},
		},
		Logger:              logger,
		MaxMessageSizeBytes: MaxMessageSizeBytesDefault,
		TLSClient: &http.Client{
			Transport: &http2.Transport{
				TLSClientConfig: tlsClientConfig,
			},
		},
	}
}

func (i *plainPostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("content-type") != "text/plain" {
		http.Error(w, "", http.StatusUnsupportedMediaType)
		return
	}
	bodyBytes, err := io.ReadAll(
		base64.NewDecoder(
			i.B64Encoding,
			http.MaxBytesReader(w, r.Body, i.MaxMessageSizeBytes),
		),
	)
	if err != nil {
		if b64Err := new(base64.CorruptInputError); errors.As(err, &b64Err) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusRequestEntityTooLarge)
		return
	}
	envelopeRequest := &studiov1alpha1.InvokeRequest{}
	if err := protoencoding.NewWireUnmarshaler(nil).Unmarshal(bodyBytes, envelopeRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	request := connect.NewRequest(bytes.NewBuffer(envelopeRequest.GetBody()))
	for _, header := range envelopeRequest.Headers {
		if _, ok := i.DisallowedHeaders[textproto.CanonicalMIMEHeaderKey(header.Key)]; ok {
			http.Error(w, fmt.Sprintf("header %q disallowed by agent", header.Key), http.StatusBadRequest)
			return
		}
		for _, value := range header.Value {
			request.Header().Add(header.Key, value)
		}
	}
	for fromHeader, toHeader := range i.ForwardHeaders {
		headerValues := r.Header.Values(fromHeader)
		if len(headerValues) > 0 {
			request.Header().Del(toHeader)
			for _, headerValue := range headerValues {
				request.Header().Add(toHeader, headerValue)
			}
		}
	}
	targetURL, err := url.Parse(envelopeRequest.GetTarget())
	if err != nil {
		http.Er