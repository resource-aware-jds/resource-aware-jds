package handler

import (
	"context"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type WorkerGRPCSocketHandler struct {
	proto.UnimplementedWorkerNodeContainerReceiverServer
}

func ProvideWorkerGRPCSocketHandler(grpcSocketServer grpc.SocketServer) WorkerGRPCSocketHandler {
	handler := WorkerGRPCSocketHandler{}
	proto.RegisterWorkerNodeContainerReceiverServer(grpcSocketServer.GetGRPCServer(), &handler)
	return handler
}

func (w *WorkerGRPCSocketHandler) GetTaskFromQueue(ctx context.Context, _ *emptypb.Empty) (*proto.Task, error) {
	return &proto.Task{
		ID:             "ABC:123",
		TaskAttributes: make([]byte, 0),
	}, nil
}
