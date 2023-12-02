//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/resource-aware-jds/resource-aware-jds/cmd/controlplane/handler"
	configdi "github.com/resource-aware-jds/resource-aware-jds/config/di"
	pkgdi "github.com/resource-aware-jds/resource-aware-jds/pkg/di"
)

//go:generate wire

func InitializeApplication() (ControlPlaneApp, func(), error) {
	panic(
		wire.Build(
			configdi.ConfigWireSet,
			pkgdi.PKGWireSet,
			handler.ProvideControlPlaneGRPCHandler,
			ProvideControlPlaneApp,
		),
	)
}
