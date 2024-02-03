package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/resource-aware-jds/resource-aware-jds/cmd/controlplane/http/requestmodel"
	"github.com/resource-aware-jds/resource-aware-jds/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type HttpHandler struct {
	jobService  service.Job
	taskService service.Task
}

func ProvideHTTPHandler(jobService service.Job, taskService service.Task) HttpHandler {
	return HttpHandler{
		jobService:  jobService,
		taskService: taskService,
	}
}

func (j *HttpHandler) ListJob(c *gin.Context) {
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

func (j *HttpHandler) CreateJob(c *gin.Context) {
	ctx := c.Request.Context()
	var req requestmodel.CreateJobRequest
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("req binding error: %v", err)})
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

func (j *HttpHandler) GetJobDetail(c *gin.Context) {
	ctx := c.Request.Context()
	jobID := c.Param("jobID")
	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid jobID",
		})
		return
	}

	jobIDParsed, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid jobID",
		})
		return
	}

	job, err := j.jobService.GetJob(ctx, jobIDParsed)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "job not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("can't get job: %v", err),
		})
		return
	}

	tasks, err := j.taskService.GetTaskByJob(ctx, job)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("can't get job: %v", err),
		})
		return
	}

	res := requestmodel.JobDetailResponse{
		Job:   *job,
		Tasks: tasks,
	}

	c.JSON(http.StatusOK, res)
}
