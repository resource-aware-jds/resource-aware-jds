package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/resource-aware-jds/resource-aware-jds/config"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/taskqueue"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

type Worker struct {
	dockerClient *client.Client
	config       config.WorkerConfigModel
	taskQueue    taskqueue.Queue
}

type IWorker interface {
	RemoveContainer(containerID string) error
	SubmitTask(containerImage string, taskId string, input []byte) error
}

func ProvideWorker(dockerClient *client.Client, config config.WorkerConfigModel, taskQueue taskqueue.Queue) IWorker {
	return &Worker{
		dockerClient: dockerClient,
		config:       config,
		taskQueue:    taskQueue,
	}
}

func (w *Worker) SubmitTask(containerImage string, taskId string, input []byte) error {
	if !w.isContainerExist(containerImage) {
		containerName := "rajds-" + taskId
		err := w.startContainer(containerImage, containerName, types.ImagePullOptions{}, taskId)
		if err != nil {
			return err
		}
	}

	task := models.Task{
		ImageUrl: containerImage,
		TaskId:   taskId,
		Input:    input,
	}
	w.taskQueue.StoreTask(&task)
	return nil
}

func (w *Worker) startContainer(dockerImage string, name string, options types.ImagePullOptions, taskId string) error {
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

	// Get mount directory path
	mountDirParh := w.getBindPath(taskId)

	// Create container
	resp, err := w.dockerClient.ContainerCreate(ctx, w.getContainerConfig(dockerImage, taskId), w.getHostConfig(mountDirParh), nil, nil, name)
	if err != nil {
		logrus.Warn("Create container error: ", err)
	}

	// Start container
	if err := w.dockerClient.ContainerStart(ctx, name, types.ContainerStartOptions{}); err != nil {
		logrus.Error(err)
		return err
	}

	fmt.Println(resp.ID)
	return nil
}

func (w *Worker) getHostConfig(sourceMountPath string) *container.HostConfig {
	return &container.HostConfig{
		ExtraHosts: []string{"host.docker.internal:host-gateway"},
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: sourceMountPath,
				Target: "/tmp",
			},
		},
	}
}

func (w *Worker) getBindPath(taskId string) string {
	hostPath := "/tmp"
	path := filepath.Join(hostPath, "rajds", taskId)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		logrus.Error("Unable to create folder on", path, err)
	}
	return path
}

func (w *Worker) isContainerExist(imageUrl string) bool {
	ctx := context.Background()
	opt := types.ContainerListOptions{All: true}
	opt.Filters = filters.NewArgs()
	opt.Filters.Add("status", "running")

	containers, err := w.dockerClient.ContainerList(ctx, opt)
	logrus.Info("Container list:", containers)
	if err != nil {
		logrus.Error("Failed to retrieve containers: ", err)
	}

	for _, container := range containers {
		if container.Image == imageUrl {
			logrus.Info("Container with image: ", imageUrl)
			return true
		}
	}
	return false
}

func (w *Worker) getContainerConfig(dockerImage string, taskId string) *container.Config {
	return &container.Config{
		Image:      dockerImage,
		Env:        []string{"MAXIMUM_CONCURRENT=" + "3", "TASK_ID=" + taskId, "WORKER_NODE_GRPC_SERVER_UNIX_SOCKET_PATH=/tmp"},
		Entrypoint: []string{"/bin/sh", "-c", "sleep infinity"},
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
