package models

import (
	"context"
	"sync"
)

type TaskWithContext struct {
	mutex                sync.Mutex
	IsReportedBackToCP   bool
	Task                 Task
	Ctx                  context.Context
	CancelFunc           func()
	ContainerId          string
	AverageResourceUsage AverageResourceUsage
}

func (t *TaskWithContext) SetTaskReportedBackToCP() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.IsReportedBackToCP = true
}

func (t *TaskWithContext) GetTaskReportStatus() bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	return t.IsReportedBackToCP
}
