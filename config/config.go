package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
	"github.com/sirupsen/logrus"
	"os"
)

type Config struct {
	CommonConfigModel
	ControlPlaneConfig ControlPlaneConfigModel `envconfig:"CONTROL_PLANE"`
	WorkerConfig       WorkerConfigModel
}

func ProvideConfig() (*Config, error) {
	// Load the .env file only for
	EnvConfigLocation, ok := os.LookupEnv("ENV_CONFIG")
	if !ok {
		EnvConfigLocation = "./.env"
	}

	err := godotenv.Load(EnvConfigLocation)
	if err != nil {
		logrus.Warn("Can't load env file")
		return nil, err
	}

	var config Config
	envconfig.MustProcess("RAJDS", &config)

	return &config, nil
}

func ProvideGRPCConfig(config *Config) grpc.Config {
	return grpc.Config{
		Port: config.ControlPlaneConfig.GRPCServerPort,
	}
}

func ProvideControlPlaneConfigModel(config *Config) ControlPlaneConfigModel {
	return config.ControlPlaneConfig
}

func ProvideWorkerConfigModel(config *Config) WorkerConfigModel {
	return config.WorkerConfig
}
