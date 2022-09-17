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
// source: buf/alpha/registry/v1alpha1/user.proto

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

type UserState int32

const (
	UserState_USER_STATE_UNSPECIFIED UserState = 0
	UserState_USER_STATE_ACTIVE      UserState = 1
	UserState_USER_STATE_DEACTIVATED UserState = 2
)

// Enum value maps for UserState.
var (
	UserState_name = map[int32]string{
		0: "USER_STATE_UNSPECIFIED",
		1: "USER_STATE_ACTIVE",
		2: "USER_STATE_DEACTIVATED",
	}
	UserState_value = map[string]int32{
		"USER_STATE_UNSPECIFIED": 0,
		"USER_STATE_ACTIVE":      1,
		"USER_STATE_DEACTIVATED": 2,
	}
)

func (x UserState) Enum() *UserState {
	p := new(UserState)
	*p = x
	return p
}

func (x UserState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (UserState) Descriptor() protoreflect.EnumDescriptor {
	return file_buf_alpha_registry_v1alpha1_user_proto_enumTypes[0].Descriptor()
}

func (UserState) Type() protoreflect.EnumType {
	return &file_buf_alpha_registry_v1alpha1_user_proto_enumTypes[0]
}

func (x UserState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use UserState.Descriptor instead.
func (UserState) EnumDescriptor() ([]byte, []int) {
	return file_buf_alpha_registry_v1alpha1_user_proto_rawDescGZIP(), []int{0}
}

type UserType int32

const (
	UserType_USER_TYPE_UNSPECIFIED UserType = 0
	UserType_USER_TYPE_PERSONAL    UserType = 1
	UserType_USER_TYPE_MACHINE     UserType = 2
	UserType_USER_TYPE_SYSTEM      UserType = 3
)

// Enum value maps for UserType.
var (
	UserType_name = map[int32]string{
		0: "USER_TYPE_UNSPECIFIED",
		1: "USER_TYPE_PERSONAL",
		2: "USER_TYPE_MACHINE",
		3: "USER_TYPE_SYSTEM",
	}
	UserType_value = map[string]int32{
		"USER_TYPE_UNSPECIFIED": 0,
		"USER_TYPE_PERSONAL":    1,
		"USER_TYPE_MACHINE":     2,
		"USER_TYPE_SYSTEM":      3,
	}
)

func (x UserType) Enum() *UserType {
	p := new(UserType)
	*p = x
	return p
}

func (x UserType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (UserType) Descriptor() protoreflect.EnumDescriptor {
	return file_buf_alpha_registry_v1alpha1_user_proto_enumTypes[1].Descriptor()
}

func (UserType) Type() protoreflect.EnumType {
	return &file_buf_alpha_registry_v1alpha1_user_proto_enumTypes[1]
}

func (x UserType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use UserType.Descriptor instead.
func (UserType) EnumDescriptor() ([]byte, []int) {
	return file_buf_alpha_registry_v1alpha1_user_proto_rawDescGZIP(), []int{1}
}

type User struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// primary key, unique, immutable
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// immutable
	CreateTime *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=create_time,json=createTime,proto3" json:"create_time,omitempty"`
	// mutable
	UpdateTime *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=update_time,json=updateTime,proto3" json:"update_time,omitempty"`
	// unique, mutable
	Username string `protobuf:"bytes,4,opt,name=username,proto3" json:"username,omitempty"`
	// mutable
	Deactivated bool `protobuf:"varint,5,opt,name=deactivated,proto3" json:"deactivated,omitempty"`
	// description is the user configurable description of the user.
	Description string `protobuf:"bytes,6,opt,name=description,proto3" json:"description,omitempty"`
	// url is the user configurable URL in the description of the user,
	// always included the scheme and will not have a #fragment suffix.
	Url string `protobuf:"bytes,7,opt,name=url,proto3" json:"url,omitempty"`
	// verification status of the user, configurable by server admin.
	VerificationStatus VerificationStatus `protobuf:"varint,8,opt,name=verification_status,json=verificationStatus,proto3,enum=buf.alpha.registry.v1alpha1.VerificationStatus" json:"verification_status,omitempty"`
	// user type of the user, depends on how the user was created.
	UserType UserType `protobuf:"varint,9,opt,name=user_type,json=userType,proto3,enum=buf.alpha.registry.v1alpha1.UserType" json:"user_type,omitempty"`
}

func (x *User) Reset() {
	*x = User{}
	if protoimpl.UnsafeEnabled {
		mi := &file_buf_alpha_registry_v1alpha1_user_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_buf_alpha_registry_v1alpha1_user_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_buf_alpha_registry_v1alpha1_user_proto_rawDescGZIP(), []int{0}
}

func (x *User) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *User) GetCreateTime() *timestamppb.Timestamp {
	if x != nil {
		return x.CreateTime
	}
	return nil
}

func (x *User) GetUpdateTime() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdateTime
	}
	return nil
}

func (x *User) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *User) GetDeactivated() bool {
	if x != nil {
		return x.Deactivated
	}
	return false
}

func (x *User) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *User) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *User) GetVerificationStatus() VerificationStatus {
	if x != nil {
		return x.VerificationStatus
	}
	return VerificationStatus_VERIFICATION_STATUS_UNSPECIFIED
}

func (x *User) GetUserType() UserType {
	if x != nil {
		return x.UserType
	}
	return UserType_USER_TYPE_UNSPECIFIED
}

// TODO: #663 move this to organization service
type OrganizationUser struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	User *User `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	// The ID of the organization for which the role belongs to.
	OrganizationId string `protobuf:"bytes,2,opt,name=organization_id,json=organizationId,proto3" json:"organization_id,omitempty"`
	// The role that the user has in the organization above.
	OrganizationRole OrganizationRole `protobuf:"varint,3,opt,name=organization_role,json=organizationRole,proto3,enum=buf.alpha.registry.v1alpha1.OrganizationRole" json:"organization_role,omitempty"`
	// The source of the user's role in the organization above.
	OrganizationRoleSource OrganizationRoleSource `protobuf:"varint,4,opt,name=organization_role_source,json=organizationRoleSource,proto3,enum=buf.alpha.registry.v1alpha1.OrganizationRoleSource" json:"organization_role_source,omitempty"`
}

func (x *OrganizationUser) Reset() {
	*x = OrganizationUser{}
	if protoimpl.UnsafeEnabled {
		mi := &file_buf_alpha_registry_v1alpha1_user_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OrganizationUser) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OrganizationUser) ProtoMessage() {}

func (x *OrganizationUser) ProtoReflect() protoreflect.Message {
	mi := &file_buf_alpha_registry_v1alpha1_user_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OrganizationUser.ProtoReflect.Descriptor instead.
func (*OrganizationUser) Descriptor() ([]byte, []int) {
	return file_buf_alpha_registry_v1alpha1_user_proto_rawDescGZIP(), []int{1}
}

func (x *OrganizationUser) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

func (x *OrganizationUser) GetOrganizationId() string {
	if x != nil {
		return x.OrganizationId
	}
	return ""
}

func (x *OrganizationUser) GetOrganizationRole() OrganizationRole {
	if x != nil {
		return x.OrganizationRole
	}
	return OrganizationRole_ORGANIZATION_ROLE_UNSPECIFIED
}

func (x *OrganizationUser) GetOrganizationRoleSource() OrganizationRoleSource {
	if x != nil {
		return x.OrganizationRoleSource
	}
	return OrganizationRoleSource_ORGANIZATION_ROLE_SOURCE_UNSPECIFIED
}

type CreateUserRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
}

func (x *CreateUserRequest) Reset() {
	*x = CreateUserRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_buf_alpha_registry_v1alpha1_user_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateUserRequest) ProtoMessage() {}

func (x *CreateUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_buf_alpha_registry_v1alpha1_user_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateUserRequest.ProtoReflect.Descriptor instead.
func (*CreateUserRequest) Descriptor() ([]byte, []int) {
	return file_buf_alpha_registry_v1alpha1_user_proto_rawDescGZIP(), []int{2}
}

func (x *CreateUserRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

type CreateUserResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	User *User `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
}

func (x *CreateUserResponse) Reset() {
	*x = CreateUserResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_buf_alpha_registry_v1alpha1_user_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateUserResponse) ProtoMessage() {}

func (x *CreateUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_buf_alpha_registry_v1alpha1_user_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateUserResponse.ProtoReflect.Descriptor instead.
func (*CreateUserResponse) Descriptor() ([]byte, []int) {
	return file_buf_alpha_registry_v1alpha1_user_proto_rawDescGZIP(), []int{3}
}

func (x *CreateUserResponse) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

type GetUserRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetUserRequest) Reset() {
	*x = GetUserRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_buf_alpha_registry_v1alpha1_user_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserRequest) ProtoMessage() {}

func (x *GetUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_buf_alpha_registry_v1alpha1_user_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserRequest.ProtoReflect.Descriptor instead.
func (*GetUserRequest) Descriptor() ([]byte, []int) {
	return file_buf_alpha_registry_v1alpha1_user_proto_rawDescGZIP(), []int{4}
}

func (x *GetUserRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type GetUserResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	User *User `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
}

func (x *GetUserResponse) Reset() {
	*x = GetUserResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_buf_alpha_registry_v1alpha1_user_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserResponse) ProtoMessage() {}

func (x *GetUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_buf_alpha_registry_v1alpha1_user_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserResponse.ProtoReflect.Descriptor instead.
func (*GetUserResponse) Descriptor() ([]byte, []int) {
	return file_buf_alpha_registry_v1alpha1_user_proto_rawDescGZIP(), []int{5}
}

func (x *GetUserResponse) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

type GetUserByUsernameRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
}

func (x *GetUserByUsernameRequest) Reset() {
	*x = GetUserByUsernameRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_buf_alpha_registry_v1alpha1_user_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserByUsernameRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserByUsernameRequest) ProtoMessage() {}

func (x *GetUserByUsernameRequest) ProtoReflect() protoreflect.Message {
	mi := &file_buf_alpha_registry_v1alpha1_user_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserByUsernameRequest.ProtoReflect.Descriptor instead.
func (*GetUserByUsernameRequest) Descriptor() ([]byte, []int) {
	return file_buf_alpha_registry_v1alpha1_user_proto_rawDescGZIP(), []int{6}
}

func (x *GetUserByUsernameRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

type GetUserByUsernameResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	User *User `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
}

func (x *GetUserByUsernameResponse) Reset() {
	*x = GetUserByUsernameResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_buf_alpha_registry_v1alpha1_user_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserByUsernameResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserByUsernameResponse) ProtoMessage() {}

func (x *GetUserByUsernameResponse) ProtoReflect() protoreflect.Message {
	mi := &file_buf_alpha_registry_v1alpha1_user_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserByUsernameResponse.ProtoReflect.Descriptor instead.
func (*GetUserByUsernameResponse) Descriptor() ([]byte, []int) {
	return file_buf_alpha_registry_v1alpha1_user_proto_rawDescGZIP(), []int{7}
}

func (x *GetUserByUsernameResponse) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

type ListUsersRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PageSize uint32 `protobuf:"varint,1,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	// The first page is returned if this is empty.
	PageToken string `protobuf:"bytes,2,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
	Reverse   bool   `protobuf:"varint,3,opt,name=reverse,proto3" json:"reverse,omitempty"`
	// If the user_state_filter is unspecified, users of all states are included.
	UserStateFilter UserState `protobuf:"varint,4,opt,name=user_state_filter,json=userStateFilter,proto3,enum=buf.alpha.registry.v1alpha1.UserState" json:"user_state_filter,omitempty"`
	// If the user_type_filters is empty, users of all types are included.
	UserTypeFilters []UserType `protobuf:"varint,5,rep,packed,name=user_type_filters,json=userTypeFilters,proto3,enum=buf.alpha.registry.v1alpha1.UserType" json:"user_type_filters,omitempty"`
}

func (x *ListUsersRequest) Reset() {
	*x = ListUsersRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_buf_alpha_registry_v1alpha1_user_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListUsersRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListUsersRequest) ProtoMessage() {}

func (x *ListUsersRequest) ProtoReflect() protoreflect.Message {
	mi := &file_buf_alpha_registry_v1alpha1_user_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListUsersRequest.ProtoReflect.Descriptor instead.
func (*ListUsersRequest) Descriptor() ([]byte, []int) {
	return file_buf_alpha_registry_v1alpha1_user_proto_rawDescGZIP(), []int{8}
}

func (x *ListUsersRequest) GetPageSize() uint32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *ListUsersRequest) GetPageToken() string {
	if x != nil {
		return x.PageToken
	}
	return ""
}

func (x *ListUsersRequest) GetReverse() bool {
	if x != nil {
		return x.Reverse
	}
	return false
}

func (x *ListUsersRequest) GetUserStateFilter() UserState {
	if x != nil {
		return x.UserStateFilter
	}
	return UserState_USER_STATE_UNSPECIFIED
}

func (x *ListUsersRequest) GetUserTypeFilters() []UserType {
	if x != nil {
		return x.UserTypeFilters
	}
	return nil
}

type ListUsersResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Users []*User `protobuf:"bytes,1,rep,name=users,proto3" json:"users,omitempty"`
	// There are no more pages if this is empty.
	NextPageToken string `protobuf:"bytes,2,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"`
}

func (x *ListUsersResponse) Reset() {
	*x = ListUsersResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_buf_alpha_registry_v1alpha1_user_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListUsersResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListUsersResponse) ProtoMessage() {}

func (x *ListUsersResponse) ProtoReflect() protoreflect.Message {
	mi := &file_buf_alpha_registry_v1alpha1_user_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListUsersResponse.ProtoReflect.Descriptor instead.
func (*ListUsersResponse) Descriptor() ([]byte, []int) {
	return file_buf_alpha_registry_v1alpha1_user_proto_rawDescGZIP(), []int{9}
}

func (x *ListUsersResponse) GetUsers() []*User {
	if x != nil {
		return x.Users
	}
	return nil
}

func (x *ListUsersResponse) GetNextPageToken() string {
	if x != nil {
		return x.NextPageToken
	}
	return ""
}

type ListOrganizationUsersRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OrganizationId string `protobuf:"bytes,1,opt,name=organization_id,json=organizationId,proto3" json:"organization_id,omitempty"`
	PageSize       uint32 `protobuf:"varint,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	// The first page is returned if this is empty.
	PageToken string `protobuf:"bytes,3,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
	Reverse   bool   `protobuf:"varint,4,opt,name=reverse,proto3" json:"reverse,omitempty"`
}

func (x *ListOrganizationUsersRequest) Reset() {
	*x = ListOrganizationUsersRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_buf_alpha_registry_v1alpha1_user_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListOrganizationUsersRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListOrganizationUsersRequest) ProtoMessage() {}

func (x *ListOrganizationUsersRequest) ProtoReflect() protoreflect.Message {
	mi := &file_buf_alpha_registry_v1alpha1_user_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListOrganizationUsersRequest.ProtoReflect.Descriptor instead.
func (*ListOrganizationUsersRequest) Descriptor() ([]byte, []int) {
	return file_buf_alpha_registry_v1alpha1_user_proto_rawDescGZIP(), []int{10}
}

func (x *ListOrganizationUsersRequest) GetOrganizationId() string {
	if x != nil {
		return x.OrganizationId
	}
	return ""
}

func (x *ListOrganizationUsersRequest) GetPageSize() uint32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *ListOrganizationUsersRequest) GetPageToken() string {
	if x != nil {
		return x.PageToken
	}
	return ""
}

func (x *ListOrganizationUsersRequest) GetReverse() bool {
	if x != nil {
		return x.Reverse
	}
	return false
}

type ListOrganizationUsersResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Users []*OrganizationUser `protobuf:"bytes,1,rep,name=users,proto3" json:"users,omitempty"`
	// There are no more pages if this is empty.
	NextPageToken string `protobuf:"bytes,2,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"`
}

func (x *ListOrganizationUsersResponse) Reset() {
	*x = ListOrganizationUsersResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_buf_alpha_registry_v1alpha1_user_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListOrganizationUsersResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListOrganizationUsersResponse) ProtoMessage() {}

func (x *ListOrganizationUsersResponse) ProtoReflect() protoreflect.Message {
	mi := &file_buf_alpha_registry_v1alpha1_user_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListOrganizationUsersResponse.ProtoReflect.Descriptor instead.
func (*ListOrganizationUsersResponse) Descriptor() ([]byte, []int) {
	return file_buf_alpha_registry_v1alpha1_user_proto_rawDescGZIP(), []int{11}
}

func (x *ListOrganizationUsersResponse) GetUsers() []*OrganizationUser {
	if x != nil {
		return x.Users
	}
	return nil
}

func (x *ListOrganizationUsersResponse) GetNextPageToken() string {
	if x != nil {
		return x.NextPageToken
	}
	return ""
}

type DeleteUserRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeleteUserRequest) Reset() {
	*x = DeleteUserRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_buf_alpha_registry_v1alpha1_user_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteUserRequest) ProtoMessage() {}

func (x *DeleteUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_buf_alpha_registry_v1alpha1_user_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteUserRequest.ProtoReflect.Descriptor instead.
func (*DeleteUserRequest) Descriptor() ([]byte, []int) {
	return file_buf_alpha_registry_v1alpha1_user_proto_rawDescGZIP(), []int{12}
}

type DeleteUserResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeleteUserResponse) Reset() {
	*x = DeleteUserResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_buf_alpha_registry_v1alpha1_user_proto_msgTypes[13]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteUserResponse) ProtoMessage() {}

func (x *DeleteUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_buf_alpha_registry_v1alpha1_user_proto_msgTypes[13]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteUserResponse.ProtoReflect.Descriptor instead.
func (*DeleteUserResponse) Descriptor() ([]byte, []int) {
	return file_buf_alpha_registry_v1alpha1_user_proto_rawDescGZIP(), []int{13}
}

type DeactivateUserRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *DeactivateUserRequest) Reset() {
	*x = DeactivateUserRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_buf_alpha_registry_v1alpha1_user_proto_msgTypes[14]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeactivateUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeactivateUserRequest) ProtoMessage() {}

func (x *DeactivateUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_buf_alpha_registry_v1alpha1_user_proto_msgTypes[14]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeactivateUserRequest.ProtoReflect.Descriptor instead.
func (*DeactivateUserRequest) Descriptor() ([]byte, []int) {
	return file_buf_alpha_registry_v1alpha1_user_proto_rawDescGZIP(), []int{14}
}

func (x *DeactivateUserRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type DeactivateUserResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeactivateUserResponse) Reset() {
	*x = DeactivateUserResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_buf_alpha_registry_v1alpha1_user_proto_msgTypes[15]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeactivateUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeactivateUserResponse) ProtoMessage() {}

func (x *DeactivateUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_buf_alpha_registry_v1alpha1_user_proto_msgTypes[15]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeactivateUserResponse.ProtoReflect.Descriptor instead.
func (*DeactivateUserResponse) Descriptor() ([]byte, []int) {
	return file_buf_alpha_registry_v1alpha1_user_proto_rawDescGZIP(), []int{15}
}

type UpdateUserServerRoleRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The ID of the user for which to be updated a role.
	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	// The new role of the user in the server.
	ServerRole ServerRole `protobuf:"varint,2,opt,name=server_role,json=serverRole,proto3,enum=buf.alpha.registry.v1alpha1.ServerRole" json:"server_role,omitempty"`
}

func (x *UpdateUserServerRoleRequest) Reset() {
	*x = UpdateUserServerRoleRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_buf_alpha_registry_v1alpha1_user_proto_msgTypes[16]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateUserServerRoleRequest) String() string {
	return pr