package requestmodel

import (
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateJobRequest struct {
	Name           string                   `json:"name" binding:"required"`
	ImageURL       string                   `json:"imageURL" binding:"required"`
	TaskAttributes []map[string]interface{} `json:"taskAttributes" binding:"required"`
}

func (c *CreateJobRequest) ToJobModel() models.Job {
	return models.Job{
		Name:     c.Name,
		ImageURL: c.ImageURL,
	}
}

type JobDetailResponse struct {
	models.Job
	Tasks []TaskJobDetailResponse `json:"tasks"`
}

type TaskJobDetailResponse struct {
	ID                      *primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Status                  models.TaskStatus   `bson:"task_status" json:"status" `
	ImageUrl                string              `bson:"image_url" json:"imageURL"`
	JobID                   *primitive.ObjectID `bson:"job_id" json:"jobID"`
	LatestDistributedNodeID string              `bson:"latest_distributed_node_id,omitempty" json:"latestDistributedNodeID,omitempty"`
}
