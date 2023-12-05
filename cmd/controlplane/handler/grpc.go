package handler

import (
	"context"
	"fmt"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
	"github.com/resource-aware-jds/resource-aware-jds/service"
	"google.golang.org/grpc/metadata"
)

type GRPCHandler struct {
	proto.UnimplementedControlPlaneServer
	controlPlaneService service.IControlPlane
}

func ProvideControlPlaneGRPCHandler(grpcServer grpc.RAJDSGrpc, controlPlaneService service.IControlPlane) GRPCHandler {
	handler := GRPCHandler{
		controlPlaneService: controlPlaneService,
	}
	proto.RegisterControlPlaneServer(grpcServer.GetGRPCServer(), &handler)
	return handler
}

func (g *GRPCHandler) WorkerRegistration(ctx context.Context, req *proto.ComputeNodeRegistrationRequest) (*proto.ComputeNodeRegistrationResponse, error) {
	result, _ := metadata.FromIncomingContext(ctx)
	fmt.Println(result)

	return &proto.ComputeNodeRegistrationResponse{
		Id: "Test",
	}, nil
}
