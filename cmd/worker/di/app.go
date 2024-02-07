package di

import (
	"github.com/resource-aware-jds/resource-aware-jds/cmd/worker/handler"
	"github.com/resource-aware-jds/resource-aware-jds/cmd/worker/http"
	"github.com/resource-aware-jds/resource-aware-jds/daemon"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
	httpServer "github.com/resource-aware-jds/resource-aware-jds/pkg/http"
)

type WorkerApp struct {
	GRPCServer                   grpc.RAJDSGrpcServer
	WorkerGRPCHandler            handler.GRPCHandler
	WorkerNodeReceiverGRPCServer grpc.WorkerNodeReceiverGRPCServer
	WorkerGRPCSocketHandler      handler.WorkerNodeReceiverGRPCHandler
	WorkerNodeDaemon             daemon.WorkerNode
	ControlPlaneGRPCClient       proto.ControlPlaneClient
	WorkerHTTPServer             httpServer.Server
	RouterResult                 http.RouterResult
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
	workerHTTPServer httpServer.Server,
	routerResult http.RouterResult,
) WorkerApp {
	return WorkerApp{
		GRPCServer:                   grpcServer,
		WorkerGRPCHandler:            workerGRPCHandler,
		WorkerNodeReceiverGRPCServer: workerNodeReceiverGRPCServer,
		WorkerGRPCSocketHandler:      workerGRPCSocketHandler,
		WorkerNodeDaemon:             workerNodeDaemon,
		WorkerHTTPServer:             workerHTTPServer,
		RouterResult:                 routerResult,
	}
}
