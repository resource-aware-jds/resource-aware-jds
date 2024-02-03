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
	ID                      *primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Status                  TaskStatus          `bson:"task_status" json:"status" `
	ImageUrl                string              `bson:"image_url" json:"imageURL"`
	JobID                   *primitive.ObjectID `bson:"job_id" json:"jobID"`
	TaskAttributes          []byte              `bson:"task_attributes" json:"taskAttributes"`
	LatestDistributedNodeID string              `bson:"latest_distributed_node_id,omitempty" json:"latestDistributedNodeID,omitempty"`
	Logs                    []TaskLog           `bson:"logs,omitempty" json:"logs"`
	CreatedAt               time.Time           `bson:"created_at" json:"createdAt"`
	UpdatedAt               time.Time           `bson:"updated_at" json:"updatedAt"`
}

func (t *Task) DistributionSuccess(nodeID string) {
	t.Status = DistributedTaskStatus
	t.LatestDistributedNodeID = nodeID
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
	Severity   LogSeverity       `bson:"severity" json:"severity"`
	Parameters map[string]string `bson:"parameters" json:"parameters"`
	Message    string            `bson:"message" json:"message"`
	Timestamp  time.Time         `bson:"timestamp" json:"timestamp"`
}
