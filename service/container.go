package service

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/container"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/datastructure"
	"github.com/sirupsen/logrus"
)

type IContainer interface {
	StartContainer(ctx context.Context, imageUrl string) (container.IContainer, error)
	GetContainerIdShort() []string
}

type ContainerService struct {
	dockerClient    *client.Client
	containerBuffer datastructure.Buffer[string, container.IContainer]
}

func ProvideContainer(dockerClient *client.Client) IContainer {
	return &ContainerService{
		dockerClient:    dockerClient,
		containerBuffer: make(datastructure.Buffer[string, container.IContainer]),
	}
}

func (c *ContainerService) StartContainer(ctx context.Context, imageUrl string) (container.IContainer, error) {
	logrus.Info("Starting container with image:", imageUrl)
	containerInstance := container.ProvideContainer(c.dockerClient, imageUrl, types.ImagePullOptions{})
	err := containerInstance.Start(ctx)
	if err != nil {
		return nil, err
	}
	containerID, ok := containerInstance.GetContainerID()
	if !ok {
		logrus.Error("Unable to get container id")
		return nil, fmt.Errorf("Unable to get container id")
	}
	c.containerBuffer.Store(containerID, containerInstance)
	return containerInstance, err
}

func (c *ContainerService) GetContainerIdShort() []string {
	containerIds := c.containerBuffer.GetKeys()
	return datastructure.Map(containerIds, func(id string) string {
		if len(id) <= ContainerIdShortSize {
			return id
		}
		return id[:ContainerIdShortSize]
	})
}
