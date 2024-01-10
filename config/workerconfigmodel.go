package config

import (
	"github.com/resource-aware-jds/resource-aware-jds/pkg/cert"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
)

type WorkerConfigModel struct {
	GRPCServerPort                            int    `envconfig:"GRPC_SERVER_PORT" default:"31236"`
	WorkerNodeReceiverGRPCServerListeningPort int    `envconfig:"WORKER_NODE_RECEIVER_GRPC_SERVER_LISTENING_PORT" default:"31237"`
	ControlPlaneHost                          string `envconfig:"CONTROL_PLANE_HOST"`
	CACertificatePath                         string `envconfig:"CA_CERTIFICATE_PATH"`
	CertificatePath                           string `envconfig:"CERTIFICATE_PATH"`
	CertificatePrivateKeyPath                 string `envconfig:"CERTIFICATE_PRIVATE_KEY_PATH"`
}

func ProvideWorkerNodeReceiverConfig(config WorkerConfigModel) grpc.WorkerNodeReceiverConfig {
	return grpc.WorkerNodeReceiverConfig{
		Port: config.WorkerNodeReceiverGRPCServerListeningPort,
	}

}

func ProvideClientCATLSCertificateConfig(config WorkerConfigModel) cert.WorkerNodeCACertificateConfig {
	return cert.WorkerNodeCACertificateConfig{
		CACertificateFilePath: config.CACertificatePath,
	}
}

func ProvideGRPCClientConfig(config WorkerConfigModel, clientCACertificate cert.WorkerNodeCACertificate) grpc.ClientConfig {
	return grpc.ClientConfig{
		Target:        config.ControlPlaneHost,
		CACertificate: clientCACertificate,
	}
}

func ProvideWorkerNodeTransportCertificate(config WorkerConfigModel) cert.WorkerNodeTransportCertificateConfig {
	return cert.WorkerNodeTransportCertificateConfig{
		CertificateFileLocation: config.CertificatePath,
		PrivateKeyFileLocation:  config.CertificatePrivateKeyPath,
	}
}
