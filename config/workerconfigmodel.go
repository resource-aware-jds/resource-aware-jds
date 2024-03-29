package config

import (
	"fmt"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/cert"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
	httpServer "github.com/resource-aware-jds/resource-aware-jds/pkg/http"
	"time"
)

type WorkerConfigModel struct {
	GRPCServerPort                            int           `envconfig:"GRPC_SERVER_PORT" default:"31236"`
	WorkerNodeReceiverGRPCServerListeningPort int           `envconfig:"WORKER_NODE_RECEIVER_GRPC_SERVER_LISTENING_PORT" default:"31237"`
	ControlPlaneHost                          string        `envconfig:"CONTROL_PLANE_HOST"`
	CACertificatePath                         string        `envconfig:"CA_CERTIFICATE_PATH"`
	CertificatePath                           string        `envconfig:"CERTIFICATE_PATH"`
	CertificatePrivateKeyPath                 string        `envconfig:"CERTIFICATE_PRIVATE_KEY_PATH"`
	ContainerStartDelayTimeSeconds            time.Duration `envconfig:"CONTAINER_START_DELAY_TIME_SECONDS" default:"10s"`
	HTTPServerPort                            int           `envconfig:"HTTP_SERVER_PORT" default:"30001"`
	MaxMemoryUsage                            string        `envconfig:"MAX_MEMORY_USAGE" default:"16GiB"`
	MemoryBufferSize                          string        `envconfig:"MEMORY_BUFFER_SIZE" default:"1GiB"`
	MaxCpuUsagePercentage                     int           `envconfig:"MAX_CPU_USAGE_PERCENTAGE" default:"100"`
	CpuBufferSize                             int           `envconfig:"CPU_BUFFER_SIZE" default:"20"`
	DockerCoreLimit                           int           `envconfig:"DOCKER_CORE_LIMIT" required:"true"`
	TotalContainerLimit                       int           `envconfig:"TOTAL_CONTAINER_CONFIG" default:"1"`
	TaskBufferTimeout                         time.Duration `envconfig:"TASK_BUFFER_TIMEOUT" default:"30s"`
	ContainerBufferTimeout                    time.Duration `envconfig:"CONTAINER_BUFFER_TIMEOUT" default:"30s"`
	ContainerLogDir                           string        `envconfig:"CONTAINER_LOG_DIR" default:"/tmp/rajds-log"`
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

func ProvideGRPCClientConfig(config WorkerConfigModel, clientCACertificate cert.WorkerNodeCACertificate, grpcResolver grpc.RAJDSGRPCResolver) grpc.ClientConfig {
	// Add ControlPlaneHost is in the Local DNS Resolver
	focusedDomainName := fmt.Sprintf("cp.%s", cert.GetDefaultDomainName())
	grpcResolver.AddHost(focusedDomainName, config.ControlPlaneHost)

	return grpc.ClientConfig{
		Target:        "rajds://cp.rajds",
		CACertificate: clientCACertificate,
		ServerName:    "cp.rajds",
	}
}

func ProvideWorkerNodeTransportCertificate(config WorkerConfigModel) cert.WorkerNodeTransportCertificateConfig {
	return cert.WorkerNodeTransportCertificateConfig{
		CertificateFileLocation: config.CertificatePath,
		PrivateKeyFileLocation:  config.CertificatePrivateKeyPath,
	}
}

func ProvideWorkerHTTPServerConfig(config WorkerConfigModel) httpServer.ServerConfig {
	return httpServer.ServerConfig{
		Port: config.HTTPServerPort,
	}
}
