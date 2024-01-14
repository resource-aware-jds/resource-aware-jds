package buffer

import (
	"github.com/resource-aware-jds/resource-aware-jds/models"
)

type TaskBuffer interface {
	Store(id string, object models.Task)
	Pop(id string) *models.Task
}

func ProvideTaskBuffer() TaskBuffer {
	buffer := &Buffer[string, models.Task]{}
	buffer.InitializeMap()
	return buffer
}
