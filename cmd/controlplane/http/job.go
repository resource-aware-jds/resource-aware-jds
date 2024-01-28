package http

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/resource-aware-jds/resource-aware-jds/cmd/controlplane/http/requestmodel"
	"github.com/resource-aware-jds/resource-aware-jds/service"
	"net/http"
)

type JobHandler struct {
	jobService  service.Job
	taskService service.Task
}

func ProvideJobHandler(jobService service.Job, taskService service.Task) JobHandler {
	return JobHandler{
		jobService:  jobService,
		taskService: taskService,
	}
}

func (j *JobHandler) ListJob(c *gin.Context) {
	results, err := j.jobService.ListJob(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": results,
	})
}

func (j *JobHandler) CreateJob(c *gin.Context) {
	ctx := c.Request.Context()
	var req requestmodel.CreateJobRequest
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("req binding error: %e", err)})
		return
	}

	job, err := j.jobService.CreateJob(ctx, req.Name, req.ImageURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("create job error: %e", err)})
		return
	}

	taskAttributesConverted := make([][]byte, 0, len(req.TaskAttributes))
	for _, taskAttribute := range req.TaskAttributes {
		marshalredTaskAttributes, err := json.Marshal(taskAttribute)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("failed to convert a task to bytes: %e", err),
			})
			return
		}

		taskAttributesConverted = append(taskAttributesConverted, marshalredTaskAttributes)
	}
	_, err = j.taskService.CreateTask(ctx, job, taskAttributesConverted)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("create task error: %e", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": job,
	})
}
