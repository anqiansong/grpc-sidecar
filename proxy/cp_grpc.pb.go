// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: cp.proto

package proxy

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

// CPServiceClient is the client API for CPService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CPServiceClient interface {
	SyncConfig(ctx context.Context, in *CPRequest, opts ...grpc.CallOption) (*CPResponse, error)
}

type cPServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCPServiceClient(cc grpc.ClientConnInterface) CPServiceClient {
	return &cPServiceClient{cc}
}

func (c *cPServiceClient) SyncConfig(ctx context.Context, in *CPRequest, opts ...grpc.CallOption) (*CPResponse, error) {
	out := new(CPResponse)
	err := c.cc.Invoke(ctx, "/cp.CPService/SyncConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CPServiceServer is the server API for CPService service.
// All implementations must embed UnimplementedCPServiceServer
// for forward compatibility
type CPServiceServer interface {
	SyncConfig(context.Context, *CPRequest) (*CPResponse, error)
	mustEmbedUnimplementedCPServiceServer()
}

// UnimplementedCPServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCPServiceServer struct {
}

func (UnimplementedCPServiceServer) SyncConfig(context.Context, *CPRequest) (*CPResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SyncConfig not implemented")
}
func (UnimplementedCPServiceServer) mustEmbedUnimplementedCPServiceServer() {}

// UnsafeCPServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CPServiceServer will
// result in compilation errors.
type UnsafeCPServiceServer interface {
	mustEmbedUnimplementedCPServiceServer()
}

func RegisterCPServiceServer(s grpc.ServiceRegistrar, srv CPServiceServer) {
	s.RegisterService(&CPService_ServiceDesc, srv)
}

func _CPService_SyncConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CPRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CPServiceServer).SyncConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cp.CPService/SyncConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CPServiceServer).SyncConfig(ctx, req.(*CPRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CPService_ServiceDesc is the grpc.ServiceDesc for CPService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CPService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "cp.CPService",
	HandlerType: (*CPServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SyncConfig",
			Handler:    _CPService_SyncConfig_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cp.proto",
}