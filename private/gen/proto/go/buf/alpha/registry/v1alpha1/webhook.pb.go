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
// source: buf/alpha/registry/v1alpha1/webhook.proto

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

// WebhookEvent contains the currently supported webhook event types.
type WebhookEvent int32

const (
	// WEBHOOK_EVENT_UNSPECIFIED is a safe noop default for webhook events
	// subscription. It will trigger an error if trying to register a webhook with
	// this event.
	WebhookEvent_WEBHOOK_EVENT_UNSPECIFIED WebhookEvent = 0
	// WEBHOOK_EVENT_REPOSITORY_PUSH is emitted whenever a successful buf push is
	// completed for a specific repository.
	WebhookEvent_WEBHOOK_EVENT_REPOSITORY_PUSH WebhookEvent = 1
)

// Enum value maps for WebhookEvent.
var (
	WebhookEvent_name = map[int32]string{
		0: "WEBHOOK_EVENT_UNSPECIFIED",
		1: "WEBHOOK_EVENT_REPOSITORY_PUSH",
	}
	WebhookEvent_value = map[string]int32{
		"WEBHOOK_EVENT_UNSPECIFIED":     0,
		"WEBHOOK_EVENT_REPOSITORY_PUSH": 1,
	}
)

func (x WebhookEvent) Enum() *WebhookEvent {
	p := new(WebhookEvent)
	*p = x
	return p
}

func (x WebhookEvent) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (WebhookEvent) Descriptor() protoreflect.EnumDescriptor {
	return file_buf_alpha_registry_v1alpha1_webhook_proto_enumTypes[0].Descriptor()
}

func (WebhookEvent) Type() protoreflect.EnumType {
	return &file_buf_alpha_registry_v1alpha1_webhook_proto_enumTypes[0]
}

func (x WebhookEvent) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use WebhookEvent.Descriptor instead.
func (WebhookEvent) EnumDescriptor() ([]byte, []int) {
	return file_buf_alpha_registry_v1alpha1_webhook_proto_rawDescGZIP(), []int{0}
}

// CreateWebhookRequest is the proto request representation of a
// webhook request body.
type CreateWebhookRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The event to subscribe to for the given repository.
	WebhookEvent WebhookEvent `protobuf:"varint,1,opt,name=webhook_event,json=webhookEvent,proto3,enum=buf.alpha.registry.v1alpha1.WebhookEvent" json:"webhook_event,omitempty"`
	// The owner name of the repository in the corresponding subscription request.
	OwnerName string `protobuf:"bytes,2,opt,name=owner_name,json=ownerName,proto3" json:"owner_name,omitempty"`
	// The repository name that the subscriber wishes create a subscription for.
	RepositoryName string `protobuf:"bytes,3,opt,name=repository_name,json=repositoryName,proto3" json:"repository_name,omitempty"`
	// The subscriber's callback URL where notifications should be delivered.
	CallbackUrl string `protobuf:"bytes,4,opt,name=callback_url,json=callbackUrl,proto3" json:"callback_url,omitempty"`
}

func (x *CreateWebhookRequest) Reset() {
	*x = CreateWebhookRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_buf_alpha_registry_v1alpha1_webhook_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateWebhookRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateWebhookRequest) ProtoMessage() {}

func (x *CreateWebhookRequest) ProtoReflect() protoreflect.Message {
	mi := &file_buf_alpha_registry_v1alpha1_webhook_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateWebhookRequest.ProtoReflect.Descriptor instead.
func (*CreateWebhookRe