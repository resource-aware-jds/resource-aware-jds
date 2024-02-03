package requestmodel

import "github.com/resource-aware-jds/resource-aware-jds/models"

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
	Tasks []models.Task `json:"tasks"`
}
