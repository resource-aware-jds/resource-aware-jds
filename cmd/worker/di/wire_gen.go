// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"github.com/resource-aware-jds/resource-aware-jds/cmd/worker/handler"
	"github.com/resource-aware-jds/resource-aware-jds/config"
	"github.com/resource-aware-jds/resource-aware-jds/daemon"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/cert"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/dockerclient"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/taskBuffer"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/taskqueue"
	"github.com/resource-aware-jds/resource-aware-jds/service"
)

// Injectors from wire.go:

func InitializeApplication() (WorkerApp, func(), error) {
	configConfig, err := config.ProvideConfig()
	if err != nil {
		return WorkerApp{}, nil, err
	}
	grpcConfig := config.ProvideWorkerGRPCConfig(configConfig)
	workerConfigModel := config.ProvideWorkerConfigModel(configConfig)
	workerNodeTransportCertificateConfig := config.ProvideWorkerNodeTransportCertificate(workerConfigModel)
	workerNodeCACertificateConfig := config.ProvideClientCATLSCertificateConfig(workerConfigModel)
	workerNodeCACertificate, err := cert.ProvideWorkerNodeCACertificate(workerNodeCACertificateConfig)
	if err != nil {
		return WorkerApp{}, nil, err
	}
	clientConfig := config.ProvideGRPCClientConfig(workerConfigModel, workerNodeCACertificate)
	rajdsGrpcClient, err := grpc.ProvideRAJDSGrpcClient(clientConfig)
	if err != nil {
		return WorkerApp{}, nil, err
	}
	controlPlaneClient := ProvideControlPlaneGRPCClient(rajdsGrpcClient)
	transportCertificate, err := cert.ProvideWorkerNodeTransportCertificate(workerNodeTransportCertificateConfig, controlPlaneClient)
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
	queue := taskqueue.ProvideTaskQueue()
	taskBufferTaskBuffer := taskBuffer.ProvideTaskBuffer()
	iWorker := service.ProvideWorker(client, workerConfigModel, queue, taskBufferTaskBuffer)
	grpcHandler := handler.ProvideWorkerGRPCHandler(rajdsGrpcServer, iWorker)
	socketServerConfig := config.ProvideGRPCSocketServerConfig(workerConfigModel)
	socketServer, cleanup3, err := grpc.ProvideGRPCSocketServer(socketServerConfig)
	if err != nil {
		cleanup2()
		cleanup()
		return WorkerApp{}, nil, err
	}
	workerGRPCSocketHandler := handler.ProvideWorkerGRPCSocketHandler(socketServer, iWorker)
	workerNode := daemon.ProvideWorkerNodeDaemon(controlPlaneClient, iWorker, queue, transportCertificate, workerConfigModel)
	workerApp := ProvideWorkerApp(rajdsGrpcServer, grpcHandler, socketServer, workerGRPCSocketHandler, workerNode)
	return workerApp, func() {
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}
