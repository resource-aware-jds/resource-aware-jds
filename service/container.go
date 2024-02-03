package service

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/resource-aware-jds/resource-aware-jds/config"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/container"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/datastructure"
	"github.com/sirupsen/logrus"
	"time"
)

type IContainer interface {
	StartContainer(ctx context.Context, imageUrl string) (container.IContainer, error)
	GetContainerIdShort() []string
	GetContainerCoolDownState() datastructure.Buffer[string, time.Time]
}

type ContainerService struct {
	dockerClient           *client.Client
	config                 config.WorkerConfigModel
	containerBuffer        datastructure.Buffer[string, container.IContainer]
	containerCoolDownState datastructure.Buffer[string, time.Time]
}

func ProvideContainer(dockerClient *client.Client, config config.WorkerConfigModel) IContainer {
	return &ContainerService{
		dockerClient:           dockerClient,
		config:                 config,
		containerBuffer:        make(datastructure.Buffer[string, container.IContainer]),
		containerCoolDownState: make(datastructure.Buffer[string, time.Time]),
	}
}

func (c *ContainerService) StartContainer(ctx context.Context, imageUrl string) (container.IContainer, error) {
	// Remove ContainerCoolDownState
	delete(c.containerCoolDownState, imageUrl)

	logrus.Info("Starting container with image:", imageUrl)
	containerInstance := container.ProvideContainer(c.dockerClient, imageUrl, types.ImagePullOptions{})
	err := containerInstance.Start(ctx)
	if err != nil {
		return nil, err
	}
	containerID, ok := containerInstance.GetContainerID()
	if !ok {
		logrus.Error("Unable to get container id")
		return nil, fmt.Errorf("unable to get container id")
	}

	c.containerBuffer.Store(containerID, containerInstance)
	c.containerCoolDownState[imageUrl] = time.Now().Add(c.config.ContainerStartDelayTimeSeconds)
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

func (c *ContainerService) GetContainerCoolDownState() datastructure.Buffer[string, time.Time] {
	return c.containerCoolDownState
}
