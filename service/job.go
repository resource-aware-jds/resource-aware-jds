package service

import (
	"context"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/repository"
)

type job struct {
	jobRepository repository.IJob
}

type Job interface {
	CreateJob(ctx context.Context, name, imageURL string) (*models.Job, error)
	ListJob(ctx context.Context) ([]models.Job, error)
}

func ProvideJobService(jobRepository repository.IJob) Job {
	return &job{
		jobRepository: jobRepository,
	}
}

func (j *job) CreateJob(ctx context.Context, name, imageURL string) (*models.Job, error) {
	// Create Job
	jobResult := models.Job{
		Name:     name,
		Status:   models.PendingJobStatus,
		ImageURL: imageURL,
	}
	insertedJobID, err := j.jobRepository.Insert(ctx, jobResult)
	if err != nil {
		return nil, err
	}
	jobResult.ID = insertedJobID

	return &jobResult, nil
}

func (j *job) ListJob(ctx context.Context) ([]models.Job, error) {
	return j.jobRepository.FindAll(ctx)
}
