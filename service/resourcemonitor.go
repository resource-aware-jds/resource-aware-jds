package service

import (
	"context"
	"github.com/docker/docker/client"
	"github.com/resource-aware-jds/dockerstats"
	"github.com/resource-aware-jds/go-osstat/cpu"
	"github.com/resource-aware-jds/go-osstat/memory"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/datastructure"
	"github.com/sirupsen/logrus"
	"time"
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

func (r *ResourceMonitor) GetUserMemUsage() (*models.MemoryUsage, error) {
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

func sleepContext(ctx context.Context, d time.Duration) {
	timer := time.NewTimer(d)
	select {
	case <-ctx.Done():
		timer.Stop()
	case <-timer.C:

	}
}

func (r *ResourceMonitor) GetUserCpuUsage(ctx context.Context) (*models.CpuUsage, error) {
	before, err := cpu.Get()
	if err != nil {
		logrus.Errorf("Unable to get os cpu usage: %e", err)
		return nil, err
	}
	sleepContext(ctx, time.Duration(1)*time.Second)
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
