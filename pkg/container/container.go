package container

import (
	"context"
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/util"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
	"time"
)

var (
	ErrContainerIsNotStarted = errors.New("container is not started")
)

type containerSvc struct {
	dockerClient     *client.Client
	imageURL         string
	containerName    string
	imagePullOptions types.ImagePullOptions
	containerID      *string
	startupTime      time.Time
}

type IContainer interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	RemoveContainer(ctx context.Context) error
	GetContainerID() (string, bool)
	GetContainerName() string
}

func ProvideContainer(dockerClient *client.Client, imageURL string, imagePullOptions types.ImagePullOptions) IContainer {
	randomId := rand.Intn(50000-10000) + 10000
	containerName := "rads-" + strconv.Itoa(randomId)

	return &containerSvc{
		dockerClient:     dockerClient,
		imageURL:         imageURL,
		containerName:    containerName,
		imagePullOptions: imagePullOptions,
	}
}

func (c *containerSvc) Start(ctx context.Context) error {
	logrus.Info("Creating container: ", c.containerName, " with image: ", c.imageURL)

	// Pull image
	logrus.Info("Pulling docker image")
	out, err := c.dockerClient.ImagePull(ctx, c.imageURL, c.imagePullOptions)
	if err != nil {
		logrus.Error("Pull image error: ", err)
	}
	defer out.Close()

	// Create container
	resp, err := c.dockerClient.ContainerCreate(
		ctx,
		c.getContainerConfig(c.imageURL),
		c.getHostConfig(),
		nil,
		nil,
		c.containerName,
	)
	if err != nil {
		logrus.Warn("Create container error: ", err)
	}

	// Start container
	if err := c.dockerClient.ContainerStart(ctx, c.containerName, types.ContainerStartOptions{}); err != nil {
		logrus.Error(err)
		return err
	}

	logrus.Info("Create container ", c.containerName, " success with id: ", resp.ID)

	c.startupTime = time.Now()
	c.containerID = util.ToPointer(resp.ID)
	return nil
}

func (c *containerSvc) Stop(ctx context.Context) error {
	if c.containerID == nil {
		return ErrContainerIsNotStarted
	}

	return c.dockerClient.ContainerStop(ctx, *c.containerID, container.StopOptions{})
}

func (c *containerSvc) RemoveContainer(ctx context.Context) error {
	if c.containerID == nil {
		return ErrContainerIsNotStarted
	}

	responseCh, errCh := c.dockerClient.ContainerWait(ctx, *c.containerID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			logrus.Error(err)
			return err
		}
	case response := <-responseCh:
		if response.Error != nil {
			logrus.Error(response.Error)
			return errors.New(response.Error.Message)
		}
		err := c.dockerClient.ContainerRemove(ctx, *c.containerID, types.ContainerRemoveOptions{})
		if err != nil {
			logrus.Error(err)
			return err
		}
	}
	return nil
}

func (c *containerSvc) getHostConfig() *container.HostConfig {
	return &container.HostConfig{
		ExtraHosts: []string{"host.docker.internal:host-gateway"},
	}
}

func (c *containerSvc) getContainerConfig(dockerImage string) *container.Config {
	return &container.Config{
		Image: dockerImage,
		Env: []string{
			"INITIAL_TASK_RUNNER=3",
			"IMAGE_URL=" + dockerImage,
		},
		// For testing
		//Entrypoint: []string{"/bin/sh", "-c", "sleep infinity"},
	}
}

func (c *containerSvc) GetContainerID() (string, bool) {
	if c.containerID == nil {
		return "", false
	}

	return *c.containerID, true
}

func (c *containerSvc) GetContainerName() string {
	return c.containerName
}
