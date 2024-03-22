package service

import (
	"context"
	"github.com/resource-aware-jds/resource-aware-jds/config"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/datastructure"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
	"time"
)

type cpTaskWatcher struct {
	config      config.TaskWatcherConfigModel
	mutex       sync.Mutex
	taskService Task
	taskBuffer  datastructure.Buffer[primitive.ObjectID, time.Time]
}

// CPTaskWatcher is a service to watch for a task that haven't heard back from the WorkerNode after being distributed
// for a specified time.
type CPTaskWatcher interface {
	datastructure.Observer[models.TaskEventBus]

	AddTaskToWatch(taskID primitive.ObjectID)
	WatcherLoop(ctx context.Context)
}

func ProvideCPTaskWatcher(taskService Task) CPTaskWatcher {
	return &cpTaskWatcher{
		taskService: taskService,
		taskBuffer:  make(datastructure.Buffer[primitive.ObjectID, time.Time]),
	}
}

func (c *cpTaskWatcher) AddTaskToWatch(taskID primitive.ObjectID) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.taskBuffer[taskID] = time.Now().Add(c.config.Timeout)
}

func (c *cpTaskWatcher) OnEvent(_ context.Context, t models.TaskEventBus) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.taskBuffer, t.TaskID)

	return nil
}

func (c *cpTaskWatcher) WatcherLoop(ctx context.Context) {
	if len(c.taskBuffer) == 0 {
		return
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	taskIDToCallTimeouts := make([]primitive.ObjectID, 0, len(c.taskBuffer))
	for x, deadline := range c.taskBuffer {
		if deadline.After(time.Now()) {
			// Remove task from the watcher and update the status as failed.
			taskIDToCallTimeouts = append(taskIDToCallTimeouts, x)
			delete(c.taskBuffer, x)
		}
	}

	go func() {
		for _, taskIDToCallTimeout := range taskIDToCallTimeouts {
			err := c.taskService.UpdateTaskWaitTimeout(ctx, taskIDToCallTimeout)
			if err != nil {
				logrus.Error("Failed to Update Task Wait Timeout: ", err)

				// Add it back to the check loop and let the next loop handle the task
				c.mutex.Lock()
				c.taskBuffer[taskIDToCallTimeout] = time.Now()
				c.mutex.Unlock()
			}
		}
	}()
}
