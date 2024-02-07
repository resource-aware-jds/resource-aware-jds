package handler

import (
	"context"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
	"github.com/resource-aware-jds/resource-aware-jds/service"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/protobuf/types/known/emptypb"
)

type WorkerNodeReceiverGRPCHandler struct {
	proto.UnimplementedWorkerNodeContainerReceiverServer
	workerService      service.IWorker
	successTaskCounter metric.Int64Counter
	failureTaskCounter metric.Int64Counter
}

func ProvideWorkerGRPCSocketHandler(grpcSocketServer grpc.WorkerNodeReceiverGRPCServer, workerService service.IWorker, meter metric.Meter) WorkerNodeReceiverGRPCHandler {
	successTaskCounter, err := meter.Int64Counter("success_task_total")
	if err != nil {
		return WorkerNodeReceiverGRPCHandler{}
	}

	failureTaskCounter, err := meter.Int64Counter("failure_task_total")
	if err != nil {
		return WorkerNodeReceiverGRPCHandler{}
	}

	handler := WorkerNodeReceiverGRPCHandler{
		workerService:      workerService,
		failureTaskCounter: failureTaskCounter,
		successTaskCounter: successTaskCounter,
	}
	proto.RegisterWorkerNodeContainerReceiverServer(grpcSocketServer.GetGRPCServer(), &handler)
	return handler
}

func (w *WorkerNodeReceiverGRPCHandler) GetTaskFromQueue(ctx context.Context, payload *proto.GetTaskPayload) (*proto.Task, error) {
	task, err := w.workerService.GetTask(payload.ImageUrl)
	return task, err
}

func (w *WorkerNodeReceiverGRPCHandler) SubmitSuccessTask(ctx context.Context, payload *proto.SubmitSuccessTaskRequest) (*emptypb.Empty, error) {
	w.successTaskCounter.Add(ctx, 1)
	err := w.workerService.SubmitSuccessTask(payload.ID, payload.Results)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

func (w *WorkerNodeReceiverGRPCHandler) ReportTaskFailure(ctx context.Context, payload *proto.ReportTaskFailureRequest) (*emptypb.Empty, error) {
	w.failureTaskCounter.Add(ctx, 1)
	err := w.workerService.ReportFailTask(ctx, payload.GetID(), payload.GetErrorDetail())
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}
