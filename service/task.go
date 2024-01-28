package service

import (
	"context"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/distribution"
	"github.com/resource-aware-jds/resource-aware-jds/repository"
)

type task struct {
	taskRepository repository.ITask
}

type Task interface {
	GetAvailableTask(ctx context.Context) ([]models.Task, error)
	UpdateTaskAfterDistribution(ctx context.Context, successTasks []models.Task, errorTasks []distribution.DistributeError) error
}

func ProvideTaskService(taskRepository repository.ITask) Task {
	return &task{
		taskRepository: taskRepository,
	}
}

func (t *task) GetAvailableTask(ctx context.Context) ([]models.Task, error) {
	return t.taskRepository.GetTaskToDistribute(ctx)
}

func (t *task) UpdateTaskAfterDistribution(ctx context.Context, successTasks []models.Task, errorTasks []distribution.DistributeError) error {
	taskToUpdate := make([]models.Task, 0, len(successTasks)+len(errorTasks))
	taskToUpdate = append(taskToUpdate, successTasks...)

	// Add failure task
	for _, errorTask := range errorTasks {
		taskToUpdate = append(taskToUpdate, errorTask.Task)
	}

	return t.taskRepository.BulkWriteStatusAndLogByID(ctx, taskToUpdate)
}
