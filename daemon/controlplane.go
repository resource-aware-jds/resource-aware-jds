package daemon

import (
	"context"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/pool"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	MaximumUnavailableCount        = 3
	AvailabilityCheckSleepDuration = 5 * time.Second
)

type availableWorkerNodeGRPCMapper struct {
	nodeEntry        models.NodeEntry
	grpcConnection   proto.WorkerNodeClient
	unavailableCount uint
}

type controlPlane struct {
	ctx            context.Context
	cancelFunc     func()
	workerNodePool pool.WorkerNode
}

type IControlPlane interface {
	Start()
	GracefullyShutdown()
}

func ProvideControlPlaneDaemon(workerNodePool pool.WorkerNode) (IControlPlane, func()) {
	ctx := context.Background()
	ctxWithCancel, cancelFunc := context.WithCancel(ctx)

	cp := controlPlane{
		ctx:            ctxWithCancel,
		cancelFunc:     cancelFunc,
		workerNodePool: workerNodePool,
	}

	cleanup := func() {
		cp.GracefullyShutdown()
	}

	return &cp, cleanup
}

func (c *controlPlane) Start() {
	logrus.Info("[ControlPlane Daemon] Starting the CP Daemon loop")
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				c.loop(ctx)
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
			}
		}
	}(c.ctx)
}

func (c *controlPlane) loop(ctx context.Context) {
	// TODO: Get some available task from database

	// TODO: Call Distribute function
}

func (c *controlPlane) GracefullyShutdown() {
	logrus.Info("[ControlPlane Daemon] Received gracefully shutdown command")
	c.cancelFunc()
	logrus.Info("[ControlPlane Daemon] Gracefully Shutdown success.")
}
