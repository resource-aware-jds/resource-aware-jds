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

// ComputeNodeClient is the client API for ComputeNode service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ComputeNodeClient interface {
	SendTask(ctx context.Context, in *Task, opts ...grpc.CallOption) (*emptypb.Empty, error)
	HealthCheck(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
	ReportJob(ctx context.Context, in *ReportTaskRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type computeNodeClient struct {
	cc grpc.ClientConnInterface
}

func NewComputeNodeClient(cc grpc.ClientConnInterface) ComputeNodeClient {
	return &computeNodeClient{cc}
}

func (c *computeNodeClient) SendTask(ctx context.Context, in *Task, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/computenode.ComputeNode/SendTask", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *computeNodeClient) HealthCheck(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/computenode.ComputeNode/HealthCheck", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *computeNodeClient) ReportJob(ctx context.Context, in *ReportTaskRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/computenode.ComputeNode/ReportJob", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ComputeNodeServer is the server API for ComputeNode service.
// All implementations must embed UnimplementedComputeNodeServer
// for forward compatibility
type ComputeNodeServer interface {
	SendTask(context.Context, *Task) (*emptypb.Empty, error)
	HealthCheck(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	ReportJob(context.Context, *ReportTaskRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedComputeNodeServer()
}

// UnimplementedComputeNodeServer must be embedded to have forward compatible implementations.
type UnimplementedComputeNodeServer struct {
}

func (UnimplementedComputeNodeServer) SendTask(context.Context, *Task) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendTask not implemented")
}
func (UnimplementedComputeNodeServer) HealthCheck(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HealthCheck not implemented")
}
func (UnimplementedComputeNodeServer) ReportJob(context.Context, *ReportTaskRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReportJob not implemented")
}
func (UnimplementedComputeNodeServer) mustEmbedUnimplementedComputeNodeServer() {}

// UnsafeComputeNodeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ComputeNodeServer will
// result in compilation errors.
type UnsafeComputeNodeServer interface {
	mustEmbedUnimplementedComputeNodeServer()
}

func RegisterComputeNodeServer(s grpc.ServiceRegistrar, srv ComputeNodeServer) {
	s.RegisterService(&ComputeNode_ServiceDesc, srv)
}

func _ComputeNode_SendTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Task)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ComputeNodeServer).SendTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/computenode.ComputeNode/SendTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ComputeNodeServer).SendTask(ctx, req.(*Task))
	}
	return interceptor(ctx, in, info, handler)
}

func _ComputeNode_HealthCheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ComputeNodeServer).HealthCheck(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/computenode.ComputeNode/HealthCheck",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ComputeNodeServer).HealthCheck(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ComputeNode_ReportJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReportTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ComputeNodeServer).ReportJob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/computenode.ComputeNode/ReportJob",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ComputeNodeServer).ReportJob(ctx, req.(*ReportTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ComputeNode_ServiceDesc is the grpc.ServiceDesc for ComputeNode service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ComputeNode_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "computenode.ComputeNode",
	HandlerType: (*ComputeNodeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendTask",
			Handler:    _ComputeNode_SendTask_Handler,
		},
		{
			MethodName: "HealthCheck",
			Handler:    _ComputeNode_HealthCheck_Handler,
		},
		{
			MethodName: "ReportJob",
			Handler:    _ComputeNode_ReportJob_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/workernode.proto",
}
