package taskBuffer

import (
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/sirupsen/logrus"
)

type taskBuffer struct {
	taskMap map[string]*models.Task
}

type TaskBuffer interface {
	Store(task *models.Task)
	Pop(id string) *models.Task
}

func ProvideTaskBuffer() TaskBuffer {
	var buffer taskBuffer
	buffer.taskMap = make(map[string]*models.Task)
	return &buffer
}

func (t *taskBuffer) Store(task *models.Task) {
	logrus.Info("Buffer task: " + task.ID.Hex())
	t.taskMap[task.ID.Hex()] = task
}

func (t *taskBuffer) Pop(id string) *models.Task {
	logrus.Info("Remove task from buffer: " + id)
	task, ok := t.taskMap[id]
	if !ok {
		return nil
	}
	delete(t.taskMap, id)
	return task
}

func (t *taskBuffer) isTaskRunning(id string) bool {
	_, ok := t.taskMap[id]
	return ok
}
