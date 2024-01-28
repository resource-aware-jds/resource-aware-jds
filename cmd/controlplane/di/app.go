package di

import (
	grpcHandler "github.com/resource-aware-jds/resource-aware-jds/cmd/controlplane/grpc"
	"github.com/resource-aware-jds/resource-aware-jds/daemon"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
)

type ControlPlaneApp struct {
	GRPCServer              grpc.RAJDSGrpcServer
	ControlPlaneGRPCHandler grpcHandler.GRPCHandler
	ControlPlaneDaemon      daemon.IControlPlane
}

func ProvideControlPlaneApp(
	grpcServer grpc.RAJDSGrpcServer,
	controlPlaneGRPCHandler grpcHandler.GRPCHandler,
	controlPlaneDaemon daemon.IControlPlane,
) ControlPlaneApp {
	return ControlPlaneApp{
		GRPCServer:              grpcServer,
		ControlPlaneGRPCHandler: controlPlaneGRPCHandler,
		ControlPlaneDaemon:      controlPlaneDaemon,
	}
}
