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

// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: buf/alpha/registry/v1alpha1/reference.proto

package registryv1alpha1connect

import (
	context "context"
	errors "errors"
	v1alpha1 "github.com/bufbuild/buf/private/gen/proto/go/buf/alpha/registry/v1alpha1"
	connect_go "github.com/bufbuild/connect-go"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect_go.IsAtLeastVersion0_1_0

const (
	// ReferenceServiceName is the fully-qualified name of the ReferenceService service.
	ReferenceServiceName = "buf.alpha.registry.v1alpha1.ReferenceService"
)

// ReferenceServiceClient is a client for the buf.alpha.registry.v1alpha1.ReferenceService service.
type ReferenceServiceClient interface {
	// GetReferenceByName takes a reference name and returns the
	// reference either as 'main', a tag, or commit.
	GetReferenceByName(context.Context, *connect_go.Request[v1alpha1.GetReferenceByNameRequest]) (*connect_go.Response[v1alpha1.GetReferenceByNameResponse], error)
}

// NewReferenceServiceClient constructs a client for the
// buf.alpha.registry.v1alpha1.ReferenceService service. By default, it uses the Connect protocol
// with the binary Protobuf Codec, asks for gzipped responses, and sends uncompressed requests. To
// use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or connect.WithGRPCWeb()
// options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewReferenceServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) ReferenceServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &referenceServiceClient{
		getReferenceByName: connect_go.NewClient[v1alpha1.GetReferenceByNameRequest, v1alpha1.GetReferenceByNameResponse](
			httpClient,
			baseURL+"/buf.alpha.registry.v1alpha1.ReferenceService/GetReferenceByName",
			opts...,
		),
	}
}

// referenceServiceClient implements ReferenceServiceClient.
type referenceServiceClient struct {
	getReferenceByName *connect_go.Client[v1alpha1.GetReferenceByNameRequest, v1alpha1.GetReferenceByNameResponse]
}

// GetReferenceByName calls buf.alpha.registry.v1alpha1.ReferenceService.GetReferenceByName.
func (c *referenceServiceClient) GetReferenceByName(ctx context.Context, req *connect_go.Request[v1alpha1.GetReferenceByNameRequest]) (*connect_go.Response[v1alpha1.GetReferenceByNameResponse], error) {
	return c.getReferenceByName.CallUnary(ctx, req)
}

// ReferenceServiceHandler is an implementation of the buf.alpha.registry.v1alpha1.ReferenceService
// service.
type ReferenceServiceHandler interface {
	// GetReferenceByName takes a reference name and returns the
	// reference either as 'main', a tag, or commit.
	GetReferenceByName(context.Context, *connect_go.Request[v1alpha1.GetReferenceByNameRequest]) (*connect_go.Response[v1alpha1.GetReferenceByNameResponse], error)
}

// NewReferenceServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewReferenceServiceHandler(svc ReferenceServiceHandler, opts ...connect_go.HandlerOptio