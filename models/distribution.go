package models

type DistributeError struct {
	NodeEntry NodeEntry
	Task      Task
	Error     error
}

type DistributorName string

const (
	RoundRobinDistributorName    DistributorName = "round_robin"
	ResourceAwareDistributorName DistributorName = "resource_aware"
)
