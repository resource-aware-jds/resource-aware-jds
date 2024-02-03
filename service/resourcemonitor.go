package service

import (
	"github.com/docker/docker/client"
	"github.com/nabhan-au/dockerstats"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/datastructure"
	"github.com/sirupsen/logrus"
)

type ResourceMonitor struct {
	dockerClient     *client.Client
	containerService IContainer
}

type IResourceMonitor interface {
	GetResourceUsage() ([]models.ContainerResourceUsage, error)
}

func ProvideResourcesMonitor(dockerClient *client.Client, workerService IContainer) IResourceMonitor {
	return &ResourceMonitor{
		dockerClient:     dockerClient,
		containerService: workerService,
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
