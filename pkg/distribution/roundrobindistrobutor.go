package distribution

import (
	"context"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/datastructure"
	"go.opentelemetry.io/otel/metric"
)

type RoundRobinDistributor struct {
	baseDistributor
}

func ProvideRoundRobinDistributor(meter metric.Meter) Distributor {
	return &RoundRobinDistributor{
		baseDistributor: newBaseDistributor("round_robin", meter),
	}
}

func (r RoundRobinDistributor) Distribute(ctx context.Context, nodes []NodeMapper, tasks []models.Task, _ DistributorDependency) (successTask []models.Task, distributionError []DistributeError, err error) {
	nodeRoundRobin, err := datastructure.ProvideRoundRobin[NodeMapper](nodes...)
	if err != nil {
		return nil, nil, err
	}

	for _, task := range tasks {
		focusedNode := nodeRoundRobin.Next()
		r.distributeToNode(ctx, focusedNode, task, &successTask, &distributionError)
	}
	return successTask, distributionError, nil
}
