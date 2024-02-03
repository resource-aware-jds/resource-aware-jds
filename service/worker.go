package service

import (
	"context"
	"fmt"
	"github.com/docker/docker/client"
	"github.com/resource-aware-jds/resource-aware-jds/config"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/cert"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/container"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/datastructure"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/taskqueue"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/workerdistribution"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	ContainerIdShortSize = 12
)

type Worker struct {
	controlPlaneGRPCClient proto.ControlPlaneClient
	dockerClient           *client.Client

	workerNodeCertificate cert.TransportCertificate
	config                config.WorkerConfigModel

	workerNodeDistribution workerdistribution.WorkerDistributor
	containerService       ContainerService

	taskQueue              taskqueue.Queue
	taskBuffer             datastructure.Buffer[string, models.Task]
	containerBuffer        datastructure.Buffer[string, container.IContainer]
	containerCoolDownState datastructure.Buffer[string, time.Time]
}

type IWorker interface {
	// ControlPlane related method
	CheckInWorkerNodeToControlPlane(ctx context.Context) error

	// Task related method
	SubmitTask(containerImage string, taskId string, input []byte) error
	GetTask(containerImage string) (*proto.Task, error)
	SubmitSuccessTask(id string, results [][]byte) error
	ReportFailTask(id string, errorMessage string) error

	// TaskDistributionDaemonLoop is a method allowing the daemon to call to accomplish its routine.
	TaskDistributionDaemonLoop(ctx context.Context)

	// Container management related method
	GetContainerIdShort() []string
}

func ProvideWorker(
	controlPlaneGRPCClient proto.ControlPlaneClient,
	dockerClient *client.Client,
	workerNodeCertificate cert.TransportCertificate,
	config config.WorkerConfigModel,
	taskQueue taskqueue.Queue,
	workerNodeDistribution workerdistribution.WorkerDistributor,
	containerService ContainerService,
) IWorker {
	return &Worker{
		controlPlaneGRPCClient: controlPlaneGRPCClient,
		dockerClient:           dockerClient,
		config:                 config,
		taskQueue:              taskQueue,
		workerNodeCertificate:  workerNodeCertificate,
		taskBuffer:             make(datastructure.Buffer[string, models.Task]),
		containerBuffer:        make(datastructure.Buffer[string, container.IContainer]),
		containerCoolDownState: make(datastructure.Buffer[string, time.Time]),
		workerNodeDistribution: workerNodeDistribution,
		containerService:       containerService,
	}
}

func (w *Worker) CheckInWorkerNodeToControlPlane(ctx context.Context) error {
	certificate, err := w.workerNodeCertificate.GetCertificateInPEM()
	if err != nil {
		return err
	}

	_, err = w.controlPlaneGRPCClient.WorkerCheckIn(ctx, &proto.WorkerCheckInRequest{
		Certificate: certificate,
		Port:        int32(w.config.GRPCServerPort),
	})
	return err
}

func (w *Worker) GetTask(containerImage string) (*proto.Task, error) {
	task, err := w.taskQueue.GetTask(containerImage)
	if err != nil {
		logrus.Warn(err)
		return nil, err
	}

	w.taskBuffer.Store(task.ID.Hex(), *task)
	return &proto.Task{
		ID:             task.ID.Hex(),
		TaskAttributes: task.TaskAttributes,
	}, nil
}

func (w *Worker) GetContainerIdShort() []string {
	containerIds := w.containerBuffer.GetKeys()
	return datastructure.Map(containerIds, func(id string) string {
		if len(id) <= ContainerIdShortSize {
			return id
		}
		return id[:ContainerIdShortSize]
	})
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

func (w *Worker) TaskDistributionDaemonLoop(ctx context.Context) {
	logrus.Info("run start container")
	task, ok := w.taskQueue.Pop()
	if !ok {
		return
	}
	taskDepointer := *task

	// Store the ContainerCoolDownState
	distributionResult := w.workerNodeDistribution.Distribute(ctx, taskDepointer, workerdistribution.WorkerState{
		ContainerCoolDownState: w.containerCoolDownState,
		ContainerList:          w.containerBuffer,
		WorkerNodeConfig:       w.config,
	})

	if !distributionResult.CreateContainerToSupportTask {
		return
	}

	// Remove ContainerCoolDownState
	delete(w.containerCoolDownState, taskDepointer.ImageUrl)

	logrus.Info("Starting container with image:", taskDepointer.ImageUrl)
	_, err := w.containerService.StartContainer(ctx, taskDepointer.ImageUrl)
	if err != nil {
		logrus.Error("Unable to start container with image:", taskDepointer.ImageUrl, err)
		errorTaskList := datastructure.Filter(w.taskQueue.ReadQueue(), func(task *models.Task) bool {
			return task.ImageUrl == taskDepointer.ImageUrl
		})
		logrus.Warn("Removing these task due to unable to start container", errorTaskList)
		w.taskQueue.BulkRemove(errorTaskList)
		return
	}

	// Add the cool down state
	w.containerCoolDownState[taskDepointer.ImageUrl] = time.Now().Add(w.config.ContainerStartDelayTimeSeconds)
}
