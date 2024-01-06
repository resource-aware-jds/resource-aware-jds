package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type TaskStatus string

const (
	CreatedTaskStatus TaskStatus = "created"
)

type Task struct {
	ID             *primitive.ObjectID `bson:"_id,omitempty"`
	Status         TaskStatus          `bson:"task_status"`
	ImageUrl       string              `bson:"image_url"`
	JobID          *primitive.ObjectID `bson:"job_id"`
	TaskAttributes []byte              `bson:"task_attributes"`
}
