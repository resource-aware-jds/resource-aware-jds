package workerdistribution

import (
	"context"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"time"
)

type delayWorkerDistributor struct{}

func ProvideDelayWorkerDistributor() WorkerDistributor {
	return &delayWorkerDistributor{}
}

func (d delayWorkerDistributor) Distribute(ctx context.Context, task models.Task, workerState WorkerState) DistributionResult {
	// Check if container is already persisted.
	container := workerState.ContainerBuffer.Get(task.ImageUrl)
	if container == nil {
		return DistributionResult{
			CreateContainerToSupportTask: true,
		}
	}

	coolDownStatus := workerState.ContainerCoolDownState.Get(task.ImageUrl)
	if coolDownStatus != nil {
		coolDownStatusDepointer := *coolDownStatus
		if coolDownStatusDepointer.After(time.Now()) {
			return DistributionResult{
				CreateContainerToSupportTask: false,
			}
		}

	}
	return DistributionResult{
		CreateContainerToSupportTask: true,
	}
}
