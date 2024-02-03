package http

import (
	httpServer "github.com/resource-aware-jds/resource-aware-jds/pkg/http"
)

type RouterResult bool

func ProvideHTTPRouter(handler Handler, server httpServer.Server) RouterResult {
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

	return true
}
