package daemon

import (
	"context"
	"fmt"
	"github.com/docker/docker/client"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/timeutil"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/workerlogic"
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
	dockerClient           *client.Client
	workerService          service.IWorker
	resourceMonitor        service.IResourceMonitor
	dynamicScaling         service.IDynamicScaling
	containerTakeDownLogic workerlogic.ContainerTakeDown
	containerService       service.IContainer
}

type WorkerNode interface {
	Start()
}

func ProvideWorkerNodeDaemon(
	dockerClient *client.Client,
	workerService service.IWorker,
	resourceMonitor service.IResourceMonitor,
	dynamicScaling service.IDynamicScaling,
	containerTakeDownLogic workerlogic.ContainerTakeDown,
	containerService service.IContainer,
) WorkerNode {
	ctx := context.Background()
	ctxWithCancel, cancelFunc := context.WithCancel(ctx)
	return &workerNode{
		dockerClient:           dockerClient,
		ctx:                    ctxWithCancel,
		cancelFunc:             cancelFunc,
		workerService:          workerService,
		resourceMonitor:        resourceMonitor,
		dynamicScaling:         dynamicScaling,
		containerTakeDownLogic: containerTakeDownLogic,
		containerService:       containerService,
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
				report, err := w.dynamicScaling.CheckResourceUsageLimit(ctx)
				if err != nil {
					logrus.Error(err)
					continue
				}

				if report.CpuUsageExceed == 0 && report.MemoryUsageExceed.Size == 0 {
					continue
				}
				logrus.Warn("CPU Usage or Memory Usage exceeded the limit, taking down the container")
				w.containerTakeDownLogic.Calculate(workerlogic.ContainerTakeDownState{
					ContainerBuffer: w.containerService.GetContainerBuffer(),
					Report:          report,
				})
				timeutil.SleepWithContext(ctx, ResourceMonitorDuration)
			}
		}
	}(w.ctx)
}
