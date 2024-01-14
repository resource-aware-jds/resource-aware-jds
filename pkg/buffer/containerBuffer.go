package buffer

import (
	"github.com/resource-aware-jds/resource-aware-jds/service"
)

type ContainerBuffer interface {
	Store(id string, object service.ContainerSvc)
	Pop(id string) *service.ContainerSvc
}

func ProvideContainerBuffer() ContainerBuffer {
	buffer := &Buffer[string, service.ContainerSvc]{}
	buffer.InitializeMap()
	return buffer
}
