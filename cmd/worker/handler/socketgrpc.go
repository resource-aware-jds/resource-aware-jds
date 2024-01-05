package handler

import (
	"context"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
	"github.com/resource-aware-jds/resource-aware-jds/service"
	"google.golang.org/protobuf/types/known/emptypb"
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

func (w *WorkerGRPCSocketHandler) SubmitSuccessTask(ctx context.Context, payload *proto.SubmitSuccessTaskRequest) (*emptypb.Empty, error) {
	err := w.workerService.SubmitSuccessTask(payload.ID, payload.Results)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

func (w *WorkerGRPCSocketHandler) ReportTaskFailure(ctx context.Context, payload *proto.ReportTaskFailureRequest) (*emptypb.Empty, error) {
	err := w.workerService.ReportFailTask(payload.ID, payload.ErrorDetail)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}
