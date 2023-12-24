package grpc

import (
	"crypto/tls"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type rajdsGRPCClient struct {
	connection *grpc.ClientConn
}

type RAJDSGrpcClient interface {
	GetConnection() *grpc.ClientConn
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair("/Users/sirateek/.rajds/controlplane/transport/cert.pem", "/Users/sirateek/.rajds/controlplane/transport/key.pem")
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}

	return credentials.NewTLS(config), nil
}

func ProvideRAJDSGrpcClient() (RAJDSGrpcClient, error) {
	cert, err := loadTLSCredentials()
	if err != nil {
		return &rajdsGRPCClient{}, err
	}

	grpcConnection, err := grpc.Dial(
		"localhost:31234",
		grpc.WithTransportCredentials(cert),
	)
	if err != nil {
		return nil, err
	}
	return &rajdsGRPCClient{
		connection: grpcConnection,
	}, nil
}

func (r *rajdsGRPCClient) GetConnection() *grpc.ClientConn {
	return r.connection
}
