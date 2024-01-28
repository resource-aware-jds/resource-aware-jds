package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type JobStatus string

const (
	PendingJobStatus      JobStatus = "pending"
	DistributingJobStatus JobStatus = "distributing"
	SuccessJobStatus      JobStatus = "success"
)

type Job struct {
	ID       *primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Status   JobStatus           `bson:"status" json:"status"`
	ImageURL string              `bson:"image_url" json:"imageURL"`
}
