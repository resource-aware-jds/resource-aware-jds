package di

import (
	"github.com/resource-aware-jds/resource-aware-jds/cmd/worker/handler"
	"github.com/resource-aware-jds/resource-aware-jds/daemon"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
)

type WorkerApp struct {
	GRPCServer                   grpc.RAJDSGrpcServer
	WorkerGRPCHandler            handler.GRPCHandler
	WorkerNodeReceiverGRPCServer grpc.WorkerNodeReceiverGRPCServer
	WorkerGRPCSocketHandler      handler.WorkerNodeReceiverGRPCHandler
	WorkerNodeDaemon             daemon.WorkerNode
	ControlPlaneGRPCClient       proto.ControlPlaneClient
}

func ProvideControlPlaneGRPCClient(grpcClient grpc.RAJDSGrpcClient) proto.ControlPlaneClient {
	return proto.NewControlPlaneClient(grpcClient.GetConnection())
}

func ProvideWorkerApp(
	grpcServer grpc.RAJDSGrpcServer,
	workerGRPCHandler handler.GRPCHandler,
	workerNodeReceiverGRPCServer grpc.WorkerNodeReceiverGRPCServer,
	workerGRPCSocketHandler handler.WorkerNodeReceiverGRPCHandler,
	workerNodeDaemon daemon.WorkerNode,
) WorkerApp {
	return WorkerApp{
		GRPCServer:                   grpcServer,
		WorkerGRPCHandler:            workerGRPCHandler,
		WorkerNodeReceiverGRPCServer: workerNodeReceiverGRPCServer,
		WorkerGRPCSocketHandler:      workerGRPCSocketHandler,
		WorkerNodeDaemon:             workerNodeDaemon,
	}
}
