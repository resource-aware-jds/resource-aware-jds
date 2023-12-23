package handler

import (
	"context"
	"crypto/x509"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/cert"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
	"github.com/resource-aware-jds/resource-aware-jds/service"
)

type GRPCHandler struct {
	proto.UnimplementedControlPlaneServer
	controlPlaneService service.IControlPlane
}

func ProvideControlPlaneGRPCHandler(grpcServer grpc.RAJDSGrpc, controlPlaneService service.IControlPlane) GRPCHandler {
	handler := GRPCHandler{
		controlPlaneService: controlPlaneService,
	}
	proto.RegisterControlPlaneServer(grpcServer.GetGRPCServer(), &handler)
	return handler
}

func (g *GRPCHandler) WorkerRegistration(ctx context.Context, req *proto.ComputeNodeRegistrationRequest) (*proto.ComputeNodeRegistrationResponse, error) {
	parsedNodePublicKey, err := x509.ParsePKCS1PublicKey(req.NodePublicKey)
	if err != nil {
		return nil, err
	}

	parsedRAJDSPublicKey, err := cert.ParseToRAJDSPublicKey(parsedNodePublicKey)
	if err != nil {
		return nil, err
	}

	certificate, err := g.controlPlaneService.RegisterWorker(ctx, req.Ip, req.Port, parsedRAJDSPublicKey)
	if err != nil {
		return nil, err
	}

	return &proto.ComputeNodeRegistrationResponse{
		Id:          certificate.GetCertificate().Subject.SerialNumber,
		Certificate: certificate.GetCertificate().Raw,
	}, nil
}
