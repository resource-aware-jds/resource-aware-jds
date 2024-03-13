// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.1
// source: proto/controlplane.proto

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

const (
	ControlPlane_WorkerRegistration_FullMethodName = "/controlplane.ControlPlane/WorkerRegistration"
	ControlPlane_CreateJob_FullMethodName          = "/controlplane.ControlPlane/CreateJob"
	ControlPlane_WorkerCheckIn_FullMethodName      = "/controlplane.ControlPlane/WorkerCheckIn"
	ControlPlane_ReportFailureTask_FullMethodName  = "/controlplane.ControlPlane/ReportFailureTask"
	ControlPlane_ReportSuccessTask_FullMethodName  = "/controlplane.ControlPlane/ReportSuccessTask"
)

// ControlPlaneClient is the client API for ControlPlane service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ControlPlaneClient interface {
	WorkerRegistration(ctx context.Context, in *ComputeNodeRegistrationRequest, opts ...grpc.CallOption) (*ComputeNodeRegistrationResponse, error)
	CreateJob(ctx context.Context, in *CreateJobRequest, opts ...grpc.CallOption) (*CreateJobResponse, error)
	WorkerCheckIn(ctx context.Context, in *WorkerCheckInRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	ReportFailureTask(ctx context.Context, in *ReportFailureTaskRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	ReportSuccessTask(ctx context.Context, in *ReportSuccessTaskRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type controlPlaneClient struct {
	cc grpc.ClientConnInterface
}

func NewControlPlaneClient(cc grpc.ClientConnInterface) ControlPlaneClient {
	return &controlPlaneClient{cc}
}

func (c *controlPlaneClient) WorkerRegistration(ctx context.Context, in *ComputeNodeRegistrationRequest, opts ...grpc.CallOption) (*ComputeNodeRegistrationResponse, error) {
	out := new(ComputeNodeRegistrationResponse)
	err := c.cc.Invoke(ctx, ControlPlane_WorkerRegistration_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *controlPlaneClient) CreateJob(ctx context.Context, in *CreateJobRequest, opts ...grpc.CallOption) (*CreateJobResponse, error) {
	out := new(CreateJobResponse)
	err := c.cc.Invoke(ctx, ControlPlane_CreateJob_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *controlPlaneClient) WorkerCheckIn(ctx context.Context, in *WorkerCheckInRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, ControlPlane_WorkerCheckIn_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *controlPlaneClient) ReportFailureTask(ctx context.Context, in *ReportFailureTaskRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, ControlPlane_ReportFailureTask_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *controlPlaneClient) ReportSuccessTask(ctx context.Context, in *ReportSuccessTaskRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, ControlPlane_ReportSuccessTask_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ControlPlaneServer is the server API for ControlPlane service.
// All implementations must embed UnimplementedControlPlaneServer
// for forward compatibility
type ControlPlaneServer interface {
	WorkerRegistration(context.Context, *ComputeNodeRegistrationRequest) (*ComputeNodeRegistrationResponse, error)
	CreateJob(context.Context, *CreateJobRequest) (*CreateJobResponse, error)
	WorkerCheckIn(context.Context, *WorkerCheckInRequest) (*emptypb.Empty, error)
	ReportFailureTask(context.Context, *ReportFailureTaskRequest) (*emptypb.Empty, error)
	ReportSuccessTask(context.Context, *ReportSuccessTaskRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedControlPlaneServer()
}

// UnimplementedControlPlaneServer must be embedded to have forward compatible implementations.
type UnimplementedControlPlaneServer struct {
}

func (UnimplementedControlPlaneServer) WorkerRegistration(context.Context, *ComputeNodeRegistrationRequest) (*ComputeNodeRegistrationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WorkerRegistration not implemented")
}
func (UnimplementedControlPlaneServer) CreateJob(context.Context, *CreateJobRequest) (*CreateJobResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateJob not implemented")
}
func (UnimplementedControlPlaneServer) WorkerCheckIn(context.Context, *WorkerCheckInRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WorkerCheckIn not implemented")
}
func (UnimplementedControlPlaneServer) ReportFailureTask(context.Context, *ReportFailureTaskRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReportFailureTask not implemented")
}
func (UnimplementedControlPlaneServer) ReportSuccessTask(context.Context, *ReportSuccessTaskRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReportSuccessTask not implemented")
}
func (UnimplementedControlPlaneServer) mustEmbedUnimplementedControlPlaneServer() {}

// UnsafeControlPlaneServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ControlPlaneServer will
// result in compilation errors.
type UnsafeControlPlaneServer interface {
	mustEmbedUnimplementedControlPlaneServer()
}

func RegisterControlPlaneServer(s grpc.ServiceRegistrar, srv ControlPlaneServer) {
	s.RegisterService(&ControlPlane_ServiceDesc, srv)
}

func _ControlPlane_WorkerRegistration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ComputeNodeRegistrationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ControlPlaneServer).WorkerRegistration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ControlPlane_WorkerRegistration_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ControlPlaneServer).WorkerRegistration(ctx, req.(*ComputeNodeRegistrationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ControlPlane_CreateJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateJobRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ControlPlaneServer).CreateJob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ControlPlane_CreateJob_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ControlPlaneServer).CreateJob(ctx, req.(*CreateJobRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ControlPlane_WorkerCheckIn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WorkerCheckInRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ControlPlaneServer).WorkerCheckIn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ControlPlane_WorkerCheckIn_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ControlPlaneServer).WorkerCheckIn(ctx, req.(*WorkerCheckInRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ControlPlane_ReportFailureTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReportFailureTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ControlPlaneServer).ReportFailureTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ControlPlane_ReportFailureTask_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ControlPlaneServer).ReportFailureTask(ctx, req.(*ReportFailureTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ControlPlane_ReportSuccessTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReportSuccessTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ControlPlaneServer).ReportSuccessTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ControlPlane_ReportSuccessTask_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ControlPlaneServer).ReportSuccessTask(ctx, req.(*ReportSuccessTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ControlPlane_ServiceDesc is the grpc.ServiceDesc for ControlPlane service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ControlPlane_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "controlplane.ControlPlane",
	HandlerType: (*ControlPlaneServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "WorkerRegistration",
			Handler:    _ControlPlane_WorkerRegistration_Handler,
		},
		{
			MethodName: "CreateJob",
			Handler:    _ControlPlane_CreateJob_Handler,
		},
		{
			MethodName: "WorkerCheckIn",
			Handler:    _ControlPlane_WorkerCheckIn_Handler,
		},
		{
			MethodName: "ReportFailureTask",
			Handler:    _ControlPlane_ReportFailureTask_Handler,
		},
		{
			MethodName: "ReportSuccessTask",
			Handler:    _ControlPlane_ReportSuccessTask_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/controlplane.proto",
}
