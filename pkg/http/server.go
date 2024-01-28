package http

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

type httpServer struct {
	engine *gin.Engine
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

	server := &httpServer{
		engine: router,
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

	return server, func() {
		server.GracefullyShutdown()
	}
}

func (h *httpServer) Serve() {
	go func() {
		h.engine.Run(fmt.Sprintf(":%d", h.config.Port))
	}()
}

func (h *httpServer) GracefullyShutdown() {
	// TODO: Gracefully shutdown the server
}

func (h *httpServer) Engine() *gin.Engine {
	return h.engine
}
