package daemon

import (
	"context"
	"errors"
	"github.com/resource-aware-jds/resource-aware-jds/config"
	"github.com/resource-aware-jds/resource-aware-jds/handlerservice"
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
	controlPlaneService handlerservice.IControlPlane
	cpTaskWatcher       service.CPTaskWatcher
	taskService         service.Task
	jobService          service.Job
	config              config.ControlPlaneConfigModel
}

type IControlPlane interface {
	Start()
	GracefullyShutdown()
	CheckTheDistributedTask()
}

func ProvideControlPlaneDaemon(
	workerNodePool pool.WorkerNode,
	controlPlaneService handlerservice.IControlPlane,
	taskService service.Task,
	jobService service.Job,
	cpTaskWatcher service.CPTaskWatcher,
	config config.ControlPlaneConfigModel,
) (IControlPlane, func()) {
	ctx := context.Background()
	ctxWithCancel, cancelFunc := context.WithCancel(ctx)

	cp := controlPlane{
		ctx:                 ctxWithCancel,
		cancelFunc:          cancelFunc,
		workerNodePool:      workerNodePool,
		controlPlaneService: controlPlaneService,
		taskService:         taskService,
		jobService:          jobService,
		cpTaskWatcher:       cpTaskWatcher,
		config:              config,
	}

	cleanup := func() {
		cp.GracefullyShutdown()
	}

	return &cp, cleanup
}

func (c *controlPlane) Start() {
	logrus.Info("[WorkerNode Pool] Get all available worker node from registry")
	nodes, err := c.controlPlaneService.GetAllWorkerNodeFromRegistry(c.ctx)
	if err != nil {
		logrus.Warnf("[WorkerNode Pool] Failed to get all available worker node from registry with error (%s)", err.Error())
	}

	logrus.Info("[ControlPlane Daemon] Starting the CP Daemon loop")
	c.workerNodePool.InitializePool(c.ctx, nodes)

	c.CheckTheDistributedTask()

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

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				c.cpTaskWatcher.WatcherLoop(ctx)
				timeutil.SleepWithContext(ctx, c.config.TaskWatcherConfig.SleepTime)
			}
		}
	}(c.ctx)
}

func (c *controlPlane) CheckTheDistributedTask() {
	logrus.Info("[ControlPlane Daemon] Recovering the Distributed tasks")
	runningTaskIDs := c.workerNodePool.CheckRunningTaskInEachWorkerNode(c.ctx)
	distributedTasks, err := c.taskService.GetAllDistributedTask(c.ctx)
	if err != nil {
		logrus.Error("Failed to Get DistributedTasks", err)
		return
	}

	recoveredTask := 0
	workOnFailureTask := 0
	for _, distributedTask := range distributedTasks {
		if ok := runningTaskIDs[*distributedTask.ID]; ok {
			recoveredTask++
			c.cpTaskWatcher.AddTaskToWatch(*distributedTask.ID)
			continue
		}

		// Update task to make it possible to be distributed again
		err = c.taskService.UpdateTaskWorkOnFailure(c.ctx, *distributedTask.ID, "", "Control Plane startup process detect no worker in the pool running this task")
		if err != nil {
			logrus.Error("Failed to UpdateTaskTo WorkOnFailure", err)
		}
		workOnFailureTask++
	}
	logrus.Infof("[ControlPlane Daemon] Total Recovered task: %d / Total WorkOnFailure task: %d", recoveredTask, workOnFailureTask)
}

func (c *controlPlane) taskScanLoop(ctx context.Context) {
	// Check if there is any available worker node first
	if !c.workerNodePool.IsAvailableWorkerNode() {
		logrus.Warn("[ControlPlane Daemon] No available worker node in the pool, skipping task scan loop")
		return
	}

	jobList, err := c.jobService.ListJobReadyToDistribute(ctx)
	if err != nil {
		logrus.Errorf("[ControlPlane Daemon] Failed to get available job (%s)", err.Error())
		return
	}
	if len(jobList) == 0 {
		logrus.Error("[ControlPlane Daemon] No job available to be distributed, skipping")
		return
	}

	job, tasks, err := c.taskService.GetAvailableTask(ctx, jobList)
	if err != nil {
		logrus.Errorf("[ControlPlane Daemon] Failed to get available task (%s)", err.Error())
		return
	}

	if len(tasks) == 0 {
		logrus.Warn("[ControlPlane Daemon] No task available in database, skipping")
		return
	}

	// Call Distribute function
	successTask, failureTask, err := c.workerNodePool.DistributeWork(ctx, *job, tasks)
	if err != nil {
		logrus.Warnf("[ControlPlane Daemon] Failed to distribute work to any worker nodes in the pool (%s)", err.Error())
		if errors.Is(err, pool.ErrNoAvailableDistributor) {
			err := c.jobService.UpdateJobToFailed(ctx, job.ID, "No distribution solution", err)
			if err != nil {
				logrus.Warnf("[ControlPlane Daemon] Fail to update job status (%s)", err.Error())
			}

			// Update all the task to be failed
			err = c.taskService.UpdateAllTaskToWorkOnFailure(ctx, job, "No distribution solution")
			if err != nil {
				logrus.Warnf("[ControlPlane Daemon] Fail to update tasks status (%s)", err.Error())
			}
		}
		return
	}

	for _, eachSuccessTask := range successTask {
		c.cpTaskWatcher.AddTaskToWatch(*eachSuccessTask.ID)
	}

	err = c.taskService.UpdateTaskAfterDistribution(ctx, successTask, failureTask)
	if err != nil {
		logrus.Warnf("[ControlPlane Daemon] Failed to update task status (%s)", err.Error())
	}
}

func (c *controlPlane) GracefullyShutdown() {
	logrus.Info("[ControlPlane Daemon] Received gracefully shutdown command")
	c.cancelFunc()
	logrus.Info("[ControlPlane Daemon] Gracefully Shutdown success.")
}
