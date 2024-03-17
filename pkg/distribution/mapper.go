package distribution

import (
	"github.com/resource-aware-jds/resource-aware-jds/config"
	"github.com/resource-aware-jds/resource-aware-jds/service"
	"go.opentelemetry.io/otel/metric"
)

type DistributorName string

const (
	RoundRobinDistributorName    DistributorName = "round_robin"
	ResourceAwareDistributorName DistributorName = "resource_aware"
)

type distributorMapper struct {
	distributorList map[DistributorName]Distributor
}

type DistributorMapper interface {
	GetDistributor(name DistributorName) (Distributor, bool)
}

func ProvideDistributorMapper(resourceAwareDistributorConfig config.ResourceAwareDistributorConfigModel, metric metric.Meter, taskService service.Task) DistributorMapper {
	return &distributorMapper{
		distributorList: map[DistributorName]Distributor{
			RoundRobinDistributorName:    ProvideRoundRobinDistributor(metric),
			ResourceAwareDistributorName: ProvideResourceAwareDistributor(resourceAwareDistributorConfig, metric, taskService),
		},
	}
}

func (d *distributorMapper) GetDistributor(name DistributorName) (Distributor, bool) {
	result, ok := d.distributorList[name]
	return result, ok
}
