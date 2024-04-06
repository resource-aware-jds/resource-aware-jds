package pool

import (
	"context"
	"errors"
	"fmt"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/cert"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/distribution"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/grpc"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/metrics"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/util"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
	"strconv"
	"sync"
	"time"
)

//go:generate mockgen -source=./workernodepool.go -destination=./mock_pool/mock_workernodepool.go -package=mock_pool

const (
	MaximumUnavailableCount  = 3
	AvailabilityCheckTimeout = 5 * time.Second
)

var (
	ErrNoAvailableWorkerNode  = errors.New("no available worker node in the pool")
	ErrNoAvailableDistributor = errors.New("no distributor available")
)

type InitialWorkerNodeSet []models.NodeEntry

type workerNodePoolMapper struct {
	nodeEntry         models.NodeEntry
	availableResource models.AvailableResource
	grpcConnection    proto.WorkerNodeClient
	unavailableCount  uint
	logger            *logrus.Entry
}

type workerNode struct {
	pool              map[string]workerNodePoolMapper
	caCertificate     cert.CACertificate
	distributorMapper distribution.DistributorMapper
	poolMu            sync.Mutex
	grpcResolver      grpc.RAJDSGRPCResolver
	metricCounter     metric.Int64ObservableCounter
	logger            *logrus.Entry
}

type WorkerNode interface {
	InitializePool(ctx context.Context, nodeEntries []models.NodeEntry)
	AddWorkerNode(ctx context.Context, node models.NodeEntry) error
	WorkerNodeAvailabilityCheck(ctx context.Context)
	DistributeWork(ctx context.Context, jobID models.Job, tasks []models.Task) ([]models.Task, []models.DistributeError, error)
	IsAvailableWorkerNode() bool
	RemoveNodeFromPool(ctx context.Context, nodeID string)
	CheckRunningTaskInEachWorkerNode(ctx context.Context) map[primitive.ObjectID]bool
	PoolSize() int
	GetAllWorkerNode() []models.NodeEntry
}

func ProvideWorkerNode(caCertificate cert.CACertificate, distributorMapper distribution.DistributorMapper, grpcResolver grpc.RAJDSGRPCResolver, meter metric.Meter) WorkerNode {
	pool := make(map[string]workerNodePoolMapper)
	counter, _ := meter.Int64ObservableCounter(
		metrics.GenerateControlPlaneMetric("connected_worker_nodes"),
		metric.WithInt64Callback(func(ctx context.Context, observer metric.Int64Observer) error {
			observer.Observe(int64(len(pool)))
			return nil
		}),
		metric.WithUnit("Node"),
		metric.WithDescription("The total alive Worker Node connected to this Control Plane"),
	)

	logger := logrus.WithFields(logrus.Fields{
		"role":      "Control Plane",
		"component": "node_pool",
	})

	meter.Float64ObservableGauge( // nolint:errcheck
		metrics.GenerateControlPlaneMetric("node_available_cpu"),
		metric.WithFloat64Callback(func(ctx context.Context, observer metric.Float64Observer) error {
			for _, node := range pool {
				observer.Observe(
					float64(node.availableResource.AvailableCpuPercentage),
					metric.WithAttributes(
						attribute.String("nodeID", node.nodeEntry.NodeID),
						attribute.String("ip", node.nodeEntry.IP),
					),
				)
			}
			return nil
		}),
	)
	meter.Float64ObservableGauge( // nolint:errcheck
		metrics.GenerateControlPlaneMetric("node_available_cpu"),
		metric.WithUnit("mb"),
		metric.WithFloat64Callback(func(ctx context.Context, observer metric.Float64Observer) error {
			for _, node := range pool {
				observer.Observe(
					util.ConvertToMib(node.availableResource.AvailableMemory).Size,
					metric.WithAttributes(
						attribute.String("nodeID", node.nodeEntry.NodeID),
						attribute.String("ip", node.nodeEntry.IP),
					),
				)
			}
			return nil
		}),
	)
	return &workerNode{
		caCertificate:     caCertificate,
		pool:              pool,
		distributorMapper: distributorMapper,
		grpcResolver:      grpcResolver,
		metricCounter:     counter,
		logger:            logger,
	}
}

func (w *workerNode) InitializePool(ctx context.Context, nodeEntries []models.NodeEntry) {
	var wg sync.WaitGroup
	for _, node := range nodeEntries {
		wg.Add(1)
		go func(nodeToAdd models.NodeEntry) {
			w.AddWorkerNode(ctx, nodeToAdd) // nolint:errcheck
			wg.Done()
		}(node)
	}

	wg.Wait()

	logrus.Infof("[WorkerNode Pool] Added %d available worker node to the pool", len(w.pool))
}

func (w *workerNode) AddWorkerNode(ctx context.Context, node models.NodeEntry) error {
	logger := logrus.WithFields(logrus.Fields{
		"nodeID": node.NodeID,
		"ip":     node.IP,
		"port":   node.Port,
	})

	// Check if /etc/host/ already contain the host and domain
	focusedHost := fmt.Sprintf("%s.%s", node.NodeID, cert.GetDefaultDomainName())
	target := fmt.Sprintf("rajds://%s", focusedHost)

	w.grpcResolver.AddHost(focusedHost, net.JoinHostPort(node.IP, strconv.Itoa(int(node.Port))))

	// Create gRPC connection
	client, err := grpc.ProvideRAJDSGrpcClient(grpc.ClientConfig{
		Target:        target,
		CACertificate: w.caCertificate,
		ServerName:    focusedHost,
	})
	if err != nil {
		logger.Warnf("Failed add worker node to the pool with error (%s)", err.Error())
		return err
	}

	clientProto := proto.NewWorkerNodeClient(client.GetConnection())
	_, err = clientProto.HealthCheck(ctx, &emptypb.Empty{})
	if err != nil {
		logger.Warnf("Failed add worker node to the pool with error (%s)", err.Error())
		return err
	}

	w.poolMu.Lock()
	w.pool[node.NodeID] = workerNodePoolMapper{
		nodeEntry:      node,
		grpcConnection: clientProto,
		logger:         logger,
	}
	w.poolMu.Unlock()

	logger.Infof("A Worker has been added to the pool")
	return nil
}

func (w *workerNode) WorkerNodeAvailabilityCheck(ctx context.Context) {
	ok := w.poolMu.TryLock()
	if !ok {
		w.logger.Warn("Skipping the worker node availability check due to distribution is performing")
		return
	}
	defer w.poolMu.Unlock()
	// Check for all available worker node.
	for key, focusedNode := range w.pool {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, AvailabilityCheckTimeout)
		resource, err := focusedNode.grpcConnection.HealthCheck(ctxWithTimeout, &emptypb.Empty{})
		cancel()

		focusedNodeLogger := focusedNode.logger.WithFields(w.logger.Data)
		if err != nil {
			focusedNode.unavailableCount++
			focusedNodeLogger.Warnf("Worker node didn't response to the ping command (%d/%d)", focusedNode.unavailableCount, MaximumUnavailableCount)
			if focusedNode.unavailableCount+1 > MaximumUnavailableCount {
				focusedNodeLogger.Warnf("Worker node has been deleted from the available worker node pool due to unresponsive has been detected.")
				delete(w.pool, key)
				continue
			}
		} else {
			// If the node become available again, reset it unavailable stat.
			focusedNode.unavailableCount = 0
		}
		focusedNode.availableResource = models.AvailableResource{
			CpuCores:               resource.GetCpuCores(),
			AvailableCpuPercentage: resource.GetAvailableCpuPercentage(),
			AvailableMemory:        util.ExtractMemoryUsageString(resource.GetAvailableMemory()),
		}
		w.pool[key] = focusedNode
	}
}

func (w *workerNode) DistributeWork(ctx context.Context, job models.Job, tasks []models.Task) ([]models.Task, []models.DistributeError, error) {
	w.poolMu.Lock()
	defer w.poolMu.Unlock()

	if len(w.pool) == 0 {
		return nil, nil, ErrNoAvailableWorkerNode
	}

	nodeMapper := make([]distribution.NodeMapper, 0, len(w.pool))
	for _, node := range w.pool {
		nodeMapper = append(nodeMapper, distribution.NodeMapper{
			NodeEntry:         node.nodeEntry,
			AvailableResource: node.availableResource,
			GRPCConnection:    node.grpcConnection,
			Logger:            node.logger,
		})
	}
	var dist distribution.Distributor
	var ok bool

	if job.Status == models.ExperimentingJobStatus && len(tasks) == 1 {
		w.logger.Info("Job in experimenting, Temporary switch to use RoundRobin distributor for experimenting task")
		dist, ok = w.distributorMapper.GetDistributor(models.RoundRobinDistributorName)
	} else {
		dist, ok = w.distributorMapper.GetDistributor(job.DistributorLogic)
	}

	if !ok {
		w.logger.Errorf("No Distributing solution found in the distributor logics (%s)", job.DistributorLogic)
		return nil, nil, ErrNoAvailableDistributor
	}

	w.logger.Infof("Distributing the Task using the %s distributor", job.DistributorLogic)
	return dist.Distribute(ctx, nodeMapper, tasks)
}

func (w *workerNode) IsAvailableWorkerNode() bool {
	w.poolMu.Lock()
	defer w.poolMu.Unlock()

	return len(w.pool) != 0
}

func (w *workerNode) RemoveNodeFromPool(_ context.Context, nodeID string) {
	w.poolMu.Lock()
	defer w.poolMu.Unlock()
	w.logger.Infof("Remove Node %s from the pool", nodeID)

	delete(w.pool, nodeID)
}

func (w *workerNode) CheckRunningTaskInEachWorkerNode(ctx context.Context) map[primitive.ObjectID]bool {
	responseObjectIDs := make(map[primitive.ObjectID]bool, 0)
	for _, nodeMapper := range w.pool {
		logger := w.logger.WithFields(logrus.Fields{
			"nodeID": nodeMapper.nodeEntry.NodeID,
		})
		response, err := nodeMapper.grpcConnection.GetAllTasks(ctx, &emptypb.Empty{})
		if err != nil {
			logger.Warnf("Fail ot get all running task on node : %e", err)
			continue
		}

		for _, taskID := range response.GetTaskIDs() {
			parsedTaskID, err := primitive.ObjectIDFromHex(taskID)
			if err != nil {
				logger.Warnf("Node Response invalid ObjectID: %e", err)
				continue
			}
			responseObjectIDs[parsedTaskID] = true
		}
	}
	return responseObjectIDs
}

func (w *workerNode) PoolSize() int {
	return len(w.pool)
}

func (w *workerNode) GetAllWorkerNode() []models.NodeEntry {
	w.poolMu.Lock()
	defer w.poolMu.Unlock()

	result := make([]models.NodeEntry, 0, len(w.pool))
	for _, data := range w.pool {
		result = append(result, data.nodeEntry)
	}

	return result
}
