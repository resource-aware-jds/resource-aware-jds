package container

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/util"
	"github.com/sirupsen/logrus"
	"io"
	"math/rand"
	"os"
	"path/filepath"
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
	GetImageUrl() string
	ExportLog(ctx context.Context) error
}

func ProvideContainer(dockerClient *client.Client, imageURL string, imagePullOptions types.ImagePullOptions) IContainer {
	randomId := rand.Intn(50000-10000) + 10000
	containerName := "rajds-" + strconv.Itoa(randomId)

	return &containerSvc{
		dockerClient:     dockerClient,
		imageURL:         imageURL,
		containerName:    containerName,
		imagePullOptions: imagePullOptions,
	}
}

func (c *containerSvc) Start(ctx context.Context) error {
	logrus.Info("Creating container: ", c.containerName, " with image: ", c.imageURL)

	// Check if image is already exists in the local machine
	_, _, err := c.dockerClient.ImageInspectWithRaw(ctx, c.imageURL)
	if err != nil {
		logrus.Info("Pulling docker image")
		_, err = c.dockerClient.ImagePull(ctx, c.imageURL, c.imagePullOptions)
		if err != nil {
			logrus.Error("Pull image error: ", err)
			return err
		}
	} else {
		logrus.Debug("Using the cached docker image")
	}

	time.Sleep(1 * time.Second)
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
	return &container.HostConfig{}
}

func (c *containerSvc) getContainerConfig(dockerImage string) *container.Config {
	return &container.Config{
		Image: dockerImage,
		Env: []string{
			"INITIAL_TASK_RUNNER=1",
			"IMAGE_URL=" + dockerImage,
			"CONTAINER_GRPC_LISTENING_URL=0.0.0.0:30000",
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

func (c *containerSvc) GetImageUrl() string { return c.imageURL }

func (c *containerSvc) ExportLog(ctx context.Context) error {
	containerId, ok := c.GetContainerID()
	if !ok {
		return fmt.Errorf("[Export log failed] Unable to get container id")
	}

	// Set options for log output
	options := types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Follow: false, Details: true}

	// Get the container logs
	out, err := c.dockerClient.ContainerLogs(ctx, containerId, options)
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
	defer out.Close()

	path := filepath.Join("logs", containerId+"-container-logs.txt")

	// Create or open the file where the logs will be stored
	err = os.MkdirAll("/"+path, os.ModeDir)
	if err != nil {
		panic(err)
	}
	logFile, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	// Copy logs from out to logFile
	_, err = io.Copy(logFile, out)
	if err != nil {
		panic(err)
	}

	return nil
}
