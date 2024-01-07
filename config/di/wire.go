package di

import (
	"github.com/google/wire"
	"github.com/resource-aware-jds/resource-aware-jds/config"
)

var ControlPlaneConfigWireSet = wire.NewSet(
	config.ProvideConfig,
	config.ProvideControlPlaneGRPCConfig,
	config.ProvideControlPlaneConfigModel,
	config.ProvideMongoConfig,
	config.ProvideCACertificateConfig,
	config.ProvideTransportCertificateConfig,
	config.ProvideWorkerConfigModel,
	config.ProvideGRPCSocketServerConfig,
)

var WorkerConfigWireSet = wire.NewSet(
	config.ProvideConfig,
	config.ProvideWorkerGRPCConfig,
	config.ProvideControlPlaneConfigModel,
	config.ProvideMongoConfig,
	config.ProvideCACertificateConfig,
	config.ProvideTransportCertificateConfig,
	config.ProvideWorkerConfigModel,
	config.ProvideGRPCSocketServerConfig,
	config.ProvideGRPCClientConfig,
	config.ProvideClientCATLSCertificateConfig,
)
