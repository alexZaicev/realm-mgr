// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: realm_mgr/v1/service.proto

package realm_mgr_v1

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

// RealmManagerServiceClient is the client API for RealmManagerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RealmManagerServiceClient interface {
	// Realm RPCs
	//
	// Get a single active realm
	GetRealm(ctx context.Context, in *GetRealmRequest, opts ...grpc.CallOption) (*GetRealmResponse, error)
	// Create a new realm
	CreateRealm(ctx context.Context, in *CreateRealmRequest, opts ...grpc.CallOption) (*CreateRealmResponse, error)
	// Release existing draft copy of the realm
	ReleaseRealm(ctx context.Context, in *ReleaseRealmRequest, opts ...grpc.CallOption) (*ReleaseRealmResponse, error)
	// Update single realm
	UpdateRealm(ctx context.Context, in *UpdateRealmRequest, opts ...grpc.CallOption) (*UpdateRealmResponse, error)
}

type realmManagerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRealmManagerServiceClient(cc grpc.ClientConnInterface) RealmManagerServiceClient {
	return &realmManagerServiceClient{cc}
}

func (c *realmManagerServiceClient) GetRealm(ctx context.Context, in *GetRealmRequest, opts ...grpc.CallOption) (*GetRealmResponse, error) {
	out := new(GetRealmResponse)
	err := c.cc.Invoke(ctx, "/realm_mgr.v1.RealmManagerService/GetRealm", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *realmManagerServiceClient) CreateRealm(ctx context.Context, in *CreateRealmRequest, opts ...grpc.CallOption) (*CreateRealmResponse, error) {
	out := new(CreateRealmResponse)
	err := c.cc.Invoke(ctx, "/realm_mgr.v1.RealmManagerService/CreateRealm", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *realmManagerServiceClient) ReleaseRealm(ctx context.Context, in *ReleaseRealmRequest, opts ...grpc.CallOption) (*ReleaseRealmResponse, error) {
	out := new(ReleaseRealmResponse)
	err := c.cc.Invoke(ctx, "/realm_mgr.v1.RealmManagerService/ReleaseRealm", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *realmManagerServiceClient) UpdateRealm(ctx context.Context, in *UpdateRealmRequest, opts ...grpc.CallOption) (*UpdateRealmResponse, error) {
	out := new(UpdateRealmResponse)
	err := c.cc.Invoke(ctx, "/realm_mgr.v1.RealmManagerService/UpdateRealm", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RealmManagerServiceServer is the server API for RealmManagerService service.
// All implementations must embed UnimplementedRealmManagerServiceServer
// for forward compatibility
type RealmManagerServiceServer interface {
	// Realm RPCs
	//
	// Get a single active realm
	GetRealm(context.Context, *GetRealmRequest) (*GetRealmResponse, error)
	// Create a new realm
	CreateRealm(context.Context, *CreateRealmRequest) (*CreateRealmResponse, error)
	// Release existing draft copy of the realm
	ReleaseRealm(context.Context, *ReleaseRealmRequest) (*ReleaseRealmResponse, error)
	// Update single realm
	UpdateRealm(context.Context, *UpdateRealmRequest) (*UpdateRealmResponse, error)
	mustEmbedUnimplementedRealmManagerServiceServer()
}

// UnimplementedRealmManagerServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRealmManagerServiceServer struct {
}

func (UnimplementedRealmManagerServiceServer) GetRealm(context.Context, *GetRealmRequest) (*GetRealmResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRealm not implemented")
}
func (UnimplementedRealmManagerServiceServer) CreateRealm(context.Context, *CreateRealmRequest) (*CreateRealmResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateRealm not implemented")
}
func (UnimplementedRealmManagerServiceServer) ReleaseRealm(context.Context, *ReleaseRealmRequest) (*ReleaseRealmResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReleaseRealm not implemented")
}
func (UnimplementedRealmManagerServiceServer) UpdateRealm(context.Context, *UpdateRealmRequest) (*UpdateRealmResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateRealm not implemented")
}
func (UnimplementedRealmManagerServiceServer) mustEmbedUnimplementedRealmManagerServiceServer() {}

// UnsafeRealmManagerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RealmManagerServiceServer will
// result in compilation errors.
type UnsafeRealmManagerServiceServer interface {
	mustEmbedUnimplementedRealmManagerServiceServer()
}

func RegisterRealmManagerServiceServer(s grpc.ServiceRegistrar, srv RealmManagerServiceServer) {
	s.RegisterService(&RealmManagerService_ServiceDesc, srv)
}

func _RealmManagerService_GetRealm_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRealmRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RealmManagerServiceServer).GetRealm(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/realm_mgr.v1.RealmManagerService/GetRealm",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RealmManagerServiceServer).GetRealm(ctx, req.(*GetRealmRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RealmManagerService_CreateRealm_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRealmRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RealmManagerServiceServer).CreateRealm(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/realm_mgr.v1.RealmManagerService/CreateRealm",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RealmManagerServiceServer).CreateRealm(ctx, req.(*CreateRealmRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RealmManagerService_ReleaseRealm_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReleaseRealmRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RealmManagerServiceServer).ReleaseRealm(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/realm_mgr.v1.RealmManagerService/ReleaseRealm",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RealmManagerServiceServer).ReleaseRealm(ctx, req.(*ReleaseRealmRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RealmManagerService_UpdateRealm_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRealmRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RealmManagerServiceServer).UpdateRealm(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/realm_mgr.v1.RealmManagerService/UpdateRealm",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RealmManagerServiceServer).UpdateRealm(ctx, req.(*UpdateRealmRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RealmManagerService_ServiceDesc is the grpc.ServiceDesc for RealmManagerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RealmManagerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "realm_mgr.v1.RealmManagerService",
	HandlerType: (*RealmManagerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetRealm",
			Handler:    _RealmManagerService_GetRealm_Handler,
		},
		{
			MethodName: "CreateRealm",
			Handler:    _RealmManagerService_CreateRealm_Handler,
		},
		{
			MethodName: "ReleaseRealm",
			Handler:    _RealmManagerService_ReleaseRealm_Handler,
		},
		{
			MethodName: "UpdateRealm",
			Handler:    _RealmManagerService_UpdateRealm_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "realm_mgr/v1/service.proto",
}
