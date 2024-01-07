//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/resource-aware-jds/resource-aware-jds/cmd/controlplane/handler"
	configDI "github.com/resource-aware-jds/resource-aware-jds/config/di"
	daemonDI "github.com/resource-aware-jds/resource-aware-jds/daemon/di"
	pkgDI "github.com/resource-aware-jds/resource-aware-jds/pkg/di"
	repositoryDI "github.com/resource-aware-jds/resource-aware-jds/repository/di"
	serviceDI "github.com/resource-aware-jds/resource-aware-jds/service/di"
)

//go:generate wire

func InitializeApplication() (ControlPlaneApp, func(), error) {
	panic(
		wire.Build(
			configDI.ControlPlaneConfigWireSet,
			pkgDI.PKGWireSet,
			handler.ProvideControlPlaneGRPCHandler,
			repositoryDI.RepositoryWireSet,
			serviceDI.ServiceWireSet,
			daemonDI.DaemonWireSet,
			ProvideControlPlaneApp,
		),
	)
}
