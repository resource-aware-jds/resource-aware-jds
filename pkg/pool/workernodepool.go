package pool

import (
	"context"
	"fmt"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/cert"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/distribution"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/timeutil"
	"github.com/resource-aware-jds/resource-aware-jds/service"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
	"sync"
	"time"
)

const (
	MaximumUnavailableCount        = 3
	AvailabilityCheckSleepDuration = 5 * time.Second
	AvailabilityCheckTimeout       = 5 * time.Second
)

type workerNodePoolMapper struct {
	nodeEntry        models.NodeEntry
	grpcConnection   proto.WorkerNodeClient
	unavailableCount uint
	logger           *logrus.Entry
}

type workerNode struct {
	pool                map[string]workerNodePoolMapper
	caCertificate       cert.CACertificate
	controlPlaneService service.IControlPlane
	distributor         distribution.Distributor
	poolMu              sync.Mutex
}

type WorkerNode interface {
	InitializePool(ctx context.Context)
	AddWorkerNode(ctx context.Context, node models.NodeEntry) error
	WorkerNodeAvailabilityCheck(ctx context.Context)
}

func ProvideWorkerNode(caCertificate cert.CACertificate, controlPlaneService service.IControlPlane, distributor distribution.Distributor) WorkerNode {
	return &workerNode{
		caCertificate:       caCertificate,
		pool:                make(map[string]workerNodePoolMapper),
		controlPlaneService: controlPlaneService,
		distributor:         distributor,
	}
}

func (w *workerNode) InitializePool(ctx context.Context) {
	logrus.Info("[WorkerNode Pool] Get all available worker node from registry")
	nodes, err := w.controlPlaneService.GetAllWorkerNodeFromRegistry(ctx)
	if err != nil {
		logrus.Warnf("[WorkerNode Pool] Failed to get all available worker node from registry with error (%s)", err.Error())
	}

	for _, node := range nodes {
		err = w.AddWorkerNode(ctx, node)
		if err != nil {
			continue
		}
	}

	logrus.Infof("[WorkerNode Pool] Added %d available worker node to the pool", len(w.pool))

}

func (w *workerNode) AddWorkerNode(ctx context.Context, node models.NodeEntry) error {
	logger := logrus.WithFields(logrus.Fields{
		"nodeID": node.NodeID,
		"ip":     node.IP,
		"port":   node.Port,
	})

	// Create gRPC connection
	target := fmt.Sprintf("%s:%d", node.IP, node.Port)
	client, err := grpc.ProvideRAJDSGrpcClient(target, w.caCertificate)
	if err != nil {
		logger.Warnf("[WorkerNode Pool] Failed add worker node to the pool with error (%s)", err.Error())
		return err
	}

	clientProto := proto.NewWorkerNodeClient(client.GetConnection())
	_, err = clientProto.HealthCheck(ctx, &emptypb.Empty{})
	if err != nil {
		logger.Warnf("[WorkerNode Pool] Failed add worker node to the pool with error (%s)", err.Error())
		return err
	}

	w.pool[node.NodeID] = workerNodePoolMapper{
		nodeEntry:      node,
		grpcConnection: clientProto,
		logger:         logger,
	}

	logger.Infof("[WorkerNode Pool] A Worker has been added to the pool")
	return nil
}

func (w *workerNode) WorkerNodeAvailabilityCheck(ctx context.Context) {
	logrus.Info("[WorkerNode Pool] Performing on worker node availability check")
	ok := w.poolMu.TryLock()
	if !ok {
		logrus.Warn("[WorkerNode Pool] Skipping the worker node availability check due to distribution is performing")
		return
	}
	defer w.poolMu.Unlock()
	// Check for all available worker node.
	for key, focusedNode := range w.pool {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, AvailabilityCheckTimeout)
		_, err := focusedNode.grpcConnection.HealthCheck(ctxWithTimeout, &emptypb.Empty{})
		cancel()
		if err != nil {
			focusedNode.unavailableCount++
			focusedNode.logger.Warnf("[ControlPlane Daemon] Worker node didn't response to the ping command (%d/%d)", focusedNode.unavailableCount, MaximumUnavailableCount)
			if focusedNode.unavailableCount+1 > MaximumUnavailableCount {
				focusedNode.logger.Warnf("[ControlPlane Daemon] Worker node has been deleted from the available worker node pool due to unresponsive has been detected.")
				delete(w.pool, key)
				continue
			}
		} else {
			// If the node become available again, reset it unavailable stat.
			focusedNode.unavailableCount = 0
		}
		w.pool[key] = focusedNode
	}

	timeutil.SleepWithContext(ctx, AvailabilityCheckSleepDuration)
}

func (w *workerNode) DistributeWork(ctx context.Context, tasks []models.Task) ([]distribution.DistributeError, error) {
	w.poolMu.Lock()
	defer w.poolMu.Unlock()

	nodeMapper := make([]distribution.NodeMapper, 0, len(w.pool))
	for _, node := range w.pool {
		nodeMapper = append(nodeMapper, distribution.NodeMapper{
			NodeEntry:      node.nodeEntry,
			GRPCConnection: node.grpcConnection,
			Logger:         node.logger,
		})
	}

	return w.distributor.Distribute(ctx, nodeMapper, tasks)
}
