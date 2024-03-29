package config

import (
	"github.com/resource-aware-jds/resource-aware-jds/pkg/cert"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/http"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/mongo"
	"time"
)

type ControlPlaneConfigModel struct {
	GRPCServerPort            int          `envconfig:"GRPC_SERVER_PORT" default:"31234"`
	HTTPServerPort            int          `envconfig:"HTTP_SERVER_PORT" default:"31313"`
	MongoConfig               mongo.Config `envconfig:"MONGO"`
	CACertificatePath         string       `envconfig:"CA_CERTIFICATE_PATH"`
	CAPrivateKeyPath          string       `envconfig:"CA_PRIVATE_KEY_PATH"`
	CertificatePath           string       `envconfig:"CERTIFICATE_PATH"`
	CertificatePrivateKeyPath string       `envconfig:"CERTIFICATE_PRIVATE_KEY_PATH"`

	ResourceAwareDistributorConfig ResourceAwareDistributorConfigModel `envconfig:"RESOURCE_AWARE_DISTRIBUTOR"`
	TaskWatcherConfig              TaskWatcherConfigModel              `envconfig:"TASK_WATCHER_CONFIG"`
}

type TaskWatcherConfigModel struct {
	SleepTime time.Duration `envconfig:"SLEEP_TIME" default:"1s"`
	Timeout   time.Duration `envconfig:"TIMEOUT" default:"30s"`
}

type ResourceAwareDistributorConfigModel struct {
	AvailableResourceClearanceThreshold float32 `envconfig:"AVAILABLE_RESOURCE_CLEARANCE_THRESHOLD" default:"80"`
}

func ProvideMongoConfig(config ControlPlaneConfigModel) mongo.Config {
	return config.MongoConfig
}

func ProvideCACertificateConfig(config ControlPlaneConfigModel) cert.CACertificateConfig {
	return cert.CACertificateConfig{
		CertificateFileLocation: config.CACertificatePath,
		PrivateKeyFileLocation:  config.CAPrivateKeyPath,
	}
}

func ProvideTransportCertificateConfig(config ControlPlaneConfigModel) cert.TransportCertificateConfig {
	return cert.TransportCertificateConfig{
		CertificateFileLocation: config.CertificatePath,
		PrivateKeyFileLocation:  config.CertificatePrivateKeyPath,
		ValidDuration:           365 * 24 * time.Hour,
		CommonName:              "Resource Aware Job Distribution Transport",
	}
}

func ProvideHTTPServerConfig(config ControlPlaneConfigModel) http.ServerConfig {
	return http.ServerConfig{
		Port: config.HTTPServerPort,
	}
}

func ProvideResourceAwareDistributorConfigMode(config ControlPlaneConfigModel) ResourceAwareDistributorConfigModel {
	return config.ResourceAwareDistributorConfig
}

func ProvideTaskWatcherConfigModel(config ControlPlaneConfigModel) TaskWatcherConfigModel {
	return config.TaskWatcherConfig
}
