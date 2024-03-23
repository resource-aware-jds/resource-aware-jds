package handlerservice

import (
	"context"
	"errors"
	"github.com/docker/docker/client"
	"github.com/resource-aware-jds/resource-aware-jds/config"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/cert"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/datastructure"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/metrics"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/taskqueue"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/util"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/workerlogic"
	"github.com/resource-aware-jds/resource-aware-jds/service"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"strings"
	"time"
)

type Worker struct {
	controlPlaneGRPCClient proto.ControlPlaneClient
	dockerClient           *client.Client

	workerNodeCertificate cert.TransportCertificate
	config                config.WorkerConfigModel

	workerNodeDistribution workerlogic.WorkerDistributor
	containerService       service.IContainer

	taskQueue  taskqueue.Queue
	taskBuffer datastructure.Buffer[string, models.TaskWithContext]

	containerSubmitTask metric.Int64Counter
}

type IWorker interface {
	// ControlPlane related method
	CheckInWorkerNodeToControlPlane(ctx context.Context) error

	// Task related method
	StoreTaskInQueue(containerImage string, taskId string, input []byte) error
	GetTask(containerImage string, containerId string) (*proto.Task, error)
	SubmitSuccessTask(ctx context.Context, id string, results []byte) error
	ReportFailTask(ctx context.Context, id string, errorMessage string) error
	CalculateAverageContainerResourceUsage(usage []models.ContainerResourceUsage) error
	GetRunningTask() []string
	GetAvailableTaskSlot() int
	GetQueuedTask() []string

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
	containerService service.IContainer,
	meter metric.Meter,
) (IWorker, error) {
	containerSubmitTask, err := meter.Int64Counter("container_submit_task")

	return &Worker{
		controlPlaneGRPCClient: controlPlaneGRPCClient,
		dockerClient:           dockerClient,
		config:                 config,
		taskQueue:              taskQueue,
		workerNodeCertificate:  workerNodeCertificate,
		taskBuffer: datastructure.ProvideBuffer[string, models.TaskWithContext](
			datastructure.WithBufferMetrics(
				meter,
				metrics.GenerateWorkerNodeMetric("task_buffer"),
				metric.WithUnit("Task"),
				metric.WithDescription("The total task that currently running in the container"),
			),
		),
		workerNodeDistribution: workerNodeDistribution,
		containerService:       containerService,
		containerSubmitTask:    containerSubmitTask,
	}, err
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

func (w *Worker) GetTask(containerImage string, containerId string) (*proto.Task, error) {
	task, err := w.taskQueue.GetTask(containerImage)
	if err != nil {
		addError := w.containerService.AddContainerTakeDownTimer(containerId)
		if addError != nil {
			logrus.Error(addError)
		}
		logrus.Warn(err)
		return nil, err
	}
	w.containerService.RemoveContainerTakeDownTimer(containerId)

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
		Task:        *task,
		Ctx:         ctx,
		CancelFunc:  cancelFunc,
		ContainerId: containerId,
		AverageResourceUsage: models.AverageResourceUsage{
			IsInitialized:   false,
			AverageCpuUsage: 0,
			AverageMemoryUsage: models.MemorySize{
				Size: 0,
				Unit: "GiB",
			}},
	})
	return &proto.Task{
		ID:             task.ID.Hex(),
		TaskAttributes: task.TaskAttributes,
	}, nil
}

func (w *Worker) GetRunningTask() []string {
	return w.taskBuffer.GetKeys()
}

func (w *Worker) GetQueuedTask() []string {
	allQueuedTask := w.taskQueue.ReadQueue()
	return datastructure.Map(allQueuedTask, func(task *models.Task) string {
		return task.ID.Hex()
	})
}

func (w *Worker) GetAvailableTaskSlot() int {
	return w.config.TotalContainerLimit - len(w.taskBuffer)
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
		TaskResourceUsage: &proto.TaskResourceUsage{
			AverageMemoryUsage: util.MemoryToString(task.AverageResourceUsage.AverageMemoryUsage),
			AverageCpuUsage:    float32(task.AverageResourceUsage.AverageCpuUsage),
		},
	})
	w.containerSubmitTask.Add(ctx, 1, metric.WithAttributes(attribute.String("status", "success")))
	return err
}

func (w *Worker) ReportFailTask(ctx context.Context, id string, errorMessage string) error {
	task := w.taskBuffer.Pop(id)
	if task == nil {
		logrus.Error("Task is not running")
		//return errors.New("task not found in task buffer, maybe it already timeout and get sent back to cp")
	}

	logrus.Error("Task failed with id: " + id)
	w.containerSubmitTask.Add(ctx, 1, metric.WithAttributes(attribute.String("status", "failure")))
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

func (w *Worker) CalculateAverageContainerResourceUsage(usage []models.ContainerResourceUsage) error {
	tasks := w.taskBuffer
	for key, task := range tasks {
		containerIdShort := task.ContainerId
		for _, container := range usage {
			if strings.TrimSpace(container.ContainerIdShort) != strings.TrimSpace(containerIdShort) {
				continue
			}
			cpuUsage, err := util.ExtractCpuUsage(container)
			if err != nil {
				return err
			}
			cpuUsagePerCore := cpuUsage / float64(w.config.DockerCoreLimit)
			memoryUsage := util.ExtractMemoryUsageFromModel(container)
			if !task.AverageResourceUsage.IsInitialized {
				task.AverageResourceUsage = models.AverageResourceUsage{
					AverageCpuUsage:    cpuUsagePerCore,
					AverageMemoryUsage: memoryUsage,
					IsInitialized:      true,
				}
			} else {
				task.AverageResourceUsage.AverageCpuUsage = (task.AverageResourceUsage.AverageCpuUsage + cpuUsagePerCore) / 2
				task.AverageResourceUsage.AverageMemoryUsage = util.DivideBy(
					util.SumInGb(task.AverageResourceUsage.AverageMemoryUsage, memoryUsage),
					2,
				)
			}
			w.taskBuffer[key] = task
		}
	}
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
		for _, removedTask := range errorTaskList {
			err := w.ReportFailTask(ctx, removedTask.ID.Hex(), "Fail to start the container with the provided image")
			if err != nil {
				logrus.Error("Submit the Fail Task fail: ", err)
			}
		}
		return
	}
}
