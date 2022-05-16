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
// buf/alpha/registry/v1alpha1/plugin.proto is a deprecated file.

package registryv1alpha1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// PluginVisibility defines the visibility options available
// for Plugins and Templates.
//
// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
type PluginVisibility int32

const (
	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	PluginVisibility_PLUGIN_VISIBILITY_UNSPECIFIED PluginVisibility = 0
	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	PluginVisibility_PLUGIN_VISIBILITY_PUBLIC PluginVisibility = 1
	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	PluginVisibility_PLUGIN_VISIBILITY_PRIVATE PluginVisibility = 2
)

// Enum value maps for PluginVisibility.
var (
	PluginVisibility_name = map[int32]string{
		0: "PLUGIN_VISIBILITY_UNSPECIFIED",
		1: "PLUGIN_VISIBILITY_PUBLIC",
		2: "PLUGIN_VISIBILITY_PRIVATE",
	}
	PluginVisibility_value = map[string]int32{
		"PLUGIN_VISIBILITY_UNSPECIFIED": 0,
		"PLUGIN_VISIBILITY_PUBLIC":      1,
		"PLUGIN_VISIBILITY_PRIVATE":     2,
	}
)

func (x PluginVisibility) Enum() *PluginVisibility {
	p := new(PluginVisibility)
	*p = x
	return p
}

func (x PluginVisibility) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PluginVisibility) Descriptor() protoreflect.EnumDescriptor {
	return file_buf_alpha_registry_v1alpha1_plugin_proto_enumTypes[0].Descriptor()
}

func (PluginVisibility) Type() protoreflect.EnumType {
	return &file_buf_alpha_registry_v1alpha1_plugin_proto_enumTypes[0]
}

func (x PluginVisibility) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PluginVisibility.Descriptor instead.
func (PluginVisibility) EnumDescriptor() ([]byte, []int) {
	return file_buf_alpha_registry_v1alpha1_plugin_proto_rawDescGZIP(), []int{0}
}

// Plugin represents a protoc plugin, such as protoc-gen-go.
//
// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
type Plugin struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The ID of the plugin, which uniquely identifies the plugin.
	//
	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// The name of the plugin, i.e. "protoc-gen-go".
	//
	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// The name of the owner of the plugin. Either a username or
	// organization name.
	//
	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	Owner string `protobuf:"bytes,3,opt,name=owner,proto3" json:"owner,omitempty"`
	// The visibility of the plugin.
	//
	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	Visibility PluginVisibility `protobuf:"varint,4,opt,name=visibility,proto3,enum=buf.alpha.registry.v1alpha1.PluginVisibility" json:"visibility,omitempty"`
	// deprecated means this plugin is deprecated.
	//
	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	Deprecated bool `protobuf:"varint,5,opt,name=deprecated,proto3" json:"deprecated,omitempty"`
	// deprecation_message is the message shown if the plugin is deprecated.
	//
	// Deprecat