// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: user_mgmt/user_mgmt.proto

package user_mgmt

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

// UserMgmtClient is the client API for UserMgmt service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserMgmtClient interface {
	AddUser(ctx context.Context, in *AddUserRequest, opts ...grpc.CallOption) (*UserResponse, error)
	GetAllUsers(ctx context.Context, in *GetAllUsersRequest, opts ...grpc.CallOption) (*GetAllUsersResponse, error)
}

type userMgmtClient struct {
	cc grpc.ClientConnInterface
}

func NewUserMgmtClient(cc grpc.ClientConnInterface) UserMgmtClient {
	return &userMgmtClient{cc}
}

func (c *userMgmtClient) AddUser(ctx context.Context, in *AddUserRequest, opts ...grpc.CallOption) (*UserResponse, error) {
	out := new(UserResponse)
	err := c.cc.Invoke(ctx, "/user_mgmt.UserMgmt/AddUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userMgmtClient) GetAllUsers(ctx context.Context, in *GetAllUsersRequest, opts ...grpc.CallOption) (*GetAllUsersResponse, error) {
	out := new(GetAllUsersResponse)
	err := c.cc.Invoke(ctx, "/user_mgmt.UserMgmt/GetAllUsers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserMgmtServer is the server API for UserMgmt service.
// All implementations must embed UnimplementedUserMgmtServer
// for forward compatibility
type UserMgmtServer interface {
	AddUser(context.Context, *AddUserRequest) (*UserResponse, error)
	GetAllUsers(context.Context, *GetAllUsersRequest) (*GetAllUsersResponse, error)
	mustEmbedUnimplementedUserMgmtServer()
}

// UnimplementedUserMgmtServer must be embedded to have forward compatible implementations.
type UnimplementedUserMgmtServer struct {
}

func (UnimplementedUserMgmtServer) AddUser(context.Context, *AddUserRequest) (*UserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddUser not implemented")
}
func (UnimplementedUserMgmtServer) GetAllUsers(context.Context, *GetAllUsersRequest) (*GetAllUsersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllUsers not implemented")
}
func (UnimplementedUserMgmtServer) mustEmbedUnimplementedUserMgmtServer() {}

// UnsafeUserMgmtServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserMgmtServer will
// result in compilation errors.
type UnsafeUserMgmtServer interface {
	mustEmbedUnimplementedUserMgmtServer()
}

func RegisterUserMgmtServer(s grpc.ServiceRegistrar, srv UserMgmtServer) {
	s.RegisterService(&UserMgmt_ServiceDesc, srv)
}

func _UserMgmt_AddUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserMgmtServer).AddUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user_mgmt.UserMgmt/AddUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserMgmtServer).AddUser(ctx, req.(*AddUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserMgmt_GetAllUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllUsersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserMgmtServer).GetAllUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user_mgmt.UserMgmt/GetAllUsers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserMgmtServer).GetAllUsers(ctx, req.(*GetAllUsersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserMgmt_ServiceDesc is the grpc.ServiceDesc for UserMgmt service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserMgmt_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "user_mgmt.UserMgmt",
	HandlerType: (*UserMgmtServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddUser",
			Handler:    _UserMgmt_AddUser_Handler,
		},
		{
			MethodName: "GetAllUsers",
			Handler:    _UserMgmt_GetAllUsers_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "user_mgmt/user_mgmt.proto",
}
