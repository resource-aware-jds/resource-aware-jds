package di

import (
	"github.com/resource-aware-jds/resource-aware-jds/cmd/worker/handler"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
)

type WorkerApp struct {
	GRPCServer        grpc.RAJDSGrpcServer
	WorkerGRPCHandler handler.GRPCHandler
}

func ProvideWorkerApp(
	grpcServer grpc.RAJDSGrpcServer,
	workerGRPCHandler handler.GRPCHandler,
) WorkerApp {
	return WorkerApp{
		GRPCServer:        grpcServer,
		WorkerGRPCHandler: workerGRPCHandler,
	}
}
