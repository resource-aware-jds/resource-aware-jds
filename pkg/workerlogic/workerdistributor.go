package workerlogic

import (
	"context"
	"github.com/resource-aware-jds/resource-aware-jds/config"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/container"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/datastructure"
	"time"
)

type WorkerState struct {
	ContainerCoolDownState datastructure.Buffer[string, time.Time]
	WorkerNodeConfig       config.WorkerConfigModel
	ContainerBuffer        datastructure.Buffer[string, container.IContainer]
}

type DistributionResult struct {
	CreateContainerToSupportTask bool
}

type WorkerDistributor interface {
	Distribute(ctx context.Context, task models.Task, workerState WorkerState) DistributionResult
}
