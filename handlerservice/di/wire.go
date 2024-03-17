package di

import (
	"github.com/google/wire"
	"github.com/resource-aware-jds/resource-aware-jds/handlerservice"
)

var HandlerServiceWireSet = wire.NewSet(
	handlerservice.ProvideWorker,
	handlerservice.ProvideControlPlane,
)
