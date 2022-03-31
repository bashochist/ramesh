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
// Source: buf/alpha/registry/v1alpha1/organization.proto

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
	// OrganizationServiceName is the fully-qualified name of the OrganizationService service.
	OrganizationServiceName = "buf.alpha.registry.v1alpha1.OrganizationService"
)

// OrganizationServiceClient is a client for the buf.alpha.registry.v1alpha1.OrganizationService
// service.
type OrganizationServiceClient interface {
	// GetOrganization gets a organization by ID.
	GetOrganization(context.Context, *connect_go.Request[v1alpha1.GetOrganizationRequest]) (*connect_go.Response[v1alpha1.GetOrganizationResponse], error)
	// GetOrganizationByName gets a organization by name.
	GetOrganizationByName(context.Context, *connect_go.Request[v1alpha1.GetOrganizationByNameRequest]) (*connect_go.Response[v1alpha1.GetOrganizationByNameResponse], error)
	// ListOrganizations lists all organizations.
	ListOrganizations(context.Context, *connect_go.Request[v1alpha1.ListOrganizationsRequest]) (*connect_go.Response[v1alpha1.ListOrganizationsResponse], error)
	// ListUserOrganizations lists all organizations a user is member of.
	ListUserOrganizations(context.Context, *connect_go.Request[v1alpha1.ListUserOrganizationsRequest]) (*connect_go.Response[v1alpha1.ListUserOrganizationsResponse], error)
	// CreateOrganization creates a new organization.
	CreateOrganization(context.Context, *connect_go.Request[v1alpha1.CreateOrganizationRequest]) (*connect_go.Response[v1alpha1.CreateOrganizationResponse], error)
	// DeleteOrganization deletes a organization.
	DeleteOrganization(context.Context, *connect_go.Request[v1alpha1.DeleteOrganizationRequest]) (*connect_go.Response[v1alpha1.DeleteOrganizationResponse], error)
	// DeleteOrganizationByName deletes a organization by name.
	DeleteOrganizationByName(context.Context, *connect_go.Request[v1alpha1.DeleteOrganizationByNameRequest]) (*connect_go.Response[v1alpha1.DeleteOrganizationByNam