package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type JobStatus string

const (
	PendingJobStatus      JobStatus = "pending"
	DistributingJobStatus JobStatus = "distributing"
	SuccessJobStatus      JobStatus = "success"
)

type Job struct {
	ID       *primitive.ObjectID `bson:"_id,omitempty"`
	Status   JobStatus           `bson:"status"`
	ImageURL string              `bson:"image_url"`
}
