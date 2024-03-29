// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.2
// source: database.proto

package protos

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Database_Ping_FullMethodName       = "/ch.unibas.ub.mediaserver.Database/Ping"
	Database_GetCache_FullMethodName   = "/ch.unibas.ub.mediaserver.Database/GetCache"
	Database_CreateItem_FullMethodName = "/ch.unibas.ub.mediaserver.Database/CreateItem"
)

// DatabaseClient is the client API for Database service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DatabaseClient interface {
	// ping the database
	Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// get data from a cache entry
	GetCache(ctx context.Context, in *CacheRequest, opts ...grpc.CallOption) (*CacheResult, error)
	// create new item
	CreateItem(ctx context.Context, in *CreateItemRequest, opts ...grpc.CallOption) (*Item, error)
}

type databaseClient struct {
	cc grpc.ClientConnInterface
}

func NewDatabaseClient(cc grpc.ClientConnInterface) DatabaseClient {
	return &databaseClient{cc}
}

func (c *databaseClient) Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Database_Ping_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) GetCache(ctx context.Context, in *CacheRequest, opts ...grpc.CallOption) (*CacheResult, error) {
	out := new(CacheResult)
	err := c.cc.Invoke(ctx, Database_GetCache_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) CreateItem(ctx context.Context, in *CreateItemRequest, opts ...grpc.CallOption) (*Item, error) {
	out := new(Item)
	err := c.cc.Invoke(ctx, Database_CreateItem_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DatabaseServer is the server API for Database service.
// All implementations must embed UnimplementedDatabaseServer
// for forward compatibility
type DatabaseServer interface {
	// ping the database
	Ping(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	// get data from a cache entry
	GetCache(context.Context, *CacheRequest) (*CacheResult, error)
	// create new item
	CreateItem(context.Context, *CreateItemRequest) (*Item, error)
	mustEmbedUnimplementedDatabaseServer()
}

// UnimplementedDatabaseServer must be embedded to have forward compatible implementations.
type UnimplementedDatabaseServer struct {
}

func (UnimplementedDatabaseServer) Ping(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedDatabaseServer) GetCache(context.Context, *CacheRequest) (*CacheResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCache not implemented")
}
func (UnimplementedDatabaseServer) CreateItem(context.Context, *CreateItemRequest) (*Item, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateItem not implemented")
}
func (UnimplementedDatabaseServer) mustEmbedUnimplementedDatabaseServer() {}

// UnsafeDatabaseServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DatabaseServer will
// result in compilation errors.
type UnsafeDatabaseServer interface {
	mustEmbedUnimplementedDatabaseServer()
}

func RegisterDatabaseServer(s grpc.ServiceRegistrar, srv DatabaseServer) {
	s.RegisterService(&Database_ServiceDesc, srv)
}

func _Database_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Database_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServer).Ping(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_GetCache_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CacheRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServer).GetCache(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Database_GetCache_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServer).GetCache(ctx, req.(*CacheRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_CreateItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServer).CreateItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Database_CreateItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServer).CreateItem(ctx, req.(*CreateItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Database_ServiceDesc is the grpc.ServiceDesc for Database service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Database_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ch.unibas.ub.mediaserver.Database",
	HandlerType: (*DatabaseServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _Database_Ping_Handler,
		},
		{
			MethodName: "GetCache",
			Handler:    _Database_GetCache_Handler,
		},
		{
			MethodName: "CreateItem",
			Handler:    _Database_CreateItem_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "database.proto",
}
