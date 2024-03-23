// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.1
// source: proto/workernode.proto

package proto

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

// WorkerNodeClient is the client API for WorkerNode service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type WorkerNodeClient interface {
	SendTask(ctx context.Context, in *RecievedTask, opts ...grpc.CallOption) (*emptypb.Empty, error)
	HealthCheck(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Resource, error)
	GetAllTasks(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*TaskResponse, error)
}

type workerNodeClient struct {
	cc grpc.ClientConnInterface
}

func NewWorkerNodeClient(cc grpc.ClientConnInterface) WorkerNodeClient {
	return &workerNodeClient{cc}
}

func (c *workerNodeClient) SendTask(ctx context.Context, in *RecievedTask, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/WorkerNode.WorkerNode/SendTask", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *workerNodeClient) HealthCheck(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Resource, error) {
	out := new(Resource)
	err := c.cc.Invoke(ctx, "/WorkerNode.WorkerNode/HealthCheck", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *workerNodeClient) GetAllTasks(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*TaskResponse, error) {
	out := new(TaskResponse)
	err := c.cc.Invoke(ctx, "/WorkerNode.WorkerNode/GetAllTasks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WorkerNodeServer is the server API for WorkerNode service.
// All implementations must embed UnimplementedWorkerNodeServer
// for forward compatibility
type WorkerNodeServer interface {
	SendTask(context.Context, *RecievedTask) (*emptypb.Empty, error)
	HealthCheck(context.Context, *emptypb.Empty) (*Resource, error)
	GetAllTasks(context.Context, *emptypb.Empty) (*TaskResponse, error)
	mustEmbedUnimplementedWorkerNodeServer()
}

// UnimplementedWorkerNodeServer must be embedded to have forward compatible implementations.
type UnimplementedWorkerNodeServer struct {
}

func (UnimplementedWorkerNodeServer) SendTask(context.Context, *RecievedTask) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendTask not implemented")
}
func (UnimplementedWorkerNodeServer) HealthCheck(context.Context, *emptypb.Empty) (*Resource, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HealthCheck not implemented")
}
func (UnimplementedWorkerNodeServer) GetAllTasks(context.Context, *emptypb.Empty) (*TaskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllTasks not implemented")
}
func (UnimplementedWorkerNodeServer) mustEmbedUnimplementedWorkerNodeServer() {}

// UnsafeWorkerNodeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to WorkerNodeServer will
// result in compilation errors.
type UnsafeWorkerNodeServer interface {
	mustEmbedUnimplementedWorkerNodeServer()
}

func RegisterWorkerNodeServer(s grpc.ServiceRegistrar, srv WorkerNodeServer) {
	s.RegisterService(&WorkerNode_ServiceDesc, srv)
}

func _WorkerNode_SendTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecievedTask)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkerNodeServer).SendTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/WorkerNode.WorkerNode/SendTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkerNodeServer).SendTask(ctx, req.(*RecievedTask))
	}
	return interceptor(ctx, in, info, handler)
}

func _WorkerNode_HealthCheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkerNodeServer).HealthCheck(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/WorkerNode.WorkerNode/HealthCheck",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkerNodeServer).HealthCheck(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _WorkerNode_GetAllTasks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkerNodeServer).GetAllTasks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/WorkerNode.WorkerNode/GetAllTasks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkerNodeServer).GetAllTasks(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// WorkerNode_ServiceDesc is the grpc.ServiceDesc for WorkerNode service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var WorkerNode_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "WorkerNode.WorkerNode",
	HandlerType: (*WorkerNodeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendTask",
			Handler:    _WorkerNode_SendTask_Handler,
		},
		{
			MethodName: "HealthCheck",
			Handler:    _WorkerNode_HealthCheck_Handler,
		},
		{
			MethodName: "GetAllTasks",
			Handler:    _WorkerNode_GetAllTasks_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/workernode.proto",
}
