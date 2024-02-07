package http

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpServer "github.com/resource-aware-jds/resource-aware-jds/pkg/http"
)

type RouterResult bool

func ProvideHTTPRouter(server httpServer.Server) RouterResult {
	server.Engine().GET("/metrics", gin.WrapH(promhttp.Handler()))

	return true
}
