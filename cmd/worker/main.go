package main

import (
	"fmt"
	"github.com/resource-aware-jds/resource-aware-jds/cmd/worker/di"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logrus.Info("Starting up Worker GRPC server.")
	app, cleanup, err := di.InitializeApplication()
	if err != nil {
		if cleanup != nil {
			cleanup()
		}
		panic(fmt.Sprintf("failed to initialize app: %e", err))
	}
  
	app.GRPCServer.Serve()
	app.GRPCSocketServer.Serve()

	// Gracefully Shutdown
	// Make channel listen for signals from OS
	gracefulStop := make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	<-gracefulStop

	logrus.Info("Gracefully shutting down, cleaning up.")
	cleanup()
	logrus.Info("Clean up success. Good Bye")
}
