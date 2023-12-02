package handler

import (
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
)

type GRPCHandler struct {
	proto.UnimplementedControlPlaneServer
}

func ProvideControlPlaneGRPCHandler(grpcServer grpc.RAJDSGrpc) GRPCHandler {
	handler := GRPCHandler{}
	proto.RegisterControlPlaneServer(grpcServer.GetGRPCServer(), handler)
	return handler
}
