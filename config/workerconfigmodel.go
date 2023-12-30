package config

type WorkerConfigModel struct {
	GRPCServerPort int `envconfig:"GRPC_SERVER_PORT" default:"31234"`
}
