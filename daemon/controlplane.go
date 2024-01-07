package daemon

import (
	"context"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/pool"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/timeutil"
	"github.com/resource-aware-jds/resource-aware-jds/service"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	AvailabilityCheckSleepDuration = 5 * time.Second
)

type controlPlane struct {
	ctx                 context.Context
	cancelFunc          func()
	workerNodePool      pool.WorkerNode
	controlPlaneService service.IControlPlane
}

type IControlPlane interface {
	Start()
	GracefullyShutdown()
}

func ProvideControlPlaneDaemon(workerNodePool pool.WorkerNode, controlPlaneService service.IControlPlane) (IControlPlane, func()) {
	ctx := context.Background()
	ctxWithCancel, cancelFunc := context.WithCancel(ctx)

	cp := controlPlane{
		ctx:                 ctxWithCancel,
		cancelFunc:          cancelFunc,
		workerNodePool:      workerNodePool,
		controlPlaneService: controlPlaneService,
	}

	cleanup := func() {
		cp.GracefullyShutdown()
	}

	return &cp, cleanup
}

func (c *controlPlane) Start() {
	logrus.Info("[ControlPlane Daemon] Starting the CP Daemon loop")
	c.workerNodePool.InitializePool(c.ctx)

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				c.taskScanLoop(ctx)
				timeutil.SleepWithContext(ctx, 5*time.Second)
			}
		}
	}(c.ctx)

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				c.workerNodePool.WorkerNodeAvailabilityCheck(ctx)
				timeutil.SleepWithContext(ctx, AvailabilityCheckSleepDuration)
			}
		}
	}(c.ctx)
}

func (c *controlPlane) taskScanLoop(ctx context.Context) {
	// Check if there is any available worker node first
	if !c.workerNodePool.IsAvailableWorkerNode() {
		logrus.Warn("[ControlPlane Daemon] No available worker node in the pool, skipping task scan loop")
		return
	}

	tasks, err := c.controlPlaneService.GetAvailableTask(ctx)
	if err != nil {
		logrus.Errorf("[ControlPlane Daemon] Failed to get available task (%s)", err.Error())
		return
	}

	if len(tasks) == 0 {
		logrus.Warn("[ControlPlane Daemon] No task available in database, skipping")
		return
	}

	// Call Distribute function
	successTask, failureTask, err := c.workerNodePool.DistributeWork(ctx, tasks)
	if err != nil {
		logrus.Warnf("[ControlPlane Daemon] Failed to distribute work to any worker nodes in the pool (%s)", err.Error())
		return
	}

	err = c.controlPlaneService.UpdateTaskAfterDistribution(ctx, successTask, failureTask)
	if err != nil {
		logrus.Warnf("[ControlPlane Daemon] Failed to update task status (%s)", err.Error())
	}
}

func (c *controlPlane) GracefullyShutdown() {
	logrus.Info("[ControlPlane Daemon] Received gracefully shutdown command")
	c.cancelFunc()
	logrus.Info("[ControlPlane Daemon] Gracefully Shutdown success.")
}
