package workerlogic

import (
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/container"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/datastructure"
)

type ContainerTakeDownState struct {
	ContainerBuffer   datastructure.Buffer[string, container.IContainer]
	Report            *models.CheckResourceReport     // nolint:unused
	ContainerResource []models.ContainerResourceUsage // nolint:unused
}

type ContainerTakeDown interface {
	Calculate(state ContainerTakeDownState) []container.IContainer
}

func ProvideOverResourceUsageContainerTakeDown() ContainerTakeDown {
	return &OverResourceUsageContainerTakeDown{}
}

type OverResourceUsageContainerTakeDown struct{}

// Calculate check the container which one should be shutdown. Don't call stop container in this function,
// Instead, Just return it and let the caller handle it instead
func (o OverResourceUsageContainerTakeDown) Calculate(state ContainerTakeDownState) []container.IContainer {
	//TODO implement me
	panic("implement me")
}
