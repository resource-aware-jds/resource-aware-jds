package di

import (
	"github.com/google/wire"
	"github.com/resource-aware-jds/resource-aware-jds/service"
)

var ServiceWireSet = wire.NewSet(
	service.ProvideControlPlane,
	service.ProvideWorker,
	service.ProvideResourcesMonitor,
	service.ProvideJobService,
	service.ProvideTaskService,
	service.ProvideContainer,
)
