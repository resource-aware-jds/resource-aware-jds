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

func ProvideControlPlaneGRPCHandler(grpcServer grpc.RAJDSGrpcServer, controlPlaneService service.IControlPlane) GRPCHandler {
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

	parsedKeyData, err := cert.ParsePublicKeyToKeyData(parsedNodePublicKey)
	if err != nil {
		return nil, err
	}

	certificate, err := g.controlPlaneService.RegisterWorker(ctx, req.Ip, req.Port, parsedKeyData)
	if err != nil {
		return nil, err
	}

	certificateResult := make([]byte, 0)
	pemEncodedCertificateChain := certificate.GetCertificateChains(true)

	for _, item := range pemEncodedCertificateChain {
		certificateResult = append(certificateResult, item...)
	}

	return &proto.ComputeNodeRegistrationResponse{
		Id:          certificate.GetCertificateSubjectSerialNumber(),
		Certificate: certificateResult,
	}, nil
}

func (g *GRPCHandler) CreateJob(ctx context.Context, req *proto.CreateJobRequest) (*proto.CreateJobRequest, error) {
	// TODO: Create Job
	// TODO: Loop Create Task
	// TODO: Response result
	return nil, nil
}
