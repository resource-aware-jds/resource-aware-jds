package di

import (
	grpcHandler "github.com/resource-aware-jds/resource-aware-jds/cmd/controlplane/grpc"
	httpHandler "github.com/resource-aware-jds/resource-aware-jds/cmd/controlplane/http"
	"github.com/resource-aware-jds/resource-aware-jds/daemon"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/http"
)

type ControlPlaneApp struct {
	GRPCServer              grpc.RAJDSGrpcServer
	HTTPServer              http.Server
	ControlPlaneGRPCHandler grpcHandler.GRPCHandler
	ControlPlaneDaemon      daemon.IControlPlane
	httpRouterResult        httpHandler.RouterResult
}

func ProvideControlPlaneApp(
	grpcServer grpc.RAJDSGrpcServer,
	httpServer http.Server,
	controlPlaneGRPCHandler grpcHandler.GRPCHandler,
	controlPlaneDaemon daemon.IControlPlane,
	httpRouterResult httpHandler.RouterResult,
) ControlPlaneApp {
	return ControlPlaneApp{
		GRPCServer:              grpcServer,
		HTTPServer:              httpServer,
		ControlPlaneGRPCHandler: controlPlaneGRPCHandler,
		ControlPlaneDaemon:      controlPlaneDaemon,
		httpRouterResult:        httpRouterResult,
	}
}
