// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"github.com/resource-aware-jds/resource-aware-jds/cmd/worker/handler"
	http2 "github.com/resource-aware-jds/resource-aware-jds/cmd/worker/http"
	"github.com/resource-aware-jds/resource-aware-jds/config"
	"github.com/resource-aware-jds/resource-aware-jds/daemon"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/cert"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/dockerclient"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/http"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/metrics"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/taskqueue"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/workerlogic"
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
	rajdsgrpcResolver := grpc.ProvideRAJDSGRPCResolver()
	clientConfig := config.ProvideGRPCClientConfig(workerConfigModel, workerNodeCACertificate, rajdsgrpcResolver)
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
	meter, err := metrics.ProvideMeter()
	if err != nil {
		cleanup2()
		cleanup()
		return WorkerApp{}, nil, err
	}
	queue := taskqueue.ProvideTaskQueue(meter)
	workerDistributor := workerlogic.ProvideDelayWorkerDistributor(workerConfigModel)
	iContainer := service.ProvideContainer(client, workerConfigModel, meter)
	iWorker, err := service.ProvideWorker(controlPlaneClient, client, transportCertificate, workerConfigModel, queue, workerDistributor, iContainer, meter)
	if err != nil {
		cleanup2()
		cleanup()
		return WorkerApp{}, nil, err
	}
	iResourceMonitor := service.ProvideResourcesMonitor(client, iContainer, workerConfigModel)
	grpcHandler := handler.ProvideWorkerGRPCHandler(rajdsGrpcServer, iWorker, iResourceMonitor, meter, transportCertificate)
	workerNodeReceiverConfig := config.ProvideWorkerNodeReceiverConfig(workerConfigModel)
	workerNodeReceiverGRPCServer, cleanup3, err := grpc.ProvideWorkerNodeReceiverGRPCServer(workerNodeReceiverConfig)
	if err != nil {
		cleanup2()
		cleanup()
		return WorkerApp{}, nil, err
	}
	workerNodeReceiverGRPCHandler := handler.ProvideWorkerGRPCSocketHandler(workerNodeReceiverGRPCServer, iWorker, meter)
	iDynamicScaling := service.ProvideDynamicScaling(iContainer, iResourceMonitor, workerConfigModel)
	containerTakeDown := workerlogic.ProvideOverResourceUsageContainerTakeDown()
	workerNode := daemon.ProvideWorkerNodeDaemon(client, iWorker, iResourceMonitor, iDynamicScaling, containerTakeDown, iContainer)
	serverConfig := config.ProvideWorkerHTTPServerConfig(workerConfigModel)
	server, cleanup4 := http.ProvideHttpServer(serverConfig)
	routerResult := http2.ProvideHTTPRouter(server)
	workerApp := ProvideWorkerApp(rajdsGrpcServer, grpcHandler, workerNodeReceiverGRPCServer, workerNodeReceiverGRPCHandler, workerNode, server, routerResult)
	return workerApp, func() {
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}
