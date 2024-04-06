package service

import (
	"context"
	"github.com/docker/docker/client"
	"github.com/resource-aware-jds/dockerstats"
	"github.com/resource-aware-jds/go-osstat/cpu"
	"github.com/resource-aware-jds/go-osstat/memory"
	"github.com/resource-aware-jds/resource-aware-jds/config"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/datastructure"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/timeutil"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/util"
	"github.com/sirupsen/logrus"
	"time"
)

type ResourceMonitor struct {
	dockerClient     *client.Client
	containerService IContainer
	config           config.WorkerConfigModel
}

type IResourceMonitor interface {
	GetResourceUsage() ([]models.ContainerResourceUsage, error)
	GetSystemMemUsage() (*models.MemoryUsage, error)
	GetSystemCpuUsage(ctx context.Context) (*models.CpuUsage, error)
	CalculateAvailableResource(ctx context.Context) (*models.AvailableResource, error)
}

func ProvideResourcesMonitor(dockerClient *client.Client, workerService IContainer, config config.WorkerConfigModel) IResourceMonitor {
	return &ResourceMonitor{
		dockerClient:     dockerClient,
		containerService: workerService,
		config:           config,
	}
}

func (r *ResourceMonitor) GetResourceUsage() ([]models.ContainerResourceUsage, error) {
	var containerStatList []models.ContainerResourceUsage
	containerKeys := r.containerService.GetContainerIdShort()
	stats, err := dockerstats.Current()
	if err != nil {
		logrus.Errorf("Unable to collect docker stats: %e", err)
	}
	for _, s := range stats {
		if datastructure.Contains(containerKeys, s.Container) {
			containerStatList = append(
				containerStatList,
				models.ContainerResourceUsage{
					ContainerIdShort: s.Container,
					CpuUsage:         s.CPU,
					MemoryUsage:      s.Memory,
				})
		}
	}
	return containerStatList, nil
}

func (r *ResourceMonitor) GetSystemMemUsage() (*models.MemoryUsage, error) {
	memUsage, err := memory.Get()
	if err != nil {
		logrus.Errorf("Unable to get os memory usage: %e", err)
		return nil, err
	}
	return &models.MemoryUsage{
		Total:  memUsage.Total,
		Used:   memUsage.Used,
		Cached: memUsage.Cached,
		Free:   memUsage.Free,
	}, nil
}

func (r *ResourceMonitor) GetSystemCpuUsage(ctx context.Context) (*models.CpuUsage, error) {
	before, err := cpu.Get()
	if err != nil {
		logrus.Errorf("Unable to get os cpu usage: %e", err)
		return nil, err
	}
	timeutil.SleepWithContext(ctx, time.Duration(1)*time.Second)
	after, err := cpu.Get()
	if err != nil {
		logrus.Errorf("Unable to get os cpu usage: %e", err)
		return nil, err
	}
	total := float64(after.Total - before.Total)
	user := float64(after.User-before.User) / total * 100
	system := float64(after.System-before.System) / total * 100
	idle := float64(after.Idle-before.Idle) / total * 100
	return &models.CpuUsage{
		User:   user,
		System: system,
		Idle:   idle,
	}, nil
}

func (r *ResourceMonitor) CalculateAvailableResource(ctx context.Context) (*models.AvailableResource, error) {
	//Read require configuration
	memoryLimit := r.config.MaxMemoryUsage
	memoryBuffer := r.config.MemoryBufferSize
	cpuLimit := r.config.MaxCpuUsagePercentage
	cpuBuffer := r.config.CpuBufferSize
	dockerCoreLimit := r.config.DockerCoreLimit

	//Read current resource usage
	systemMemoryUsage, err := r.GetSystemMemUsage()
	if err != nil {
		logrus.Errorf("Unable to retrieve system memory usage: %e", err)
		return nil, err
	}
	systemCpuUsage, err := r.GetSystemCpuUsage(ctx)
	if err != nil {
		logrus.Errorf("Unable to retrieve system cpu usage: %e", err)
		return nil, err
	}
	containerResourceUsage, err := r.GetResourceUsage()
	if err != nil {
		logrus.Errorf("Unable to retrieve container resource usage: %e", err)
		return nil, err
	}

	memoryUsage, cpuUsage, err := calculateContainerResourceUsage(containerResourceUsage)
	if err != nil {
		return nil, err
	}

	upperBoundAvailableResourceReport := models.AvailableResource{}

	r.checkCpuUpperBound(cpuUsage, dockerCoreLimit, cpuLimit, &upperBoundAvailableResourceReport)
	r.checkMemoryUpperBound(memoryUsage, memoryLimit, &upperBoundAvailableResourceReport)

	// Check buffer size
	bufferAvailableResourceReport := models.AvailableResource{}

	r.checkCpuBuffer(systemCpuUsage, cpuBuffer, &bufferAvailableResourceReport)
	r.checkMemoryBuffer(systemMemoryUsage, memoryBuffer, &bufferAvailableResourceReport)

	return &models.AvailableResource{
		CpuCores:               int64(dockerCoreLimit),
		AvailableCpuPercentage: min(upperBoundAvailableResourceReport.AvailableCpuPercentage, bufferAvailableResourceReport.AvailableCpuPercentage),
		AvailableMemory: models.MemorySize{
			Size: min(upperBoundAvailableResourceReport.AvailableMemory.Size, bufferAvailableResourceReport.AvailableMemory.Size),
			Unit: "Gib",
		},
	}, nil
}

func (r *ResourceMonitor) checkMemoryBuffer(systemMemoryUsage *models.MemoryUsage, memoryBuffer string, report *models.AvailableResource) {
	freeMemory := systemMemoryUsage.Total - systemMemoryUsage.Used
	memoryBufferGb := util.ConvertToGb(util.ExtractMemoryUsageString(memoryBuffer)).Size
	freeMemoryGb := float64(freeMemory) / (1024 * 1024 * 1024)
	memoryDelta := freeMemoryGb - memoryBufferGb
	if memoryDelta > 0 {
		report.AvailableMemory = util.SumInGb(
			report.AvailableMemory,
			models.MemorySize{
				Size: report.AvailableMemory.Size + (memoryDelta),
				Unit: "GiB",
			})
	}

}

func (r *ResourceMonitor) checkCpuBuffer(systemCpuUsage *models.CpuUsage, cpuBuffer int, report *models.AvailableResource) {
	cpuDelta := float32(systemCpuUsage.Idle) - float32(cpuBuffer)
	if cpuDelta > 0 {
		report.AvailableCpuPercentage += cpuDelta
	}

}

func (r *ResourceMonitor) checkCpuUpperBound(cpuUsage float64, dockerCoreLimit int, cpuLimit int, report *models.AvailableResource) {
	currentCpuPercentage := cpuUsage / float64(dockerCoreLimit)
	cpuDelta := float64(cpuLimit) - currentCpuPercentage
	if cpuDelta > 0 {
		report.AvailableCpuPercentage += float32(cpuDelta)
	}
}

func (r *ResourceMonitor) checkMemoryUpperBound(memoryUsage models.MemorySize, memoryLimit string, report *models.AvailableResource) {
	currentMemoryUsageGb := util.ConvertToGb(memoryUsage).Size
	memoryLimitGb := util.ConvertToGb(util.ExtractMemoryUsageString(memoryLimit)).Size
	memoryDelta := memoryLimitGb - currentMemoryUsageGb
	if memoryDelta > 0 {
		report.AvailableMemory = util.SumInGb(
			report.AvailableMemory,
			models.MemorySize{Size: report.AvailableMemory.Size + memoryDelta, Unit: "GiB"},
		)
	}

}

func calculateContainerResourceUsage(containerResourceUsage []models.ContainerResourceUsage) (models.MemorySize, float64, error) {
	//Calculate all container usage
	memoryUsageSlice := datastructure.Map(containerResourceUsage,
		func(containerUsage models.ContainerResourceUsage) models.MemorySize {
			return util.ExtractMemoryUsageFromModel(containerUsage)
		})

	memoryUsage := datastructure.SumAny(memoryUsageSlice, util.SumInGb, models.MemorySize{Size: 0, Unit: "GiB"})

	cpuUsage := datastructure.SumFloat(containerResourceUsage, func(containerUsage models.ContainerResourceUsage) float64 {
		percentageFloat, err := util.ExtractCpuUsage(containerUsage)
		if err != nil {
			logrus.Errorf("Unabel to extract percentage: %e", err)
		}
		return percentageFloat
	})
	return memoryUsage, cpuUsage, nil
}
