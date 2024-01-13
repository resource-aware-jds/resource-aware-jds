package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/client"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/buffer"
	"github.com/sirupsen/logrus"
	"runtime"
	"time"
)

type ResourceMonitor struct {
	m               runtime.MemStats
	dockerClient    *client.Client
	containerBuffer buffer.ContainerBuffer
}

type IResourceMonitor interface {
	GetMemoryUsage(ctx context.Context)
}

func ProvideResourcesMonitor(dockerClient *client.Client, containerBuffer buffer.ContainerBuffer) IResourceMonitor {
	return &ResourceMonitor{
		dockerClient:    dockerClient,
		containerBuffer: containerBuffer,
	}
}

func (r *ResourceMonitor) GetMemoryUsage(ctx context.Context) {
	//runtime.ReadMemStats(&r.m)
	//logrus.Infof("Alloc = %v MiB", bToMb(r.m.Alloc))
	//logrus.Infof("\tTotalAlloc = %v MiB", bToMb(r.m.TotalAlloc))
	//logrus.Infof("\tSys = %v MiB", bToMb(r.m.Sys))
	//logrus.Infof("\tNumGC = %v\n", r.m.NumGC)
	containerIds := r.containerBuffer.GetKeys()
	for _, containerId := range containerIds {
		r.getDockerContainerStat(ctx, containerId)
	}
}

func (r *ResourceMonitor) getCPUUsage() {

}

type myStruct struct {
	Id       string `json:"id"`
	Read     string `json:"read"`
	Preread  string `json:"preread"`
	CpuStats cpu    `json:"cpu_stats"`
}

type cpu struct {
	Usage cpuUsage `json:"cpu_usage"`
}

type cpuUsage struct {
	Total float64 `json:"total_usage"`
}

func (r *ResourceMonitor) getDockerContainerStat(ctx context.Context, containerId string) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	stats, e := r.dockerClient.ContainerStats(ctxWithTimeout, containerId, true)
	if stats.Body == nil {
		return
	}
	if e != nil {
		fmt.Errorf("%s", e.Error())
	}
	for {
		select {
		case <-ctxWithTimeout.Done():
			stats.Body.Close()
			fmt.Println("Stop logging")
			return
		default:
			var containerStats map[string]interface{}
			logrus.Infof(stats.OSType)
			logrus.Info(stats.Body)
			err := json.NewDecoder(stats.Body).Decode(&containerStats)
			if err != nil {
				cancel()
			}
			logrus.Info("==================================")
			logrus.Info(containerStats)
		}
	}
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
