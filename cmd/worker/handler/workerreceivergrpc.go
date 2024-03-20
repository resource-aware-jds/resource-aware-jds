package handler

import (
	"context"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/handlerservice"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/protobuf/types/known/emptypb"
)

type WorkerNodeReceiverGRPCHandler struct {
	proto.UnimplementedWorkerNodeContainerReceiverServer
	workerService handlerservice.IWorker
}

func ProvideWorkerGRPCSocketHandler(grpcSocketServer grpc.WorkerNodeReceiverGRPCServer, workerService handlerservice.IWorker, meter metric.Meter) WorkerNodeReceiverGRPCHandler {
	handler := WorkerNodeReceiverGRPCHandler{
		workerService: workerService,
	}
	proto.RegisterWorkerNodeContainerReceiverServer(grpcSocketServer.GetGRPCServer(), &handler)
	return handler
}

func (w *WorkerNodeReceiverGRPCHandler) GetTaskFromQueue(ctx context.Context, payload *proto.GetTaskPayload) (*proto.Task, error) {
	task, err := w.workerService.GetTask(payload.ImageUrl, payload.ContainerId)
	return task, err
}

func (w *WorkerNodeReceiverGRPCHandler) SubmitSuccessTask(ctx context.Context, payload *proto.SubmitSuccessTaskRequest) (*emptypb.Empty, error) {
	err := w.workerService.SubmitSuccessTask(ctx, payload.ID, payload.Results)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

func (w *WorkerNodeReceiverGRPCHandler) ReportTaskFailure(ctx context.Context, payload *proto.ReportTaskFailureRequest) (*emptypb.Empty, error) {
	err := w.workerService.ReportFailTask(ctx, payload.GetID(), payload.GetErrorDetail())
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}
