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

type DistributeError struct {
	NodeEntry models.NodeEntry
	Task      models.Task
	Error     error
}

type Distributor interface {
	Distribute(ctx context.Context, nodes []NodeMapper, tasks []models.Task, dependency DistributorDependency) ([]models.Task, []DistributeError, error)
}

type DistributorDependency struct {
	TaskResourceUsage models.TaskResourceUsage
}
