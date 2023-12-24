package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/resource-aware-jds/resource-aware-jds/config"
	"github.com/sirupsen/logrus"
	"strconv"
)

type Worker struct {
	dockerClient *client.Client
	config       config.WorkerConfigModel
}

type IWorker interface {
	RunJob(dockerImage string, name string, options types.ImagePullOptions, jobIdStr string) error
	RemoveContainer(containerID string) error
}

func ProvideWorker(dockerClient *client.Client, config config.WorkerConfigModel) IWorker {
	return &Worker{
		dockerClient: dockerClient,
		config:       config,
	}
}

func (w *Worker) RunJob(dockerImage string, name string, options types.ImagePullOptions, jobIdStr string) error {
	logrus.Info("Creating container: ", name, " with image: ", dockerImage)
	defer logrus.Info("Create container ", name, " success")

	ctx := context.Background()

	//Pull image
	out, err := w.dockerClient.ImagePull(ctx, dockerImage, options)
	if err != nil {
		logrus.Warn("Pull image error: ", err)
	} else {
		defer out.Close()
	}
	// Create container
	resp, err := w.dockerClient.ContainerCreate(ctx, w.getContainerConfig(dockerImage, strconv.Itoa(w.config.GRPCServerPort), jobIdStr), w.getHostConfig(), nil, nil, name)
	if err != nil {
		logrus.Error("Create container error: ", err)
		return err
	}

	// Start container
	if err := w.dockerClient.ContainerStart(ctx, name, types.ContainerStartOptions{}); err != nil {
		logrus.Error(err)
		return err
	}

	fmt.Println(resp.ID)
	return nil
}

func (w *Worker) getHostConfig() *container.HostConfig {
	return &container.HostConfig{
		ExtraHosts: []string{"host.docker.internal:host-gateway"},
	}
}

func (w *Worker) getContainerConfig(dockerImage string, hostPort string, jobID string) *container.Config {
	return &container.Config{
		Image: dockerImage,
		Env:   []string{"HOST_PORT=" + hostPort, "JOB_ID=" + jobID},
	}
}

func (w *Worker) RemoveContainer(containerID string) error {
	ctx := context.Background()
	responseCh, errCh := w.dockerClient.ContainerWait(ctx, containerID, container.WaitConditionNotRunning)
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
		err := w.dockerClient.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{})
		if err != nil {
			logrus.Error(err)
			return err
		}
	}
	return nil
}
