package service

import (
	"context"
	"fmt"
	"github.com/resource-aware-jds/resource-aware-jds/config"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/datastructure"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/util"
	"github.com/sirupsen/logrus"
	"regexp"
	"strconv"
	"strings"
)

type DynamicScalingService struct {
	containerService          IContainer
	resourceMonitoringService IResourceMonitor
	config                    config.WorkerConfigModel
}

type IDynamicScaling interface {
	CheckResourceUsageLimit(ctx context.Context) models.CheckResourceReport
}

func ProvideDynamicScaling(containerService IContainer, resourceMonitoringService IResourceMonitor, config config.WorkerConfigModel) IDynamicScaling {
	return &DynamicScalingService{
		containerService,
		resourceMonitoringService,
		config,
	}
}

func (d *DynamicScalingService) CheckResourceUsageLimit(ctx context.Context) models.CheckResourceReport {
	//Read require configuration
	fmt.Println("Get config")
	memoryLimit := d.config.MaxMemoryUsage
	memoryBuffer := d.config.MemoryBufferSize
	cpuLimit := d.config.MaxCpuUsagePercentage
	cpuBuffer := d.config.CpuBufferSize
	dockerCoreLimit := d.config.DockerCoreLimit

	//Read current resource usage
	fmt.Println("Get resource usage")
	containerResourceUsage, err := d.resourceMonitoringService.GetResourceUsage()
	if err != nil {
		logrus.Errorf("Unable to retrieve container resource usage: %e", err)
	}
	systemMemoryUsage, err := d.resourceMonitoringService.GetSystemMemUsage()
	if err != nil {
		logrus.Errorf("Unable to retrieve system memory usage: %e", err)
	}
	systemCpuUsage, err := d.resourceMonitoringService.GetSystemCpuUsage(ctx)
	if err != nil {
		logrus.Errorf("Unable to retrieve system cpu usage: %e", err)
	}

	//Calculate all container usage
	fmt.Println("Get container usage")
	memoryUsageSlice := datastructure.Map(containerResourceUsage,
		func(containerUsage models.ContainerResourceUsage) models.MemoryWithUnit {
			fmt.Println(d.extractMemoryUsage(containerUsage.MemoryUsage.Raw))
			return d.extractMemoryUsage(containerUsage.MemoryUsage.Raw)
		})

	memoryUsage := datastructure.SumAny(memoryUsageSlice, util.SumInGb, models.MemoryWithUnit{Size: 0, Unit: "GiB"})

	cpuUsage := datastructure.SumFloat(containerResourceUsage, func(containerUsage models.ContainerResourceUsage) float64 {
		trimmedStr := strings.TrimSuffix(containerUsage.CpuUsage, "%")
		percentageFloat, err := strconv.ParseFloat(trimmedStr, 64)

		if err != nil {
			fmt.Printf("There was an error converting the string to a float:  %v\n", err)
			// TODO add error handler
			return 0
		}
		return percentageFloat
	})

	fmt.Printf("Docker usage: ")
	fmt.Println(containerResourceUsage)
	fmt.Printf("Docker memory: ")
	fmt.Println(memoryUsage)

	// Check resource upper bound
	upperBoundReport := models.CheckResourceReport{
		CpuUsageExceed:    0,
		MemoryUsageExceed: models.MemoryWithUnit{Size: 0, Unit: "GiB"},
	}
	d.checkCpuUpperBound(cpuUsage, dockerCoreLimit, cpuLimit, upperBoundReport)
	d.checkMemoryUpperBound(memoryUsage, memoryLimit, upperBoundReport)

	// Check buffer size
	bufferReport := models.CheckResourceReport{
		CpuUsageExceed:    0,
		MemoryUsageExceed: models.MemoryWithUnit{Size: 0, Unit: "GiB"},
	}
	d.checkCpuBuffer(systemCpuUsage, cpuBuffer, bufferReport)
	d.checkMemoryBuffer(systemMemoryUsage, memoryBuffer, bufferReport)

	return models.CheckResourceReport{
		CpuUsageExceed: max(upperBoundReport.CpuUsageExceed, bufferReport.CpuUsageExceed),
		MemoryUsageExceed: models.MemoryWithUnit{
			Size: max(upperBoundReport.MemoryUsageExceed.Size, bufferReport.MemoryUsageExceed.Size),
			Unit: "GiB",
		},
	}
}

func (d *DynamicScalingService) checkMemoryBuffer(systemMemoryUsage *models.MemoryUsage, memoryBuffer string, report models.CheckResourceReport) {
	freeMemory := systemMemoryUsage.Total - systemMemoryUsage.Used
	memoryBufferGb := util.ConvertToGb(d.extractMemoryUsage(memoryBuffer)).Size
	freeMemoryGb := float64(freeMemory) / (1024 * 1024 * 1024)
	if freeMemoryGb < memoryBufferGb {
		report.MemoryUsageExceed = util.SumInGb(
			report.MemoryUsageExceed,
			models.MemoryWithUnit{
				Size: report.MemoryUsageExceed.Size + float64(memoryBufferGb-freeMemoryGb),
				Unit: "GiB",
			})
	}
}

func (d *DynamicScalingService) checkCpuBuffer(systemCpuUsage *models.CpuUsage, cpuBuffer int, report models.CheckResourceReport) {
	if systemCpuUsage.Idle < float64(cpuBuffer) {
		report.CpuUsageExceed += float64(cpuBuffer) - systemCpuUsage.Idle
	}
}

func (d *DynamicScalingService) checkCpuUpperBound(cpuUsage float64, dockerCoreLimit int, cpuLimit int, report models.CheckResourceReport) {
	currentCpuPercentage := cpuUsage / float64(dockerCoreLimit)
	if cpuUsage/float64(dockerCoreLimit) > float64(cpuLimit) {
		cpuDelta := currentCpuPercentage - float64(cpuLimit)
		if cpuDelta > 0 {
			report.CpuUsageExceed += cpuDelta
		}
	}
}

func (d *DynamicScalingService) checkMemoryUpperBound(memoryUsage models.MemoryWithUnit, memoryLimit string, report models.CheckResourceReport) {
	currentMemoryUsageGb := util.ConvertToGb(memoryUsage).Size
	memoryLimitGb := util.ConvertToGb(d.extractMemoryUsage(memoryLimit)).Size
	if currentMemoryUsageGb > memoryLimitGb {
		memoryDelta := currentMemoryUsageGb - memoryLimitGb
		if memoryDelta > 0 {
			report.MemoryUsageExceed = util.SumInGb(
				report.MemoryUsageExceed,
				models.MemoryWithUnit{Size: report.MemoryUsageExceed.Size + memoryDelta, Unit: "GiB"},
			)
		}
	}
}

func (d *DynamicScalingService) extractMemoryUsage(input string) models.MemoryWithUnit {
	regex := regexp.MustCompile(`(\d+(\.\d+)?)([a-zA-Z]+)`)
	match := regex.FindStringSubmatch(input)

	if match != nil {
		number, _ := strconv.ParseFloat(match[1], 64)
		unit := match[3]

		result := models.MemoryWithUnit{
			Size: number,
			Unit: unit,
		}

		return result
	}

	return models.MemoryWithUnit{}
}
