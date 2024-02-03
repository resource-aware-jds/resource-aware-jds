package http

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type httpServer struct {
	engine *gin.Engine
	server *http.Server
	config ServerConfig
}

type Server interface {
	Serve()
	GracefullyShutdown()
	Engine() *gin.Engine
}

type ServerConfig struct {
	Port int
}

func ProvideHttpServer(config ServerConfig) (Server, func()) {
	router := gin.Default()

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: router,
	}

	result := &httpServer{
		engine: router,
		server: &server,
		config: config,
	}

	// Allow CORs Policy
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	return result, func() {
		result.GracefullyShutdown()
	}
}

func (h *httpServer) Serve() {
	go func() {
		h.server.ListenAndServe() //nolint:errcheck
	}()
}

func (h *httpServer) GracefullyShutdown() {
	logrus.Info("Gracefully shutting down the HTTP Server")
	h.server.Shutdown(context.Background()) //nolint:errcheck
}

func (h *httpServer) Engine() *gin.Engine {
	return h.engine
}
