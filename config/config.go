package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"runtime"
	"strconv"
)

type Config struct {
	CommonConfigModel
	ControlPlaneConfig ControlPlaneConfigModel `envconfig:"CONTROL_PLANE"`
	WorkerConfig       WorkerConfigModel       `envconfig:"WORKER_NODE"`
	WriteLogToFile     bool                    `envconfig:"WRITE_LOG_TO_FILE"`
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
	}

	var config Config
	envconfig.MustProcess("RAJDS", &config)

	logrus.SetReportCaller(true)

	var formatter logrus.Formatter

	formatter = &logrus.JSONFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			fileName := path.Base(frame.File) + ":" + strconv.Itoa(frame.Line)
			return "", fileName
		},
	}
	if config.DebugMode {
		formatter = &logrus.TextFormatter{
			CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
				fileName := path.Base(frame.File) + ":" + strconv.Itoa(frame.Line)
				//return frame.Function, fileName
				return "", fileName
			},
		}
	}

	logWriter := &doubleLogWrite{}
	if config.WriteLogToFile {
		file, err := os.OpenFile("out.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}

		logWriter.file = file
	}
	logrus.SetOutput(logWriter)
	logrus.SetFormatter(formatter)

	return &config, nil
}

func ProvideControlPlaneGRPCConfig(config *Config) grpc.Config {
	return grpc.Config{
		Port: config.ControlPlaneConfig.GRPCServerPort,
	}
}

func ProvideWorkerGRPCConfig(config *Config) grpc.Config {
	return grpc.Config{
		Port: config.WorkerConfig.GRPCServerPort,
	}
}

func ProvideControlPlaneConfigModel(config *Config) ControlPlaneConfigModel {
	return config.ControlPlaneConfig
}

func ProvideWorkerConfigModel(config *Config) WorkerConfigModel {
	return config.WorkerConfig
}

type doubleLogWrite struct {
	file io.Writer
}

func (d *doubleLogWrite) Write(p []byte) (n int, err error) {
	write, err := os.Stdout.Write(p)
	if err != nil || d.file == nil {
		return write, err
	}

	// Write to the file
	return d.file.Write(p)
}
