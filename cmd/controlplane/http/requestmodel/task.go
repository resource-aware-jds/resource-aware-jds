package requestmodel

import (
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type TaskJobDetailResponse struct {
	ID                      *primitive.ObjectID `json:"id"`
	Status                  models.TaskStatus   `json:"status" `
	ImageUrl                string              `json:"imageURL"`
	JobID                   *primitive.ObjectID `json:"jobID"`
	LatestDistributedNodeID string              `json:"latestDistributedNodeID,omitempty"`
	CreatedAt               time.Time           `json:"createdAt"`
	UpdatedAt               time.Time           `json:"updatedAt"`
}

type TaskJobFullDetailResponse struct {
	TaskJobDetailResponse
	Logs           []models.TaskLog       `json:"logs"`
	TaskAttributes map[string]interface{} `json:"taskAttributes"`
}
