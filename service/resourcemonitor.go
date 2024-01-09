package service

import (
	"github.com/sirupsen/logrus"
	"runtime"
)

type ResourceMonitor struct {
	m runtime.MemStats
}

type IResourceMonitor interface {
	GetMemoryUsage()
}

func ProvideResourcesMonitor() IResourceMonitor {
	return &ResourceMonitor{}
}

func (r *ResourceMonitor) GetMemoryUsage() {
	runtime.ReadMemStats(&r.m)
	logrus.Infof("Alloc = %v MiB", bToMb(r.m.Alloc))
	logrus.Infof("\tTotalAlloc = %v MiB", bToMb(r.m.TotalAlloc))
	logrus.Infof("\tSys = %v MiB", bToMb(r.m.Sys))
	logrus.Infof("\tNumGC = %v\n", r.m.NumGC)
}

func (r *ResourceMonitor) getCPUUsage() {

}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
