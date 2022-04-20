
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
// buf/alpha/registry/v1alpha1/plugin.proto is a deprecated file.

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
	// PluginServiceName is the fully-qualified name of the PluginService service.
	PluginServiceName = "buf.alpha.registry.v1alpha1.PluginService"
)

// PluginServiceClient is a client for the buf.alpha.registry.v1alpha1.PluginService service.
type PluginServiceClient interface {
	// ListPlugins returns all the plugins available to the user. This includes
	// public plugins, those uploaded to organizations the user is part of,
	// and any plugins uploaded directly by the user.
	ListPlugins(context.Context, *connect_go.Request[v1alpha1.ListPluginsRequest]) (*connect_go.Response[v1alpha1.ListPluginsResponse], error)
	// ListUserPlugins lists all plugins belonging to a user.
	ListUserPlugins(context.Context, *connect_go.Request[v1alpha1.ListUserPluginsRequest]) (*connect_go.Response[v1alpha1.ListUserPluginsResponse], error)
	// ListOrganizationPlugins lists all plugins for an organization.
	ListOrganizationPlugins(context.Context, *connect_go.Request[v1alpha1.ListOrganizationPluginsRequest]) (*connect_go.Response[v1alpha1.ListOrganizationPluginsResponse], error)
	// GetPluginVersion returns the plugin version, if found.
	GetPluginVersion(context.Context, *connect_go.Request[v1alpha1.GetPluginVersionRequest]) (*connect_go.Response[v1alpha1.GetPluginVersionResponse], error)
	// ListPluginVersions lists all the versions available for the specified plugin.
	ListPluginVersions(context.Context, *connect_go.Request[v1alpha1.ListPluginVersionsRequest]) (*connect_go.Response[v1alpha1.ListPluginVersionsResponse], error)
	// CreatePlugin creates a new plugin.
	CreatePlugin(context.Context, *connect_go.Request[v1alpha1.CreatePluginRequest]) (*connect_go.Response[v1alpha1.CreatePluginResponse], error)
	// GetPlugin returns the plugin, if found.
	GetPlugin(context.Context, *connect_go.Request[v1alpha1.GetPluginRequest]) (*connect_go.Response[v1alpha1.GetPluginResponse], error)
	// DeletePlugin deletes the plugin, if it exists. Note that deleting
	// a plugin may cause breaking changes for templates using that plugin,
	// and should be done with extreme care.
	DeletePlugin(context.Context, *connect_go.Request[v1alpha1.DeletePluginRequest]) (*connect_go.Response[v1alpha1.DeletePluginResponse], error)
	// SetPluginContributor sets the role of a user in the plugin.
	SetPluginContributor(context.Context, *connect_go.Request[v1alpha1.SetPluginContributorRequest]) (*connect_go.Response[v1alpha1.SetPluginContributorResponse], error)
	// ListPluginContributors returns the list of contributors that has an explicit role against the plugin.
	// This does not include users who have implicit roles against the plugin, unless they have also been
	// assigned a role explicitly.
	ListPluginContributors(context.Context, *connect_go.Request[v1alpha1.ListPluginContributorsRequest]) (*connect_go.Response[v1alpha1.ListPluginContributorsResponse], error)
	// DeprecatePlugin deprecates the plugin, if found.
	DeprecatePlugin(context.Context, *connect_go.Request[v1alpha1.DeprecatePluginRequest]) (*connect_go.Response[v1alpha1.DeprecatePluginResponse], error)
	// UndeprecatePlugin makes the plugin not deprecated and removes any deprecation_message.
	UndeprecatePlugin(context.Context, *connect_go.Request[v1alpha1.UndeprecatePluginRequest]) (*connect_go.Response[v1alpha1.UndeprecatePluginResponse], error)
	// GetTemplate returns the template, if found.
	GetTemplate(context.Context, *connect_go.Request[v1alpha1.GetTemplateRequest]) (*connect_go.Response[v1alpha1.GetTemplateResponse], error)
	// ListTemplates returns all the templates available to the user. This includes
	// public templates, those owned by organizations the user is part of,
	// and any created directly by the user.
	ListTemplates(context.Context, *connect_go.Request[v1alpha1.ListTemplatesRequest]) (*connect_go.Response[v1alpha1.ListTemplatesResponse], error)
	// ListTemplatesUserCanAccess is like ListTemplates, but does not return
	// public templates.
	ListTemplatesUserCanAccess(context.Context, *connect_go.Request[v1alpha1.ListTemplatesUserCanAccessRequest]) (*connect_go.Response[v1alpha1.ListTemplatesUserCanAccessResponse], error)
	// ListUserPlugins lists all templates belonging to a user.
	ListUserTemplates(context.Context, *connect_go.Request[v1alpha1.ListUserTemplatesRequest]) (*connect_go.Response[v1alpha1.ListUserTemplatesResponse], error)
	// ListOrganizationTemplates lists all templates for an organization.
	ListOrganizationTemplates(context.Context, *connect_go.Request[v1alpha1.ListOrganizationTemplatesRequest]) (*connect_go.Response[v1alpha1.ListOrganizationTemplatesResponse], error)
	// GetTemplateVersion returns the template version, if found.
	GetTemplateVersion(context.Context, *connect_go.Request[v1alpha1.GetTemplateVersionRequest]) (*connect_go.Response[v1alpha1.GetTemplateVersionResponse], error)
	// ListTemplateVersions lists all the template versions available for the specified template.
	ListTemplateVersions(context.Context, *connect_go.Request[v1alpha1.ListTemplateVersionsRequest]) (*connect_go.Response[v1alpha1.ListTemplateVersionsResponse], error)
	// CreateTemplate creates a new template.
	CreateTemplate(context.Context, *connect_go.Request[v1alpha1.CreateTemplateRequest]) (*connect_go.Response[v1alpha1.CreateTemplateResponse], error)
	// DeleteTemplate deletes the template, if it exists.
	DeleteTemplate(context.Context, *connect_go.Request[v1alpha1.DeleteTemplateRequest]) (*connect_go.Response[v1alpha1.DeleteTemplateResponse], error)
	// CreateTemplateVersion creates a new template version.
	CreateTemplateVersion(context.Context, *connect_go.Request[v1alpha1.CreateTemplateVersionRequest]) (*connect_go.Response[v1alpha1.CreateTemplateVersionResponse], error)
	// SetTemplateContributor sets the role of a user in the template.
	SetTemplateContributor(context.Context, *connect_go.Request[v1alpha1.SetTemplateContributorRequest]) (*connect_go.Response[v1alpha1.SetTemplateContributorResponse], error)
	// ListTemplateContributors returns the list of contributors that has an explicit role against the template.
	// This does not include users who have implicit roles against the template, unless they have also been
	// assigned a role explicitly.
	ListTemplateContributors(context.Context, *connect_go.Request[v1alpha1.ListTemplateContributorsRequest]) (*connect_go.Response[v1alpha1.ListTemplateContributorsResponse], error)
	// DeprecateTemplate deprecates the template, if found.
	DeprecateTemplate(context.Context, *connect_go.Request[v1alpha1.DeprecateTemplateRequest]) (*connect_go.Response[v1alpha1.DeprecateTemplateResponse], error)
	// UndeprecateTemplate makes the template not deprecated and removes any deprecation_message.
	UndeprecateTemplate(context.Context, *connect_go.Request[v1alpha1.UndeprecateTemplateRequest]) (*connect_go.Response[v1alpha1.UndeprecateTemplateResponse], error)
}

// NewPluginServiceClient constructs a client for the buf.alpha.registry.v1alpha1.PluginService
// service. By default, it uses the Connect protocol with the binary Protobuf Codec, asks for
// gzipped responses, and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply
// the connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewPluginServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) PluginServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &pluginServiceClient{
		listPlugins: connect_go.NewClient[v1alpha1.ListPluginsRequest, v1alpha1.ListPluginsResponse](
			httpClient,
			baseURL+"/buf.alpha.registry.v1alpha1.PluginService/ListPlugins",
			opts...,
		),
		listUserPlugins: connect_go.NewClient[v1alpha1.ListUserPluginsRequest, v1alpha1.ListUserPluginsResponse](
			httpClient,
			baseURL+"/buf.alpha.registry.v1alpha1.PluginService/ListUserPlugins",
			opts...,
		),
		listOrganizationPlugins: connect_go.NewClient[v1alpha1.ListOrganizationPluginsRequest, v1alpha1.ListOrganizationPluginsResponse](
			httpClient,
			baseURL+"/buf.alpha.registry.v1alpha1.PluginService/ListOrganizationPlugins",
			opts...,
		),
		getPluginVersion: connect_go.NewClient[v1alpha1.GetPluginVersionRequest, v1alpha1.GetPluginVersionResponse](
			httpClient,
			baseURL+"/buf.alpha.registry.v1alpha1.PluginService/GetPluginVersion",
			opts...,
		),
		listPluginVersions: connect_go.NewClient[v1alpha1.ListPluginVersionsRequest, v1alpha1.ListPluginVersionsResponse](
			httpClient,
			baseURL+"/buf.alpha.registry.v1alpha1.PluginService/ListPluginVersions",
			opts...,
		),
		createPlugin: connect_go.NewClient[v1alpha1.CreatePluginRequest, v1alpha1.CreatePluginResponse](
			httpClient,
			baseURL+"/buf.alpha.registry.v1alpha1.PluginService/CreatePlugin",
			opts...,
		),
		getPlugin: connect_go.NewClient[v1alpha1.GetPluginRequest, v1alpha1.GetPluginResponse](
			httpClient,
			baseURL+"/buf.alpha.registry.v1alpha1.PluginService/GetPlugin",
			opts...,
		),
		deletePlugin: connect_go.NewClient[v1alpha1.DeletePluginRequest, v1alpha1.DeletePluginResponse](
			httpClient,
			baseURL+"/buf.alpha.registry.v1alpha1.PluginService/DeletePlugin",
			opts...,
		),
		setPluginContributor: connect_go.NewClient[v1alpha1.SetPluginContributorRequest, v1alpha1.SetPluginContributorResponse](
			httpClient,
			baseURL+"/buf.alpha.registry.v1alpha1.PluginService/SetPluginContributor",
			opts...,
		),
		listPluginContributors: connect_go.NewClient[v1alpha1.ListPluginContributorsRequest, v1alpha1.ListPluginContributorsResponse](
			httpClient,
			baseURL+"/buf.alpha.registry.v1alpha1.PluginService/ListPluginContributors",
			opts...,
		),
		deprecatePlugin: connect_go.NewClient[v1alpha1.DeprecatePluginRequest, v1alpha1.DeprecatePluginResponse](
			httpClient,
			baseURL+"/buf.alpha.registry.v1alpha1.PluginService/DeprecatePlugin",
			opts...,
		),
		undeprecatePlugin: connect_go.NewClient[v1alpha1.UndeprecatePluginRequest, v1alpha1.UndeprecatePluginResponse](
			httpClient,
			baseURL+"/buf.alpha.registry.v1alpha1.PluginService/UndeprecatePlugin",
			opts...,
		),
		getTemplate: connect_go.NewClient[v1alpha1.GetTemplateRequest, v1alpha1.GetTemplateResponse](
			httpClient,
			baseURL+"/buf.alpha.registry.v1alpha1.PluginService/GetTemplate",
			opts...,
		),
		listTemplates: connect_go.NewClient[v1alpha1.ListTemplatesRequest, v1alpha1.ListTemplatesResponse](
			httpClient,
			baseURL+"/buf.alpha.registry.v1alpha1.PluginService/ListTemplates",
			opts...,
		),
		listTemplatesUserCanAccess: connect_go.NewClient[v1alpha1.ListTemplatesUserCanAccessRequest, v1alpha1.ListTemplatesUserCanAccessResponse](
			httpClient,
			baseURL+"/buf.alpha.registry.v1alpha1.PluginService/ListTemplatesUserCanAccess",
			opts...,
		),
		listUserTemplates: connect_go.NewClient[v1alpha1.ListUserTemplatesRequest, v1alpha1.ListUserTemplatesResponse](
			httpClient,
			baseURL+"/buf.alpha.registry.v1alpha1.PluginService/ListUserTemplates",
			opts...,
		),
		listOrganizationTemplates: connect_go.NewClient[v1alpha1.ListOrganizationTemplatesRequest, v1alpha1.ListOrganizationTemplatesResponse](
			httpClient,
			baseURL+"/buf.alpha.registry.v1alpha1.PluginService/ListOrganizationTemplates",
			opts...,
		),
		getTemplateVersion: connect_go.NewClient[v1alpha1.GetTemplateVersionRequest, v1alpha1.GetTemplateVersionResponse](
			httpClient,
			baseURL+"/buf.alpha.registry.v1alpha1.PluginService/GetTemplateVersion",
			opts...,
		),