// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: raftdkv/v1/kv.proto

package v1_raftdkv_v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	KVService_Set_FullMethodName    = "/raftdkv.v1.KVService/Set"
	KVService_Get_FullMethodName    = "/raftdkv.v1.KVService/Get"
	KVService_Delete_FullMethodName = "/raftdkv.v1.KVService/Delete"
)

// KVServiceClient is the client API for KVService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KVServiceClient interface {
	Set(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*SetResponse, error)
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
	Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error)
}

type kVServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewKVServiceClient(cc grpc.ClientConnInterface) KVServiceClient {
	return &kVServiceClient{cc}
}

func (c *kVServiceClient) Set(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*SetResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SetResponse)
	err := c.cc.Invoke(ctx, KVService_Set_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kVServiceClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, KVService_Get_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kVServiceClient) Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteResponse)
	err := c.cc.Invoke(ctx, KVService_Delete_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KVServiceServer is the server API for KVService service.
// All implementations must embed UnimplementedKVServiceServer
// for forward compatibility.
type KVServiceServer interface {
	Set(context.Context, *SetRequest) (*SetResponse, error)
	Get(context.Context, *GetRequest) (*GetResponse, error)
	Delete(context.Context, *DeleteRequest) (*DeleteResponse, error)
	mustEmbedUnimplementedKVServiceServer()
}

// UnimplementedKVServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedKVServiceServer struct{}

func (UnimplementedKVServiceServer) Set(context.Context, *SetRequest) (*SetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Set not implemented")
}
func (UnimplementedKVServiceServer) Get(context.Context, *GetRequest) (*GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedKVServiceServer) Delete(context.Context, *DeleteRequest) (*DeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedKVServiceServer) mustEmbedUnimplementedKVServiceServer() {}
func (UnimplementedKVServiceServer) testEmbeddedByValue()                   {}

// UnsafeKVServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to KVServiceServer will
// result in compilation errors.
type UnsafeKVServiceServer interface {
	mustEmbedUnimplementedKVServiceServer()
}

func RegisterKVServiceServer(s grpc.ServiceRegistrar, srv KVServiceServer) {
	// If the following call pancis, it indicates UnimplementedKVServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&KVService_ServiceDesc, srv)
}

func _KVService_Set_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KVServiceServer).Set(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KVService_Set_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KVServiceServer).Set(ctx, req.(*SetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KVService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KVServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KVService_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KVServiceServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KVService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KVServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KVService_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KVServiceServer).Delete(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// KVService_ServiceDesc is the grpc.ServiceDesc for KVService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var KVService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "raftdkv.v1.KVService",
	HandlerType: (*KVServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Set",
			Handler:    _KVService_Set_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _KVService_Get_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _KVService_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "raftdkv/v1/kv.proto",
}
