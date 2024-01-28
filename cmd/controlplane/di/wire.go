//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/resource-aware-jds/resource-aware-jds/cmd/controlplane/grpc"
	configDI "github.com/resource-aware-jds/resource-aware-jds/config/di"
	daemonDI "github.com/resource-aware-jds/resource-aware-jds/daemon/di"
	certDI "github.com/resource-aware-jds/resource-aware-jds/pkg/cert/di"
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
			grpc.ProvideControlPlaneGRPCHandler,
			repositoryDI.RepositoryWireSet,
			serviceDI.ServiceWireSet,
			daemonDI.DaemonWireSet,
			certDI.ControlPlaneCertWireSet,
			ProvideControlPlaneApp,
		),
	)
}
