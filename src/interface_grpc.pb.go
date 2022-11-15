// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.9
// source: src/interface.proto

package p2p

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

// ReqAccessToCSClient is the client API for ReqAccessToCS service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ReqAccessToCSClient interface {
	ReqAccessToCS(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Reply, error)
}

type reqAccessToCSClient struct {
	cc grpc.ClientConnInterface
}

func NewReqAccessToCSClient(cc grpc.ClientConnInterface) ReqAccessToCSClient {
	return &reqAccessToCSClient{cc}
}

func (c *reqAccessToCSClient) ReqAccessToCS(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Reply, error) {
	out := new(Reply)
	err := c.cc.Invoke(ctx, "/p2p.ReqAccessToCS/reqAccessToCS", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ReqAccessToCSServer is the server API for ReqAccessToCS service.
// All implementations must embed UnimplementedReqAccessToCSServer
// for forward compatibility
type ReqAccessToCSServer interface {
	ReqAccessToCS(context.Context, *Request) (*Reply, error)
	mustEmbedUnimplementedReqAccessToCSServer()
}

// UnimplementedReqAccessToCSServer must be embedded to have forward compatible implementations.
type UnimplementedReqAccessToCSServer struct {
}

func (UnimplementedReqAccessToCSServer) ReqAccessToCS(context.Context, *Request) (*Reply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReqAccessToCS not implemented")
}
func (UnimplementedReqAccessToCSServer) mustEmbedUnimplementedReqAccessToCSServer() {}

// UnsafeReqAccessToCSServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ReqAccessToCSServer will
// result in compilation errors.
type UnsafeReqAccessToCSServer interface {
	mustEmbedUnimplementedReqAccessToCSServer()
}

func RegisterReqAccessToCSServer(s grpc.ServiceRegistrar, srv ReqAccessToCSServer) {
	s.RegisterService(&ReqAccessToCS_ServiceDesc, srv)
}

func _ReqAccessToCS_ReqAccessToCS_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReqAccessToCSServer).ReqAccessToCS(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/p2p.ReqAccessToCS/reqAccessToCS",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReqAccessToCSServer).ReqAccessToCS(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

// ReqAccessToCS_ServiceDesc is the grpc.ServiceDesc for ReqAccessToCS service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ReqAccessToCS_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "p2p.ReqAccessToCS",
	HandlerType: (*ReqAccessToCSServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "reqAccessToCS",
			Handler:    _ReqAccessToCS_ReqAccessToCS_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "src/interface.proto",
}
