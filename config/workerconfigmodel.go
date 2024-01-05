package config

import "github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"

type WorkerConfigModel struct {
	GRPCServerPort                     int    `envconfig:"GRPC_SERVER_PORT" default:"31234"`
	WorkerNodeGRPCServerUnixSocketPath string `envconfig:"WORKER_NODE_GRPC_SERVER_UNIX_SOCKET_PATH" default:"/tmp/rajds_workernode.sock"`
}

func ProvideGRPCSocketServerConfig(config WorkerConfigModel) grpc.SocketServerConfig {
	return grpc.SocketServerConfig{
		UnixSocketPath: config.WorkerNodeGRPCServerUnixSocketPath,
	}
}
