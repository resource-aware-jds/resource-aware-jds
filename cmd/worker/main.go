package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"github.com/resource-aware-jds/resource-aware-jds/cmd/worker/di"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/cert"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logrus.Info("\"Starting up the Worker.\"")
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

	logrus.Info("\"Starting up Worker GRPC server.\"")
	app, cleanup, err := di.InitializeApplication()
	app.GRPCServer.Serve()

	// Gracefully Shutdown
	// Make channel listen for signals from OS
	gracefulStop := make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	<-gracefulStop

	logrus.Info("Gracefully shutting down, cleaning up.")
	cleanup()
	logrus.Info("Clean up success. Good Bye")
}
