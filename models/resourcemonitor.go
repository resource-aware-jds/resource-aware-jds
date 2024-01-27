package models

import (
	"github.com/nabhan-au/dockerstats"
)

type ContainerResourceUsage struct {
	ContainerId string
	CpuUsage    string
	MemoryUsage dockerstats.MemoryStats
}

type MemoryUsage struct {
	Usage float64
	Limit float64
}
