package grpc

import (
	"context"
	"crypto/x509"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/cert"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
	"github.com/resource-aware-jds/resource-aware-jds/service"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
)

type GRPCHandler struct {
	proto.UnimplementedControlPlaneServer
	controlPlaneService service.IControlPlane
	jobService          service.Job
	taskService         service.Task
}

func ProvideControlPlaneGRPCHandler(grpcServer grpc.RAJDSGrpcServer, controlPlaneService service.IControlPlane, jobService service.Job, taskService service.Task) GRPCHandler {
	handler := GRPCHandler{
		controlPlaneService: controlPlaneService,
		jobService:          jobService,
		taskService:         taskService,
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
	job, err := g.jobService.CreateJob(ctx, req.GetName(), req.GetImageURL())
	if err != nil {
		return nil, err
	}

	tasks, err := g.taskService.CreateTask(ctx, job, req.GetTaskAttributes())
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

func (g *GRPCHandler) ReportFailureTask(ctx context.Context, req *proto.ReportFailureTaskRequest) (*emptypb.Empty, error) {
	id := req.GetId()
	parsedTaskID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logrus.Errorf("parsing task id error %v", err)
		return nil, err
	}

	err = g.taskService.UpdateTaskWorkOnFailure(ctx, parsedTaskID, req.GetNodeID(), req.GetMessage())
	return &emptypb.Empty{}, err
}

func (g *GRPCHandler) ReportSuccessTask(ctx context.Context, req *proto.ReportSuccessTaskRequest) (*emptypb.Empty, error) {
	id := req.GetId()
	parsedTaskID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logrus.Errorf("parsing task id error %v", err)
		return nil, err
	}

	err = g.taskService.UpdateTaskSuccess(ctx, parsedTaskID, req.GetNodeID(), req.GetResult())
	return &emptypb.Empty{}, err
}
