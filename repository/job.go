package repository

import (
	"context"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	JobCollectionName = "job"
)

type job struct {
	database   *mongo.Database
	collection *mongo.Collection
}

type IJob interface {
	Insert(ctx context.Context, job models.Job) (insertedJobID *primitive.ObjectID, err error)
	FindAll(ctx context.Context) ([]models.Job, error)
	FindOneByDocumentID(ctx context.Context, id primitive.ObjectID) (*models.Job, error)
	FindJobToDistribute(ctx context.Context) ([]models.Job, error)
}

func ProvideJob(database *mongo.Database) IJob {
	return &job{
		database:   database,
		collection: database.Collection(JobCollectionName),
	}
}

func (j *job) Insert(ctx context.Context, job models.Job) (insertedJobID *primitive.ObjectID, err error) {
	result, err := j.collection.InsertOne(ctx, job)
	if err != nil {
		return nil, err
	}

	objID := result.InsertedID.(primitive.ObjectID)
	return &objID, nil
}

func (j *job) FindAll(ctx context.Context) ([]models.Job, error) {
	var result []models.Job
	data, err := j.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	err = data.All(ctx, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (j *job) FindOneByDocumentID(ctx context.Context, id primitive.ObjectID) (*models.Job, error) {
	result := j.collection.FindOne(ctx, bson.M{
		"_id": id,
	})
	if result.Err() != nil {
		return nil, result.Err()
	}

	var jobResult models.Job
	err := result.Decode(&jobResult)
	return &jobResult, err
}

func (j *job) FindJobToDistribute(ctx context.Context) ([]models.Job, error) {
	result, err := j.collection.Find(ctx, bson.M{
		"status": bson.M{
			"$in": []models.JobStatus{
				models.DistributingJobStatus,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	var res []models.Job
	err = result.All(ctx, &res)
	return res, err
}
