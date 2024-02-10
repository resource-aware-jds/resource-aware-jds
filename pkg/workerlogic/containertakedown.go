package workerlogic

import (
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/container"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/datastructure"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/util"
	"github.com/sirupsen/logrus"
	"slices"
	"sort"
	"strings"
)

type ContainerTakeDownState struct {
	ContainerBuffer datastructure.Buffer[string, container.IContainer]
	Report          *models.CheckResourceReport // nolint:unused
}

type ContainerTakeDown interface {
	Calculate(state ContainerTakeDownState) []container.IContainer
}

func ProvideOverResourceUsageContainerTakeDown() ContainerTakeDown {
	return &OverResourceUsageContainerTakeDown{}
}

type OverResourceUsageContainerTakeDown struct{}

// Calculate check the container which one should be shutdown. Don't call stop container in this function,
// Instead, Just return it and let the caller handle it instead
func (o OverResourceUsageContainerTakeDown) Calculate(state ContainerTakeDownState) (removeContainer []container.IContainer) {
	containerResourceUsage := state.Report.ContainerResourceUsages
	containerIdRemoveList := make([]string, len(containerResourceUsage))
	cpuExceed := state.Report.CpuUsageExceed
	memoryExceed := state.Report.MemoryUsageExceed

	// Sort by CPU
	sort.Slice(containerResourceUsage, func(i, j int) bool {
		return containerResourceUsage[i].CpuUsage < containerResourceUsage[j].CpuUsage
	})

	for i := 0; i < len(containerResourceUsage); i++ {
		containerUsage := containerResourceUsage[i]
		if cpuExceed <= 0 {
			break
		}
		memoryUsage := util.ConvertToGb(util.ExtractMemoryUsage(containerUsage.CpuUsage))
		cpuUsage, err := util.ExtractCpuUsage(containerUsage)
		if err != nil {
			return nil
		}
		cpuExceed -= cpuUsage
		memoryExceed = util.SubtractInGb(memoryExceed, memoryUsage)
		containerIdRemoveList = append(containerIdRemoveList, containerUsage.ContainerIdShort)
	}

	// Sort by Memory
	sort.Slice(containerResourceUsage, func(i, j int) bool {
		first := util.ConvertToGb(util.ExtractMemoryUsage(containerResourceUsage[i].MemoryUsage.Raw))
		second := util.ConvertToGb(util.ExtractMemoryUsage(containerResourceUsage[j].MemoryUsage.Raw))
		return first.Size < second.Size
	})

	for i := 0; i < len(containerResourceUsage); i++ {
		containerUsage := containerResourceUsage[i]
		if memoryExceed.Size <= 0 {
			break
		}
		if slices.Contains(containerIdRemoveList, containerUsage.ContainerIdShort) {
			continue
		}
		memoryUsage := util.ConvertToGb(util.ExtractMemoryUsage(containerUsage.CpuUsage))
		cpuUsage, err := util.ExtractCpuUsage(containerUsage)
		if err != nil {
			return nil
		}
		cpuExceed -= cpuUsage
		memoryExceed = util.SubtractInGb(memoryExceed, memoryUsage)
		containerIdRemoveList = append(containerIdRemoveList, containerUsage.ContainerIdShort)
	}

	containerKeys := state.ContainerBuffer.GetKeys()
	removeContainer = datastructure.Map(containerIdRemoveList, func(shortKey string) container.IContainer {
		id := getCompleteId(shortKey, containerKeys)
		containerInstance := state.ContainerBuffer.Get(id)
		if containerInstance == nil {
			logrus.Errorf("Unable to get container with id: %s", shortKey)
			// Return nil to skip item in the slice
			return nil
		}
		return *containerInstance
	})
	return datastructure.Filter(removeContainer, func(iContainer container.IContainer) bool {
		return iContainer != nil
	})
}

func getCompleteId(prefix string, data []string) string {
	for _, str := range data {
		if strings.HasPrefix(str, prefix) {
			return str
		}
	}
	return ""
}
