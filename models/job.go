package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type JobStatus string

const (
	CreatedJobStatus       JobStatus = "created"
	ExperimentingJobStatus JobStatus = "experimenting"
	DistributingJobStatus  JobStatus = "distributing"
	SuccessJobStatus       JobStatus = "success"
)

type Job struct {
	ID               *primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name             string              `bson:"name" json:"name"`
	Status           JobStatus           `bson:"status" json:"status"`
	IsExperiment     bool                `bson:"is_experiment" json:"isExperiment"`
	ImageURL         string              `bson:"image_url" json:"imageURL"`
	DistributorLogic DistributorName     `bson:"distributor_logic" json:"distributorLogic"`
	CreatedAt        time.Time           `bson:"created_at" json:"createdAt"`
	UpdatedAt        time.Time           `bson:"updated_at" json:"updatedAt"`
}
