package distribution

import (
	"context"
	"errors"
	"github.com/resource-aware-jds/resource-aware-jds/config"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/util"
	"github.com/resource-aware-jds/resource-aware-jds/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel/metric"
	"sort"
)

var (
	ErrResourceAwareDistributorTaskDifferenceJob = errors.New("distribution failed due to jobID is not the same")
)

type ResourceAwareDistributor struct {
	baseDistributor

	config      config.ResourceAwareDistributorConfigModel
	taskService service.Task
}

func ProvideResourceAwareDistributor(config config.ResourceAwareDistributorConfigModel, meter metric.Meter, taskService service.Task) Distributor {
	return &ResourceAwareDistributor{
		baseDistributor: newBaseDistributor(ResourceAwareDistributorName, meter),
		config:          config,
		taskService:     taskService,
	}
}

type maximumTaskForNode struct {
	node      NodeMapper
	totalTask int
}

func (r ResourceAwareDistributor) Distribute(ctx context.Context, nodes []NodeMapper, tasks []models.Task) (successTask []models.Task, distributionError []models.DistributeError, err error) {
	err = r.checkTaskWithSameJobID(tasks)
	if err != nil {
		return nil, nil, err
	}

	// Get the average resource usage.
	averageResourceUsage, err := r.taskService.GetAverageResourceUsage(ctx, tasks[0].JobID)
	if err != nil {
		r.logger.Error("Fail to ResourceAware distribute task since no average resource usage info available.")
		return nil, nil, err
	}

	// Calculate maximum resource for each node that can take
	nodeWithMaximumTasks := make([]maximumTaskForNode, 0, len(nodes))
	taskToDistribute := len(tasks)
	for _, node := range nodes {
		maximumTask := r.calculateTheMaximumTaskOnNode(*averageResourceUsage, node.AvailableResource, taskToDistribute)
		nodeWithMaximumTasks = append(nodeWithMaximumTasks, maximumTaskForNode{
			node:      node,
			totalTask: maximumTask,
		})
	}

	// Sort by the totalTask desc
	sort.Slice(nodeWithMaximumTasks, func(i, j int) bool {
		return nodeWithMaximumTasks[i].totalTask > nodeWithMaximumTasks[j].totalTask
	})

	for _, nodeWithMaximumTask := range nodeWithMaximumTasks {
		toBeDistributedTasks := tasks[0:nodeWithMaximumTask.totalTask]
		tasks = tasks[nodeWithMaximumTask.totalTask:]
		for _, toBeDistributedTask := range toBeDistributedTasks {
			r.distributeToNode(ctx, nodeWithMaximumTask.node, toBeDistributedTask, &successTask, &distributionError)
		}
	}
	return
}

func (r ResourceAwareDistributor) checkTaskWithSameJobID(tasks []models.Task) error {
	// Expect that all the task should have the same job id
	var jobID *primitive.ObjectID
	for _, task := range tasks {
		if jobID == nil {
			jobID = task.JobID
		} else if jobID != task.JobID {
			return ErrResourceAwareDistributorTaskDifferenceJob
		}
	}
	return nil
}

func (r ResourceAwareDistributor) calculateTheMaximumTaskOnNode(averageResourceUsage models.TaskResourceUsage, nodeResource models.AvailableResource, toBeDistributedTask int) int {
	memory := averageResourceUsage.Memory
	cpu := averageResourceUsage.CPU

	maximumMemory := util.ConvertToMib(nodeResource.AvailableMemory).Size * float64(r.config.AvailableResourceClearanceThreshold/100)
	maximumCPU := nodeResource.AvailableCpuPercentage * (r.config.AvailableResourceClearanceThreshold / 100)

	totalTask := 1
	for {
		desiredMemory := memory * float64(totalTask)
		desiredCPU := cpu * float32(totalTask)

		if desiredMemory > maximumMemory || desiredCPU > maximumCPU || totalTask > toBeDistributedTask {
			return totalTask - 1
		}
		totalTask++
	}
}
