// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.4
// source: learnfs.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// LearnFSClient is the client API for LearnFS service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LearnFSClient interface {
	CreateFile(ctx context.Context, in *CreateFileRequest, opts ...grpc.CallOption) (*CreateFileResponse, error)
	Flush(ctx context.Context, in *FlushRequest, opts ...grpc.CallOption) (*FlushResponse, error)
	GetAttr(ctx context.Context, in *GetAttrRequest, opts ...grpc.CallOption) (*GetAttrResponse, error)
	SetAttr(ctx context.Context, in *SetAttrRequest, opts ...grpc.CallOption) (*SetAttrResponse, error)
	Lookup(ctx context.Context, in *LookupRequest, opts ...grpc.CallOption) (*LookupResponse, error)
	OpenDir(ctx context.Context, opts ...grpc.CallOption) (LearnFS_OpenDirClient, error)
	RemoveFile(ctx context.Context, in *RemoveFileRequest, opts ...grpc.CallOption) (*RemoveFileResponse, error)
}

type learnFSClient struct {
	cc grpc.ClientConnInterface
}

func NewLearnFSClient(cc grpc.ClientConnInterface) LearnFSClient {
	return &learnFSClient{cc}
}

func (c *learnFSClient) CreateFile(ctx context.Context, in *CreateFileRequest, opts ...grpc.CallOption) (*CreateFileResponse, error) {
	out := new(CreateFileResponse)
	err := c.cc.Invoke(ctx, "/learnfs.LearnFS/CreateFile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *learnFSClient) Flush(ctx context.Context, in *FlushRequest, opts ...grpc.CallOption) (*FlushResponse, error) {
	out := new(FlushResponse)
	err := c.cc.Invoke(ctx, "/learnfs.LearnFS/Flush", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *learnFSClient) GetAttr(ctx context.Context, in *GetAttrRequest, opts ...grpc.CallOption) (*GetAttrResponse, error) {
	out := new(GetAttrResponse)
	err := c.cc.Invoke(ctx, "/learnfs.LearnFS/GetAttr", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *learnFSClient) SetAttr(ctx context.Context, in *SetAttrRequest, opts ...grpc.CallOption) (*SetAttrResponse, error) {
	out := new(SetAttrResponse)
	err := c.cc.Invoke(ctx, "/learnfs.LearnFS/SetAttr", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *learnFSClient) Lookup(ctx context.Context, in *LookupRequest, opts ...grpc.CallOption) (*LookupResponse, error) {
	out := new(LookupResponse)
	err := c.cc.Invoke(ctx, "/learnfs.LearnFS/Lookup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *learnFSClient) OpenDir(ctx context.Context, opts ...grpc.CallOption) (LearnFS_OpenDirClient, error) {
	stream, err := c.cc.NewStream(ctx, &LearnFS_ServiceDesc.Streams[0], "/learnfs.LearnFS/OpenDir", opts...)
	if err != nil {
		return nil, err
	}
	x := &learnFSOpenDirClient{stream}
	return x, nil
}

type LearnFS_OpenDirClient interface {
	Send(*NextDirEntryRequest) error
	Recv() (*DirEntryResponse, error)
	grpc.ClientStream
}

type learnFSOpenDirClient struct {
	grpc.ClientStream
}

func (x *learnFSOpenDirClient) Send(m *NextDirEntryRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *learnFSOpenDirClient) Recv() (*DirEntryResponse, error) {
	m := new(DirEntryResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *learnFSClient) RemoveFile(ctx context.Context, in *RemoveFileRequest, opts ...grpc.CallOption) (*RemoveFileResponse, error) {
	out := new(RemoveFileResponse)
	err := c.cc.Invoke(ctx, "/learnfs.LearnFS/RemoveFile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LearnFSServer is the server API for LearnFS service.
// All implementations must embed UnimplementedLearnFSServer
// for forward compatibility
type LearnFSServer interface {
	CreateFile(context.Context, *CreateFileRequest) (*CreateFileResponse, error)
	Flush(context.Context, *FlushRequest) (*FlushResponse, error)
	GetAttr(context.Context, *GetAttrRequest) (*GetAttrResponse, error)
	SetAttr(context.Context, *SetAttrRequest) (*SetAttrResponse, error)
	Lookup(context.Context, *LookupRequest) (*LookupResponse, error)
	OpenDir(LearnFS_OpenDirServer) error
	RemoveFile(context.Context, *RemoveFileRequest) (*RemoveFileResponse, error)
	mustEmbedUnimplementedLearnFSServer()
}

// UnimplementedLearnFSServer must be embedded to have forward compatible implementations.
type UnimplementedLearnFSServer struct {
}

func (UnimplementedLearnFSServer) CreateFile(context.Context, *CreateFileRequest) (*CreateFileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateFile not implemented")
}
func (UnimplementedLearnFSServer) Flush(context.Context, *FlushRequest) (*FlushResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Flush not implemented")
}
func (UnimplementedLearnFSServer) GetAttr(context.Context, *GetAttrRequest) (*GetAttrResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAttr not implemented")
}
func (UnimplementedLearnFSServer) SetAttr(context.Context, *SetAttrRequest) (*SetAttrResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetAttr not implemented")
}
func (UnimplementedLearnFSServer) Lookup(context.Context, *LookupRequest) (*LookupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Lookup not implemented")
}
func (UnimplementedLearnFSServer) OpenDir(LearnFS_OpenDirServer) error {
	return status.Errorf(codes.Unimplemented, "method OpenDir not implemented")
}
func (UnimplementedLearnFSServer) RemoveFile(context.Context, *RemoveFileRequest) (*RemoveFileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveFile not implemented")
}
func (UnimplementedLearnFSServer) mustEmbedUnimplementedLearnFSServer() {}

// UnsafeLearnFSServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LearnFSServer will
// result in compilation errors.
type UnsafeLearnFSServer interface {
	mustEmbedUnimplementedLearnFSServer()
}

func RegisterLearnFSServer(s grpc.ServiceRegistrar, srv LearnFSServer) {
	s.RegisterService(&LearnFS_ServiceDesc, srv)
}

func _LearnFS_CreateFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LearnFSServer).CreateFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/learnfs.LearnFS/CreateFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LearnFSServer).CreateFile(ctx, req.(*CreateFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LearnFS_Flush_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FlushRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LearnFSServer).Flush(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/learnfs.LearnFS/Flush",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LearnFSServer).Flush(ctx, req.(*FlushRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LearnFS_GetAttr_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAttrRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LearnFSServer).GetAttr(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/learnfs.LearnFS/GetAttr",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LearnFSServer).GetAttr(ctx, req.(*GetAttrRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LearnFS_SetAttr_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetAttrRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LearnFSServer).SetAttr(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/learnfs.LearnFS/SetAttr",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LearnFSServer).SetAttr(ctx, req.(*SetAttrRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LearnFS_Lookup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LookupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LearnFSServer).Lookup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/learnfs.LearnFS/Lookup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LearnFSServer).Lookup(ctx, req.(*LookupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LearnFS_OpenDir_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(LearnFSServer).OpenDir(&learnFSOpenDirServer{stream})
}

type LearnFS_OpenDirServer interface {
	Send(*DirEntryResponse) error
	Recv() (*NextDirEntryRequest, error)
	grpc.ServerStream
}

type learnFSOpenDirServer struct {
	grpc.ServerStream
}

func (x *learnFSOpenDirServer) Send(m *DirEntryResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *learnFSOpenDirServer) Recv() (*NextDirEntryRequest, error) {
	m := new(NextDirEntryRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _LearnFS_RemoveFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LearnFSServer).RemoveFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/learnfs.LearnFS/RemoveFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LearnFSServer).RemoveFile(ctx, req.(*RemoveFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// LearnFS_ServiceDesc is the grpc.ServiceDesc for LearnFS service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LearnFS_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "learnfs.LearnFS",
	HandlerType: (*LearnFSServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateFile",
			Handler:    _LearnFS_CreateFile_Handler,
		},
		{
			MethodName: "Flush",
			Handler:    _LearnFS_Flush_Handler,
		},
		{
			MethodName: "GetAttr",
			Handler:    _LearnFS_GetAttr_Handler,
		},
		{
			MethodName: "SetAttr",
			Handler:    _LearnFS_SetAttr_Handler,
		},
		{
			MethodName: "Lookup",
			Handler:    _LearnFS_Lookup_Handler,
		},
		{
			MethodName: "RemoveFile",
			Handler:    _LearnFS_RemoveFile_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "OpenDir",
			Handler:       _LearnFS_OpenDir_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "learnfs.proto",
}
