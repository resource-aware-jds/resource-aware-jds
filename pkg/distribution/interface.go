package distribution

import (
	"context"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/sirupsen/logrus"
)

type NodeMapper struct {
	NodeEntry      models.NodeEntry
	GRPCConnection proto.WorkerNodeClient
	Logger         *logrus.Entry
}

type DistributeError struct {
	NodeEntry models.NodeEntry
	Task      models.Task
	Error     error
}

type Distributor interface {
	Distribute(ctx context.Context, nodes []NodeMapper, tasks []models.Task) ([]models.Task, []DistributeError, error)
}
