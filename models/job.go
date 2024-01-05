package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type JobStatus string

const (
	PendingJobStatus      JobStatus = "pending"
	DistributingJobStatus JobStatus = "distributing"
	SuccessJobStatus      JobStatus = "success"
)

type Job struct {
	ID         *primitive.ObjectID `bson:"_id,omitempty"`
	Status     JobStatus           `bson:"status"`
	Parameters map[string]any      `bson:"parameters"`
	Tasks      []map[string]any    `bson:"distribution_parameters"`
	ImageURL   string              `bson:"image_url"`
	Logs       []JobLog            `bson:"logs"`
}

type JobLogType string

const (
	CreateJobLogType  JobLogType = "create"
	DistributedToNode JobLogType = "distributed_to_node"
)

type JobLog struct {
	ID     *primitive.ObjectID `bson:"_id,omitempty"`
	Type   JobLogType          `bson:"type"`
	NodeID string              `bson:"node_id,omitempty"`
}
