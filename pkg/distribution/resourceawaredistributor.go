package distribution

import (
	"context"
	"errors"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"math/rand"
)

type ResourceAwareDistributor struct {
	distributedTaskCounter metric.Int64Counter
}

func ProvideResourceAwareDistributor(meter metric.Meter) Distributor {
	counter, err := meter.Int64Counter("resource_aware_distributor_task")
	if err != nil {
		panic(err)
	}
	return &ResourceAwareDistributor{
		distributedTaskCounter: counter,
	}
}

func (r ResourceAwareDistributor) Distribute(ctx context.Context, nodes []NodeMapper, tasks []models.Task, dependency DistributorDependency) (successTask []models.Task, distributionError []DistributeError, err error) {
	// Expect that all the task should have the same job id
	var jobID *primitive.ObjectID
	for _, task := range tasks {
		if jobID == nil {
			jobID = task.JobID
		} else if jobID != task.JobID {
			return nil, nil, errors.New("[ResourceAwareDistributor] distribution failed due to jobID is not the same")
		}
	}

	if err != nil {
		return nil, nil, errors.New("[ResourceAwareDistributor] fail to get average resource usage")
	}

	for _, task := range tasks {
		isDistributed := false
		for _, node := range nodes {
			if float32(util.ConvertToMib(node.AvailableResource.AvailableMemory).Size)-dependency.TaskResourceUsage.Memory <= 0 ||
				node.AvailableResource.AvailableCpuPercentage-dependency.TaskResourceUsage.CPU <= 0 {
				continue
			}
			err = r.distributeToNode(ctx, node, task)
			isDistributed = err == nil
		}
		if !isDistributed {
			if err == nil {
				err = errors.New("no available node")
			}
			distributionError = append(distributionError, DistributeError{
				Task:  task,
				Error: err,
			})
			continue
		}

		successTask = append(successTask, task)
		rand.Shuffle(len(nodes), func(i, j int) {
			nodes[i], nodes[j] = nodes[j], nodes[i]
		})
	}
	return
}

func (r ResourceAwareDistributor) distributeToNode(ctx context.Context, node NodeMapper, task models.Task) error {
	logger := node.Logger.WithField("taskID", task.ID.Hex())
	logger.Info("[ResourceAwareDistributor] Sending task to the worker node")
	_, err := node.GRPCConnection.SendTask(ctx, &proto.RecievedTask{
		ID:             task.ID.Hex(),
		TaskAttributes: task.TaskAttributes,
		DockerImage:    task.ImageUrl,
	})
	if err != nil {
		logger.Warnf("[ResourceAwareDistributor] Fail to distribute task to worker node (%s)", err.Error())
		task.DistributionFailure(node.NodeEntry.NodeID, err)
		return err
	}
	// Add log to success task
	task.DistributionSuccess(node.NodeEntry.NodeID)
	r.distributedTaskCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("nodeID", node.NodeEntry.NodeID)))
	logger.Info("[Distributor] Worker Node accepted the task")
	return nil
}
