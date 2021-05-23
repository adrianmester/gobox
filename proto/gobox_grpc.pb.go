// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

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

// GoBoxClient is the client API for GoBox service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GoBoxClient interface {
	GetLastUpdateTime(ctx context.Context, in *Null, opts ...grpc.CallOption) (*GetLastUpdateTimeResult, error)
	SendFileInfo(ctx context.Context, in *SendFileInfoInput, opts ...grpc.CallOption) (*SendFileInfoResponse, error)
	SendFileChunks(ctx context.Context, opts ...grpc.CallOption) (GoBox_SendFileChunksClient, error)
	SendFileChunksData(ctx context.Context, opts ...grpc.CallOption) (GoBox_SendFileChunksDataClient, error)
	InitialSyncComplete(ctx context.Context, in *Null, opts ...grpc.CallOption) (*Null, error)
	DeleteFile(ctx context.Context, in *DeleteFileInput, opts ...grpc.CallOption) (*Null, error)
}

type goBoxClient struct {
	cc grpc.ClientConnInterface
}

func NewGoBoxClient(cc grpc.ClientConnInterface) GoBoxClient {
	return &goBoxClient{cc}
}

func (c *goBoxClient) GetLastUpdateTime(ctx context.Context, in *Null, opts ...grpc.CallOption) (*GetLastUpdateTimeResult, error) {
	out := new(GetLastUpdateTimeResult)
	err := c.cc.Invoke(ctx, "/gobox.GoBox/GetLastUpdateTime", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *goBoxClient) SendFileInfo(ctx context.Context, in *SendFileInfoInput, opts ...grpc.CallOption) (*SendFileInfoResponse, error) {
	out := new(SendFileInfoResponse)
	err := c.cc.Invoke(ctx, "/gobox.GoBox/SendFileInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *goBoxClient) SendFileChunks(ctx context.Context, opts ...grpc.CallOption) (GoBox_SendFileChunksClient, error) {
	stream, err := c.cc.NewStream(ctx, &GoBox_ServiceDesc.Streams[0], "/gobox.GoBox/SendFileChunks", opts...)
	if err != nil {
		return nil, err
	}
	x := &goBoxSendFileChunksClient{stream}
	return x, nil
}

type GoBox_SendFileChunksClient interface {
	Send(*SendFileChunksInput) error
	CloseAndRecv() (*Null, error)
	grpc.ClientStream
}

type goBoxSendFileChunksClient struct {
	grpc.ClientStream
}

func (x *goBoxSendFileChunksClient) Send(m *SendFileChunksInput) error {
	return x.ClientStream.SendMsg(m)
}

func (x *goBoxSendFileChunksClient) CloseAndRecv() (*Null, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Null)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *goBoxClient) SendFileChunksData(ctx context.Context, opts ...grpc.CallOption) (GoBox_SendFileChunksDataClient, error) {
	stream, err := c.cc.NewStream(ctx, &GoBox_ServiceDesc.Streams[1], "/gobox.GoBox/SendFileChunksData", opts...)
	if err != nil {
		return nil, err
	}
	x := &goBoxSendFileChunksDataClient{stream}
	return x, nil
}

type GoBox_SendFileChunksDataClient interface {
	Send(*SendFileChunksDataInput) error
	Recv() (*SendFileChunksDataRequest, error)
	grpc.ClientStream
}

type goBoxSendFileChunksDataClient struct {
	grpc.ClientStream
}

func (x *goBoxSendFileChunksDataClient) Send(m *SendFileChunksDataInput) error {
	return x.ClientStream.SendMsg(m)
}

func (x *goBoxSendFileChunksDataClient) Recv() (*SendFileChunksDataRequest, error) {
	m := new(SendFileChunksDataRequest)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *goBoxClient) InitialSyncComplete(ctx context.Context, in *Null, opts ...grpc.CallOption) (*Null, error) {
	out := new(Null)
	err := c.cc.Invoke(ctx, "/gobox.GoBox/InitialSyncComplete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *goBoxClient) DeleteFile(ctx context.Context, in *DeleteFileInput, opts ...grpc.CallOption) (*Null, error) {
	out := new(Null)
	err := c.cc.Invoke(ctx, "/gobox.GoBox/DeleteFile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GoBoxServer is the server API for GoBox service.
// All implementations must embed UnimplementedGoBoxServer
// for forward compatibility
type GoBoxServer interface {
	GetLastUpdateTime(context.Context, *Null) (*GetLastUpdateTimeResult, error)
	SendFileInfo(context.Context, *SendFileInfoInput) (*SendFileInfoResponse, error)
	SendFileChunks(GoBox_SendFileChunksServer) error
	SendFileChunksData(GoBox_SendFileChunksDataServer) error
	InitialSyncComplete(context.Context, *Null) (*Null, error)
	DeleteFile(context.Context, *DeleteFileInput) (*Null, error)
	mustEmbedUnimplementedGoBoxServer()
}

// UnimplementedGoBoxServer must be embedded to have forward compatible implementations.
type UnimplementedGoBoxServer struct {
}

func (UnimplementedGoBoxServer) GetLastUpdateTime(context.Context, *Null) (*GetLastUpdateTimeResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLastUpdateTime not implemented")
}
func (UnimplementedGoBoxServer) SendFileInfo(context.Context, *SendFileInfoInput) (*SendFileInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendFileInfo not implemented")
}
func (UnimplementedGoBoxServer) SendFileChunks(GoBox_SendFileChunksServer) error {
	return status.Errorf(codes.Unimplemented, "method SendFileChunks not implemented")
}
func (UnimplementedGoBoxServer) SendFileChunksData(GoBox_SendFileChunksDataServer) error {
	return status.Errorf(codes.Unimplemented, "method SendFileChunksData not implemented")
}
func (UnimplementedGoBoxServer) InitialSyncComplete(context.Context, *Null) (*Null, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InitialSyncComplete not implemented")
}
func (UnimplementedGoBoxServer) DeleteFile(context.Context, *DeleteFileInput) (*Null, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteFile not implemented")
}
func (UnimplementedGoBoxServer) mustEmbedUnimplementedGoBoxServer() {}

// UnsafeGoBoxServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GoBoxServer will
// result in compilation errors.
type UnsafeGoBoxServer interface {
	mustEmbedUnimplementedGoBoxServer()
}

func RegisterGoBoxServer(s grpc.ServiceRegistrar, srv GoBoxServer) {
	s.RegisterService(&GoBox_ServiceDesc, srv)
}

func _GoBox_GetLastUpdateTime_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Null)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GoBoxServer).GetLastUpdateTime(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gobox.GoBox/GetLastUpdateTime",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GoBoxServer).GetLastUpdateTime(ctx, req.(*Null))
	}
	return interceptor(ctx, in, info, handler)
}

func _GoBox_SendFileInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendFileInfoInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GoBoxServer).SendFileInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gobox.GoBox/SendFileInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GoBoxServer).SendFileInfo(ctx, req.(*SendFileInfoInput))
	}
	return interceptor(ctx, in, info, handler)
}

func _GoBox_SendFileChunks_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GoBoxServer).SendFileChunks(&goBoxSendFileChunksServer{stream})
}

type GoBox_SendFileChunksServer interface {
	SendAndClose(*Null) error
	Recv() (*SendFileChunksInput, error)
	grpc.ServerStream
}

type goBoxSendFileChunksServer struct {
	grpc.ServerStream
}

func (x *goBoxSendFileChunksServer) SendAndClose(m *Null) error {
	return x.ServerStream.SendMsg(m)
}

func (x *goBoxSendFileChunksServer) Recv() (*SendFileChunksInput, error) {
	m := new(SendFileChunksInput)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _GoBox_SendFileChunksData_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GoBoxServer).SendFileChunksData(&goBoxSendFileChunksDataServer{stream})
}

type GoBox_SendFileChunksDataServer interface {
	Send(*SendFileChunksDataRequest) error
	Recv() (*SendFileChunksDataInput, error)
	grpc.ServerStream
}

type goBoxSendFileChunksDataServer struct {
	grpc.ServerStream
}

func (x *goBoxSendFileChunksDataServer) Send(m *SendFileChunksDataRequest) error {
	return x.ServerStream.SendMsg(m)
}

func (x *goBoxSendFileChunksDataServer) Recv() (*SendFileChunksDataInput, error) {
	m := new(SendFileChunksDataInput)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _GoBox_InitialSyncComplete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Null)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GoBoxServer).InitialSyncComplete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gobox.GoBox/InitialSyncComplete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GoBoxServer).InitialSyncComplete(ctx, req.(*Null))
	}
	return interceptor(ctx, in, info, handler)
}

func _GoBox_DeleteFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteFileInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GoBoxServer).DeleteFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gobox.GoBox/DeleteFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GoBoxServer).DeleteFile(ctx, req.(*DeleteFileInput))
	}
	return interceptor(ctx, in, info, handler)
}

// GoBox_ServiceDesc is the grpc.ServiceDesc for GoBox service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GoBox_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "gobox.GoBox",
	HandlerType: (*GoBoxServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetLastUpdateTime",
			Handler:    _GoBox_GetLastUpdateTime_Handler,
		},
		{
			MethodName: "SendFileInfo",
			Handler:    _GoBox_SendFileInfo_Handler,
		},
		{
			MethodName: "InitialSyncComplete",
			Handler:    _GoBox_InitialSyncComplete_Handler,
		},
		{
			MethodName: "DeleteFile",
			Handler:    _GoBox_DeleteFile_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SendFileChunks",
			Handler:       _GoBox_SendFileChunks_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "SendFileChunksData",
			Handler:       _GoBox_SendFileChunksData_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "proto/gobox.proto",
}
