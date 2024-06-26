package daemon

import (
	"context"
	"fmt"
	"github.com/docker/docker/client"
	"github.com/resource-aware-jds/resource-aware-jds/handlerservice"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/metrics"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/timeutil"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/util"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/workerlogic"
	"github.com/resource-aware-jds/resource-aware-jds/service"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/metric"
	"time"
)

const (
	StartContainerDuration  = 15 * time.Second
	ResourceMonitorDuration = 1 * time.Second
)

type workerNode struct {
	ctx                    context.Context
	cancelFunc             func()
	dockerClient           *client.Client
	workerService          handlerservice.IWorker
	resourceMonitor        service.IResourceMonitor
	dynamicScaling         service.IDynamicScaling
	containerTakeDownLogic workerlogic.ContainerTakeDown
	containerService       service.IContainer
	exceedValue            *models.CheckResourceReport
}

type WorkerNode interface {
	Start()
}

func ProvideWorkerNodeDaemon(
	dockerClient *client.Client,
	workerService handlerservice.IWorker,
	resourceMonitor service.IResourceMonitor,
	dynamicScaling service.IDynamicScaling,
	containerTakeDownLogic workerlogic.ContainerTakeDown,
	containerService service.IContainer,
	meter metric.Meter,
) WorkerNode {
	ctx := context.Background()
	ctxWithCancel, cancelFunc := context.WithCancel(ctx)

	res := &workerNode{
		dockerClient:           dockerClient,
		ctx:                    ctxWithCancel,
		cancelFunc:             cancelFunc,
		workerService:          workerService,
		resourceMonitor:        resourceMonitor,
		dynamicScaling:         dynamicScaling,
		containerTakeDownLogic: containerTakeDownLogic,
		containerService:       containerService,
	}

	meter.Float64ObservableCounter( //nolint:errcheck
		metrics.GenerateWorkerNodeMetric("cpu_exceed"),
		metric.WithFloat64Callback(func(ctx context.Context, observer metric.Float64Observer) error {
			if res.exceedValue == nil {
				observer.Observe(0)
				return nil
			}
			observer.Observe((*res.exceedValue).CpuUsageExceed)
			return nil
		}),
	)

	meter.Float64ObservableCounter( //nolint:errcheck
		metrics.GenerateWorkerNodeMetric("memory_exceed"),
		metric.WithFloat64Callback(func(ctx context.Context, observer metric.Float64Observer) error {
			if res.exceedValue == nil {
				observer.Observe(0)
				return nil
			}

			result := util.ConvertToMib((*res.exceedValue).MemoryUsageExceed)
			observer.Observe(result.Size)
			return nil
		}),
	)
	return res
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
				report, err := w.dynamicScaling.CheckResourceUsageLimitWithTimeBuffer(ctx)
				if err != nil {
					logrus.Error(err)
					continue
				}
				err = w.workerService.CalculateAverageContainerResourceUsage(report.ContainerResourceUsages)
				if err != nil {
					logrus.Error(err)
					continue
				}

				w.exceedValue = report

				if report.CpuUsageExceed == 0 && report.MemoryUsageExceed.Size == 0 {
					continue
				}
				logrus.Warnf("CPU Usage or Memory Usage exceeded the limit, taking down the container / CPU Exceed size: %f / MemoryUsageExceed: %f %s", report.CpuUsageExceed, report.MemoryUsageExceed.Size, report.MemoryUsageExceed.Unit)
				containerToBeTakeDowns := w.containerTakeDownLogic.Calculate(workerlogic.ContainerTakeDownState{
					ContainerBuffer: w.containerService.GetContainerBuffer(),
					Report:          report,
				})

				for _, containerToBeTakeDown := range containerToBeTakeDowns {
					err = w.containerService.DownContainer(ctx, containerToBeTakeDown)
					if err != nil {
						logrus.Error("Take Down Container error: ", err)
					}
				}
				timeutil.SleepWithContext(ctx, ResourceMonitorDuration)
			}
		}
	}(w.ctx)
}
