package config

import (
	"github.com/resource-aware-jds/resource-aware-jds/pkg/cert"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
)

type WorkerConfigModel struct {
	GRPCServerPort                     int    `envconfig:"GRPC_SERVER_PORT" default:"31236"`
	WorkerNodeGRPCServerUnixSocketPath string `envconfig:"WORKER_NODE_GRPC_SERVER_UNIX_SOCKET_PATH" default:"/tmp/rajds_workernode.sock"`
	ControlPlaneHost                   string `envconfig:"CONTROL_PLANE_HOST"`
	CACertificatePath                  string `envconfig:"CA_CERTIFICATE_PATH"`
}

func ProvideGRPCSocketServerConfig(config WorkerConfigModel) grpc.SocketServerConfig {
	return grpc.SocketServerConfig{
		UnixSocketPath: config.WorkerNodeGRPCServerUnixSocketPath,
	}
}

func ProvideClientCATLSCertificateConfig(config WorkerConfigModel) cert.ClientCATLSCertificateConfig {
	return cert.ClientCATLSCertificateConfig{
		CACertificateFilePath: config.CACertificatePath,
	}
}

func ProvideGRPCClientConfig(config WorkerConfigModel, clientCACertificate cert.ClientCATLSCertificate) grpc.ClientConfig {
	return grpc.ClientConfig{
		Target:        config.ControlPlaneHost,
		CACertificate: clientCACertificate,
	}
}
