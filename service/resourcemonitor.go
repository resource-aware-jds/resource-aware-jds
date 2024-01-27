package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/nabhan-au/dockerstats"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/sirupsen/logrus"
	"io"
	"time"
)

const MEGABYTE_SIZE = 1024 * 1024

type ResourceMonitor struct {
	dockerClient  *client.Client
	workerService IWorker
}

type IResourceMonitor interface {
	GetResourceUsage(ctx context.Context) ([]*models.ContainerResourceUsage, error)
}

func ProvideResourcesMonitor(dockerClient *client.Client, workerService IWorker) IResourceMonitor {
	return &ResourceMonitor{
		dockerClient:  dockerClient,
		workerService: workerService,
	}
}

func (r *ResourceMonitor) GetResourceUsage(ctx context.Context) ([]*models.ContainerResourceUsage, error) {
	//containerIds := r.containerBuffer.GetKeys()
	var containerStatList []*models.ContainerResourceUsage
	//for _, containerId := range containerIds {
	//	previousCpu, currentCpu, memoryStats, err := r.getDockerContainerStat(ctx, containerId)
	//	// TODO Shouldn't this part return error??? Fix later
	//	if err != nil {
	//		return nil, err
	//	}
	//	cpuUsages := extractCpuUsage(previousCpu, currentCpu)
	//	totalMemoryUsage, memoryLimit := extractMemoryUsage(memoryStats)
	//	logrus.Info("Memory usage: ", totalMemoryUsage, ", Limit: ", memoryLimit)
	//	memoryUsage := models.MemoryUsage{
	//		Usage: totalMemoryUsage,
	//		Limit: memoryLimit,
	//	}
	//	containerResourceUsage := models.ContainerResourceUsage{
	//		ContainerId: containerId,
	//		CpuUsage:    cpuUsages,
	//		MemoryUsage: memoryUsage,
	//	}
	//	containerStatList = append(containerStatList, &containerResourceUsage)
	//}
	stats, err := dockerstats.Current()
	if err != nil {
		panic(err)
	}

	for _, s := range stats {
		fmt.Println(s.Container)
		fmt.Println(s.Memory)
		fmt.Println(s.CPU)
	}
	return containerStatList, nil
}

func (r *ResourceMonitor) getDockerContainerStat(ctx context.Context, containerId string) (*types.CPUStats, *types.CPUStats, *types.MemoryStats, error) {
	var (
		previousCPU    uint64
		previousSystem uint64
		u              = make(chan error, 1)
	)

	response, err := r.dockerClient.ContainerStats(ctx, containerId, true)
	if err != nil {
		return nil, nil, nil, err
	}
	defer response.Body.Close()
	if response.Body == nil {
		return nil, nil, nil, fmt.Errorf("unable to get container stats from daemon")
	}
	dec := json.NewDecoder(response.Body)

	go func() {
		for {
			var (
				v                *types.StatsJSON
				memPercent       = 0.0
				cpuPercent       = 0.0 // Only used on Linux
				mem              = 0.0
				memLimit         = 0.0
				memPerc          = 0.0
				pidsStatsCurrent uint64
			)

			if err := dec.Decode(&v); err != nil {
				dec = json.NewDecoder(io.MultiReader(dec.Buffered(), response.Body))
				u <- err
				if err == io.EOF {
					break
				}
				time.Sleep(100 * time.Millisecond)
				continue
			}
			if v.MemoryStats.Limit != 0 {
				memPercent = float64(v.MemoryStats.Usage) / float64(v.MemoryStats.Limit) * 100.0
			}
			previousCPU = v.PreCPUStats.CPUUsage.TotalUsage
			previousSystem = v.PreCPUStats.SystemUsage
			cpuPercent = calculateCPUPercentUnix(previousCPU, previousSystem, v)
			mem = float64(v.MemoryStats.Usage)
			memLimit = float64(v.MemoryStats.Limit)
			memPerc = memPercent
			pidsStatsCurrent = v.PidsStats.Current

			logrus.Info(cpuPercent)
			logrus.Info(mem)
			logrus.Info(memLimit)
			logrus.Info(memPerc)
			logrus.Info(pidsStatsCurrent)
			u <- nil
		}
	}()

	var stats types.Stats
	err = json.NewDecoder(response.Body).Decode(&stats)
	if err != nil {
		return nil, nil, nil, err
	}
	return &stats.PreCPUStats, &stats.CPUStats, &stats.MemoryStats, nil
}

func calculateCPUPercentUnix(previousCPU, previousSystem uint64, v *types.StatsJSON) float64 {
	var (
		cpuPercent = 0.0
		// calculate the change for the cpu usage of the container in between readings
		cpuDelta = float64(v.CPUStats.CPUUsage.TotalUsage) - float64(previousCPU)
		// calculate the change for the entire system between readings
		systemDelta = float64(v.CPUStats.SystemUsage) - float64(previousSystem)
	)

	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * float64(len(v.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	}
	return cpuPercent
}

func extractCpuUsage(previousCpu *types.CPUStats, currentCpu *types.CPUStats) float64 {
	var (
		cpuPercent = 0.0
		// calculate the change for the cpu usage of the container in between readings
		cpuDelta = float64(currentCpu.CPUUsage.TotalUsage) - float64(previousCpu.CPUUsage.TotalUsage)
		// calculate the change for the entire system between readings
		systemDelta = float64(currentCpu.SystemUsage) - float64(previousCpu.CPUUsage.TotalUsage)
	)

	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * float64(5) * 100.0
	}
	logrus.Info("Total: ", currentCpu.CPUUsage.TotalUsage, " prev: ", previousCpu.CPUUsage.TotalUsage)
	logrus.Info("System: ", currentCpu.SystemUsage, " prev: ", previousCpu.CPUUsage.TotalUsage)
	logrus.Info("PerCPU: ", currentCpu.CPUUsage.PercpuUsage)
	logrus.Info("CPU usage: ", cpuPercent)
	return cpuPercent
}

func extractMemoryUsage(memoryStats *types.MemoryStats) (float64, float64) {
	return byteToMegaByte(float64(memoryStats.Usage)), byteToMegaByte(float64(memoryStats.Limit))
}

func byteToMegaByte(byte float64) float64 {
	return byte / MEGABYTE_SIZE
}
