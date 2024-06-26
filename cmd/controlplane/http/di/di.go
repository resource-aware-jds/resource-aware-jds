package di

import (
	"github.com/google/wire"
	"github.com/resource-aware-jds/resource-aware-jds/cmd/controlplane/http"
)

var HTTPWireSet = wire.NewSet(
	http.ProvideHTTPRouter,
	http.ProvideHTTPHandler,
	http.ProvideHandler,
	http.ProvideNodeHandler,
)
