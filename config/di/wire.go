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
	config.ProvideHTTPServerConfig,
)

var WorkerConfigWireSet = wire.NewSet(
	config.ProvideConfig,
	config.ProvideWorkerGRPCConfig,
	config.ProvideMongoConfig,
	config.ProvideTransportCertificateConfig,
	config.ProvideWorkerConfigModel,
	config.ProvideWorkerNodeReceiverConfig,
	config.ProvideGRPCClientConfig,
	config.ProvideClientCATLSCertificateConfig,
	config.ProvideWorkerNodeTransportCertificate,
	config.ProvideWorkerHTTPServerConfig,
)
