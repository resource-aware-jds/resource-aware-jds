// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"github.com/resource-aware-jds/resource-aware-jds/cmd/worker/handler"
	"github.com/resource-aware-jds/resource-aware-jds/config"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/cert"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/dockerclient"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/taskqueue"
	"github.com/resource-aware-jds/resource-aware-jds/service"
)

// Injectors from wire.go:

func InitializeApplication() (WorkerApp, func(), error) {
	configConfig, err := config.ProvideConfig()
	if err != nil {
		return WorkerApp{}, nil, err
	}
	grpcConfig := config.ProvideGRPCConfig(configConfig)
	controlPlaneConfigModel := config.ProvideControlPlaneConfigModel(configConfig)
	transportCertificateConfig := config.ProvideTransportCertificateConfig(controlPlaneConfigModel)
	caCertificateConfig := config.ProvideCACertificateConfig(controlPlaneConfigModel)
	caCertificate, err := cert.ProvideCACertificate(caCertificateConfig)
	if err != nil {
		return WorkerApp{}, nil, err
	}
	transportCertificate, err := cert.ProvideTransportCertificate(transportCertificateConfig, caCertificate)
	if err != nil {
		return WorkerApp{}, nil, err
	}
	rajdsGrpcServer, cleanup, err := grpc.ProvideGRPCServer(grpcConfig, transportCertificate)
	if err != nil {
		return WorkerApp{}, nil, err
	}
	client, cleanup2, err := dockerclient.ProvideDockerClient()
	if err != nil {
		cleanup()
		return WorkerApp{}, nil, err
	}
	workerConfigModel := config.ProvideWorkerConfigModel(configConfig)
	queue := taskqueue.ProvideTaskQueue()
	iWorker := service.ProvideWorker(client, workerConfigModel, queue)
	grpcHandler := handler.ProvideWorkerGRPCHandler(rajdsGrpcServer, iWorker)
	socketServerConfig := config.ProvideGRPCSocketServerConfig(workerConfigModel)
	socketServer, cleanup3, err := grpc.ProvideGRPCSocketServer(socketServerConfig)
	if err != nil {
		cleanup2()
		cleanup()
		return WorkerApp{}, nil, err
	}
	workerGRPCSocketHandler := handler.ProvideWorkerGRPCSocketHandler(socketServer)
	workerApp := ProvideWorkerApp(rajdsGrpcServer, grpcHandler, socketServer, workerGRPCSocketHandler)
	return workerApp, func() {
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}
