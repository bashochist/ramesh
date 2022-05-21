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
	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	DeprecationMessage string `protobuf:"bytes,6,opt,name=deprecation_message,json=deprecationMessage,proto3" json:"deprecation_message,omitempty"`
	// The creation time of the plugin.
	//
	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	CreateTime *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=create_time,json=createTime,proto3" json:"create_time,omitempty"`
	// The last update time of the plugin object.
	//
	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	UpdateTime *timestamppb.Timestamp `protobuf:"bytes,8,opt,name=update_time,json=updateTime,proto3" json:"update_time,omitempty"`
}

func (x *Plugin) Reset() {
	*x = Plugin{}
	if protoimpl.UnsafeEnabled {
		mi := &file_buf_alpha_registry_v1alpha1_plugin_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Plugin) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Plugin) ProtoMessage() {}

func (x *Plugin) ProtoReflect() protoreflect.Message {
	mi := &file_buf_alpha_registry_v1alpha1_plugin_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Plugin.ProtoReflect.Descriptor instead.
func (*Plugin) Descriptor() ([]byte, []int) {
	return file_buf_alpha_registry_v1alpha1_plugin_proto_rawDescGZIP(), []int{0}
}

// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
func (x *Plugin) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
func (x *Plugin) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
func (x *Plugin) GetOwner() string {
	if x != nil {
		return x.Owner
	}
	return ""
}

// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
func (x *Plugin) GetVisibility() PluginVisibility {
	if x != nil {
		return x.Visibility
	}
	return PluginVisibility_PLUGIN_VISIBILITY_UNSPECIFIED
}

// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
func (x *Plugin) GetDeprecated() bool {
	if x != nil {
		return x.Deprecated
	}
	return false
}

// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
func (x *Plugin) GetDeprecationMessage() string {
	if x != nil {
		return x.DeprecationMessage
	}
	return ""
}

// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
func (x *Plugin) GetCreateTime() *timestamppb.Timestamp {
	if x != nil {
		return x.CreateTime
	}
	return nil
}

// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
func (x *Plugin) GetUpdateTime() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdateTime
	}
	return nil
}

// PluginVersion represents a specific build of a plugin,
// such as protoc-gen-go v1.4.0.
//
// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
type PluginVersion struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The ID of the plugin version, which uniquely identifies the plugin version.
	// Mostly used for pagination.
	//
	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// The name of the version, i.e. "v1.4.0".
	//
	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// The name of the plugin to which this version relates.
	//
	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	PluginName string `protobuf:"bytes,3,opt,name=plugin_name,json=pluginName,proto3" json:"plugin_name,omitempty"`
	// The owner of the plugin to which this version relates.
	//
	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	PluginOwner string `protobuf:"bytes,4,opt,name=plugin_owner,json=pluginOwner,proto3" json:"plugin_owner,omitempty"`
	// The full container image digest associated with this plugin version including
	// the algorithm.
	// Ref: https://github.com/opencontainers/image-spec/blob/main/descriptor.md#digests
	//
	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	ContainerImageDigest string `protobuf:"bytes,5,opt,name=container_image_digest,json=containerImageDigest,proto3" json:"container_image_digest,omitempty"`
	// Optionally define the runtime libraries.
	//
	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	RuntimeLibraries []*RuntimeLibrary `protobuf:"bytes,6,rep,name=runtime_libraries,json=runtimeLibraries,proto3" json:"runtime_libraries,omitempty"`
}

func (x *PluginVersion) Reset() {
	*x = PluginVersion{}
	if protoimpl.UnsafeEnabled {
		mi := &file_buf_alpha_registry_v1alpha1_plugin_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PluginVersion) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PluginVersion) ProtoMessage() {}

func (x *PluginVersion) ProtoReflect() protoreflect.Message {
	mi := &file_buf_alpha_registry_v1alpha1_plugin_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PluginVersion.ProtoReflect.Descriptor instead.
func (*PluginVersion) Descriptor() ([]byte, []int) {
	return file_buf_alpha_registry_v1alpha1_plugin_proto_rawDescGZIP(), []int{1}
}

// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
func (x *PluginVersion) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
func (x *PluginVersion) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
func (x *PluginVersion) GetPluginName() string {
	if x != nil {
		return x.PluginName
	}
	return ""
}

// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
func (x *PluginVersion) GetPluginOwner() string {
	if x != nil {
		return x.PluginOwner
	}
	return ""
}

// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
func (x *PluginVersion) GetContainerImageDigest() string {
	if x != nil {
		return x.ContainerImageDigest
	}
	return ""
}

// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
func (x *PluginVersion) GetRuntimeLibraries() []*RuntimeLibrary {
	if x != nil {
		return x.RuntimeLibraries
	}
	return nil
}

// Template defines a set of plugins that should be used together
// i.e. "go-grpc" would include protoc-gen-go and protoc-gen-go-grpc.
//
// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
type Template struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The ID of the template, which uniquely identifies the template.
	//
	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// The name of the template, i.e. "grpc-go".
	//
	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// The name of the owner of the template. Either a
	// username or organization name.
	//
	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	Owner string `protobuf:"bytes,3,opt,name=owner,proto3" json:"owner,omitempty"`
	// Must not contain duplicate plugins. Order of plugin configs
	// dictates insertion point order. Note that we're
	// intentionally putting most of the plugin configuration
	// in the template, so that template versions are
	// less likely to cause breakages for users.
	//
	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	PluginConfigs []*PluginConfig `protobuf:"bytes,4,rep,name=plugin_configs,json=pluginConfigs,proto3" json:"plugin_configs,omitempty"`
	// The visibility of the template.
	//
	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	Visibility PluginVisibility `protobuf:"varint,5,opt,name=visibility,proto3,enum=buf.alpha.registry.v1alpha1.PluginVisibility" json:"visibility,omitempty"`
	// deprecated means this template is deprecated.
	//
	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	Deprecated bool `protobuf:"varint,8,opt,name=deprecated,proto3" json:"deprecated,omitempty"`
	// deprecation_message is the message shown if the template is deprecated.
	//
	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	DeprecationMessage string `protobuf:"bytes,9,opt,name=deprecation_message,json=deprecationMessage,proto3" json:"deprecation_message,omitempty"`
	// The creation time of the template.
	//
	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	CreateTime *timestamppb.Timestamp `protobuf:"bytes,10,opt,name=create_time,json=createTime,proto3" json:"create_time,omitempty"`
	// The last update time of the template object.
	//
	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	UpdateTime *timestamppb.Timestamp `protobuf:"bytes,11,opt,name=update_time,json=updateTime,proto3" json:"update_time,omitempty"`
}

func (x *Template) Reset() {
	*x = Template{}
	if protoimpl.UnsafeEnabled {
		mi := &file_buf_alpha_registry_v1alpha1_plugin_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Template) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Template) ProtoMessage() {}

func (x *Template) ProtoReflect() protoreflect.Message {
	mi := &file_buf_alpha_registry_v1alpha1_plugin_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Template.ProtoReflect.Descriptor instead.
func (*Template) Descriptor() ([]byte, []int) {
	return file_buf_alpha_registry_v1alpha1_plugin_proto_rawDescGZIP(), []int{2}
}

// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
func (x *Template) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
func (x *Template) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
func (x *Template) GetOwner() string {
	if x != nil {
		return x.Owner
	}
	return ""
}

// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
func (x *Template) GetPluginConfigs() []*PluginConfig {
	if x != nil {
		return x.PluginConfigs
	}
	return nil
}

// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
func (x *Template) GetVisibility() PluginVisibility {
	if x != nil {
		return x.Visibility
	}
	return PluginVisibility_PLUGIN_VISIBILITY_UNSPECIFIED
}

// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
func (x *Template) GetDeprecated() bool {
	if x != nil {
		return x.Deprecated
	}
	return false
}

// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
func (x *Template) GetDeprecationMessage() string {
	if x != nil {
		return x.DeprecationMessage
	}
	return ""
}

// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
func (x *Template) GetCreateTime() *timestamppb.Timestamp {
	if x != nil {
		return x.CreateTime
	}
	return nil
}

// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
func (x *Template) GetUpdateTime() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdateTime
	}
	return nil
}

// PluginConfig defines a runtime configuration for a plugin.
//
// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
type PluginConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The owner of the plugin to which this config relates.
	//
	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	PluginOwner string `protobuf:"bytes,1,opt,name=plugin_owner,json=pluginOwner,proto3" json:"plugin_owner,omitempty"`
	// The name of the plugin to which this config relates.
	//
	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	PluginName string `protobuf:"bytes,2,opt,name=plugin_name,json=pluginName,proto3" json:"plugin_name,omitempty"`
	// Parameters that should be provided to the plugin. These are
	// joined with a "," before being provided to the plugin at runtime.
	//
	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	Parameters []string `protobuf:"bytes,3,rep,name=parameters,proto3" json:"parameters,omitempty"`
	// True if the source plugin is inaccessible by the user.
	//
	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	Inaccessible bool `protobuf:"varint,5,opt,name=inaccessible,proto3" json:"inaccessible,omitempty"`
}

func (x *PluginConfig) Reset() {
	*x = PluginConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_buf_alpha_registry_v1alpha1_plugin_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PluginConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PluginConfig) ProtoMessage() {}

func (x *PluginConfig) ProtoReflect() protoreflect.Message {
	mi := &file_buf_alpha_registry_v1alpha1_plugin_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PluginConfig.ProtoReflect.Descriptor instead.
func (*PluginConfig) Descriptor() ([]byte, []int) {
	return file_buf_alpha_registry_v1alpha1_plugin_proto_rawDescGZIP(), []int{3}
}

// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
func (x *PluginConfig) GetPluginOwner() string {
	if x != nil {
		return x.PluginOwner
	}
	return ""
}

// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
func (x *PluginConfig) GetPluginName() string {
	if x != nil {
		return x.PluginName
	}
	return ""
}

// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
func (x *PluginConfig) GetParameters() []string {
	if x != nil {
		return x.Parameters
	}
	return nil
}

// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
func (x *PluginConfig) GetInaccessible() bool {
	if x != nil {
		return x.Inaccessible
	}
	return false
}

// TemplateVersion defines a template at a
// specific set of versions for the contained plugins.
//
// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
type TemplateVersion struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The ID of the template version, which uniquely identifies the template version.
	// Mostly used for pagination.
	//
	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// The name of the template version, i.e. "v1".
	//
	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// The owner of the template to which this version relates.
	//
	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	TemplateOwner string `protobuf:"bytes,3,opt,name=template_owner,json=templateOwner,proto3" json:"template_owner,omitempty"`
	// The name of the template to which this version relates.
	//
	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	TemplateName string `protobuf:"bytes,4,opt,name=template_name,json=templateName,proto3" json:"template_name,omitempty"`
	// A map from plugin owner and name to version for the plugins
	// defined in the template. Every plugin in the template
	// must have a corresponding version in this array.
	//
	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	PluginVersions []*PluginVersionMapping `protobuf:"bytes,5,rep,name=plugin_versions,json=pluginVersions,proto3" json:"plugin_versions,omitempty"`
}

func (x *TemplateVersion) Reset() {
	*x = TemplateVersion{}
	if protoimpl.UnsafeEnabled {
		mi := &file_buf_alpha_registry_v1alpha1_plugin_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TemplateVersion) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TemplateVersion) ProtoMessage() {}

func (x *TemplateVersion) ProtoReflect() protoreflect.Message {
	mi := &file_buf_alpha_registry_v1alpha1_plugin_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TemplateVersion.ProtoReflect.Descriptor instead.
func (*TemplateVersion) Descriptor() ([]byte, []int) {
	return file_buf_alpha_registry_v1alpha1_plugin_proto_rawDescGZIP(), []int{4}
}

// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
func (x *TemplateVersion) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
func (x *TemplateVersion) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
func (x *TemplateVersion) GetTemplateOwner() string {
	if x != nil {
		return x.TemplateOwner
	}
	return ""
}

// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
func (x *TemplateVersion) GetTemplateName() string {
	if x != nil {
		return x.TemplateName
	}
	return ""
}

// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
func (x *TemplateVersion) GetPluginVersions() []*PluginVersionMapping {
	if x != nil {
		return x.PluginVersions
	}
	return nil
}

// PluginVersionMapping maps a plugin_id to a version.
//
// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
type PluginVersionMapping struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The owner of the plugin to which this mapping relates.
	//
	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	PluginOwner string `protobuf:"bytes,1,opt,name=plugin_owner,json=pluginOwner,proto3" json:"plugin_owner,omitempty"`
	// The name of the plugin to which this mapping relates.
	//
	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	PluginName string `protobuf:"bytes,2,opt,name=plugin_name,json=pluginName,proto3" json:"plugin_name,omitempty"`
	// The version of the plugin to use, i.e. "v1.4.0".
	//
	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	Version string `protobuf:"bytes,3,opt,name=version,proto3" json:"version,omitempty"`
	// True if the source plugin is inaccessible by the user.
	//
	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	Inaccessible bool `protobuf:"varint,5,opt,name=inaccessible,proto3" json:"inaccessible,omitempty"`
}

func (x *PluginVersionMapping) Reset() {
	*x = PluginVersionMapping{}
	if protoimpl.UnsafeEnabled {
		mi := &file_buf_alpha_registry_v1alpha1_plugin_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PluginVersionMapping) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PluginVersionMapping) ProtoMessage() {}

func (x *PluginVersionMapping) ProtoReflect() protoreflect.Message {
	mi := &file_buf_alpha_registry_v1alpha1_plugin_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PluginVersionMapping.ProtoReflect.Descriptor instead.
func (*PluginVersionMapping) Descriptor() ([]byte, []int) {
	return file_buf_alpha_registry_v1alpha1_plugin_proto_rawDescGZIP(), []int{5}
}

// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
func (x *PluginVersionMapping) GetPluginOwner() string {
	if x != nil {
		return x.PluginOwner
	}
	return ""
}

// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
func (x *PluginVersionMapping) GetPluginName() string {
	if x != nil {
		return x.PluginName
	}
	return ""
}

// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
func (x *PluginVersionMapping) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
func (x *PluginVersionMapping) GetInaccessible() bool {
	if x != nil {
		return x.Inaccessible
	}
	return false
}

// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
type PluginContributor struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	User *User `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	// The ID of the plugin which the role belongs to.
	//
	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	PluginId string `protobuf:"bytes,2,opt,name=plugin_id,json=pluginId,proto3" json:"plugin_id,omitempty"`
	// The role that the user has been explicitly assigned against the plugin.
	//
	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	ExplicitRole PluginRole `protobuf:"varint,3,opt,name=explicit_role,json=explicitRole,proto3,enum=buf.alpha.registry.v1alpha1.PluginRole" json:"explicit_role,omitempty"`
	// Optionally defines the role that the user has implicitly against the plugin through the owning organization.
	// If the plugin does not belong to an organization or the user is not part of the owning organization, this is unset.
	//
	// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
	ImplicitRole PluginRole `protobuf:"varint,4,opt,name=implicit_role,json=implicitRole,proto3,enum=buf.alpha.registry.v1alpha1.PluginRole" json:"implicit_role,omitempty"`
}

func (x *PluginContributor) Reset() {
	*x = PluginContributor{}
	if protoimpl.UnsafeEnabled {
		mi := &file_buf_alpha_registry_v1alpha1_plugin_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PluginContributor) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PluginContributor) ProtoMessage() {}

func (x *PluginContributor) ProtoReflect() protoreflect.Message {
	mi := &file_buf_alpha_registry_v1alpha1_plugin_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PluginContributor.ProtoReflect.Descriptor instead.
func (*PluginContributor) Descriptor() ([]byte, []int) {
	return file_buf_alpha_registry_v1alpha1_plugin_proto_rawDescGZIP(), []int{6}
}

// Deprecated: The entire proto file buf/alpha/registry/v1alpha1/plugin.proto is marked as deprecated.
func (x *PluginContributor) GetUser() *User {
	if x != nil {
		return x.User
	}
	