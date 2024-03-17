package distribution

import (
	"context"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/sirupsen/logrus"
)

type NodeMapper struct {
	NodeEntry         models.NodeEntry
	GRPCConnection    proto.WorkerNodeClient
	AvailableResource models.AvailableResource
	Logger            *logrus.Entry
}

type Distributor interface {
	Distribute(ctx context.Context, nodes []NodeMapper, tasks []models.Task) ([]models.Task, []models.DistributeError, error)
}
