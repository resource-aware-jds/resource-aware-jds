package daemon

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/resource-aware-jds/resource-aware-jds/config"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/buffer"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/cert"
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
	ctx                    context.Context
	cancelFunc             func()
	controlPlaneGRPCClient proto.ControlPlaneClient
	workerNodeCertificate  cert.TransportCertificate
	workerService          service.IWorker
	taskQueue              taskqueue.Queue
	workerNodeConfig       config.WorkerConfigModel
	resourceMonitor        service.IResourceMonitor
	containerBuffer        buffer.ContainerBuffer
}

type WorkerNode interface {
	Start()
}

func ProvideWorkerNodeDaemon(controlPlaneGRPCClient proto.ControlPlaneClient, workerService service.IWorker, taskQueue taskqueue.Queue, workerNodeCertificate cert.TransportCertificate, workerNodeConfig config.WorkerConfigModel, resourceMonitor service.IResourceMonitor, containerBuffer buffer.ContainerBuffer) WorkerNode {
	ctx := context.Background()
	ctxWithCancel, cancelFunc := context.WithCancel(ctx)
	return &workerNode{
		ctx:                    ctxWithCancel,
		cancelFunc:             cancelFunc,
		controlPlaneGRPCClient: controlPlaneGRPCClient,
		workerNodeCertificate:  workerNodeCertificate,
		workerService:          workerService,
		taskQueue:              taskQueue,
		workerNodeConfig:       workerNodeConfig,
		resourceMonitor:        resourceMonitor,
		containerBuffer:        containerBuffer,
	}
}

func (w *workerNode) Start() {
	err := w.checkInNodeToControlPlane()
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

func (w *workerNode) checkInNodeToControlPlane() error {
	certificate, err := w.workerNodeCertificate.GetCertificateInPEM()
	if err != nil {
		return err
	}

	_, err = w.controlPlaneGRPCClient.WorkerCheckIn(w.ctx, &proto.WorkerCheckInRequest{
		Certificate: certificate,
		Port:        int32(w.workerNodeConfig.GRPCServerPort),
	})
	return err
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
