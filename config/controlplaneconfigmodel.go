package config

import "github.com/resource-aware-jds/resource-aware-jds/pkg/mongo"

type ControlPlaneConfigModel struct {
	GRPCServerPort int          `envconfig:"GRPC_SERVER_PORT" default:"31234"`
	MongoConfig    mongo.Config `envconfig:"MONGO"`
}
