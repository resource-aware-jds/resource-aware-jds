package models

type ContainerResourceUsage struct {
	ContainerId string
	CpuUsage    float64
	MemoryUsage MemoryUsage
}

type MemoryUsage struct {
	Usage float64
	Limit float64
}
