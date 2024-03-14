package distribution

import "go.opentelemetry.io/otel/metric"

type DistributorName string

const (
	RoundRobinDistributorName    DistributorName = "round-robin-distributor"
	ResourceAwareDistributorName DistributorName = "resource-aware-distributor"
)

type distributorMapper struct {
	distributorList map[DistributorName]Distributor
}

type DistributorMapper interface {
	GetDistributor(name DistributorName) (Distributor, bool)
}

func ProvideDistributorMapper(metric metric.Meter) DistributorMapper {
	return &distributorMapper{
		distributorList: map[DistributorName]Distributor{
			RoundRobinDistributorName:    ProvideRoundRobinDistributor(metric),
			ResourceAwareDistributorName: ProvideResourceAwareDistributor(metric),
		},
	}
}

func (d *distributorMapper) GetDistributor(name DistributorName) (Distributor, bool) {
	result, ok := d.distributorList[name]
	return result, ok
}
