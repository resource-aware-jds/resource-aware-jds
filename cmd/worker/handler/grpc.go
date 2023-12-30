package handler

import (
	"context"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
	"github.com/resource-aware-jds/resource-aware-jds/service"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
	"strconv"
)

type GRPCHandler struct {
	proto.UnimplementedComputeNodeServer
	workerService service.IWorker
}

func ProvideWorkerGRPCHandler(grpcServer grpc.RAJDSGrpcServer, workerService service.IWorker) GRPCHandler {
	handler := GRPCHandler{
		workerService: workerService,
	}
	proto.RegisterComputeNodeServer(grpcServer.GetGRPCServer(), &handler)
	return handler
}

func (j *GRPCHandler) SendTask(context context.Context, task *proto.Task) (*emptypb.Empty, error) {
	taskId := strconv.Itoa(int(task.TaskId))
	err := j.workerService.SubmitTask(task.DockerImage, taskId)
	return &emptypb.Empty{}, err
}

func (j *GRPCHandler) ReportJob(context context.Context, report *proto.ReportTaskRequest) (*emptypb.Empty, error) {
	if report.TotalJob == report.CurrentJob {
		jobIdStr := strconv.Itoa(int(report.TaskId))
		containerName := "rajds-" + jobIdStr
		go j.workerService.RemoveContainer(containerName)
	}
	logrus.Info("Job id: ", report.TaskId, " Current: ", report.CurrentJob, " Total: ", report.TotalJob)
	return &emptypb.Empty{}, nil
}
