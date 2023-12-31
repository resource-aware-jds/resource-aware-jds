package pool

import (
	"context"
	"errors"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/cert"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/distribution"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
	"strconv"
	"sync"
	"time"
)

const (
	MaximumUnavailableCount  = 3
	AvailabilityCheckTimeout = 5 * time.Second
)

var (
	ErrNoAvailableWorkerNode = errors.New("no available worker node in the pool")
)

type InitialWorkerNodeSet []models.NodeEntry

type workerNodePoolMapper struct {
	nodeEntry        models.NodeEntry
	grpcConnection   proto.WorkerNodeClient
	unavailableCount uint
	logger           *logrus.Entry
}

type workerNode struct {
	pool          map[string]workerNodePoolMapper
	caCertificate cert.CACertificate
	distributor   distribution.Distributor
	poolMu        sync.Mutex
}

type WorkerNode interface {
	InitializePool(ctx context.Context, nodeEntries []models.NodeEntry)
	AddWorkerNode(ctx context.Context, node models.NodeEntry) error
	WorkerNodeAvailabilityCheck(ctx context.Context)
	DistributeWork(ctx context.Context, tasks []models.Task) ([]models.Task, []distribution.DistributeError, error)
	IsAvailableWorkerNode() bool
}

func ProvideWorkerNode(caCertificate cert.CACertificate, distributor distribution.Distributor) WorkerNode {
	return &workerNode{
		caCertificate: caCertificate,
		pool:          make(map[string]workerNodePoolMapper),
		distributor:   distributor,
	}
}

func (w *workerNode) InitializePool(ctx context.Context, nodeEntries []models.NodeEntry) {
	for _, node := range nodeEntries {
		w.AddWorkerNode(ctx, node)
	}

	logrus.Infof("[WorkerNode Pool] Added %d available worker node to the pool", len(w.pool))
}

func (w *workerNode) AddWorkerNode(ctx context.Context, node models.NodeEntry) error {
	logger := logrus.WithFields(logrus.Fields{
		"nodeID": node.NodeID,
		"ip":     node.IP,
		"port":   node.Port,
	})

	joinedHostPort := net.JoinHostPort(node.IP, strconv.Itoa(int(node.Port)))

	// Create gRPC connection
	client, err := grpc.ProvideRAJDSGrpcClient(grpc.ClientConfig{
		Target:        joinedHostPort,
		CACertificate: w.caCertificate,
	})
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
}

func (w *workerNode) DistributeWork(ctx context.Context, tasks []models.Task) ([]models.Task, []distribution.DistributeError, error) {
	w.poolMu.Lock()
	defer w.poolMu.Unlock()

	if len(w.pool) == 0 {
		return nil, nil, ErrNoAvailableWorkerNode
	}

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

func (w *workerNode) IsAvailableWorkerNode() bool {
	w.poolMu.Lock()
	defer w.poolMu.Unlock()

	return len(w.pool) != 0
}
