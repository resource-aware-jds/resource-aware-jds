package http

import (
	httpServer "github.com/resource-aware-jds/resource-aware-jds/pkg/http"
)

type RouterResult bool

func ProvideHTTPRouter(handler Handler, server httpServer.Server) RouterResult {
	job := server.Engine().Group("/job")
	{
		job.GET("/", handler.JobHandler.ListJob)
		job.POST("/", handler.JobHandler.CreateJob)
	}

	return true
}
