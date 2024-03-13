package service

import (
	"context"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/distribution"
	"github.com/resource-aware-jds/resource-aware-jds/repository"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type task struct {
	taskRepository repository.ITask
}

type Task interface {
	GetAvailableTask(ctx context.Context) ([]models.Task, error)
	UpdateTaskAfterDistribution(ctx context.Context, successTasks []models.Task, errorTasks []distribution.DistributeError) error
	UpdateTaskWorkOnFailure(ctx context.Context, taskID primitive.ObjectID, nodeID string, errMessage string) error
	UpdateTaskSuccess(ctx context.Context, taskID primitive.ObjectID, nodeID string, result []byte) error
	CreateTask(ctx context.Context, job *models.Job, taskAttributes [][]byte, isExperiment bool) ([]models.Task, error)
	GetTaskByJob(ctx context.Context, job *models.Job) ([]models.Task, error)
	GetTaskByID(ctx context.Context, taskID primitive.ObjectID) (*models.Task, error)
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

func (t *task) CreateTask(ctx context.Context, job *models.Job, taskAttributes [][]byte, isExperiment bool) ([]models.Task, error) {
	// Create Tasks
	tasks := make([]models.Task, 0, len(taskAttributes))
	for _, taskAttribute := range taskAttributes {
		newTask := models.Task{
			JobID:          job.ID,
			Status:         models.CreatedTaskStatus,
			ImageUrl:       job.ImageURL,
			TaskAttributes: taskAttribute,
		}
		tasks = append(tasks, newTask)
	}

	if isExperiment {
		tasks[0].ExperimentTask()
	} else {
		for _, taskData := range tasks {
			taskData.SkipExperimentTask()
		}
	}

	err := t.taskRepository.InsertMany(ctx, tasks)
	if err != nil {
		return nil, err
	}

	return t.taskRepository.FindManyByJobID(ctx, job.ID)
}

func (t *task) GetTaskByJob(ctx context.Context, job *models.Job) ([]models.Task, error) {
	return t.taskRepository.FindManyByJobID(ctx, job.ID)
}

func (t *task) GetTaskByID(ctx context.Context, taskID primitive.ObjectID) (*models.Task, error) {
	return t.taskRepository.FindOneByID(ctx, taskID)
}

func (t *task) UpdateTaskWorkOnFailure(ctx context.Context, taskID primitive.ObjectID, nodeID string, errMessage string) error {
	taskResult, err := t.GetTaskByID(ctx, taskID)
	if err != nil {
		logrus.Errorf("get task error %v", err)
		return err
	}

	taskResult.WorkOnTaskFailure(nodeID, errMessage)
	return t.taskRepository.BulkWriteStatusAndLogByID(ctx, []models.Task{*taskResult})
}

func (t *task) UpdateTaskSuccess(ctx context.Context, taskID primitive.ObjectID, nodeID string, result []byte) error {
	taskResult, err := t.GetTaskByID(ctx, taskID)
	if err != nil {
		logrus.Errorf("get task error %v", err)
		return err
	}

	taskResult.SuccessTask(nodeID, result)
	return t.taskRepository.WriteTaskResult(ctx, *taskResult)
}
