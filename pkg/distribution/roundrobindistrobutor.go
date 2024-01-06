package distribution

import (
	"context"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/datastructure"
)

type RoundRobinDistributor struct {
}

func ProvideRoundRobinDistributor() Distributor {
	return &RoundRobinDistributor{}
}

func (r RoundRobinDistributor) Distribute(ctx context.Context, nodes []NodeMapper, tasks []models.Task) (distributionError []DistributeError, err error) {
	nodeRoundRobin, err := datastructure.ProvideRoundRobin[NodeMapper](nodes...)
	if err != nil {
		return nil, err
	}
	for _, task := range tasks {
		focusedNode := nodeRoundRobin.Next()
		logger := focusedNode.Logger.WithField("taskID", task.ID.Hex())
		logger.Info("[Distributor] Sending task to the worker node")
		_, err = focusedNode.GRPCConnection.SendTask(ctx, &proto.RecievedTask{})
		if err != nil {
			logger.Warnf("[Distributor] Fail to distribute task to worker node (%s)", err.Error())
			distributionError = append(distributionError, DistributeError{
				NodeEntry: focusedNode.NodeEntry,
				Task:      task,
				Error:     err,
			})
			continue
		}
		logger.Info("[Distributor] Worker Node accepted the task")
	}
	return distributionError, nil
}
