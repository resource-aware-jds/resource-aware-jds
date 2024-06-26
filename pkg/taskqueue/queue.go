package taskqueue

import (
	"fmt"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/datastructure"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/metrics"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/metric"
)

type queue struct {
	queue datastructure.Queue[*models.Task]
}

// Queue is a special datastructure.Queue implementation that will not fully remove models.Task
// once popped. Instead, It will still be stored in internal datastructure.Buffer to make sure
// that no task is lost during the task distribution process.
type Queue interface {
	StoreTask(task *models.Task)
	GetTask(imageUrl string) (*models.Task, error)
	ReadQueue() []*models.Task
	GetDistinctImageList() []string
	BulkRemove(tasks []*models.Task)
	Pop() (*models.Task, bool)
	PeakForNextTask() (*models.Task, bool)
	GetSize() int
}

func ProvideTaskQueue(meter metric.Meter) Queue {
	return &queue{
		queue: datastructure.ProvideQueue[*models.Task](
			datastructure.WithQueueMetrics(
				meter,
				metrics.GenerateWorkerNodeMetric("task_queue_size"),
				metric.WithUnit("Task"),
				metric.WithDescription("The total task in queue"),
			),
		),
	}
}

func (q *queue) Pop() (*models.Task, bool) {
	result, ok := q.queue.Pop()
	if result == nil {
		return nil, ok
	}
	return *result, ok
}

func (q *queue) StoreTask(task *models.Task) {
	q.queue.Push(task)
}

func (q *queue) GetTask(imageUrl string) (*models.Task, error) {
	filter := func(t *models.Task) bool {
		return t.ImageUrl == imageUrl
	}
	data, isSuccess := q.queue.PopWithFilter(filter)
	if !isSuccess {
		return nil, fmt.Errorf("unable to get task, queue empty")
	}
	return *data, nil
}

func (q *queue) ReadQueue() []*models.Task {
	return q.queue.ReadQueue()
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
	q.queue.RemoveWithCondition(func(task *models.Task) bool {
		return !datastructure.Contains(tasks, task)
	})
	logrus.Info("Task removed:", tasks)
}

func (q *queue) PeakForNextTask() (*models.Task, bool) {
	if len(q.queue.ReadQueue()) == 0 {
		return nil, false
	}

	return q.queue.ReadQueue()[0], true
}

func (q *queue) GetSize() int {
	return q.queue.Size()
}
