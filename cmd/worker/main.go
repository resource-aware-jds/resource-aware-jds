package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/cert"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
	"github.com/sirupsen/logrus"
)

func main() {
	caCertificate, err := cert.ProvideClientCATLSCertificate(cert.ClientCATLSCertificateConfig{
		CACertificateFilePath: "/Users/sirateek/.rajds/controlplane/ca/cert.pem",
	})
	if err != nil {
		panic(err)
	}

	grpcConn, err := grpc.ProvideRAJDSGrpcClient("localhost:31234", caCertificate)
	if err != nil {
		panic(err)
	}

	controlPlaneClient := proto.NewControlPlaneClient(grpcConn.GetConnection())

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)

	if err != nil {
		panic(err)
	}

	result, err := controlPlaneClient.WorkerRegistration(context.Background(), &proto.ComputeNodeRegistrationRequest{
		Ip:            "1234",
		Port:          1234,
		NodePublicKey: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
	})
	if err != nil {
		logrus.Error(err)
		panic(err)
	}

	logrus.Info(result)

}
