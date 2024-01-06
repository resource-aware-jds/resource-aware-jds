package di

import (
	"github.com/google/wire"
	"github.com/resource-aware-jds/resource-aware-jds/repository"
)

var RepositoryWireSet = wire.NewSet(
	repository.ProvideControlPlane,
	repository.ProvideJob,
	repository.ProvideTask,
)
