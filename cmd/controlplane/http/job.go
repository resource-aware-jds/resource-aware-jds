package http

import (
	"github.com/gin-gonic/gin"
)

type JobHandler struct {
}

func ProvideJobHandler() JobHandler {
	return JobHandler{}
}

func (j *JobHandler) ListJob(c *gin.Context) {

}
