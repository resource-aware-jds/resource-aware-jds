package service

import (
	"context"
	"errors"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/repository"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//go:generate mockgen -source=./task.go -destination=./mock_service/mock_task.go -package=mock_service

type task struct {
	taskRepository repository.ITask
}

type Task interface {
	GetAvailableTask(ctx context.Context, jobIDs []models.Job) (*models.Job, []models.Task, error)
	UpdateTaskAfterDistribution(ctx context.Context, successTasks []models.Task, errorTasks []models.DistributeError) error
	UpdateTaskWorkOnFailure(ctx context.Context, taskID primitive.ObjectID, nodeID string, errMessage string) error
	UpdateTaskSuccess(ctx context.Context, taskID primitive.ObjectID, nodeID string, result []byte, averageCPUUsage float32, averageMemoryUsage float64) error
	UpdateTaskWaitTimeout(ctx context.Context, taskID primitive.ObjectID) error
	CreateTask(ctx context.Context, job *models.Job, taskAttributes [][]byte, isExperiment bool) ([]models.Task, error)
	GetTaskByJob(ctx context.Context, job *models.Job) ([]models.Task, error)
	GetTaskByID(ctx context.Context, taskID primitive.ObjectID) (*models.Task, error)
	GetAverageResourceUsage(ctx context.Context, jobID *primitive.ObjectID) (*models.TaskResourceUsage, error)
	UpdateTaskToBeReadyToBeDistributed(ctx context.Context, jobID *primitive.ObjectID) error
	CountUnfinishedTaskByJobID(ctx context.Context, jobID *primitive.ObjectID) (int64, error)
}

func ProvideTaskService(taskRepository repository.ITask) Task {
	return &task{
		taskRepository: taskRepository,
	}
}

func (t *task) GetAvailableTask(ctx context.Context, jobs []models.Job) (*models.Job, []models.Task, error) {
	// Distribute Based on Job
	for _, job := range jobs {
		tasks, err := t.taskRepository.GetTaskToDistributeForJob(ctx, job.ID)
		if len(tasks) != 0 || err != nil {
			return &job, tasks, err
		}
	}
	return nil, nil, nil
}

func (t *task) UpdateTaskAfterDistribution(ctx context.Context, successTasks []models.Task, errorTasks []models.DistributeError) error {
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
		for index, taskData := range tasks {
			taskData.SkipExperimentTask()
			tasks[index] = taskData
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

func (t *task) UpdateTaskSuccess(ctx context.Context, taskID primitive.ObjectID, nodeID string, result []byte, averageCPUUsage float32, averageMemoryUsage float64) error {
	taskResult, err := t.GetTaskByID(ctx, taskID)
	if err != nil {
		logrus.Errorf("get task error %v", err)
		return err
	}

	taskResult.SuccessTask(nodeID, result)
	taskResult.ResourceUsage = models.TaskResourceUsage{
		Memory: averageMemoryUsage,
		CPU:    averageCPUUsage,
	}
	return t.taskRepository.WriteTaskResult(ctx, *taskResult)
}

func (t *task) GetAverageResourceUsage(ctx context.Context, jobID *primitive.ObjectID) (*models.TaskResourceUsage, error) {
	finishedTasks, err := t.taskRepository.FindFinishedTask(ctx, jobID)
	if err != nil {
		return nil, err
	}

	if len(finishedTasks) == 0 {
		return nil, errors.New("no finished task hence no resource usage data")
	}

	if len(finishedTasks) == 1 {
		return &finishedTasks[0].ResourceUsage, nil
	}

	result := finishedTasks[0].ResourceUsage
	for _, finishedTask := range finishedTasks[1:] {
		result.AverageWithOther(finishedTask.ResourceUsage)
	}
	return &result, nil
}

func (t *task) UpdateTaskToBeReadyToBeDistributed(ctx context.Context, jobID *primitive.ObjectID) error {
	tasks, err := t.taskRepository.FindManyByJobID(ctx, jobID)
	if err != nil {
		return err
	}

	for _, task := range tasks {
		if task.Status == models.CreatedTaskStatus {
			task.DoneExperimentTask()
		}
	}

	return t.taskRepository.BulkWriteStatusAndLogByID(ctx, tasks)
}

func (t *task) CountUnfinishedTaskByJobID(ctx context.Context, jobID *primitive.ObjectID) (int64, error) {
	return t.taskRepository.CountUnfinishedTaskByJobID(ctx, jobID)
}

func (t *task) UpdateTaskWaitTimeout(ctx context.Context, taskID primitive.ObjectID) error {
	task, err := t.GetTaskByID(ctx, taskID)
	if err != nil {
		return err
	}

	task.CPWaitTimeout()
	return t.taskRepository.WriteTaskResult(ctx, *task)
}
