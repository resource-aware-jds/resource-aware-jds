package distribution

import (
	"context"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/datastructure"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/metric"
)

type RoundRobinDistributor struct {
	baseDistributor
}

func ProvideRoundRobinDistributor(meter metric.Meter) Distributor {
	return &RoundRobinDistributor{
		baseDistributor: newBaseDistributor(models.RoundRobinDistributorName, meter),
	}
}

func (r RoundRobinDistributor) Distribute(ctx context.Context, nodes []NodeMapper, tasks []models.Task) (successTask []models.Task, distributionError []models.DistributeError, err error) {
	nodeRoundRobin, err := datastructure.ProvideRoundRobin[NodeMapper](nodes...)
	if err != nil {
		return nil, nil, err
	}

	for _, node := range nodes {
		logrus.Info("[Debug] NodeList: ", node.NodeEntry.NodeID)
	}

	for _, task := range tasks {
		focusedNode := nodeRoundRobin.Next()
		logrus.Info("[Debug] Focused Node: ", focusedNode.NodeEntry.NodeID)
		r.distributeToNode(ctx, focusedNode, task, &successTask, &distributionError)
	}
	return successTask, distributionError, nil
}
