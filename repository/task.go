package repository

import (
	"context"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	TaskCollectionName = "task"
)

type task struct {
	database   *mongo.Database
	collection *mongo.Collection
}

type ITask interface {
	FindOneByID(ctx context.Context, taskID primitive.ObjectID) (*models.Task, error)
	FindManyByJobID(ctx context.Context, jobID *primitive.ObjectID) ([]models.Task, error)
	CountUnfinishedTaskByJobID(ctx context.Context, jobID *primitive.ObjectID) (int64, error)
	InsertMany(ctx context.Context, tasks []models.Task) error
	GetTaskToDistributeForJob(ctx context.Context, jobID *primitive.ObjectID) ([]models.Task, error)
	BulkWriteStatusAndLogByID(ctx context.Context, tasks []models.Task) error
	WriteTaskResult(ctx context.Context, task models.Task) error
	FindFinishedTask(ctx context.Context, jobID *primitive.ObjectID) ([]models.Task, error)
	UpdateTaskStatusByJobID(ctx context.Context, jobID *primitive.ObjectID, status models.Task) error
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
		element.CreatedAt = time.Now()
		element.UpdatedAt = time.Now()
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
	err = result.All(ctx, &resultDecoded)
	return resultDecoded, err
}

func (t *task) GetTaskToDistributeForJob(ctx context.Context, jobID *primitive.ObjectID) ([]models.Task, error) {
	result, err := t.collection.Find(ctx, bson.M{
		"task_status": bson.M{
			"$in": []models.TaskStatus{models.ReadyToDistribute},
		},
		"job_id": jobID,
	})
	if err != nil {
		return nil, err
	}

	var resultDecoded []models.Task
	err = result.All(ctx, &resultDecoded)
	return resultDecoded, err
}

func (t *task) BulkWriteStatusAndLogByID(ctx context.Context, tasks []models.Task) error {
	var operations []mongo.WriteModel
	for _, task := range tasks {
		operation := mongo.NewUpdateOneModel()
		operation.SetFilter(bson.M{
			"_id": task.ID,
		})
		operation.SetUpdate(bson.M{
			"$set": bson.M{
				"task_status":                task.Status,
				"logs":                       task.Logs,
				"latest_distributed_node_id": task.LatestDistributedNodeID,
				"updated_at":                 time.Now(),
			},
		})

		operations = append(operations, operation)
	}

	_, err := t.collection.BulkWrite(ctx, operations)
	return err
}

func (t *task) FindOneByID(ctx context.Context, taskID primitive.ObjectID) (*models.Task, error) {
	result := t.collection.FindOne(ctx, bson.M{
		"_id": taskID,
	})

	if result.Err() != nil {
		return nil, result.Err()
	}

	var taskRes models.Task
	err := result.Decode(&taskRes)
	return &taskRes, err
}

func (t *task) WriteTaskResult(ctx context.Context, task models.Task) error {
	operation := mongo.NewUpdateOneModel()
	operation.SetFilter(bson.M{
		"_id": task.ID,
	})
	operation.SetUpdate(bson.M{
		"$set": bson.M{
			"result":      task.Result,
			"logs":        task.Logs,
			"task_status": task.Status,
			"updated_at":  time.Now(),
		},
	})

	_, err := t.collection.BulkWrite(ctx, []mongo.WriteModel{operation})
	return err
}

func (t *task) FindFinishedTask(ctx context.Context, jobID *primitive.ObjectID) ([]models.Task, error) {
	result, err := t.collection.Find(ctx, bson.M{
		"job_id":      jobID,
		"task_status": models.SuccessTaskStatus,
	})
	if err != nil {
		return nil, err
	}

	var resultDecoded []models.Task
	err = result.All(ctx, &resultDecoded)
	return resultDecoded, err
}

func (t *task) UpdateTaskStatusByJobID(ctx context.Context, jobID *primitive.ObjectID, status models.Task) error {
	_, err := t.collection.UpdateMany(
		ctx,
		bson.M{
			"job_id": jobID,
			"status": bson.M{
				"$ne": models.SuccessTaskStatus,
			},
		},
		bson.M{
			"$set": bson.M{
				"status": models.ReadyToDistribute,
			},
		},
	)

	return err
}

func (t *task) CountUnfinishedTaskByJobID(ctx context.Context, jobID *primitive.ObjectID) (int64, error) {
	return t.collection.CountDocuments(ctx, bson.M{
		"job_id": jobID,
		"status": bson.M{
			"$ne": models.SuccessTaskStatus,
		},
	})
}
