package service

import (
	"context"
	"github.com/resource-aware-jds/resource-aware-jds/config"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/timeutil"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/util"
	"github.com/sirupsen/logrus"
)

const (
	CheckResourceTimeBuffer = 3
)

type DynamicScalingService struct {
	containerService          IContainer
	resourceMonitoringService IResourceMonitor
	config                    config.WorkerConfigModel
}

type IDynamicScaling interface {
	CheckResourceUsageLimit(ctx context.Context) (*models.CheckResourceReport, error)
	CheckResourceUsageLimitWithTimeBuffer(ctx context.Context) (*models.CheckResourceReport, error)
}

func ProvideDynamicScaling(containerService IContainer, resourceMonitoringService IResourceMonitor, config config.WorkerConfigModel) IDynamicScaling {
	return &DynamicScalingService{
		containerService,
		resourceMonitoringService,
		config,
	}
}

func (d *DynamicScalingService) CheckResourceUsageLimitWithTimeBuffer(ctx context.Context) (*models.CheckResourceReport, error) {
	result, err := d.CheckResourceUsageLimit(ctx)
	if err != nil {
		return nil, err
	}
	if result.MemoryUsageExceed.Size == 0 && result.CpuUsageExceed == 0 {
		return result, nil
	}
	timeutil.SleepWithContext(ctx, CheckResourceTimeBuffer)
	result, err = d.CheckResourceUsageLimit(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (d *DynamicScalingService) CheckResourceUsageLimit(ctx context.Context) (*models.CheckResourceReport, error) {
	//Read require configuration
	memoryLimit := d.config.MaxMemoryUsage
	memoryBuffer := d.config.MemoryBufferSize
	cpuLimit := d.config.MaxCpuUsagePercentage
	cpuBuffer := d.config.CpuBufferSize
	dockerCoreLimit := d.config.DockerCoreLimit

	//Read current resource usage
	systemMemoryUsage, err := d.resourceMonitoringService.GetSystemMemUsage()
	if err != nil {
		logrus.Errorf("Unable to retrieve system memory usage: %e", err)
		return nil, err
	}
	systemCpuUsage, err := d.resourceMonitoringService.GetSystemCpuUsage(ctx)
	if err != nil {
		logrus.Errorf("Unable to retrieve system cpu usage: %e", err)
		return nil, err
	}
	containerResourceUsage, err := d.resourceMonitoringService.GetResourceUsage()
	if err != nil {
		logrus.Errorf("Unable to retrieve container resource usage: %e", err)
		return nil, err
	}

	memoryUsage, cpuUsage, err := calculateContainerResourceUsage(containerResourceUsage)
	if err != nil {
		return nil, err
	}

	// Check resource upper bound
	upperBoundReport := models.CheckResourceReport{
		CpuUsageExceed:    0,
		MemoryUsageExceed: models.MemorySize{Size: 0, Unit: "GiB"},
	}
	d.checkCpuUpperBound(cpuUsage, dockerCoreLimit, cpuLimit, &upperBoundReport)
	d.checkMemoryUpperBound(memoryUsage, memoryLimit, &upperBoundReport)

	// Check buffer size
	bufferReport := models.CheckResourceReport{
		CpuUsageExceed:    0,
		MemoryUsageExceed: models.MemorySize{Size: 0, Unit: "GiB"},
	}
	d.checkCpuBuffer(systemCpuUsage, cpuBuffer, &bufferReport)
	d.checkMemoryBuffer(systemMemoryUsage, memoryBuffer, &bufferReport)

	return &models.CheckResourceReport{
		CpuUsageExceed: max(upperBoundReport.CpuUsageExceed, bufferReport.CpuUsageExceed),
		MemoryUsageExceed: models.MemorySize{
			Size: max(upperBoundReport.MemoryUsageExceed.Size, bufferReport.MemoryUsageExceed.Size),
			Unit: "GiB",
		},
		ContainerResourceUsages: containerResourceUsage,
	}, nil
}

func (d *DynamicScalingService) checkMemoryBuffer(systemMemoryUsage *models.MemoryUsage, memoryBuffer string, report *models.CheckResourceReport) {
	freeMemory := systemMemoryUsage.Total - systemMemoryUsage.Used
	memoryBufferGb := util.ConvertToGb(util.ExtractMemoryUsageString(memoryBuffer)).Size
	freeMemoryGb := float64(freeMemory) / (1024 * 1024 * 1024)
	if freeMemoryGb < memoryBufferGb {
		report.MemoryUsageExceed = util.SumInGb(
			report.MemoryUsageExceed,
			models.MemorySize{
				Size: report.MemoryUsageExceed.Size + float64(memoryBufferGb-freeMemoryGb),
				Unit: "GiB",
			})
	}
}

func (d *DynamicScalingService) checkCpuBuffer(systemCpuUsage *models.CpuUsage, cpuBuffer int, report *models.CheckResourceReport) {
	if systemCpuUsage.Idle < float64(cpuBuffer) {
		report.CpuUsageExceed += float64(cpuBuffer) - systemCpuUsage.Idle
	}
}

func (d *DynamicScalingService) checkCpuUpperBound(cpuUsage float64, dockerCoreLimit int, cpuLimit int, report *models.CheckResourceReport) {
	currentCpuPercentage := cpuUsage / float64(dockerCoreLimit)
	if currentCpuPercentage > float64(cpuLimit) {
		cpuDelta := currentCpuPercentage - float64(cpuLimit)
		if cpuDelta > 0 {
			report.CpuUsageExceed += cpuDelta
		}
	}
}

func (d *DynamicScalingService) checkMemoryUpperBound(memoryUsage models.MemorySize, memoryLimit string, report *models.CheckResourceReport) {
	currentMemoryUsageGb := util.ConvertToGb(memoryUsage).Size
	memoryLimitGb := util.ConvertToGb(util.ExtractMemoryUsageString(memoryLimit)).Size
	if currentMemoryUsageGb > memoryLimitGb {
		memoryDelta := currentMemoryUsageGb - memoryLimitGb
		if memoryDelta > 0 {
			report.MemoryUsageExceed = util.SumInGb(
				report.MemoryUsageExceed,
				models.MemorySize{Size: report.MemoryUsageExceed.Size + memoryDelta, Unit: "GiB"},
			)
		}
	}
}
