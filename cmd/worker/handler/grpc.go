package handler

import (
	"context"
	"github.com/docker/docker/api/types"
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

func (j *GRPCHandler) SendJob(context context.Context, job *proto.Job) (*emptypb.Empty, error) {
	jobIdStr := strconv.Itoa(int(job.JobID))
	containerName := "rajds-" + jobIdStr
	err := j.workerService.RunJob(job.DockerImage, containerName, types.ImagePullOptions{}, jobIdStr)
	return &emptypb.Empty{}, err
}

func (j *GRPCHandler) ReportJob(context context.Context, report *proto.ReportJobRequest) (*emptypb.Empty, error) {
	if report.TotalJob == report.CurrentJob {
		jobIdStr := strconv.Itoa(int(report.JobID))
		containerName := "rajds-" + jobIdStr
		go j.workerService.RemoveContainer(containerName)
	}
	logrus.Info("Job id: ", report.JobID, " Current: ", report.CurrentJob, " Total: ", report.TotalJob)
	return &emptypb.Empty{}, nil
}
