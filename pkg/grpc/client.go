package grpc

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/cert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type rajdsGRPCClient struct {
	connection *grpc.ClientConn
}

type RAJDSGrpcClient interface {
	GetConnection() *grpc.ClientConn
}

func ProvideRAJDSGrpcClient(target string, certificate cert.TLSCertificate) (RAJDSGrpcClient, error) {
	// Create the trusted CA Pool
	caCertificatePool := x509.NewCertPool()

	caCertificate, err := certificate.GetCACertificate()
	if err != nil {
		return nil, err
	}

	caCertificatePool.AddCert(caCertificate)
	tlsConfig := &tls.Config{
		RootCAs: caCertificatePool,
	}

	grpcConnection, err := grpc.Dial(
		target,
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
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
