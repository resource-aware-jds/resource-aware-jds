package taskqueue

import (
	"fmt"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/datastructure"
)

type queue struct {
	runnerQueue datastructure.Queue[*models.Task]
}

type Queue interface {
	StoreTask(task *models.Task)
	GetTask(imageUrl string) (*models.Task, error)
}

func ProvideTaskQueue() Queue {
	return &queue{}
}

func (q *queue) StoreTask(task *models.Task) {
	q.runnerQueue.Push(task)
}

func (q *queue) GetTask(imageUrl string) (*models.Task, error) {
	filter := func(t *models.Task) bool {
		return t.ImageUrl == imageUrl
	}
	data, isSuccess := q.runnerQueue.PopWithFilter(filter)
	if !isSuccess {
		return nil, fmt.Errorf("unable to get task, queue empty")
	}
	return *data, nil
}
