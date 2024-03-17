package models

type DistributeError struct {
	NodeEntry NodeEntry
	Task      Task
	Error     error
}
