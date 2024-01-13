package buffer

import "github.com/resource-aware-jds/resource-aware-jds/models"

type ContainerBuffer interface {
	Store(id string, object *models.Container)
	Pop(id string) *models.Container
	GetKeys() []string
}

func ProvideContainerBuffer() ContainerBuffer {
	buffer := &Buffer[string, models.Container]{}
	buffer.InitializeMap()
	return buffer
}
