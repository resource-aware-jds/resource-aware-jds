package di

import (
	"github.com/resource-aware-jds/resource-aware-jds/cmd/controlplane/handler"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
)

type ControlPlaneApp struct {
	GRPCServer              grpc.RAJDSGrpcServer
	ControlPlaneGRPCHandler handler.GRPCHandler
}

func ProvideControlPlaneApp(
	grpcServer grpc.RAJDSGrpcServer,
	controlPlaneGRPCHandler handler.GRPCHandler,
) ControlPlaneApp {
	return ControlPlaneApp{
		GRPCServer:              grpcServer,
		ControlPlaneGRPCHandler: controlPlaneGRPCHandler,
	}
}
