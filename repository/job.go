package repository

import (
	"context"
	"github.com/resource-aware-jds/resource-aware-jds/models"
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
