package service

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/resource-aware-jds/resource-aware-jds/config"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/cert"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/datastructure"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/taskqueue"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Worker struct {
	controlPlaneGRPCClient proto.ControlPlaneClient
	dockerClient           *client.Client

	workerNodeCertificate cert.TransportCertificate
	config                config.WorkerConfigModel

	taskQueue  taskqueue.Queue
	taskBuffer datastructure.Buffer[string, models.Task]
}

type IWorker interface {
	CheckInWorkerNodeToControlPlane(ctx context.Context) error
	SubmitTask(containerImage string, taskId string, input []byte) error
	GetTask(containerImage string) (*proto.Task, error)
	SubmitSuccessTask(id string, results [][]byte) error
	ReportFailTask(id string, errorMessage string) error
	CreateContainer(image string, imagePullOptions types.ImagePullOptions) ContainerSvc
}

func ProvideWorker(controlPlaneGRPCClient proto.ControlPlaneClient, dockerClient *client.Client, workerNodeCertificate cert.TransportCertificate, config config.WorkerConfigModel, taskQueue taskqueue.Queue) IWorker {
	return &Worker{
		controlPlaneGRPCClient: controlPlaneGRPCClient,
		dockerClient:           dockerClient,
		config:                 config,
		taskQueue:              taskQueue,
		workerNodeCertificate:  workerNodeCertificate,
		taskBuffer:             make(datastructure.Buffer[string, models.Task]),
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

func (w *Worker) CreateContainer(image string, imagePullOptions types.ImagePullOptions) ContainerSvc {
	return ProvideContainer(w.dockerClient, image, imagePullOptions)
}
