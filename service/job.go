package service

import (
	"context"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type job struct {
	jobRepository repository.IJob
}

type Job interface {
	GetJob(ctx context.Context, id primitive.ObjectID) (*models.Job, error)
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
		Name:      name,
		Status:    models.CreatedJobStatus,
		ImageURL:  imageURL,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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

func (j *job) GetJob(ctx context.Context, id primitive.ObjectID) (*models.Job, error) {
	return j.jobRepository.FindOneByDocumentID(ctx, id)
}
