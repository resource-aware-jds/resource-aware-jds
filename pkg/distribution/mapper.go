package distribution

import (
	"github.com/resource-aware-jds/resource-aware-jds/config"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/service"
	"go.opentelemetry.io/otel/metric"
)

type distributorMapper struct {
	distributorList map[models.DistributorName]Distributor
}

type DistributorMapper interface {
	GetDistributor(name models.DistributorName) (Distributor, bool)
}

func ProvideDistributorMapper(resourceAwareDistributorConfig config.ResourceAwareDistributorConfigModel, metric metric.Meter, taskService service.Task) DistributorMapper {
	return &distributorMapper{
		distributorList: map[models.DistributorName]Distributor{
			models.RoundRobinDistributorName:    ProvideRoundRobinDistributor(metric),
			models.ResourceAwareDistributorName: ProvideResourceAwareDistributor(resourceAwareDistributorConfig, metric, taskService),
		},
	}
}

func (d *distributorMapper) GetDistributor(name models.DistributorName) (Distributor, bool) {
	result, ok := d.distributorList[name]
	return result, ok
}
