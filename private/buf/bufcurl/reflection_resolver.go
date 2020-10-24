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

package bufcurl

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	reflectionv1 "github.com/bufbuild/buf/private/gen/proto/go/grpc/reflection/v1"
	"github.com/bufbuild/buf/private/pkg/protoencoding"
	"github.com/bufbuild/buf/private/pkg/verbose"
	"github.com/bufbuild/connect-go"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

const (
	// ReflectProtocolUnknown represents that the server reflection protocol
	// is unknown. If given this value, the server reflection resolver will
	// cycle through the known reflection protocols from newest to oldest,
	// trying each one until a reflection protocol that works is found.
	ReflectProtocolUnknown ReflectProtocol = iota + 1
	// ReflectProtocolGRPCV1 represents the gRPC server reflection protocol
	// defined by the service grpc.reflection.v1.ServerReflection.
	ReflectProtocolGRPCV1
	// ReflectProtocolGRPCV1Alpha represents the gRPC server reflection protocol
	// defined by the service grpc.reflection.v1alpha.ServerReflection.
	ReflectProtocolGRPCV1Alpha
)

var (
	// AllKnownReflectProtocolStrings are all string values for
	// ReflectProtocol that represent known reflection protocols.
	AllKnownReflectProtocolStrings = []string{
		"grpc-v1",
		"grpc-v1alpha",
	}

	reflectProtocolToString = map[ReflectProtocol]string{
		ReflectProtocolUnknown:     "",
		ReflectProtocolGRPCV1:      "grpc-v1",
		ReflectProtocolGRPCV1Alpha: "grpc-v1alpha",
	}
	stringToReflectProtocol = map[string]ReflectProtocol{
		"":             ReflectProtocolUnknown,
		"grpc-v1":      ReflectProtocolGRPCV1,
		"grpc-v1alpha": ReflectProtocolGRPCV1Alpha,
	}
)

// ReflectProtocol is a reflection protocol.
type ReflectProtocol int

// String implements fmt.Stringer.
func (r ReflectProtocol) String() string {
	s, ok := reflectProtocolToString[r]
	if !ok {
		return strconv.Itoa(int(r))
	}
	return s
}

// ParseReflectProtocol parses the ReflectProtocol.
//
// The empty string is a parse error.
func ParseReflectProtocol(s string) (ReflectProtocol, error) {
	r, ok := stringToReflectProtocol[strings.ToLower(strings.TrimSpace(s))]
	if ok {
		return r, nil
	}
	return 0, fmt.Errorf("unknown ReflectProtocol: %q", s)
}

// NewServerReflectionResolver creates a new resolver using the given details to
// create an RPC reflection client, to ask the server for descriptors.
func NewServerReflectionResolver(
	ctx context.Context,
	httpClient connect.HTTPClient,
	opts []connect.ClientOption,
	baseURL string,
	reflectProtocol ReflectProtocol,
	headers http.Header,
	printer verbose.Printer,
) (r Resolver, closeResolver func()) {
	baseURL = strings.TrimSuffix(baseURL, "/")
	var v1Client, v1alphaClient *reflectClient
	if reflectProtocol != ReflectProtocolGRPCV1 {
		v1alphaClient = connect.NewClient[reflectionv1.ServerReflectionRequest, reflectionv1.ServerReflectionResponse](httpClient, baseURL+"/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo", opts...)
	}
	if reflectProtocol != ReflectProtocolGRPCV1Alpha {
		v1Client = connect.NewClient[reflectionv1.ServerReflectionRequest, reflectionv1.ServerReflectionResponse](httpClient, baseURL+"/grpc.reflection.v1.ServerReflection/ServerReflectionInfo", opts...)
	}
	// if version is neither "