package di

import (
	"github.com/resource-aware-jds/resource-aware-jds/cmd/worker/handler"
	"github.com/resource-aware-jds/resource-aware-jds/daemon"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
)

type WorkerApp struct {
	GRPCServer              grpc.RAJDSGrpcServer
	WorkerGRPCHandler       handler.GRPCHandler
	GRPCSocketServer        grpc.SocketServer
	WorkerGRPCSocketHandler handler.WorkerGRPCSocketHandler
	WorkerNodeDaemon        daemon.WorkerNode
	ControlPlaneGRPCClient  proto.ControlPlaneClient
}

func ProvideControlPlaneGRPCClient(grpcClient grpc.RAJDSGrpcClient) proto.ControlPlaneClient {
	return proto.NewControlPlaneClient(grpcClient.GetConnection())
}

func ProvideWorkerApp(
	grpcServer grpc.RAJDSGrpcServer,
	workerGRPCHandler handler.GRPCHandler,
	grpcSocketServer grpc.SocketServer,
	workerGRPCSocketHandler handler.WorkerGRPCSocketHandler,
	workerNodeDaemon daemon.WorkerNode,
) WorkerApp {
	return WorkerApp{
		GRPCServer:              grpcServer,
		WorkerGRPCHandler:       workerGRPCHandler,
		GRPCSocketServer:        grpcSocketServer,
		WorkerGRPCSocketHandler: workerGRPCSocketHandler,
		WorkerNodeDaemon:        workerNodeDaemon,
	}
}
