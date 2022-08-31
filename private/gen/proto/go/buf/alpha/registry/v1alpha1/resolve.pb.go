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
// source: buf/alpha/registry/v1alpha1/resolve.proto

package registryv1alpha1

import (
	v1alpha1 "github.com/bufbuild/buf/private/gen/proto/go/buf/alpha/module/v1alpha1"
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

type ResolvedReferenceType int32

const (
	ResolvedReferenceType_RESOLVED_REFERENCE_TYPE_UNSPECIFIED ResolvedReferenceType = 0
	ResolvedReferenceType_RESOLVED_REFERENCE_TYPE_COMMIT      ResolvedReferenceType = 1
	ResolvedReferenceType_RESOLVED_REFERENCE_TYPE_TAG         ResolvedReferenceType = 3
	ResolvedReferenceType_RESOLVED_REFERENCE_TYPE_DRAFT       ResolvedReferenceType = 5
)

// Enum value maps for ResolvedReferenceType.
var (
	ResolvedReferenceType_name = map[int32]string{
		0: "RESOLVED_REFERENCE_TYPE_UNSPECIFIED",
		1: "RESOLVED_REFERENCE_TYPE_COMMIT",
		3: "RESOLVED_REFERENCE_TYPE_TAG",
		5: "RESOLVED_REFERENCE_TYPE_DRAFT",
	}
	ResolvedReferenceType_value = map[string]int32{
		"RESOLVED_REFERENCE_TYPE_UNSPECIFIED": 0,
		"RESOLVED_REFERENCE_TYPE_COMMIT":      1,
		"RESOLVED_REFERENCE_TYPE_TAG":         3,
		"RESOLVED_REFERENCE_TYPE_DRAFT":       5,
	}
)

func (x ResolvedReferenceType) Enum() *ResolvedReferenceType {
	p := new(ResolvedReferenceType)
	*p = x
	return p
}

func (x ResolvedReferenceType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ResolvedReferenceType) Descriptor() protoreflect.EnumDescriptor {
	return file_buf_alpha_registry_v1alpha1_resolve_proto_enumTypes[0].Descriptor()
}

func (ResolvedReferenceType) Type() protoreflect.EnumType {
	return &file_buf_alpha_registry_v1alpha1_resolve_proto_enumTypes[0]
}

func (x ResolvedReferenceType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ResolvedReferenceType.Descriptor instead.
func (ResolvedReferenceType) EnumDescriptor() ([]byte, []int) {
	return file_buf_alpha_registry_v1alpha1_resolve_proto_rawDescGZIP(), []int{0}
}

type GetModulePinsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ModuleReferences []*v1alpha1.ModuleReference `protobuf:"bytes,1,rep,name=module_references,json=moduleReferences,proto3" json:"module_references,omitempty"`
	// current_module_pins allows for partial dependency updates by letting clients
	// send a request with the pins for their current module and only the
	// identities of the dependencies they want to update in module_references.
	//
	// When resolving, if a client supplied module pin is:
	//   - in the transitive closure of pins resolved from the module_references,
	//     the client supplied module pin will be an extra candidate for tie
	//     breaking.
	//   - NOT in the in the transitive closure of pins resolved from the
	//     module_references, it will be returned as is.
	CurrentModulePins []*v1alpha1.ModulePin `protobuf:"bytes,2,rep,name=current_module_pins,json=currentModulePins,proto3" json:"current_module_pins,omitempty"`
}

func (x *GetModulePinsRequest) Reset() {
	*x = GetModulePinsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_buf_alpha_registry_v1alpha1_resolve_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetModulePinsRequest) String() string {
	return protoimpl.X.Messa