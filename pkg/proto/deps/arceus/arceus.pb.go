// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        (unknown)
// source: proto/deps/arceus/arceus.proto

package suggest

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type SortType int32

const (
	SortType_SORT_TYPE_UNKNOWN SortType = 0
	SortType_SORT_TYPE_ASC     SortType = 1
	SortType_SORT_TYPE_DESC    SortType = 2
)

// Enum value maps for SortType.
var (
	SortType_name = map[int32]string{
		0: "SORT_TYPE_UNKNOWN",
		1: "SORT_TYPE_ASC",
		2: "SORT_TYPE_DESC",
	}
	SortType_value = map[string]int32{
		"SORT_TYPE_UNKNOWN": 0,
		"SORT_TYPE_ASC":     1,
		"SORT_TYPE_DESC":    2,
	}
)

func (x SortType) Enum() *SortType {
	p := new(SortType)
	*p = x
	return p
}

func (x SortType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SortType) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_deps_arceus_arceus_proto_enumTypes[0].Descriptor()
}

func (SortType) Type() protoreflect.EnumType {
	return &file_proto_deps_arceus_arceus_proto_enumTypes[0]
}

func (x SortType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SortType.Descriptor instead.
func (SortType) EnumDescriptor() ([]byte, []int) {
	return file_proto_deps_arceus_arceus_proto_rawDescGZIP(), []int{0}
}

type Role int32

const (
	Role_ROLE_UNKNOWN Role = 0
	Role_ROLE_BOT     Role = 1
	Role_ROLE_USER    Role = 2
)

// Enum value maps for Role.
var (
	Role_name = map[int32]string{
		0: "ROLE_UNKNOWN",
		1: "ROLE_BOT",
		2: "ROLE_USER",
	}
	Role_value = map[string]int32{
		"ROLE_UNKNOWN": 0,
		"ROLE_BOT":     1,
		"ROLE_USER":    2,
	}
)

func (x Role) Enum() *Role {
	p := new(Role)
	*p = x
	return p
}

func (x Role) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Role) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_deps_arceus_arceus_proto_enumTypes[1].Descriptor()
}

func (Role) Type() protoreflect.EnumType {
	return &file_proto_deps_arceus_arceus_proto_enumTypes[1]
}

func (x Role) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Role.Descriptor instead.
func (Role) EnumDescriptor() ([]byte, []int) {
	return file_proto_deps_arceus_arceus_proto_rawDescGZIP(), []int{1}
}

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Content string `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
	Role    Role   `protobuf:"varint,2,opt,name=role,proto3,enum=arceus.Role" json:"role,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_deps_arceus_arceus_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_proto_deps_arceus_arceus_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_proto_deps_arceus_arceus_proto_rawDescGZIP(), []int{0}
}

func (x *Message) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *Message) GetRole() Role {
	if x != nil {
		return x.Role
	}
	return Role_ROLE_UNKNOWN
}

type Conversation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       uint64     `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Messages []*Message `protobuf:"bytes,2,rep,name=messages,proto3" json:"messages,omitempty"`
}

func (x *Conversation) Reset() {
	*x = Conversation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_deps_arceus_arceus_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Conversation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Conversation) ProtoMessage() {}

func (x *Conversation) ProtoReflect() protoreflect.Message {
	mi := &file_proto_deps_arceus_arceus_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Conversation.ProtoReflect.Descriptor instead.
func (*Conversation) Descriptor() ([]byte, []int) {
	return file_proto_deps_arceus_arceus_proto_rawDescGZIP(), []int{1}
}

func (x *Conversation) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Conversation) GetMessages() []*Message {
	if x != nil {
		return x.Messages
	}
	return nil
}

type GenerateTextRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Content        string  `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
	Model          string  `protobuf:"bytes,2,opt,name=model,proto3" json:"model,omitempty"`
	ConversationId *uint64 `protobuf:"varint,3,opt,name=conversation_id,json=conversationId,proto3,oneof" json:"conversation_id,omitempty"`
}

func (x *GenerateTextRequest) Reset() {
	*x = GenerateTextRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_deps_arceus_arceus_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GenerateTextRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GenerateTextRequest) ProtoMessage() {}

func (x *GenerateTextRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_deps_arceus_arceus_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GenerateTextRequest.ProtoReflect.Descriptor instead.
func (*GenerateTextRequest) Descriptor() ([]byte, []int) {
	return file_proto_deps_arceus_arceus_proto_rawDescGZIP(), []int{2}
}

func (x *GenerateTextRequest) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *GenerateTextRequest) GetModel() string {
	if x != nil {
		return x.Model
	}
	return ""
}

func (x *GenerateTextRequest) GetConversationId() uint64 {
	if x != nil && x.ConversationId != nil {
		return *x.ConversationId
	}
	return 0
}

type GenerateTextResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Content        string                 `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
	ConversationId uint64                 `protobuf:"varint,2,opt,name=conversation_id,json=conversationId,proto3" json:"conversation_id,omitempty"`
	CreatedAt      *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
}

func (x *GenerateTextResponse) Reset() {
	*x = GenerateTextResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_deps_arceus_arceus_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GenerateTextResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GenerateTextResponse) ProtoMessage() {}

func (x *GenerateTextResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_deps_arceus_arceus_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GenerateTextResponse.ProtoReflect.Descriptor instead.
func (*GenerateTextResponse) Descriptor() ([]byte, []int) {
	return file_proto_deps_arceus_arceus_proto_rawDescGZIP(), []int{3}
}

func (x *GenerateTextResponse) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *GenerateTextResponse) GetConversationId() uint64 {
	if x != nil {
		return x.ConversationId
	}
	return 0
}

func (x *GenerateTextResponse) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

var File_proto_deps_arceus_arceus_proto protoreflect.FileDescriptor

var file_proto_deps_arceus_arceus_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x64, 0x65, 0x70, 0x73, 0x2f, 0x61, 0x72, 0x63,
	0x65, 0x75, 0x73, 0x2f, 0x61, 0x72, 0x63, 0x65, 0x75, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x06, 0x61, 0x72, 0x63, 0x65, 0x75, 0x73, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x45, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x20, 0x0a, 0x04,
	0x72, 0x6f, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0c, 0x2e, 0x61, 0x72, 0x63,
	0x65, 0x75, 0x73, 0x2e, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x22, 0x4b,
	0x0a, 0x0c, 0x43, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x2b,
	0x0a, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x0f, 0x2e, 0x61, 0x72, 0x63, 0x65, 0x75, 0x73, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x52, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x22, 0x87, 0x01, 0x0a, 0x13,
	0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x54, 0x65, 0x78, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x14, 0x0a,
	0x05, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6d, 0x6f,
	0x64, 0x65, 0x6c, 0x12, 0x2c, 0x0a, 0x0f, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x48, 0x00, 0x52, 0x0e,
	0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x88, 0x01,
	0x01, 0x42, 0x12, 0x0a, 0x10, 0x5f, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x22, 0x94, 0x01, 0x0a, 0x14, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61,
	0x74, 0x65, 0x54, 0x65, 0x78, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x27, 0x0a, 0x0f, 0x63, 0x6f, 0x6e, 0x76,
	0x65, 0x72, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x0e, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49,
	0x64, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x2a, 0x48, 0x0a, 0x08,
	0x53, 0x6f, 0x72, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x15, 0x0a, 0x11, 0x53, 0x4f, 0x52, 0x54,
	0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12,
	0x11, 0x0a, 0x0d, 0x53, 0x4f, 0x52, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x41, 0x53, 0x43,
	0x10, 0x01, 0x12, 0x12, 0x0a, 0x0e, 0x53, 0x4f, 0x52, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f,
	0x44, 0x45, 0x53, 0x43, 0x10, 0x02, 0x2a, 0x35, 0x0a, 0x04, 0x52, 0x6f, 0x6c, 0x65, 0x12, 0x10,
	0x0a, 0x0c, 0x52, 0x4f, 0x4c, 0x45, 0x5f, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00,
	0x12, 0x0c, 0x0a, 0x08, 0x52, 0x4f, 0x4c, 0x45, 0x5f, 0x42, 0x4f, 0x54, 0x10, 0x01, 0x12, 0x0d,
	0x0a, 0x09, 0x52, 0x4f, 0x4c, 0x45, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x10, 0x02, 0x32, 0x75, 0x0a,
	0x06, 0x41, 0x72, 0x63, 0x65, 0x75, 0x73, 0x12, 0x6b, 0x0a, 0x0c, 0x47, 0x65, 0x6e, 0x65, 0x72,
	0x61, 0x74, 0x65, 0x54, 0x65, 0x78, 0x74, 0x12, 0x1b, 0x2e, 0x61, 0x72, 0x63, 0x65, 0x75, 0x73,
	0x2e, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x54, 0x65, 0x78, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x61, 0x72, 0x63, 0x65, 0x75, 0x73, 0x2e, 0x47, 0x65,
	0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x54, 0x65, 0x78, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x20, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1a, 0x3a, 0x01, 0x2a, 0x22, 0x15, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x63, 0x68, 0x61, 0x74, 0x2f, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x42, 0x17, 0x5a, 0x15, 0x6d, 0x79, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63,
	0x74, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x73, 0x75, 0x67, 0x67, 0x65, 0x73, 0x74, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_deps_arceus_arceus_proto_rawDescOnce sync.Once
	file_proto_deps_arceus_arceus_proto_rawDescData = file_proto_deps_arceus_arceus_proto_rawDesc
)

func file_proto_deps_arceus_arceus_proto_rawDescGZIP() []byte {
	file_proto_deps_arceus_arceus_proto_rawDescOnce.Do(func() {
		file_proto_deps_arceus_arceus_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_deps_arceus_arceus_proto_rawDescData)
	})
	return file_proto_deps_arceus_arceus_proto_rawDescData
}

var file_proto_deps_arceus_arceus_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_proto_deps_arceus_arceus_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_proto_deps_arceus_arceus_proto_goTypes = []interface{}{
	(SortType)(0),                 // 0: arceus.SortType
	(Role)(0),                     // 1: arceus.Role
	(*Message)(nil),               // 2: arceus.Message
	(*Conversation)(nil),          // 3: arceus.Conversation
	(*GenerateTextRequest)(nil),   // 4: arceus.GenerateTextRequest
	(*GenerateTextResponse)(nil),  // 5: arceus.GenerateTextResponse
	(*timestamppb.Timestamp)(nil), // 6: google.protobuf.Timestamp
}
var file_proto_deps_arceus_arceus_proto_depIdxs = []int32{
	1, // 0: arceus.Message.role:type_name -> arceus.Role
	2, // 1: arceus.Conversation.messages:type_name -> arceus.Message
	6, // 2: arceus.GenerateTextResponse.created_at:type_name -> google.protobuf.Timestamp
	4, // 3: arceus.Arceus.GenerateText:input_type -> arceus.GenerateTextRequest
	5, // 4: arceus.Arceus.GenerateText:output_type -> arceus.GenerateTextResponse
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_proto_deps_arceus_arceus_proto_init() }
func file_proto_deps_arceus_arceus_proto_init() {
	if File_proto_deps_arceus_arceus_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_deps_arceus_arceus_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_deps_arceus_arceus_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Conversation); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_deps_arceus_arceus_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GenerateTextRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_deps_arceus_arceus_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GenerateTextResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_proto_deps_arceus_arceus_proto_msgTypes[2].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_deps_arceus_arceus_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_deps_arceus_arceus_proto_goTypes,
		DependencyIndexes: file_proto_deps_arceus_arceus_proto_depIdxs,
		EnumInfos:         file_proto_deps_arceus_arceus_proto_enumTypes,
		MessageInfos:      file_proto_deps_arceus_arceus_proto_msgTypes,
	}.Build()
	File_proto_deps_arceus_arceus_proto = out.File
	file_proto_deps_arceus_arceus_proto_rawDesc = nil
	file_proto_deps_arceus_arceus_proto_goTypes = nil
	file_proto_deps_arceus_arceus_proto_depIdxs = nil
}
