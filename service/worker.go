package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/resource-aware-jds/resource-aware-jds/config"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/buffer"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/taskqueue"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Worker struct {
	dockerClient *client.Client
	config       config.WorkerConfigModel
	taskQueue    taskqueue.Queue
	taskBuffer   buffer.TaskBuffer
}

type IWorker interface {
	RemoveContainer(ctx context.Context, containerID string) error
	SubmitTask(containerImage string, taskId string, input []byte) error
	GetTask(containerImage string) (*proto.Task, error)
	SubmitSuccessTask(id string, results [][]byte) error
	ReportFailTask(id string, errorMessage string) error
	StartContainer(ctx context.Context, dockerImage string, name string, options types.ImagePullOptions) error
	IsContainerExist(ctx context.Context, imageUrl string) bool
}

func ProvideWorker(dockerClient *client.Client, config config.WorkerConfigModel, taskQueue taskqueue.Queue, taskBuffer buffer.TaskBuffer) IWorker {
	return &Worker{
		dockerClient: dockerClient,
		config:       config,
		taskQueue:    taskQueue,
		taskBuffer:   taskBuffer,
	}
}

func (w *Worker) GetTask(containerImage string) (*proto.Task, error) {
	task, err := w.taskQueue.GetTask(containerImage)
	if err != nil {
		logrus.Warn(err)
		return nil, err
	}

	w.taskBuffer.Store(task.ID.Hex(), task)
	return &proto.Task{
		ID:             task.ID.Hex(),
		TaskAttributes: task.TaskAttributes,
	}, nil
}

func (w *Worker) SubmitSuccessTask(id string, results [][]byte) error {
	task := w.taskBuffer.Pop(id)
	if task == nil {
		logrus.Error("Task is not running")
	}
	logrus.Info("Task succeed with id: " + id)
	return nil
}

func (w *Worker) ReportFailTask(id string, errorMessage string) error {
	task := w.taskBuffer.Pop(id)
	if task == nil {
		return fmt.Errorf("Task with id:" + id + "not existed in buffer")
	}
	logrus.Error("Task failed with id: " + id)
	w.taskQueue.StoreTask(task)
	return nil
}

func (w *Worker) SubmitTask(containerImage string, taskId string, input []byte) error {
	hex, err := primitive.ObjectIDFromHex(taskId)
	if err != nil {
		return err
	}

	task := models.Task{
		ImageUrl:       containerImage,
		ID:             &hex,
		TaskAttributes: input,
	}
	w.taskQueue.StoreTask(&task)
	return nil
}

func (w *Worker) RemoveContainer(ctx context.Context, containerID string) error {
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

func (w *Worker) StartContainer(ctx context.Context, dockerImage string, name string, options types.ImagePullOptions) error {
	logrus.Info("Creating container: ", name, " with image: ", dockerImage)
	defer logrus.Info("Create container ", name, " success")

	//Pull image
	logrus.Info("Pulling docker image")
	out, err := w.dockerClient.ImagePull(ctx, dockerImage, options)
	if err != nil {
		logrus.Error("Pull image error: ", err)
	} else {
		defer out.Close()
	}

	// Create container
	resp, err := w.dockerClient.ContainerCreate(
		ctx,
		w.getContainerConfig(dockerImage),
		w.getHostConfig(),
		nil,
		nil,
		name,
	)
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

func (w *Worker) getHostConfig() *container.HostConfig {
	return &container.HostConfig{
		ExtraHosts: []string{"host.docker.internal:host-gateway"},
	}
}

func (w *Worker) IsContainerExist(ctx context.Context, imageUrl string) bool {
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

func (w *Worker) getContainerConfig(dockerImage string) *container.Config {
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
