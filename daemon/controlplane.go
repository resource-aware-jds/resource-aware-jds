package daemon

import (
	"context"
	"fmt"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/cert"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/timeutil"
	"github.com/resource-aware-jds/resource-aware-jds/repository"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
	"sync"
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
	ctx                         context.Context
	cancelFunc                  func()
	gracefullyShutdownWaitGroup sync.WaitGroup
	availableWorkerNode         map[string]availableWorkerNodeGRPCMapper
	caCertificate               cert.CACertificate
	nodeRegistryRepo            repository.INodeRegistry
}

type IControlPlane interface {
	Start()
	GracefullyShutdown()
}

func ProvideControlPlaneDaemon(caCertificate cert.CACertificate, nodeRegistryRepo repository.INodeRegistry) (IControlPlane, func()) {
	ctx := context.Background()
	ctxWithCancel, cancelFunc := context.WithCancel(ctx)

	cp := controlPlane{
		ctx:                 ctxWithCancel,
		cancelFunc:          cancelFunc,
		availableWorkerNode: map[string]availableWorkerNodeGRPCMapper{},
		caCertificate:       caCertificate,
		nodeRegistryRepo:    nodeRegistryRepo,
	}

	cleanup := func() {
		cp.GracefullyShutdown()
	}

	return &cp, cleanup
}

func (c *controlPlane) Start() {
	logrus.Info("[ControlPlane Daemon] Get all available worker node from registry")
	nodes, err := c.nodeRegistryRepo.GetAllWorkerNode(c.ctx)
	if err != nil {
		logrus.Warnf("[ControlPlane Daemon] Failed to get all available worker node from registry with error (%s)", err.Error())
	}

	for _, node := range nodes {
		err = c.AddAvailableWorkerNode(node)
		if err != nil {
			continue
		}
	}

	logrus.Infof("[ControlPlane Daemon] Added %d available worker node to the pool", len(c.availableWorkerNode))

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
				c.availableWorkerNodeLoop(ctx)
			}
		}
	}(c.ctx)
}

func (c *controlPlane) availableWorkerNodeLoop(ctx context.Context) {
	logrus.Info("[ControlPlane Daemon] Performing on worker node availability check")
	// Check for all available worker node.
	for key, availableWorkerNode := range c.availableWorkerNode {
		logger := logrus.WithFields(logrus.Fields{
			"nodeID": availableWorkerNode.nodeEntry.NodeID,
			"ip":     availableWorkerNode.nodeEntry.IP,
			"port":   availableWorkerNode.nodeEntry.Port,
		})
		_, err := availableWorkerNode.grpcConnection.HealthCheck(ctx, &emptypb.Empty{})
		if err != nil {
			availableWorkerNode.unavailableCount++
			logger.Warnf("[ControlPlane Daemon] Worker node didn't response to the ping command (%d/%d)", availableWorkerNode.unavailableCount, MaximumUnavailableCount)
			if availableWorkerNode.unavailableCount+1 > MaximumUnavailableCount {
				logger.Warnf("[ControlPlane Daemon] Worker node has been deleted from the available worker node pool due to unresponsive has been detected.")
				delete(c.availableWorkerNode, key)
				continue
			}
		} else {
			// If the node become available again, reset it unavailable stat.
			availableWorkerNode.unavailableCount = 0
		}
		c.availableWorkerNode[key] = availableWorkerNode
	}

	timeutil.SleepWithContext(ctx, AvailabilityCheckSleepDuration)
}

func (c *controlPlane) AddAvailableWorkerNode(nodeEntry models.NodeEntry) error {
	logger := logrus.WithFields(logrus.Fields{
		"nodeID": nodeEntry.NodeID,
		"ip":     nodeEntry.IP,
		"port":   nodeEntry.Port,
	})

	// Create gRPC connection
	target := fmt.Sprintf("%s:%d", nodeEntry.IP, nodeEntry.Port)
	client, err := grpc.ProvideRAJDSGrpcClient(target, c.caCertificate)
	if err != nil {
		logger.Warnf("[ControlPlane Daemon] Failed add worker node to the pool with error (%s)", err.Error())
		return err
	}

	clientProto := proto.NewWorkerNodeClient(client.GetConnection())
	_, err = clientProto.HealthCheck(c.ctx, &emptypb.Empty{})
	if err != nil {
		logger.Warnf("[ControlPlane Daemon] Failed add worker node to the pool with error (%s)", err.Error())
		return err
	}

	c.availableWorkerNode[nodeEntry.NodeID] = availableWorkerNodeGRPCMapper{
		nodeEntry:      nodeEntry,
		grpcConnection: clientProto,
	}

	logger.Infof("[ControlPlane Daemon] A Worker has been added to the pool")
	return nil
}

func (c *controlPlane) loop(ctx context.Context) {
	// TODO: Get some available task from database
	// TODO: Call Distribute function
}

func (c *controlPlane) GracefullyShutdown() {
	logrus.Info("[ControlPlane Daemon] Received gracefully shutdown command")
	c.cancelFunc()
	c.gracefullyShutdownWaitGroup.Wait()
	logrus.Info("[ControlPlane Daemon] Gracefully Shutdown success.")
}
