package models

import (
	"github.com/resource-aware-jds/dockerstats"
)

type ContainerResourceUsage struct {
	ContainerIdShort string
	CpuUsage         string
	MemoryUsage      dockerstats.MemoryStats
}
