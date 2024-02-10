package workerlogic

import (
	"context"
	"github.com/resource-aware-jds/resource-aware-jds/config"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"time"
)

type delayWorkerDistributor struct {
	config config.WorkerConfigModel
}

func ProvideDelayWorkerDistributor(config config.WorkerConfigModel) WorkerDistributor {
	return &delayWorkerDistributor{
		config: config,
	}
}

func (d delayWorkerDistributor) Distribute(_ context.Context, task models.Task, workerState WorkerState) DistributionResult {
	// Check if total container is exceeded the total container limit
	if len(workerState.ContainerBuffer) >= d.config.TotalContainerLimit {
		return DistributionResult{
			CreateContainerToSupportTask: false,
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
