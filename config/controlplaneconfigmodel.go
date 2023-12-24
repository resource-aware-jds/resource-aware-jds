package config

import (
	"github.com/resource-aware-jds/resource-aware-jds/pkg/cert"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/mongo"
	"time"
)

type ControlPlaneConfigModel struct {
	GRPCServerPort               int          `envconfig:"GRPC_SERVER_PORT" default:"31234"`
	MongoConfig                  mongo.Config `envconfig:"MONGO"`
	CACertificatePath            string       `envconfig:"CA_CERTIFICATE_PATH"`
	CAPrivateKeyPath             string       `envconfig:"CA_PRIVATE_KEY_PATH"`
	CertificatePath              string       `envconfig:"CERTIFICATE_PATH"`
	CertificatePrivateKeyPath    string       `envconfig:"CERTIFICATE_PRIVATE_KEY_PATH"`
	ClientCertificateStoragePath string       `envconfig:"CLIENT_CERTIFICATE_STORAGE_PATH"`
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
