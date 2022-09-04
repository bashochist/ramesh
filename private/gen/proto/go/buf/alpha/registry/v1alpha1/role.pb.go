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

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.29.0
// 	protoc        (unknown)
// source: buf/alpha/registry/v1alpha1/role.proto

package registryv1alpha1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// The roles that users can have in a Server.
type ServerRole int32

const (
	ServerRole_SERVER_ROLE_UNSPECIFIED ServerRole = 0
	ServerRole_SERVER_ROLE_ADMIN       ServerRole = 1
	ServerRole_SERVER_ROLE_MEMBER      ServerRole = 2
)

// Enum value maps for ServerRole.
var (
	ServerRole_name = map[int32]string{
		0: "SERVER_ROLE_UNSPECIFIED",
		1: "SERVER_ROLE_ADMIN",
		2: "SERVER_ROLE_MEMBER",
	}
	ServerRole_value = map[string]int32{
		"SERVER_ROLE_UNSPECIFIED": 0,
		"SERVER_ROLE_ADMIN":       1,
		"SERVER_ROLE_MEMBER":      2,
	}
)

func (x ServerRole) Enum() *ServerRole {
	p := new(ServerRole)
	*p = x
	return p
}

func (x ServerRole) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ServerRole) Descriptor() protoreflect.EnumDescriptor {
	return file_buf_alpha_registry_v1alpha1_role_proto_enumTypes[0].Descriptor()
}

func (ServerRole) Type() protoreflect.EnumType {
	return &file_buf_alpha_registry_v1alpha1_role_proto_enumTypes[0]
}

func (x ServerRole) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ServerRole.Descriptor instead.
func (ServerRole) EnumDescriptor() ([]byte, []int) {
	return file_buf_alpha_registry_v1alpha1_role_proto_rawDescGZIP(), []int{0}
}

// The roles that users can have in a Organization.
type OrganizationRole int32

const (
	OrganizationRole_ORGANIZATION_ROLE_UNSPECIFIED OrganizationRole = 0
	OrganizationRole_ORGANIZATION_ROLE_OWNER       OrganizationRole = 1
	OrganizationRole_ORGANIZATION_ROLE_ADMIN       OrganizationRole = 2
	OrganizationRole_ORGANIZATION_ROLE_MEMBER      OrganizationRole = 3
	OrganizationRole_ORGANIZATION_ROLE_MACHINE     OrganizationRole = 4
)

// Enum value maps for OrganizationRole.
var (
	OrganizationRole_name = map[int32]string{
		0: "ORGANIZATION_ROLE_UNSPECIFIED",
		1: "ORGANIZATION_ROLE_OWNER",
		2: "ORGANIZATION_ROLE_ADMIN",
		3: "ORGANIZATION_ROLE_MEMBER",
		4: "ORGANIZATION_ROLE_MACHINE",
	}
	OrganizationRole_value = map[string]int32{
		"ORGANIZATION_ROLE_UNSPECIFIED": 0,
		"ORGANIZATION_ROLE_OWNER":       1,
		"ORGANIZATION_ROLE_ADMIN":       2,
		"ORGANIZATION_ROLE_MEMBER":      3,
		"ORGANIZATION_ROLE_MACHINE":     4,
	}
)

func (x OrganizationRole) Enum() *OrganizationRole {
	p := new(OrganizationRole)
	*p = x
	return p
}

func (x OrganizationRole) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (OrganizationRole) Descriptor() protoreflect.EnumDescriptor {
	return file_buf_alpha_registry_v1alpha1_role_proto_enumTypes[1].Descriptor()
}

func (OrganizationRole) Type() protoreflect.EnumType {
	return &file_buf_alpha_registry_v1alpha1_role_proto_enumTypes[1]
}

func (x OrganizationRole) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use OrganizationRole.Descriptor instead.
func (OrganizationRole) EnumDescriptor() ([]byte, []int) {
	return file_buf_alpha_registry_v1alpha1_role_proto_rawDescGZIP(), []int{1}
}

// The source of a user's role in an Organization.
type OrganizationRoleSource int32

const (
	OrganizationRoleSource_ORGANIZATION_ROLE_SOURCE_UNSPECIFIED OrganizationRoleSource = 0
	OrganizationRoleSource_ORGANIZATION_ROLE_SOURCE_DIRECT      OrganizationRoleSource = 1
	OrganizationRoleSource_ORGANIZATION_ROLE_SOURCE_JIT         OrganizationRoleSource = 2
	OrganizationRoleSource_ORGANIZATION_ROLE_SOURCE_IDP_GROUP   OrganizationRoleSource = 3
)

// Enum value maps for OrganizationRoleSource.
var (
	OrganizationRoleSource_name = map[int32]string{
		0: "ORGANIZATION_ROLE_SOURCE_UNSPECIFIED",
		1: "ORGANIZATION_ROLE_SOURCE_DIRECT",
		2: "ORGANIZATION_ROLE_SOURCE_JIT",
		3: "ORGANIZATION_ROLE_SOURCE_IDP_GROUP",
	}
	OrganizationRoleSource_value = map[string]int32{
		"ORGANIZATION_ROLE_SOURCE_UNSPECIFIED": 0,
		"ORGANIZATION_ROLE_SOURCE_DIRECT":      1,
		"ORGANIZATION_ROLE_SOURCE_JIT":         2,
		"ORGANIZATION_ROLE_SOURCE_IDP_GROUP":   3,
	}
)

func (x OrganizationRoleSource) Enum() *OrganizationRoleSource {
	p := new(OrganizationRoleSource)
	*p = x
	return p
}

func (x OrganizationRoleSource) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (OrganizationRoleSource) Descriptor() protoreflect.EnumDescriptor {
	return file_buf_alpha_registry_v1alpha1_role_proto_enumTypes[2].Descriptor()
}

func (OrganizationRoleSource) Type() protoreflect.EnumType {
	return &file_buf_alpha_registry_v1alpha1_role_proto_enumTypes[2]
}

func (x OrganizationRoleSource) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use OrganizationRoleSource.Descriptor instead.
func (OrganizationRoleSource) EnumDescriptor() ([]byte, []int) {
	return file_buf_alpha_registry_v1alpha1_role_proto_rawDescGZIP(), []int{2}
}

// The roles that users can have for a Repository.
type RepositoryRole int32

const (
	RepositoryRole_REPOSITORY_ROLE_UNSPECIFIED   RepositoryRole = 0
	RepositoryRole_REPOSITORY_ROLE_OWNER         RepositoryRole = 1
	RepositoryRole_REPOSITORY_ROLE_ADMIN         RepositoryRole = 2
	RepositoryRole_REPOSITORY_ROLE_WRITE         RepositoryRole = 3
	RepositoryRole_REPOSITORY_ROLE_READ          RepositoryRole = 4
	RepositoryRole_REPOSITORY_ROLE_LIMITED_WRITE RepositoryRole = 5
)

// Enum value maps for RepositoryRole.
var (
	RepositoryRole_name = map[int32]string{
		0: "REPOSITORY_ROLE_UNSPECIFIED",
		1: "REPOSITORY_ROLE_OWNER",
		2: "REPOSITORY_ROLE_ADMIN",
		3: "REPOSITORY_ROLE_WRITE",
		4: "REPOSITORY_ROLE_READ",
		5: "REPOSITORY_ROLE_LIMITED_WRITE",
	}
	RepositoryRole_value = map[string]int32{
		"REPOSITORY_ROLE_UNSPECIFIED":   0,
		"REPOSITORY_ROLE_OWNER":         1,
		"REPOSITORY_ROLE_ADMIN":         2,
		"REPOSITORY_ROLE_WRITE":         3,
		"REPOSITORY_ROLE_READ":          4,
		"REPOSITORY_ROLE_LIMITED_WRITE": 5,
	}
)

func (x RepositoryRole) Enum() *RepositoryRole {
	p := new(RepositoryRole)
	*p = x
	return p
}

func (x RepositoryRole) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RepositoryRole) Descriptor() protoreflect.EnumDescriptor {
	return file_buf_alpha_registry_v1alpha1_role_proto_enumTypes[3].Descriptor()
}

func (RepositoryRole) Type() protoreflect.EnumType {
	return &file_buf_alpha_registry_v1alpha1_role_proto_enumTypes[3]
}

func (x RepositoryRole) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RepositoryRole.Descriptor instead.
func (RepositoryRole) EnumDescriptor() ([]byte, []int) {
	return file_buf_alpha_registry_v1alpha1_role_proto_rawDescGZIP(), []int{3}
}

// The roles that users can have for a Template.
//
// Deprecated: Marked as deprecated in buf/alpha/registry/v1alpha1/role.proto.
type TemplateRole int32

const (
	TemplateRole_TEMPLATE_ROLE_UNSPECIFIED TemplateRole = 0
	TemplateRole_TEMPLATE_ROLE_OWNER       TemplateRole = 1
	TemplateRole_TEMPLATE_ROLE_ADMIN       TemplateRole = 2
	TemplateRole_TEMPLATE_ROLE_WRITE       TemplateRole = 3
	TemplateRole_TEMPLATE_ROLE_READ        TemplateRole = 4
)

// Enum value maps for TemplateRole.
var (
	TemplateRole_name = map[int32]string{
		0: "TEMPLATE_ROLE_UNSPECIFIED",
		1: "TEMPLATE_ROLE_OWNER",
		2: "TEMPLATE_ROLE_ADMIN",
		3: "TEMPLATE_ROLE_WRITE",
		4: "TEMPLATE_ROLE_READ",
	}
	TemplateRole_value = map[string]int32{
		"TEMPLATE_ROLE_UNSPECIFIED": 0,
		"TEMPLATE_ROLE_OWNER":       1,
		"TEMPLATE_ROLE_ADMIN":       2,
		"TEMPLATE_ROLE_WRITE":       3,
		"TEMPLATE_ROLE_READ":        4,
	}
)

func (x TemplateRole) Enum() *TemplateRole {
	p := new(TemplateRole)
	*p = x
	return p
}

func (x TemplateRole) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TemplateRole) Descriptor() protoreflect.EnumDescriptor {
	return file_buf_alpha_registry_v1alpha1_role_proto_enumTypes[4].Descriptor()
}

func (TemplateRole) Type() protoreflect.EnumType {
	return &file_buf_alpha_registry_v1alpha1_role_proto_enumTypes[4]
}

func (x TemplateRole) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TemplateRole.Descriptor instead.
func (TemplateRole) EnumDescriptor() ([]byte, []int) {
	return file_buf_alpha_registry_v1alpha1_role_proto_rawDescGZIP(), []int{4}
}

// The roles that users can have for a Plugin.
//
// Deprecated: Marked as deprecated in buf/alpha/registry/v1alpha1/role.proto.
type PluginRole int32

const (
	PluginRole_PLUGIN_ROLE_UNSPECIFIED PluginRole = 0
	PluginRole_PLUGIN_ROLE_OWNER       PluginRole = 1
	P