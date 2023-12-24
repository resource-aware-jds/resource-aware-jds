package handler

import (
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
)

type GRPCHandler struct {
	proto.UnimplementedComputeNodeServer
}

func ProvideComputeNodeGRPCHandler(grpcServer grpc.RAJDSGrpc) GRPCHandler {
	handler := GRPCHandler{}
	proto.RegisterComputeNodeServer(grpcServer.GetGRPCServer(), handler)
	return handler
}
