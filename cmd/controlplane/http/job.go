package http

import (
	"github.com/gin-gonic/gin"
	"github.com/resource-aware-jds/resource-aware-jds/service"
	"net/http"
)

type JobHandler struct {
	jobService service.Job
}

func ProvideJobHandler(jobService service.Job) JobHandler {
	return JobHandler{
		jobService: jobService,
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
