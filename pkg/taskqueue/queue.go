package taskqueue

import (
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/datastructure"
)

type queue struct {
	runnerQueue datastructure.Queue[*models.Task]
}

type Queue interface {
	StoreTask(task *models.Task)
	GetTask() *models.Task
}

func ProvideTaskQueue() Queue {
	return &queue{}
}

func (q *queue) StoreTask(task *models.Task) {
	q.runnerQueue.Push(task)
}

func (q *queue) GetTask() *models.Task {
	data, _ := q.runnerQueue.Pop()
	return *data
}
