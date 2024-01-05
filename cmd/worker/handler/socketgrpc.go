package handler

import (
	"context"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
	"github.com/resource-aware-jds/resource-aware-jds/service"
)

type WorkerGRPCSocketHandler struct {
	proto.UnimplementedWorkerNodeContainerReceiverServer
	workerService service.IWorker
}

func ProvideWorkerGRPCSocketHandler(grpcSocketServer grpc.SocketServer, workerService service.IWorker) WorkerGRPCSocketHandler {
	handler := WorkerGRPCSocketHandler{
		workerService: workerService,
	}
	proto.RegisterWorkerNodeContainerReceiverServer(grpcSocketServer.GetGRPCServer(), &handler)
	return handler
}

func (w *WorkerGRPCSocketHandler) GetTaskFromQueue(ctx context.Context, payload *proto.GetTaskPayload) (*proto.Task, error) {
	task, err := w.workerService.GetTask(payload.ImageUrl)
	return task, err
}
