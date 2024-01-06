package repository

import (
	"context"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	TaskCollectionName = "task"
)

type task struct {
	database   *mongo.Database
	collection *mongo.Collection
}

type ITask interface {
	FindManyByJobID(ctx context.Context, jobID *primitive.ObjectID) ([]models.Task, error)
	InsertMany(ctx context.Context, tasks []models.Task) error
}

func ProvideTask(database *mongo.Database) ITask {
	return &task{
		database:   database,
		collection: database.Collection(TaskCollectionName),
	}
}

func (t *task) InsertMany(ctx context.Context, tasks []models.Task) error {
	iTasksSlice := make([]interface{}, 0, len(tasks))
	for _, element := range tasks {
		iTasksSlice = append(iTasksSlice, element)
	}

	_, err := t.collection.InsertMany(ctx, iTasksSlice)
	return err
}

func (t *task) FindManyByJobID(ctx context.Context, jobID *primitive.ObjectID) ([]models.Task, error) {
	result, err := t.collection.Find(ctx, bson.M{
		"job_id": jobID,
	})
	if err != nil {
		return nil, err
	}

	var resultDecoded []models.Task
	err = result.Decode(&resultDecoded)
	return resultDecoded, err
}
