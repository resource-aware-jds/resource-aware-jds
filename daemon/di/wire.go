package di

import (
	"github.com/google/wire"
	"github.com/resource-aware-jds/resource-aware-jds/daemon"
)

var (
	DaemonWireSet = wire.NewSet(
		daemon.ProvideControlPlaneDaemon,
		daemon.ProvideWorkerNodeDaemon,
	)
)
