package daemon

import (
	"context"
	"fmt"
	"github.com/docker/docker/client"
	"github.com/resource-aware-jds/resource-aware-jds/config"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/timeutil"
	"github.com/resource-aware-jds/resource-aware-jds/service"
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
	workerNodeConfig config.WorkerConfigModel
	resourceMonitor  service.IResourceMonitor
}

type WorkerNode interface {
	Start()
}

func ProvideWorkerNodeDaemon(dockerClient *client.Client, workerService service.IWorker, workerNodeConfig config.WorkerConfigModel, resourceMonitor service.IResourceMonitor) WorkerNode {
	ctx := context.Background()
	ctxWithCancel, cancelFunc := context.WithCancel(ctx)
	return &workerNode{
		dockerClient:     dockerClient,
		ctx:              ctxWithCancel,
		cancelFunc:       cancelFunc,
		workerService:    workerService,
		workerNodeConfig: workerNodeConfig,
		resourceMonitor:  resourceMonitor,
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
				w.workerService.TaskDistributionDaemonLoop(ctx)
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
