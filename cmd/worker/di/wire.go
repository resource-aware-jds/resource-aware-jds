//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/resource-aware-jds/resource-aware-jds/cmd/worker/handler"
	configDI "github.com/resource-aware-jds/resource-aware-jds/config/di"
	daemonDI "github.com/resource-aware-jds/resource-aware-jds/daemon/di"
	pkgDI "github.com/resource-aware-jds/resource-aware-jds/pkg/di"
	serviceDI "github.com/resource-aware-jds/resource-aware-jds/service/di"
)

//go:generate wire

func InitializeApplication() (WorkerApp, func(), error) {
	panic(
		wire.Build(
			configDI.WorkerConfigWireSet,
			pkgDI.PKGWireSet,
			handler.ProvideWorkerGRPCHandler,
			handler.ProvideWorkerGRPCSocketHandler,
			serviceDI.ServiceWireSet,
			ProvideWorkerApp,
			ProvideControlPlaneGRPCClient,
			daemonDI.DaemonWireSet,
		),
	)
}
