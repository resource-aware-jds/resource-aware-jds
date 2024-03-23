package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type TaskStatus string

const (
	CreatedTaskStatus     TaskStatus = "created"
	ReadyToDistribute     TaskStatus = "ready-to-distribute"
	DistributedTaskStatus TaskStatus = "distributed"
	WorkOnTaskFailure     TaskStatus = "work-on-task-failure"
	SuccessTaskStatus     TaskStatus = "success"
)

type Task struct {
	ID                      *primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Status                  TaskStatus          `bson:"task_status" json:"status" `
	ImageUrl                string              `bson:"image_url" json:"imageURL"`
	JobID                   *primitive.ObjectID `bson:"job_id" json:"jobID"`
	TaskAttributes          []byte              `bson:"task_attributes" json:"taskAttributes"`
	LatestDistributedNodeID string              `bson:"latest_distributed_node_id,omitempty" json:"latestDistributedNodeID,omitempty"`
	Logs                    []Log               `bson:"logs,omitempty" json:"logs"`
	CreatedAt               time.Time           `bson:"created_at" json:"createdAt"`
	UpdatedAt               time.Time           `bson:"updated_at" json:"updatedAt"`
	Result                  *[]byte             `bson:"result,omitempty" json:"-"`
	RetryCount              int                 `bson:"retry_count" json:"retryCount"`
	ResourceUsage           TaskResourceUsage   `bson:"resource_usage,omitempty" json:"resourceUsage"`
}

type TaskResourceUsage struct {
	Memory float64 `bson:"memory" json:"memory"`
	CPU    float32 `bson:"cpu" json:"cpu"`
}

func (t *TaskResourceUsage) AverageWithOther(other TaskResourceUsage) {
	t.Memory = (t.Memory + other.Memory) / 2
	t.CPU = (t.CPU + other.CPU) / 2
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

func (t *Task) WorkOnTaskFailure(nodeID string, message string) {
	t.LatestDistributedNodeID = ""
	t.Status = WorkOnTaskFailure
	t.AddLog(ErrorLogSeverity, "Node Report Task Failure", map[string]string{
		"nodeID":     nodeID,
		"errMessage": message,
	})
}

func (t *Task) SuccessTask(nodeID string, result []byte) {
	t.Status = SuccessTaskStatus
	if len(result) != 0 {
		t.Result = &result
	}
	t.AddLog(InfoLogSeverity, "Node Submit Successful Task", map[string]string{
		"nodeID": nodeID,
	})
}

func (t *Task) ExperimentTask() {
	t.Status = ReadyToDistribute
	t.AddLog(InfoLogSeverity, "This task has been selected to be the experiment task", nil)
}

func (t *Task) SkipExperimentTask() {
	t.Status = ReadyToDistribute
	t.AddLog(InfoLogSeverity, "This task has been skipped the experiment phrase", nil)
}

func (t *Task) DoneExperimentTask() {
	t.Status = ReadyToDistribute
	t.AddLog(InfoLogSeverity, "The experimenting has been completed and the rest of the task has been marked as ready to distribute", nil)
}

func (t *Task) CPWaitTimeout() {
	t.Status = WorkOnTaskFailure
	t.AddLog(ErrorLogSeverity, "The distributed worker node didn't response the task success in time", nil)
}

func (t *Task) AddLog(severity LogSeverity, message string, parameters map[string]string) {
	if t.Logs == nil {
		t.Logs = make([]Log, 0)
	}

	t.Logs = append(t.Logs, Log{
		Severity:   severity,
		Parameters: parameters,
		Message:    message,
		Timestamp:  time.Now(),
	})
}

type TaskResponse struct {
	ID string
}
