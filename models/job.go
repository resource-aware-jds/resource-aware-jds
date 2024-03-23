package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type JobStatus string

const (
	CreatedJobStatus       JobStatus = "created"
	ExperimentingJobStatus JobStatus = "experimenting"
	DistributingJobStatus  JobStatus = "distributing"
	SuccessJobStatus       JobStatus = "success"
	FailedJobStatus        JobStatus = "failed"
)

type Job struct {
	ID               *primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name             string              `bson:"name" json:"name"`
	Status           JobStatus           `bson:"status" json:"status"`
	IsExperiment     bool                `bson:"is_experiment" json:"isExperiment"`
	ImageURL         string              `bson:"image_url" json:"imageURL"`
	DistributorLogic DistributorName     `bson:"distributor_logic" json:"distributorLogic"`
	CreatedAt        time.Time           `bson:"created_at" json:"createdAt"`
	UpdatedAt        time.Time           `bson:"updated_at" json:"updatedAt"`
	Logs             []Log               `bson:"logs" json:"logs"`
}

func (j *Job) ExperimentingJob() {
	j.Status = ExperimentingJobStatus
	j.AddLog(InfoLogSeverity, "Job Transition to ExperimentingJob", nil)
}

func (j *Job) DistributingJob() {
	j.Status = DistributingJobStatus
	j.AddLog(InfoLogSeverity, "Job Transition to DistributingJob", nil)
}

func (j *Job) FailedJobStatus(message string, err error) {
	j.Status = FailedJobStatus
	j.AddLog(WarnLogSeverity, "Job Transition to FailedJobStatus", map[string]string{
		"error":   err.Error(),
		"message": message,
	})
}

func (j *Job) SuccessJobStatus() {
	j.Status = DistributingJobStatus
	j.AddLog(InfoLogSeverity, "Job Transition to SuccessJob", nil)
}

func (j *Job) AddLog(severity LogSeverity, message string, parameters map[string]string) {
	if j.Logs == nil {
		j.Logs = make([]Log, 0)
	}

	j.Logs = append(j.Logs, Log{
		Severity:   severity,
		Parameters: parameters,
		Message:    message,
		Timestamp:  time.Now(),
	})
}
