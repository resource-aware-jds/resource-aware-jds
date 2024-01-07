package taskqueue

import (
	"fmt"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/datastructure"
	"github.com/sirupsen/logrus"
)

type queue struct {
	runnerQueue datastructure.Queue[*models.Task]
}

type Queue interface {
	StoreTask(task *models.Task)
	GetTask(imageUrl string) (*models.Task, error)
	ReadQueue() []*models.Task
	GetDistinctImageList() []string
	BulkRemove(tasks []*models.Task)
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

func (q *queue) ReadQueue() []*models.Task {
	return q.runnerQueue.ReadQueue()
}

func (q *queue) GetDistinctImageList() []string {
	taskList := q.ReadQueue()
	var imageList []string

	distinctKeyMap := make(map[string]bool)
	for _, task := range taskList {
		if _, value := distinctKeyMap[task.ImageUrl]; !value {
			distinctKeyMap[task.ImageUrl] = true
			imageList = append(imageList, task.ImageUrl)
		}
	}
	return imageList
}

func (q *queue) BulkRemove(tasks []*models.Task) {
	q.runnerQueue.FilterWithCondition(func(task *models.Task) bool {
		return !datastructure.Contains(tasks, task)
	})
	logrus.Info("Task removed:", tasks)
	logrus.Info("Current queue:", q.runnerQueue)
}
