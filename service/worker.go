package service

import (
	"context"
	"errors"
	"github.com/docker/docker/client"
	"github.com/resource-aware-jds/resource-aware-jds/config"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/cert"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/datastructure"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/taskqueue"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/workerlogic"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel/metric"
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

	workerNodeDistribution workerlogic.WorkerDistributor
	containerService       IContainer

	taskQueue  taskqueue.Queue
	taskBuffer datastructure.Buffer[string, models.TaskWithContext]
}

type IWorker interface {
	// ControlPlane related method
	CheckInWorkerNodeToControlPlane(ctx context.Context) error

	// Task related method
	StoreTaskInQueue(containerImage string, taskId string, input []byte) error
	GetTask(containerImage string) (*proto.Task, error)
	SubmitSuccessTask(ctx context.Context, id string, results []byte) error
	ReportFailTask(ctx context.Context, id string, errorMessage string) error

	// TaskDistributionDaemonLoop is a method allowing the daemon to call to accomplish its routine.
	TaskDistributionDaemonLoop(ctx context.Context)
}

func ProvideWorker(
	controlPlaneGRPCClient proto.ControlPlaneClient,
	dockerClient *client.Client,
	workerNodeCertificate cert.TransportCertificate,
	config config.WorkerConfigModel,
	taskQueue taskqueue.Queue,
	workerNodeDistribution workerlogic.WorkerDistributor,
	containerService IContainer,
	meter metric.Meter,
) IWorker {
	return &Worker{
		controlPlaneGRPCClient: controlPlaneGRPCClient,
		dockerClient:           dockerClient,
		config:                 config,
		taskQueue:              taskQueue,
		workerNodeCertificate:  workerNodeCertificate,
		taskBuffer: datastructure.ProvideBuffer[string, models.TaskWithContext](
			datastructure.WithBufferMetrics(meter, "task_buffer_size"),
		),
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

	ctx, cancelFunc := context.WithDeadline(context.Background(), time.Now().Add(w.config.TaskBufferTimeout))

	// Start a new goroutine to check if context is timeout before calling cancel.
	// If yes, call report task failure.
	go func(innerCtx context.Context, innerW *Worker, innerTask models.Task) {
		<-ctx.Done()
		if errors.Is(ctx.Err(), context.Canceled) {
			return
		}
		err := innerW.ReportFailTask(context.Background(), innerTask.ID.Hex(), "Timeout")
		if err != nil {
			logrus.Error("report task failure fail: ", err)
		}
	}(ctx, w, *task)

	w.taskBuffer.Store(task.ID.Hex(), models.TaskWithContext{
		Task:       *task,
		Ctx:        ctx,
		CancelFunc: cancelFunc,
	})
	return &proto.Task{
		ID:             task.ID.Hex(),
		TaskAttributes: task.TaskAttributes,
	}, nil
}

func (w *Worker) SubmitSuccessTask(ctx context.Context, id string, results []byte) error {
	task := w.taskBuffer.Pop(id)
	if task == nil {
		logrus.Error("Task is not running")
		return errors.New("task not found in task buffer, maybe it already timeout and get sent back to cp")
	}
	logrus.Info("Task succeed with id: " + id)
	_, err := w.controlPlaneGRPCClient.ReportSuccessTask(ctx, &proto.ReportSuccessTaskRequest{
		Id:     id,
		NodeID: w.workerNodeCertificate.GetNodeID(),
		Result: results,
	})
	return err
}

func (w *Worker) ReportFailTask(ctx context.Context, id string, errorMessage string) error {
	task := w.taskBuffer.Pop(id)
	if task == nil {
		logrus.Error("Task is not running")
		return errors.New("task not found in task buffer, maybe it already timeout and get sent back to cp")
	}

	logrus.Error("Task failed with id: " + id)
	_, err := w.controlPlaneGRPCClient.ReportFailureTask(ctx, &proto.ReportFailureTaskRequest{
		Id:      id,
		NodeID:  w.workerNodeCertificate.GetNodeID(),
		Message: errorMessage,
	})
	return err
}

func (w *Worker) StoreTaskInQueue(containerImage string, taskId string, input []byte) error {
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
	task, ok := w.taskQueue.PeakForNextTask()
	if !ok {
		return
	}
	taskDepointer := *task

	// Store the ContainerCoolDownState
	distributionResult := w.workerNodeDistribution.Distribute(ctx, taskDepointer, workerlogic.WorkerState{
		ContainerCoolDownState: w.containerService.GetContainerCoolDownState(),
		WorkerNodeConfig:       w.config,
		ContainerBuffer:        w.containerService.GetContainerBuffer(),
	})

	if !distributionResult.CreateContainerToSupportTask {
		return
	}

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
}
