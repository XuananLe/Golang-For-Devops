// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: service_simple.proto

package pb

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
	SimpleService_CreateUser_FullMethodName = "/pb.SimpleService/CreateUser"
)

// SimpleServiceClient is the client API for SimpleService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SimpleServiceClient interface {
	CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error)
}

type simpleServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSimpleServiceClient(cc grpc.ClientConnInterface) SimpleServiceClient {
	return &simpleServiceClient{cc}
}

func (c *simpleServiceClient) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateUserResponse)
	err := c.cc.Invoke(ctx, SimpleService_CreateUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SimpleServiceServer is the server API for SimpleService service.
// All implementations must embed UnimplementedSimpleServiceServer
// for forward compatibility.
type SimpleServiceServer interface {
	CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error)
	mustEmbedUnimplementedSimpleServiceServer()
}

// UnimplementedSimpleServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedSimpleServiceServer struct{}

func (UnimplementedSimpleServiceServer) CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}
func (UnimplementedSimpleServiceServer) mustEmbedUnimplementedSimpleServiceServer() {}
func (UnimplementedSimpleServiceServer) testEmbeddedByValue()                       {}

// UnsafeSimpleServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SimpleServiceServer will
// result in compilation errors.
type UnsafeSimpleServiceServer interface {
	mustEmbedUnimplementedSimpleServiceServer()
}

func RegisterSimpleServiceServer(s grpc.ServiceRegistrar, srv SimpleServiceServer) {
	// If the following call pancis, it indicates UnimplementedSimpleServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&SimpleService_ServiceDesc, srv)
}

func _SimpleService_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SimpleServiceServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SimpleService_CreateUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SimpleServiceServer).CreateUser(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SimpleService_ServiceDesc is the grpc.ServiceDesc for SimpleService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SimpleService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.SimpleService",
	HandlerType: (*SimpleServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateUser",
			Handler:    _SimpleService_CreateUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service_simple.proto",
}
