package handler

import (
	"context"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/util"
	"github.com/resource-aware-jds/resource-aware-jds/service"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GRPCHandler struct {
	proto.UnimplementedWorkerNodeServer
	workerService             service.IWorker
	taskCounter               metric.Int64Counter
	resourceMonitoringService service.IResourceMonitor
}

func ProvideWorkerGRPCHandler(grpcServer grpc.RAJDSGrpcServer, workerService service.IWorker, resourceMonitoringService service.IResourceMonitor, meter metric.Meter) GRPCHandler {
	taskCounter, err := meter.Int64Counter("total_received_task")
	if err != nil {
		return GRPCHandler{}
	}

	handler := GRPCHandler{
		workerService:             workerService,
		taskCounter:               taskCounter,
		resourceMonitoringService: resourceMonitoringService,
	}
	proto.RegisterWorkerNodeServer(grpcServer.GetGRPCServer(), &handler)
	return handler
}

func (j *GRPCHandler) HealthCheck(ctx context.Context, req *emptypb.Empty) (*proto.Resource, error) {
	resource, err := j.resourceMonitoringService.CalculateAvailableResource(ctx)
	if err != nil {
		return nil, err
	}
	return &proto.Resource{
		CpuCores:               resource.CpuCores,
		AvailableCpuPercentage: resource.AvailableCpuPercentage,
		AvailableMemory:        util.MemoryToString(resource.AvailableMemory),
	}, nil
}

func (j *GRPCHandler) SendTask(ctx context.Context, task *proto.RecievedTask) (*emptypb.Empty, error) {
	j.taskCounter.Add(ctx, 1)
	err := j.workerService.StoreTaskInQueue(task.DockerImage, task.ID, task.TaskAttributes)
	return &emptypb.Empty{}, err
}
