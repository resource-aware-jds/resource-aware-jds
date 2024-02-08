package distribution

import (
	"context"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/datastructure"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type RoundRobinDistributor struct {
	distributedTaskCounter metric.Int64Counter
}

func ProvideRoundRobinDistributor(meter metric.Meter) Distributor {
	counter, err := meter.Int64Counter("cp_distributed_task")
	if err != nil {
		panic(err)
	}
	return &RoundRobinDistributor{
		distributedTaskCounter: counter,
	}
}

func (r RoundRobinDistributor) Distribute(ctx context.Context, nodes []NodeMapper, tasks []models.Task) (successTask []models.Task, distributionError []DistributeError, err error) {
	nodeRoundRobin, err := datastructure.ProvideRoundRobin[NodeMapper](nodes...)
	if err != nil {
		return nil, nil, err
	}

	for _, task := range tasks {
		focusedNode := nodeRoundRobin.Next()
		logger := focusedNode.Logger.WithField("taskID", task.ID.Hex())
		logger.Info("[Distributor] Sending task to the worker node")
		_, err = focusedNode.GRPCConnection.SendTask(ctx, &proto.RecievedTask{
			ID:             task.ID.Hex(),
			TaskAttributes: task.TaskAttributes,
			DockerImage:    task.ImageUrl,
		})
		if err != nil {
			logger.Warnf("[Distributor] Fail to distribute task to worker node (%s)", err.Error())
			task.DistributionFailure(focusedNode.NodeEntry.NodeID, err)
			distributionError = append(distributionError, DistributeError{
				NodeEntry: focusedNode.NodeEntry,
				Task:      task,
				Error:     err,
			})
			continue
		}
		// Add log to success task
		task.DistributionSuccess(focusedNode.NodeEntry.NodeID)
		successTask = append(successTask, task)
		r.distributedTaskCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("nodeID", focusedNode.NodeEntry.NodeID)))
		logger.Info("[Distributor] Worker Node accepted the task")
	}
	return successTask, distributionError, nil
}
