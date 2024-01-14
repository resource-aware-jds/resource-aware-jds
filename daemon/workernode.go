package daemon

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/resource-aware-jds/resource-aware-jds/config"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/datastructure"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/taskqueue"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/timeutil"
	"github.com/resource-aware-jds/resource-aware-jds/service"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	StartContainerDuration  = 15 * time.Second
	ResourceMonitorDuration = 5 * time.Second
)

type workerNode struct {
	ctx        context.Context
	cancelFunc func()

	dockerClient *client.Client

	workerService    service.IWorker
	taskQueue        taskqueue.Queue
	workerNodeConfig config.WorkerConfigModel
	resourceMonitor  service.IResourceMonitor
	containerBuffer  datastructure.Buffer[string, service.ContainerSvc]
}

type WorkerNode interface {
	Start()
}

func ProvideWorkerNodeDaemon(dockerClient *client.Client, workerService service.IWorker, taskQueue taskqueue.Queue, workerNodeConfig config.WorkerConfigModel, resourceMonitor service.IResourceMonitor) WorkerNode {
	ctx := context.Background()
	ctxWithCancel, cancelFunc := context.WithCancel(ctx)
	return &workerNode{
		dockerClient:     dockerClient,
		ctx:              ctxWithCancel,
		cancelFunc:       cancelFunc,
		workerService:    workerService,
		taskQueue:        taskQueue,
		workerNodeConfig: workerNodeConfig,
		resourceMonitor:  resourceMonitor,
		containerBuffer:  make(datastructure.Buffer[string, service.ContainerSvc]),
	}
}

func (w *workerNode) Start() {
	err := w.workerService.CheckInWorkerNodeToControlPlane(w.ctx)
	if err != nil {
		panic(fmt.Sprintf("Failed to check in worker node to control plane (%s)", err.Error()))
	}

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				w.taskStartContainer(ctx)
				timeutil.SleepWithContext(ctx, StartContainerDuration)
			}
		}
	}(w.ctx)

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				w.resourceMonitor.GetMemoryUsage()
				timeutil.SleepWithContext(ctx, ResourceMonitorDuration)
			}
		}
	}(w.ctx)
}

func (w *workerNode) taskStartContainer(ctx context.Context) {
	logrus.Info("run start container")
	imageList := w.taskQueue.GetDistinctImageList()
	taskList := w.taskQueue.ReadQueue()

	for _, image := range imageList {
		logrus.Info("Starting container with image:", image)
		container := w.workerService.CreateContainer(image, types.ImagePullOptions{})
		err := container.Start(ctx)
		if err != nil {
			logrus.Error("Unable to start container with image:", image, err)
			errorTaskList := datastructure.Filter(taskList, func(task *models.Task) bool {
				return task.ImageUrl == image
			})
			logrus.Warn("Removing these task due to unable to start container", errorTaskList)
			w.taskQueue.BulkRemove(errorTaskList)
			return
		}
		containerID, ok := container.GetContainerID()
		if !ok {
			logrus.Error("Unable to get container id")
			return
		}
		w.containerBuffer.Store(containerID, container)
	}
}
