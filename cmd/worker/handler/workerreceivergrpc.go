package handler

import (
	"context"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
	"github.com/resource-aware-jds/resource-aware-jds/service"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/protobuf/types/known/emptypb"
)

type WorkerNodeReceiverGRPCHandler struct {
	proto.UnimplementedWorkerNodeContainerReceiverServer
	workerService       service.IWorker
	containerSubmitTask metric.Int64Counter
}

func ProvideWorkerGRPCSocketHandler(grpcSocketServer grpc.WorkerNodeReceiverGRPCServer, workerService service.IWorker, meter metric.Meter) WorkerNodeReceiverGRPCHandler {
	containerSubmitTask, err := meter.Int64Counter("container_submit_task")
	if err != nil {
		return WorkerNodeReceiverGRPCHandler{}
	}

	handler := WorkerNodeReceiverGRPCHandler{
		workerService:       workerService,
		containerSubmitTask: containerSubmitTask,
	}
	proto.RegisterWorkerNodeContainerReceiverServer(grpcSocketServer.GetGRPCServer(), &handler)
	return handler
}

func (w *WorkerNodeReceiverGRPCHandler) GetTaskFromQueue(ctx context.Context, payload *proto.GetTaskPayload) (*proto.Task, error) {
	task, err := w.workerService.GetTask(payload.ImageUrl)
	return task, err
}

func (w *WorkerNodeReceiverGRPCHandler) SubmitSuccessTask(ctx context.Context, payload *proto.SubmitSuccessTaskRequest) (*emptypb.Empty, error) {
	w.containerSubmitTask.Add(ctx, 1, metric.WithAttributes(attribute.String("status", "success")))
	err := w.workerService.SubmitSuccessTask(ctx, payload.ID, payload.Results)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

func (w *WorkerNodeReceiverGRPCHandler) ReportTaskFailure(ctx context.Context, payload *proto.ReportTaskFailureRequest) (*emptypb.Empty, error) {
	w.containerSubmitTask.Add(ctx, 1, metric.WithAttributes(attribute.String("status", "failure")))
	err := w.workerService.ReportFailTask(ctx, payload.GetID(), payload.GetErrorDetail())
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}
