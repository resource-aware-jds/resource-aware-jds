package service

import (
	"context"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/repository"
)

type job struct {
	jobRepository  repository.IJob
	taskRepository repository.ITask
}

type Job interface {
	CreateJob(ctx context.Context, imageURL string, taskAttributes [][]byte) (*models.Job, []models.Task, error)
}

func ProvideJobService(jobRepository repository.IJob, taskRepository repository.ITask) Job {
	return &job{
		jobRepository:  jobRepository,
		taskRepository: taskRepository,
	}
}

func (j *job) CreateJob(ctx context.Context, imageURL string, taskAttributes [][]byte) (*models.Job, []models.Task, error) {
	// Create Job
	job := models.Job{
		Status:   models.PendingJobStatus,
		ImageURL: imageURL,
	}
	insertedJobID, err := j.jobRepository.Insert(ctx, job)
	if err != nil {
		return nil, nil, err
	}
	job.ID = insertedJobID

	// Create Tasks
	tasks := make([]models.Task, 0, len(taskAttributes))
	for _, taskAttribute := range taskAttributes {
		newTask := models.Task{
			JobID:          insertedJobID,
			Status:         models.CreatedTaskStatus,
			ImageUrl:       imageURL,
			TaskAttributes: taskAttribute,
		}
		tasks = append(tasks, newTask)
	}
	err = j.taskRepository.InsertMany(ctx, tasks)
	if err != nil {
		return nil, nil, err
	}

	tasksResponse, err := j.taskRepository.FindManyByJobID(ctx, insertedJobID)
	if err != nil {
		return nil, nil, err
	}

	return &job, tasksResponse, nil
}
