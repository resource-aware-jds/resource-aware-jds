package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/resource-aware-jds/resource-aware-jds/config"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/container"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/datastructure"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/metrics"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/metric"
	"time"
)

type IContainer interface {
	StartContainer(ctx context.Context, imageUrl string) (container.IContainer, error)
	GetContainerIdShort() []string
	GetContainerCoolDownState() datastructure.Buffer[string, time.Time]
	GetContainerBuffer() datastructure.Buffer[string, container.IContainer]
	DownContainer(ctx context.Context, container container.IContainer) error
	AddContainerTakeDownTimer(containerId string) error
	RemoveContainerTakeDownTimer(containerImage string)
}

type ContainerService struct {
	dockerClient              *client.Client
	config                    config.WorkerConfigModel
	containerBuffer           datastructure.Buffer[string, container.IContainer]
	containerImageWithContext datastructure.Buffer[string, func()]
	containerCoolDownState    datastructure.Buffer[string, time.Time]
}

func ProvideContainer(dockerClient *client.Client, config config.WorkerConfigModel, meter metric.Meter) IContainer {
	return &ContainerService{
		dockerClient: dockerClient,
		config:       config,
		containerBuffer: datastructure.ProvideBuffer[string, container.IContainer](
			datastructure.WithBufferMetrics(
				meter,
				metrics.GenerateWorkerNodeMetric("container_buffer"),
				metric.WithUnit("Container"),
				metric.WithDescription("The total container under this worker node agent supervise"),
			),
		),
		containerImageWithContext: make(datastructure.Buffer[string, func()]),
		containerCoolDownState:    make(datastructure.Buffer[string, time.Time]),
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
	c.containerCoolDownState.Store(imageUrl, time.Now().Add(c.config.ContainerStartDelayTimeSeconds))
	return containerInstance, err
}

func (c *ContainerService) DownContainer(ctx context.Context, container container.IContainer) error {
	containerId, ok := container.GetContainerID()
	if !ok {
		logrus.Error("Unable to get container id")
		return fmt.Errorf("unable to get container id")
	}
	err := container.Stop(ctx)
	if err != nil {
		return err
	}
	c.containerBuffer.Pop(containerId)
	return nil
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

func (c *ContainerService) GetContainerBuffer() datastructure.Buffer[string, container.IContainer] {
	return c.containerBuffer
}

func (c *ContainerService) AddContainerTakeDownTimer(containerImage string) error {
	if c.containerImageWithContext.IsObjectInBuffer(containerImage) {
		return nil
	}
	ctx, cancelFunc := context.WithDeadline(context.Background(), time.Now().Add(c.config.ContainerBufferTimeout))
	go func(innerCtx context.Context, innerC *ContainerService, innerContainerImage string) {
		<-ctx.Done()
		if errors.Is(ctx.Err(), context.Canceled) {
			return
		}
		bgnCtx := context.Background()
		innerC.downContainerWithImage(bgnCtx, innerContainerImage)
	}(ctx, c, containerImage)

	c.containerImageWithContext.Store(containerImage, cancelFunc)

	return nil
}

func (c *ContainerService) RemoveContainerTakeDownTimer(containerImage string) {
	funcPointer := c.containerImageWithContext.Pop(containerImage)
	if funcPointer == nil {
		logrus.Warn("Unable to pop containerImageWithContext with image:", containerImage)
		return
	}
	cancelFunc := *funcPointer
	cancelFunc()
}

func (c *ContainerService) downContainerWithImage(ctx context.Context, containerImage string) {
	containerBuffer := c.GetContainerBuffer()
	containers := containerBuffer.GetValues()
	datastructure.Map(containers, func(container container.IContainer) error {
		fmt.Println(containerImage)
		fmt.Println(container.GetImageUrl())
		if container.GetImageUrl() != containerImage {
			return nil
		}
		err := c.DownContainer(ctx, container)
		if err != nil {
			logrus.Error(err)
			return err
		}
		return nil
	})
}
