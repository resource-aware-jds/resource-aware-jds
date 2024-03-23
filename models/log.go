package models

import "time"

type LogSeverity string

const (
	InfoLogSeverity  LogSeverity = "info"
	WarnLogSeverity  LogSeverity = "warn"
	ErrorLogSeverity LogSeverity = "error"
)

type Log struct {
	Severity   LogSeverity       `bson:"severity" json:"severity"`
	Parameters map[string]string `bson:"parameters" json:"parameters"`
	Message    string            `bson:"message" json:"message"`
	Timestamp  time.Time         `bson:"timestamp" json:"timestamp"`
}
