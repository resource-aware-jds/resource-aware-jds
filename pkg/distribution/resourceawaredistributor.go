package distribution

import (
	"context"
	"errors"
	"github.com/resource-aware-jds/resource-aware-jds/config"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"sort"
)

type ResourceAwareDistributor struct {
	distributedTaskCounter metric.Int64Counter
	config                 config.ResourceAwareDistributorConfigModel
}

func ProvideResourceAwareDistributor(config config.ResourceAwareDistributorConfigModel, meter metric.Meter) Distributor {
	counter, err := meter.Int64Counter("resource_aware_distributor_task")
	if err != nil {
		panic(err)
	}
	return &ResourceAwareDistributor{
		distributedTaskCounter: counter,
		config:                 config,
	}
}

type maximumTaskForNode struct {
	node      NodeMapper
	totalTask int
}

func (r ResourceAwareDistributor) Distribute(ctx context.Context, nodes []NodeMapper, tasks []models.Task, dependency DistributorDependency) (successTask []models.Task, distributionError []DistributeError, err error) {
	// Expect that all the task should have the same job id
	var jobID *primitive.ObjectID
	for _, task := range tasks {
		if jobID == nil {
			jobID = task.JobID
		} else if jobID != task.JobID {
			return nil, nil, errors.New("[ResourceAwareDistributor] distribution failed due to jobID is not the same")
		}
	}

	if err != nil {
		return nil, nil, errors.New("[ResourceAwareDistributor] fail to get average resource usage")
	}

	// Calculate maximum resource for each node that can take
	nodeWithMaximumTasks := make([]maximumTaskForNode, 0, len(nodes))
	taskToDistribute := len(tasks)
	for _, node := range nodes {
		maximumTask := r.calculateTheMaximumTaskOnNode(dependency, node.AvailableResource, taskToDistribute)
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
			err := r.distributeToNode(ctx, nodeWithMaximumTask.node, toBeDistributedTask)
			if err != nil {
				distributionError = append(distributionError, DistributeError{
					Task:  toBeDistributedTask,
					Error: err,
				})
				continue
			}
			successTask = append(successTask, toBeDistributedTask)
		}
	}
	return
}

func (r ResourceAwareDistributor) distributeToNode(ctx context.Context, node NodeMapper, task models.Task) error {
	logger := node.Logger.WithField("taskID", task.ID.Hex())
	logger.Info("[ResourceAwareDistributor] Sending task to the worker node")
	_, err := node.GRPCConnection.SendTask(ctx, &proto.RecievedTask{
		ID:             task.ID.Hex(),
		TaskAttributes: task.TaskAttributes,
		DockerImage:    task.ImageUrl,
	})
	if err != nil {
		logger.Warnf("[ResourceAwareDistributor] Fail to distribute task to worker node (%s)", err.Error())
		task.DistributionFailure(node.NodeEntry.NodeID, err)
		return err
	}
	// Add log to success task
	task.DistributionSuccess(node.NodeEntry.NodeID)
	r.distributedTaskCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("nodeID", node.NodeEntry.NodeID)))
	logger.Info("[Distributor] Worker Node accepted the task")
	return nil
}

func (r ResourceAwareDistributor) calculateTheMaximumTaskOnNode(dependency DistributorDependency, nodeResource models.AvailableResource, toBeDistributedTask int) int {
	memory := dependency.TaskResourceUsage.Memory
	cpu := dependency.TaskResourceUsage.CPU

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
