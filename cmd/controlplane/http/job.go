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

	job, err := j.jobService.CreateJob(ctx, req.Name, req.ImageURL, req.IsExperiment, req.DistributionLogic)
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
	_, err = j.taskService.CreateTask(ctx, job, taskAttributesConverted, req.IsExperiment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("create task error: %e", err)})
		return
	}

	if !req.IsExperiment {
		err = j.jobService.UpdateJobStatusToDistributing(ctx, job.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("fail to update job status: %e", err)})
			return
		}
	} else {
		err = j.jobService.UpdateJobStatusToExperimenting(ctx, job.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("fail to update job status: %e", err)})
			return
		}
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

	taskResponse := make([]requestmodel.TaskJobDetailResponse, 0, len(tasks))

	for _, task := range tasks {
		taskResponse = append(taskResponse, requestmodel.TaskJobDetailResponse{
			ID:                      task.ID,
			Status:                  task.Status,
			LatestDistributedNodeID: task.LatestDistributedNodeID,
			JobID:                   task.JobID,
			ImageUrl:                task.ImageUrl,
			CreatedAt:               task.CreatedAt,
			UpdatedAt:               task.UpdatedAt,
		})
	}

	res := requestmodel.JobDetailResponse{
		Job:   *job,
		Tasks: taskResponse,
	}

	c.JSON(http.StatusOK, res)
}

func (j *HttpHandler) GetSpecificTaskDetail(c *gin.Context) {
	ctx := c.Request.Context()
	taskID := c.Param("taskID")

	taskIDParsed, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad task id",
		})
		return
	}

	task, err := j.taskService.GetTaskByID(ctx, taskIDParsed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	taskRes := requestmodel.TaskJobFullDetailResponse{
		TaskJobDetailResponse: requestmodel.TaskJobDetailResponse{
			ID:                      task.ID,
			JobID:                   task.JobID,
			LatestDistributedNodeID: task.LatestDistributedNodeID,
			CreatedAt:               task.CreatedAt,
			UpdatedAt:               task.UpdatedAt,
			Status:                  task.Status,
			ImageUrl:                task.ImageUrl,
		},
		Logs: task.Logs,
	}

	var taskAttributesMap map[string]interface{}
	err = json.Unmarshal(task.TaskAttributes, &taskAttributesMap)
	if err == nil {
		taskRes.TaskAttributes = taskAttributesMap
	}

	c.JSON(http.StatusOK, taskRes)
}
