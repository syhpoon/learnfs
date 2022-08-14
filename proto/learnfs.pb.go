// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.21.4
// source: learnfs.proto

package proto

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

// Inode attributes
type Attr struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Inode number
	Ino uint32 `protobuf:"varint,1,opt,name=ino,proto3" json:"ino,omitempty"`
	// Size in bytes
	Size uint32 `protobuf:"varint,2,opt,name=size,proto3" json:"size,omitempty"`
	// Size in blocks
	Blocks uint32 `protobuf:"varint,3,opt,name=blocks,proto3" json:"blocks,omitempty"`
	// Access time
	Atime *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=atime,proto3" json:"atime,omitempty"`
	// Modification time
	Mtime *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=mtime,proto3" json:"mtime,omitempty"`
	// Status change time
	Ctime *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=ctime,proto3" json:"ctime,omitempty"`
	// Inode kind
	//
	// Types that are assignable to Kind:
	//	*Attr_Unknown
	//	*Attr_Fifo
	//	*Attr_Char
	//	*Attr_Block
	//	*Attr_Directory
	//	*Attr_File
	//	*Attr_Link
	//	*Attr_Socket
	Kind isAttr_Kind `protobuf_oneof:"kind"`
	// File mode
	Mode uint32 `protobuf:"varint,8,opt,name=mode,proto3" json:"mode,omitempty"`
	/// Number of hard links
	Nlink uint32 `protobuf:"varint,9,opt,name=nlink,proto3" json:"nlink,omitempty"`
	/// User id
	Uid uint32 `protobuf:"varint,10,opt,name=uid,proto3" json:"uid,omitempty"`
	/// Group id
	Gid uint32 `protobuf:"varint,11,opt,name=gid,proto3" json:"gid,omitempty"`
}

func (x *Attr) Reset() {
	*x = Attr{}
	if protoimpl.UnsafeEnabled {
		mi := &file_learnfs_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Attr) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Attr) ProtoMessage() {}

func (x *Attr) ProtoReflect() protoreflect.Message {
	mi := &file_learnfs_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Attr.ProtoReflect.Descriptor instead.
func (*Attr) Descriptor() ([]byte, []int) {
	return file_learnfs_proto_rawDescGZIP(), []int{0}
}

func (x *Attr) GetIno() uint32 {
	if x != nil {
		return x.Ino
	}
	return 0
}

func (x *Attr) GetSize() uint32 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *Attr) GetBlocks() uint32 {
	if x != nil {
		return x.Blocks
	}
	return 0
}

func (x *Attr) GetAtime() *timestamppb.Timestamp {
	if x != nil {
		return x.Atime
	}
	return nil
}

func (x *Attr) GetMtime() *timestamppb.Timestamp {
	if x != nil {
		return x.Mtime
	}
	return nil
}

func (x *Attr) GetCtime() *timestamppb.Timestamp {
	if x != nil {
		return x.Ctime
	}
	return nil
}

func (m *Attr) GetKind() isAttr_Kind {
	if m != nil {
		return m.Kind
	}
	return nil
}

func (x *Attr) GetUnknown() bool {
	if x, ok := x.GetKind().(*Attr_Unknown); ok {
		return x.Unknown
	}
	return false
}

func (x *Attr) GetFifo() bool {
	if x, ok := x.GetKind().(*Attr_Fifo); ok {
		return x.Fifo
	}
	return false
}

func (x *Attr) GetChar() bool {
	if x, ok := x.GetKind().(*Attr_Char); ok {
		return x.Char
	}
	return false
}

func (x *Attr) GetBlock() bool {
	if x, ok := x.GetKind().(*Attr_Block); ok {
		return x.Block
	}
	return false
}

func (x *Attr) GetDirectory() bool {
	if x, ok := x.GetKind().(*Attr_Directory); ok {
		return x.Directory
	}
	return false
}

func (x *Attr) GetFile() bool {
	if x, ok := x.GetKind().(*Attr_File); ok {
		return x.File
	}
	return false
}

func (x *Attr) GetLink() bool {
	if x, ok := x.GetKind().(*Attr_Link); ok {
		return x.Link
	}
	return false
}

func (x *Attr) GetSocket() bool {
	if x, ok := x.GetKind().(*Attr_Socket); ok {
		return x.Socket
	}
	return false
}

func (x *Attr) GetMode() uint32 {
	if x != nil {
		return x.Mode
	}
	return 0
}

func (x *Attr) GetNlink() uint32 {
	if x != nil {
		return x.Nlink
	}
	return 0
}

func (x *Attr) GetUid() uint32 {
	if x != nil {
		return x.Uid
	}
	return 0
}

func (x *Attr) GetGid() uint32 {
	if x != nil {
		return x.Gid
	}
	return 0
}

type isAttr_Kind interface {
	isAttr_Kind()
}

type Attr_Unknown struct {
	Unknown bool `protobuf:"varint,20,opt,name=unknown,proto3,oneof"`
}

type Attr_Fifo struct {
	Fifo bool `protobuf:"varint,21,opt,name=fifo,proto3,oneof"`
}

type Attr_Char struct {
	Char bool `protobuf:"varint,22,opt,name=char,proto3,oneof"`
}

type Attr_Block struct {
	Block bool `protobuf:"varint,23,opt,name=block,proto3,oneof"`
}

type Attr_Directory struct {
	Directory bool `protobuf:"varint,24,opt,name=directory,proto3,oneof"`
}

type Attr_File struct {
	File bool `protobuf:"varint,25,opt,name=file,proto3,oneof"`
}

type Attr_Link struct {
	Link bool `protobuf:"varint,26,opt,name=link,proto3,oneof"`
}

type Attr_Socket struct {
	Socket bool `protobuf:"varint,27,opt,name=socket,proto3,oneof"`
}

func (*Attr_Unknown) isAttr_Kind() {}

func (*Attr_Fifo) isAttr_Kind() {}

func (*Attr_Char) isAttr_Kind() {}

func (*Attr_Block) isAttr_Kind() {}

func (*Attr_Directory) isAttr_Kind() {}

func (*Attr_File) isAttr_Kind() {}

func (*Attr_Link) isAttr_Kind() {}

func (*Attr_Socket) isAttr_Kind() {}

// RPC messages
type GetAttrRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Inode uint32 `protobuf:"varint,1,opt,name=inode,proto3" json:"inode,omitempty"`
}

func (x *GetAttrRequest) Reset() {
	*x = GetAttrRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_learnfs_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAttrRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAttrRequest) ProtoMessage() {}

func (x *GetAttrRequest) ProtoReflect() protoreflect.Message {
	mi := &file_learnfs_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAttrRequest.ProtoReflect.Descriptor instead.
func (*GetAttrRequest) Descriptor() ([]byte, []int) {
	return file_learnfs_proto_rawDescGZIP(), []int{1}
}

func (x *GetAttrRequest) GetInode() uint32 {
	if x != nil {
		return x.Inode
	}
	return 0
}

type GetAttrResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Errno uint32 `protobuf:"varint,1,opt,name=errno,proto3" json:"errno,omitempty"`
	Attr  *Attr  `protobuf:"bytes,2,opt,name=attr,proto3" json:"attr,omitempty"`
}

func (x *GetAttrResponse) Reset() {
	*x = GetAttrResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_learnfs_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAttrResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAttrResponse) ProtoMessage() {}

func (x *GetAttrResponse) ProtoReflect() protoreflect.Message {
	mi := &file_learnfs_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAttrResponse.ProtoReflect.Descriptor instead.
func (*GetAttrResponse) Descriptor() ([]byte, []int) {
	return file_learnfs_proto_rawDescGZIP(), []int{2}
}

func (x *GetAttrResponse) GetErrno() uint32 {
	if x != nil {
		return x.Errno
	}
	return 0
}

func (x *GetAttrResponse) GetAttr() *Attr {
	if x != nil {
		return x.Attr
	}
	return nil
}

type LookupRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Inode uint32 `protobuf:"varint,1,opt,name=inode,proto3" json:"inode,omitempty"`
	Name  string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *LookupRequest) Reset() {
	*x = LookupRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_learnfs_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LookupRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LookupRequest) ProtoMessage() {}

func (x *LookupRequest) ProtoReflect() protoreflect.Message {
	mi := &file_learnfs_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LookupRequest.ProtoReflect.Descriptor instead.
func (*LookupRequest) Descriptor() ([]byte, []int) {
	return file_learnfs_proto_rawDescGZIP(), []int{3}
}

func (x *LookupRequest) GetInode() uint32 {
	if x != nil {
		return x.Inode
	}
	return 0
}

func (x *LookupRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type LookupResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Errno uint32 `protobuf:"varint,1,opt,name=errno,proto3" json:"errno,omitempty"`
	Attr  *Attr  `protobuf:"bytes,2,opt,name=attr,proto3" json:"attr,omitempty"`
}

func (x *LookupResponse) Reset() {
	*x = LookupResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_learnfs_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LookupResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LookupResponse) ProtoMessage() {}

func (x *LookupResponse) ProtoReflect() protoreflect.Message {
	mi := &file_learnfs_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LookupResponse.ProtoReflect.Descriptor instead.
func (*LookupResponse) Descriptor() ([]byte, []int) {
	return file_learnfs_proto_rawDescGZIP(), []int{4}
}

func (x *LookupResponse) GetErrno() uint32 {
	if x != nil {
		return x.Errno
	}
	return 0
}

func (x *LookupResponse) GetAttr() *Attr {
	if x != nil {
		return x.Attr
	}
	return nil
}

type NextDirEntryRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Inode uint32 `protobuf:"varint,1,opt,name=inode,proto3" json:"inode,omitempty"`
}

func (x *NextDirEntryRequest) Reset() {
	*x = NextDirEntryRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_learnfs_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NextDirEntryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NextDirEntryRequest) ProtoMessage() {}

func (x *NextDirEntryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_learnfs_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NextDirEntryRequest.ProtoReflect.Descriptor instead.
func (*NextDirEntryRequest) Descriptor() ([]byte, []int) {
	return file_learnfs_proto_rawDescGZIP(), []int{5}
}

func (x *NextDirEntryRequest) GetInode() uint32 {
	if x != nil {
		return x.Inode
	}
	return 0
}

type DirEntryResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Errno uint32                     `protobuf:"varint,1,opt,name=errno,proto3" json:"errno,omitempty"`
	Entry *DirEntryResponse_DirEntry `protobuf:"bytes,2,opt,name=entry,proto3" json:"entry,omitempty"`
}

func (x *DirEntryResponse) Reset() {
	*x = DirEntryResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_learnfs_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DirEntryResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DirEntryResponse) ProtoMessage() {}

func (x *DirEntryResponse) ProtoReflect() protoreflect.Message {
	mi := &file_learnfs_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DirEntryResponse.ProtoReflect.Descriptor instead.
func (*DirEntryResponse) Descriptor() ([]byte, []int) {
	return file_learnfs_proto_rawDescGZIP(), []int{6}
}

func (x *DirEntryResponse) GetErrno() uint32 {
	if x != nil {
		return x.Errno
	}
	return 0
}

func (x *DirEntryResponse) GetEntry() *DirEntryResponse_DirEntry {
	if x != nil {
		return x.Entry
	}
	return nil
}

type DirEntryResponse_DirEntry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Inode uint32 `protobuf:"varint,1,opt,name=inode,proto3" json:"inode,omitempty"`
	Name  string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Mode  uint32 `protobuf:"varint,3,opt,name=mode,proto3" json:"mode,omitempty"`
}

func (x *DirEntryResponse_DirEntry) Reset() {
	*x = DirEntryResponse_DirEntry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_learnfs_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DirEntryResponse_DirEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DirEntryResponse_DirEntry) ProtoMessage() {}

func (x *DirEntryResponse_DirEntry) ProtoReflect() protoreflect.Message {
	mi := &file_learnfs_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DirEntryResponse_DirEntry.ProtoReflect.Descriptor instead.
func (*DirEntryResponse_DirEntry) Descriptor() ([]byte, []int) {
	return file_learnfs_proto_rawDescGZIP(), []int{6, 0}
}

func (x *DirEntryResponse_DirEntry) GetInode() uint32 {
	if x != nil {
		return x.Inode
	}
	return 0
}

func (x *DirEntryResponse_DirEntry) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *DirEntryResponse_DirEntry) GetMode() uint32 {
	if x != nil {
		return x.Mode
	}
	return 0
}

var File_learnfs_proto protoreflect.FileDescriptor

var file_learnfs_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x6c, 0x65, 0x61, 0x72, 0x6e, 0x66, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x6c, 0x65, 0x61, 0x72, 0x6e, 0x66, 0x73, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xf6, 0x03, 0x0a, 0x04, 0x41, 0x74,
	0x74, 0x72, 0x12, 0x10, 0x0a, 0x03, 0x69, 0x6e, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x03, 0x69, 0x6e, 0x6f, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x62, 0x6c, 0x6f, 0x63,
	0x6b, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x73,
	0x12, 0x30, 0x0a, 0x05, 0x61, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x05, 0x61, 0x74, 0x69,
	0x6d, 0x65, 0x12, 0x30, 0x0a, 0x05, 0x6d, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x05, 0x6d,
	0x74, 0x69, 0x6d, 0x65, 0x12, 0x30, 0x0a, 0x05, 0x63, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x05, 0x63, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x07, 0x75, 0x6e, 0x6b, 0x6e, 0x6f, 0x77,
	0x6e, 0x18, 0x14, 0x20, 0x01, 0x28, 0x08, 0x48, 0x00, 0x52, 0x07, 0x75, 0x6e, 0x6b, 0x6e, 0x6f,
	0x77, 0x6e, 0x12, 0x14, 0x0a, 0x04, 0x66, 0x69, 0x66, 0x6f, 0x18, 0x15, 0x20, 0x01, 0x28, 0x08,
	0x48, 0x00, 0x52, 0x04, 0x66, 0x69, 0x66, 0x6f, 0x12, 0x14, 0x0a, 0x04, 0x63, 0x68, 0x61, 0x72,
	0x18, 0x16, 0x20, 0x01, 0x28, 0x08, 0x48, 0x00, 0x52, 0x04, 0x63, 0x68, 0x61, 0x72, 0x12, 0x16,
	0x0a, 0x05, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x18, 0x17, 0x20, 0x01, 0x28, 0x08, 0x48, 0x00, 0x52,
	0x05, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x1e, 0x0a, 0x09, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x79, 0x18, 0x18, 0x20, 0x01, 0x28, 0x08, 0x48, 0x00, 0x52, 0x09, 0x64, 0x69, 0x72,
	0x65, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x12, 0x14, 0x0a, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x19,
	0x20, 0x01, 0x28, 0x08, 0x48, 0x00, 0x52, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x12, 0x14, 0x0a, 0x04,
	0x6c, 0x69, 0x6e, 0x6b, 0x18, 0x1a, 0x20, 0x01, 0x28, 0x08, 0x48, 0x00, 0x52, 0x04, 0x6c, 0x69,
	0x6e, 0x6b, 0x12, 0x18, 0x0a, 0x06, 0x73, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x18, 0x1b, 0x20, 0x01,
	0x28, 0x08, 0x48, 0x00, 0x52, 0x06, 0x73, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x12, 0x0a, 0x04,
	0x6d, 0x6f, 0x64, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x6d, 0x6f, 0x64, 0x65,
	0x12, 0x14, 0x0a, 0x05, 0x6e, 0x6c, 0x69, 0x6e, 0x6b, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x05, 0x6e, 0x6c, 0x69, 0x6e, 0x6b, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18, 0x0a, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x03, 0x75, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x67, 0x69, 0x64, 0x18,
	0x0b, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x03, 0x67, 0x69, 0x64, 0x42, 0x06, 0x0a, 0x04, 0x6b, 0x69,
	0x6e, 0x64, 0x22, 0x26, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x41, 0x74, 0x74, 0x72, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x05, 0x69, 0x6e, 0x6f, 0x64, 0x65, 0x22, 0x4a, 0x0a, 0x0f, 0x47, 0x65,
	0x74, 0x41, 0x74, 0x74, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a,
	0x05, 0x65, 0x72, 0x72, 0x6e, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x65, 0x72,
	0x72, 0x6e, 0x6f, 0x12, 0x21, 0x0a, 0x04, 0x61, 0x74, 0x74, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x0d, 0x2e, 0x6c, 0x65, 0x61, 0x72, 0x6e, 0x66, 0x73, 0x2e, 0x41, 0x74, 0x74, 0x72,
	0x52, 0x04, 0x61, 0x74, 0x74, 0x72, 0x22, 0x39, 0x0a, 0x0d, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x6f, 0x64, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x69, 0x6e, 0x6f, 0x64, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x22, 0x49, 0x0a, 0x0e, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6e, 0x6f, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6e, 0x6f, 0x12, 0x21, 0x0a, 0x04, 0x61, 0x74, 0x74,
	0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x6c, 0x65, 0x61, 0x72, 0x6e, 0x66,
	0x73, 0x2e, 0x41, 0x74, 0x74, 0x72, 0x52, 0x04, 0x61, 0x74, 0x74, 0x72, 0x22, 0x2b, 0x0a, 0x13,
	0x4e, 0x65, 0x78, 0x74, 0x44, 0x69, 0x72, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x05, 0x69, 0x6e, 0x6f, 0x64, 0x65, 0x22, 0xac, 0x01, 0x0a, 0x10, 0x44, 0x69,
	0x72, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x65, 0x72, 0x72, 0x6e, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x65,
	0x72, 0x72, 0x6e, 0x6f, 0x12, 0x38, 0x0a, 0x05, 0x65, 0x6e, 0x74, 0x72, 0x79, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x6c, 0x65, 0x61, 0x72, 0x6e, 0x66, 0x73, 0x2e, 0x44, 0x69,
	0x72, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x44,
	0x69, 0x72, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x65, 0x6e, 0x74, 0x72, 0x79, 0x1a, 0x48,
	0x0a, 0x08, 0x44, 0x69, 0x72, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e,
	0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x69, 0x6e, 0x6f, 0x64, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6d, 0x6f, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x04, 0x6d, 0x6f, 0x64, 0x65, 0x32, 0xca, 0x01, 0x0a, 0x07, 0x4c, 0x65, 0x61,
	0x72, 0x6e, 0x46, 0x53, 0x12, 0x3c, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x41, 0x74, 0x74, 0x72, 0x12,
	0x17, 0x2e, 0x6c, 0x65, 0x61, 0x72, 0x6e, 0x66, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x74, 0x74,
	0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x6c, 0x65, 0x61, 0x72, 0x6e,
	0x66, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x74, 0x74, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x39, 0x0a, 0x06, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x12, 0x16, 0x2e, 0x6c,
	0x65, 0x61, 0x72, 0x6e, 0x66, 0x73, 0x2e, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x6c, 0x65, 0x61, 0x72, 0x6e, 0x66, 0x73, 0x2e, 0x4c,
	0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x46, 0x0a,
	0x07, 0x4f, 0x70, 0x65, 0x6e, 0x44, 0x69, 0x72, 0x12, 0x1c, 0x2e, 0x6c, 0x65, 0x61, 0x72, 0x6e,
	0x66, 0x73, 0x2e, 0x4e, 0x65, 0x78, 0x74, 0x44, 0x69, 0x72, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x6c, 0x65, 0x61, 0x72, 0x6e, 0x66, 0x73,
	0x2e, 0x44, 0x69, 0x72, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x28, 0x01, 0x30, 0x01, 0x42, 0x0f, 0x5a, 0x0d, 0x6c, 0x65, 0x61, 0x72, 0x6e, 0x66, 0x73,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_learnfs_proto_rawDescOnce sync.Once
	file_learnfs_proto_rawDescData = file_learnfs_proto_rawDesc
)

func file_learnfs_proto_rawDescGZIP() []byte {
	file_learnfs_proto_rawDescOnce.Do(func() {
		file_learnfs_proto_rawDescData = protoimpl.X.CompressGZIP(file_learnfs_proto_rawDescData)
	})
	return file_learnfs_proto_rawDescData
}

var file_learnfs_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_learnfs_proto_goTypes = []interface{}{
	(*Attr)(nil),                      // 0: learnfs.Attr
	(*GetAttrRequest)(nil),            // 1: learnfs.GetAttrRequest
	(*GetAttrResponse)(nil),           // 2: learnfs.GetAttrResponse
	(*LookupRequest)(nil),             // 3: learnfs.LookupRequest
	(*LookupResponse)(nil),            // 4: learnfs.LookupResponse
	(*NextDirEntryRequest)(nil),       // 5: learnfs.NextDirEntryRequest
	(*DirEntryResponse)(nil),          // 6: learnfs.DirEntryResponse
	(*DirEntryResponse_DirEntry)(nil), // 7: learnfs.DirEntryResponse.DirEntry
	(*timestamppb.Timestamp)(nil),     // 8: google.protobuf.Timestamp
}
var file_learnfs_proto_depIdxs = []int32{
	8, // 0: learnfs.Attr.atime:type_name -> google.protobuf.Timestamp
	8, // 1: learnfs.Attr.mtime:type_name -> google.protobuf.Timestamp
	8, // 2: learnfs.Attr.ctime:type_name -> google.protobuf.Timestamp
	0, // 3: learnfs.GetAttrResponse.attr:type_name -> learnfs.Attr
	0, // 4: learnfs.LookupResponse.attr:type_name -> learnfs.Attr
	7, // 5: learnfs.DirEntryResponse.entry:type_name -> learnfs.DirEntryResponse.DirEntry
	1, // 6: learnfs.LearnFS.GetAttr:input_type -> learnfs.GetAttrRequest
	3, // 7: learnfs.LearnFS.Lookup:input_type -> learnfs.LookupRequest
	5, // 8: learnfs.LearnFS.OpenDir:input_type -> learnfs.NextDirEntryRequest
	2, // 9: learnfs.LearnFS.GetAttr:output_type -> learnfs.GetAttrResponse
	4, // 10: learnfs.LearnFS.Lookup:output_type -> learnfs.LookupResponse
	6, // 11: learnfs.LearnFS.OpenDir:output_type -> learnfs.DirEntryResponse
	9, // [9:12] is the sub-list for method output_type
	6, // [6:9] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_learnfs_proto_init() }
func file_learnfs_proto_init() {
	if File_learnfs_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_learnfs_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Attr); i {
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
		file_learnfs_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAttrRequest); i {
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
		file_learnfs_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAttrResponse); i {
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
		file_learnfs_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LookupRequest); i {
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
		file_learnfs_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LookupResponse); i {
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
		file_learnfs_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NextDirEntryRequest); i {
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
		file_learnfs_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DirEntryResponse); i {
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
		file_learnfs_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DirEntryResponse_DirEntry); i {
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
	file_learnfs_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Attr_Unknown)(nil),
		(*Attr_Fifo)(nil),
		(*Attr_Char)(nil),
		(*Attr_Block)(nil),
		(*Attr_Directory)(nil),
		(*Attr_File)(nil),
		(*Attr_Link)(nil),
		(*Attr_Socket)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_learnfs_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_learnfs_proto_goTypes,
		DependencyIndexes: file_learnfs_proto_depIdxs,
		MessageInfos:      file_learnfs_proto_msgTypes,
	}.Build()
	File_learnfs_proto = out.File
	file_learnfs_proto_rawDesc = nil
	file_learnfs_proto_goTypes = nil
	file_learnfs_proto_depIdxs = nil
}
