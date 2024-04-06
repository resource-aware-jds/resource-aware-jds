package handler

import (
	"context"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/handlerservice"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/cert"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/metrics"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/util"
	"github.com/resource-aware-jds/resource-aware-jds/service"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GRPCHandler struct {
	proto.UnimplementedWorkerNodeServer
	workerService             handlerservice.IWorker
	taskCounter               metric.Int64Counter
	resourceMonitoringService service.IResourceMonitor
	workerNodeCertificate     cert.TransportCertificate
}

func ProvideWorkerGRPCHandler(grpcServer grpc.RAJDSGrpcServer, workerService handlerservice.IWorker, resourceMonitoringService service.IResourceMonitor, meter metric.Meter, workerNodeCertificate cert.TransportCertificate) GRPCHandler {
	taskCounter, err := meter.Int64Counter(
		metrics.GenerateWorkerNodeMetric("total_received_task"),
		metric.WithUnit("Task"),
		metric.WithDescription("The total received task in this Worker Node"),
	)
	if err != nil {
		return GRPCHandler{}
	}

	handler := GRPCHandler{
		workerService:             workerService,
		taskCounter:               taskCounter,
		resourceMonitoringService: resourceMonitoringService,
		workerNodeCertificate:     workerNodeCertificate,
	}
	proto.RegisterWorkerNodeServer(grpcServer.GetGRPCServer(), &handler)
	return handler
}

func (j *GRPCHandler) HealthCheck(ctx context.Context, req *emptypb.Empty) (*proto.Resource, error) {
	availableTaskSlot := j.workerService.GetAvailableTaskSlot()
	if availableTaskSlot <= 0 {
		return &proto.Resource{
			CpuCores:               0,
			AvailableCpuPercentage: 0,
			AvailableMemory:        "0Mib",
		}, nil
	}
	resource, err := j.resourceMonitoringService.CalculateAvailableResource(ctx)
	if err != nil {
		return nil, err
	}

	logrus.Info("Health check report, cpu: ", resource.AvailableCpuPercentage, ", memory: ", util.MemoryToString(resource.AvailableMemory))
	return &proto.Resource{
		CpuCores:               resource.CpuCores,
		AvailableCpuPercentage: resource.AvailableCpuPercentage,
		AvailableMemory:        util.MemoryToString(resource.AvailableMemory),
	}, nil
}

func (j *GRPCHandler) SendTask(ctx context.Context, task *proto.RecievedTask) (*emptypb.Empty, error) {
	j.taskCounter.Add(
		ctx,
		1,
		metric.WithAttributes(attribute.String("nodeID", j.workerNodeCertificate.GetNodeID())),
	)
	err := j.workerService.StoreTaskInQueue(task.DockerImage, task.ID, task.TaskAttributes)
	return &emptypb.Empty{}, err
}

func (j *GRPCHandler) GetAllTasks(context.Context, *emptypb.Empty) (*proto.TaskResponse, error) {
	runningTask := j.workerService.GetRunningTask()
	queuedTask := j.workerService.GetQueuedTask()

	return &proto.TaskResponse{TaskIDs: append(runningTask, queuedTask...)}, nil
}
