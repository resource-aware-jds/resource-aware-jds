package grpc

import (
	"crypto/tls"
	"fmt"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/cert"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
)

type rajdsGRPC struct {
	Listener   net.Listener
	GRPCServer *grpc.Server
}

type RAJDSGrpc interface {
	Serve()
	GetGRPCServer() *grpc.Server
}

type Config struct {
	Port int
}

func ProvideGRPCServer(config Config, transportCertificate cert.TransportCertificate) (RAJDSGrpc, func(), error) {
	// GRPC Server Listening
	lis, err := net.Listen("tcp", fmt.Sprint(":", config.Port))
	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
		return nil, nil, err
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{
			{
				Certificate: transportCertificate.GetCertificateChains(),
				PrivateKey:  transportCertificate.GetPrivateKey(),
			},
		},
		ClientAuth: tls.NoClientCert,
	}

	grpcServer := grpc.NewServer(
		grpc.Creds(credentials.NewTLS(tlsConfig)),
	)

	cleanup := func() {
		grpcServer.GracefulStop()
	}

	return &rajdsGRPC{
		Listener:   lis,
		GRPCServer: grpcServer,
	}, cleanup, nil
}

func (r *rajdsGRPC) Serve() {
	go func() {
		logrus.Info("GRPC Server is Listening on: ", r.Listener.Addr())
		r.GRPCServer.Serve(r.Listener)
	}()
}

func (r *rajdsGRPC) GetGRPCServer() *grpc.Server {
	return r.GRPCServer
}
