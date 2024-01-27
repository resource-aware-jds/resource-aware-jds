package models

import (
	"github.com/nabhan-au/dockerstats"
)

type ContainerResourceUsage struct {
	ContainerIdShort string
	CpuUsage         string
	MemoryUsage      dockerstats.MemoryStats
}
