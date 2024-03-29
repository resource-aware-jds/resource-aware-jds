package grpc

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/cert"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"net"
)

type rajdsGRPCServer struct {
	logger     *logrus.Entry
	Listener   net.Listener
	GRPCServer *grpc.Server
}

type RAJDSGrpcServer interface {
	Serve()
	GetGRPCServer() *grpc.Server
}

type Config struct {
	Port int
}

func ProvideGRPCServer(config Config, transportCertificate cert.TransportCertificate) (RAJDSGrpcServer, func(), error) {
	// GRPC Server Listening
	lis, err := net.Listen("tcp", fmt.Sprint(":", config.Port))
	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
		return nil, nil, err
	}

	// Create Client CA Pool
	caPool := x509.NewCertPool()
	caCertificate, err := transportCertificate.GetCACertificate()
	if err != nil {
		logrus.Errorf("failed to get CA certificate: %v", err)
		return nil, nil, err
	}
	caPool.AddCert(caCertificate)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{
			{
				Certificate: transportCertificate.GetCertificateChains(false)[:1],
				PrivateKey:  transportCertificate.GetPrivateKey().GetRawKeyData(),
			},
		},
		ClientAuth: tls.NoClientCert,
	}

	grpcServer := grpc.NewServer(
		grpc.Creds(credentials.NewTLS(tlsConfig)),
		grpc.UnaryInterceptor(grpcUnaryInterceptor),
	)

	cleanup := func() {
		grpcServer.GracefulStop()
	}

	return &rajdsGRPCServer{
		logger:     logrus.WithField("type", "Main Receiver GRPC Server"),
		Listener:   lis,
		GRPCServer: grpcServer,
	}, cleanup, nil
}

func (r *rajdsGRPCServer) Serve() {
	go func() {
		r.logger.Info("GRPC Server is Listening on: ", r.Listener.Addr())
		r.GRPCServer.Serve(r.Listener) // nolint:errcheck
	}()
}

func (r *rajdsGRPCServer) GetGRPCServer() *grpc.Server {
	return r.GRPCServer
}

type WorkerNodeReceiverGRPCServer RAJDSGrpcServer

type WorkerNodeReceiverConfig Config

func ProvideWorkerNodeReceiverGRPCServer(config WorkerNodeReceiverConfig) (WorkerNodeReceiverGRPCServer, func(), error) {
	// GRPC Server Listening
	lis, err := net.Listen("tcp", fmt.Sprint(":", config.Port))
	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
		return nil, nil, err
	}

	grpcServer := grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(grpcUnaryInterceptor),
	)

	cleanup := func() {
		grpcServer.GracefulStop()
	}

	return &rajdsGRPCServer{
		logger:     logrus.WithField("type", "WorkerNodeReceiver GRPC Server"),
		Listener:   lis,
		GRPCServer: grpcServer,
	}, cleanup, nil
}
