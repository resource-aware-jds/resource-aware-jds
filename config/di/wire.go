package di

import (
	"github.com/google/wire"
	"github.com/resource-aware-jds/resource-aware-jds/config"
)

var ConfigWireSet = wire.NewSet(
	config.ProvideConfig,
	config.ProvideGRPCConfig,
	config.ProvideControlPlaneConfigModel,
	config.ProvideMongoConfig,
	config.ProvideCACertificateConfig,
	config.ProvideTransportCertificateConfig,
	config.ProvideWorkerConfigModel,
)
