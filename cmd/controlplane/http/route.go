package http

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpServer "github.com/resource-aware-jds/resource-aware-jds/pkg/http"
)

type RouterResult bool

func ProvideHTTPRouter(handler Handler, server httpServer.Server) RouterResult {
	server.Engine().GET("/metrics", gin.WrapH(promhttp.Handler()))

	job := server.Engine().Group("/job")
	{
		job.GET("/", handler.httpHandler.ListJob)
		job.POST("/", handler.httpHandler.CreateJob)
		job.GET("/:jobID/detail", handler.httpHandler.GetJobDetail)
	}

	task := server.Engine().Group("/task")
	{
		task.GET("/:taskID/detail", handler.httpHandler.GetSpecificTaskDetail)
	}

	node := server.Engine().Group("/node")
	{
		node.GET("/", handler.nodePoolHandler.GetNodes)
		node.GET("/task-watch", handler.nodePoolHandler.GetTaskWatcher)
	}
	return true
}
