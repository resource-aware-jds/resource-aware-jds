package models

import "context"

type TaskWithContext struct {
	Task       Task
	Ctx        context.Context
	CancelFunc func()
}
