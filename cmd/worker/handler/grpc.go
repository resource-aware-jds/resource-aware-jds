package handler

import (
	"context"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
	"github.com/resource-aware-jds/resource-aware-jds/service"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GRPCHandler struct {
	proto.UnimplementedWorkerNodeServer
	workerService service.IWorker
}

func ProvideWorkerGRPCHandler(grpcServer grpc.RAJDSGrpcServer, workerService service.IWorker) GRPCHandler {
	handler := GRPCHandler{
		workerService: workerService,
	}
	proto.RegisterWorkerNodeServer(grpcServer.GetGRPCServer(), &handler)
	return handler
}

func (j *GRPCHandler) SendTask(context context.Context, task *proto.RecievedTask) (*emptypb.Empty, error) {
	err := j.workerService.SubmitTask(task.DockerImage, task.ID, task.TaskAttributes)
	return &emptypb.Empty{}, err
}
