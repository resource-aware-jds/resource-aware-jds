package models

type CheckResourceReport struct {
	CpuUsageExceed    float64
	MemoryUsageExceed MemorySize
}
