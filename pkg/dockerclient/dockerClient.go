package dockerclient

import (
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

func ProvideDockerClient() (*client.Client, func(), error) {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logrus.Error("Connect Docker Client Error: ", err)
		return nil, func() {}, err
	}

	cleanup := func() {
		err = dockerClient.Close()
		if err != nil {
			logrus.Error("Docker Client Close Connection Error: ", err)
		}
	}

	return dockerClient, cleanup, err
}
