package grpc

import (
	"context"
	"crypto/x509"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/cert"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
	"github.com/resource-aware-jds/resource-aware-jds/service"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
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

	p, _ := peer.FromContext(ctx)

	peerIp := p.Addr.String()
	host, _, err := net.SplitHostPort(peerIp)
	if err != nil {
		return nil, err
	}

	certificate, err := g.controlPlaneService.RegisterWorker(ctx, host, req.Port, parsedKeyData)
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

func (g *GRPCHandler) CreateJob(ctx context.Context, req *proto.CreateJobRequest) (*proto.CreateJobResponse, error) {
	job, tasks, err := g.controlPlaneService.CreateJob(ctx, req.GetImageURL(), req.GetTaskAttributes())
	if err != nil {
		return nil, err
	}

	res := proto.CreateJobResponse{
		ID:       job.ID.Hex(),
		Status:   string(job.Status),
		ImageURL: job.ImageURL,
	}

	responseTasks := make([]*proto.ControlPlaneTask, 0, len(tasks))
	for _, task := range tasks {
		parsedTask := proto.ControlPlaneTask{
			ID:             task.ID.Hex(),
			Status:         string(task.Status),
			TaskAttributes: task.TaskAttributes,
		}

		responseTasks = append(responseTasks, &parsedTask)
	}
	res.Tasks = responseTasks

	return &res, nil
}

func (g *GRPCHandler) WorkerCheckIn(ctx context.Context, req *proto.WorkerCheckInRequest) (*emptypb.Empty, error) {
	p, _ := peer.FromContext(ctx)

	peerIp := p.Addr.String()
	host, _, err := net.SplitHostPort(peerIp)
	if err != nil {
		return nil, err
	}

	err = g.controlPlaneService.CheckInWorkerNode(ctx, host, req.Port, req.GetCertificate())
	return &emptypb.Empty{}, err
}
