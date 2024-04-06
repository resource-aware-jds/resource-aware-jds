// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	grpc2 "github.com/resource-aware-jds/resource-aware-jds/cmd/controlplane/grpc"
	http2 "github.com/resource-aware-jds/resource-aware-jds/cmd/controlplane/http"
	"github.com/resource-aware-jds/resource-aware-jds/config"
	"github.com/resource-aware-jds/resource-aware-jds/daemon"
	"github.com/resource-aware-jds/resource-aware-jds/handlerservice"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/cert"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/distribution"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/eventbus"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/http"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/metrics"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/mongo"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/pool"
	"github.com/resource-aware-jds/resource-aware-jds/repository"
	"github.com/resource-aware-jds/resource-aware-jds/service"
)

// Injectors from wire.go:

func InitializeApplication() (ControlPlaneApp, func(), error) {
	configConfig, err := config.ProvideConfig()
	if err != nil {
		return ControlPlaneApp{}, nil, err
	}
	grpcConfig := config.ProvideControlPlaneGRPCConfig(configConfig)
	controlPlaneConfigModel := config.ProvideControlPlaneConfigModel(configConfig)
	transportCertificateConfig := config.ProvideTransportCertificateConfig(controlPlaneConfigModel)
	caCertificateConfig := config.ProvideCACertificateConfig(controlPlaneConfigModel)
	caCertificate, err := cert.ProvideCACertificate(caCertificateConfig)
	if err != nil {
		return ControlPlaneApp{}, nil, err
	}
	transportCertificate, err := cert.ProvideTransportCertificate(transportCertificateConfig, caCertificate)
	if err != nil {
		return ControlPlaneApp{}, nil, err
	}
	rajdsGrpcServer, cleanup, err := grpc.ProvideGRPCServer(grpcConfig, transportCertificate)
	if err != nil {
		return ControlPlaneApp{}, nil, err
	}
	serverConfig := config.ProvideHTTPServerConfig(controlPlaneConfigModel)
	server, cleanup2 := http.ProvideHttpServer(serverConfig)
	mongoConfig := config.ProvideMongoConfig(controlPlaneConfigModel)
	database, cleanup3, err := mongo.ProvideMongoConnection(mongoConfig)
	if err != nil {
		cleanup2()
		cleanup()
		return ControlPlaneApp{}, nil, err
	}
	iNodeRegistry := repository.ProvideControlPlane(database)
	resourceAwareDistributorConfigModel := config.ProvideResourceAwareDistributorConfigMode(controlPlaneConfigModel)
	meter, err := metrics.ProvideMeter(transportCertificate)
	if err != nil {
		cleanup3()
		cleanup2()
		cleanup()
		return ControlPlaneApp{}, nil, err
	}
	iTask := repository.ProvideTask(database)
	task := service.ProvideTaskService(iTask)
	distributorMapper := distribution.ProvideDistributorMapper(resourceAwareDistributorConfigModel, meter, task)
	rajdsgrpcResolver := grpc.ProvideRAJDSGRPCResolver()
	workerNode := pool.ProvideWorkerNode(caCertificate, distributorMapper, rajdsgrpcResolver, meter)
	iControlPlane := handlerservice.ProvideControlPlane(iNodeRegistry, caCertificate, controlPlaneConfigModel, workerNode)
	iJob := repository.ProvideJob(database)
	job := service.ProvideJobService(iJob)
	taskEventBus := eventbus.ProvideTaskEventBus()
	grpcHandler := grpc2.ProvideControlPlaneGRPCHandler(rajdsGrpcServer, iControlPlane, job, task, meter, taskEventBus)
	taskWatcherConfigModel := config.ProvideTaskWatcherConfigModel(controlPlaneConfigModel)
	cpTaskWatcher := service.ProvideCPTaskWatcher(task, taskWatcherConfigModel)
	daemonIControlPlane, cleanup4 := daemon.ProvideControlPlaneDaemon(workerNode, iControlPlane, task, job, cpTaskWatcher, controlPlaneConfigModel)
	httpHandler := http2.ProvideHTTPHandler(job, task)
	nodeHandler := http2.ProvideNodeHandler(workerNode, cpTaskWatcher)
	handler := http2.ProvideHandler(httpHandler, nodeHandler)
	routerResult := http2.ProvideHTTPRouter(handler, server)
	observerInit := ProvideObserverInit(taskEventBus, cpTaskWatcher)
	controlPlaneApp := ProvideControlPlaneApp(rajdsGrpcServer, server, grpcHandler, daemonIControlPlane, routerResult, observerInit)
	return controlPlaneApp, func() {
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}
