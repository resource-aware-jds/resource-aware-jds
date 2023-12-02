package handler

import (
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
	"github.com/resource-aware-jds/resource-aware-jds/service"
)

type GRPCHandler struct {
	proto.UnimplementedControlPlaneServer
	controlPlaneService service.IControlPlane
}

func ProvideControlPlaneGRPCHandler(grpcServer grpc.RAJDSGrpc, controlPlaneService service.IControlPlane) GRPCHandler {
	handler := GRPCHandler{
		controlPlaneService: controlPlaneService,
	}
	proto.RegisterControlPlaneServer(grpcServer.GetGRPCServer(), handler)
	return handler
}
