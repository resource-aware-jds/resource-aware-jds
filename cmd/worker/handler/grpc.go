package handler

import (
	"context"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
	"github.com/resource-aware-jds/resource-aware-jds/service"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GRPCHandler struct {
	proto.UnimplementedWorkerNodeServer
	workerService service.IWorker
	taskCounter   metric.Int64Counter
}

func ProvideWorkerGRPCHandler(grpcServer grpc.RAJDSGrpcServer, workerService service.IWorker, meter metric.Meter) GRPCHandler {
	taskCounter, err := meter.Int64Counter("total_received_task")
	if err != nil {
		return GRPCHandler{}
	}

	handler := GRPCHandler{
		workerService: workerService,
		taskCounter:   taskCounter,
	}
	proto.RegisterWorkerNodeServer(grpcServer.GetGRPCServer(), &handler)
	return handler
}

func (j *GRPCHandler) HealthCheck(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (j *GRPCHandler) SendTask(ctx context.Context, task *proto.RecievedTask) (*emptypb.Empty, error) {
	j.taskCounter.Add(ctx, 1)
	err := j.workerService.StoreTaskInQueue(task.DockerImage, task.ID, task.TaskAttributes)
	return &emptypb.Empty{}, err
}
