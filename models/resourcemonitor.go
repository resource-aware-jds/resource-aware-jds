package models

import (
	"github.com/resource-aware-jds/dockerstats"
)

type ContainerResourceUsage struct {
	ContainerIdShort string
	CpuUsage         string
	MemoryUsage      dockerstats.MemoryStats
}

type MemoryUsage struct {
	Total  uint64
	Used   uint64
	Cached uint64
	Free   uint64
}

type CpuUsage struct {
	User   float64
	System float64
	Idle   float64
}

type OSResourceUsage struct {
	MemoryUsage MemoryUsage
	CpuUsage    CpuUsage
}

type MemorySize struct {
	Size float64
	Unit string
}

type AvailableResource struct {
	CpuCores               int64
	AvailableCpuPercentage float32
	AvailableMemory        MemorySize
}
