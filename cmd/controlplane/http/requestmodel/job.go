package requestmodel

import (
	"github.com/resource-aware-jds/resource-aware-jds/models"
)

type CreateJobRequest struct {
	Name           string                   `json:"name" binding:"required"`
	ImageURL       string                   `json:"imageURL" binding:"required"`
	TaskAttributes []map[string]interface{} `json:"taskAttributes" binding:"required"`
	IsExperiment   bool                     `json:"isExperiment"`
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
