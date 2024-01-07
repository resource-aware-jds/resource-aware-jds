package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type TaskStatus string

const (
	CreatedTaskStatus     TaskStatus = "created"
	DistributedTaskStatus TaskStatus = "distributed"
)

type Task struct {
	ID             *primitive.ObjectID `bson:"_id,omitempty"`
	Status         TaskStatus          `bson:"task_status"`
	ImageUrl       string              `bson:"image_url"`
	JobID          *primitive.ObjectID `bson:"job_id"`
	TaskAttributes []byte              `bson:"task_attributes"`
	Logs           []TaskLog           `bson:"logs,omitempty"`
}

func (t *Task) DistributionSuccess(nodeID string) {
	t.Status = DistributedTaskStatus
	t.AddLog(InfoLogSeverity, "Distributed to node", map[string]string{
		"nodeID": nodeID,
	})
}

func (t *Task) DistributionFailure(nodeID string, err error) {
	t.AddLog(WarnLogSeverity, "Fail to distribute task to node", map[string]string{
		"nodeID": nodeID,
		"error":  err.Error(),
	})
}

func (t *Task) AddLog(severity LogSeverity, message string, parameters map[string]string) {
	if t.Logs == nil {
		t.Logs = make([]TaskLog, 0)
	}

	t.Logs = append(t.Logs, TaskLog{
		Severity:   severity,
		Parameters: parameters,
		Message:    message,
		Timestamp:  time.Now(),
	})
}

type LogSeverity string

const (
	InfoLogSeverity LogSeverity = "info"
	WarnLogSeverity LogSeverity = "warn"
)

type TaskLog struct {
	Severity   LogSeverity       `bson:"severity"`
	Parameters map[string]string `bson:"parameters"`
	Message    string            `bson:"message"`
	Timestamp  time.Time         `bson:"timestamp"`
}
