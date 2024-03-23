package grpc

import (
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/models"
)

func ToDistributorName(in proto.DistributionLogic) models.DistributorName {
	switch in {
	case proto.DistributionLogic_RoundRobin:
		return models.RoundRobinDistributorName
	case proto.DistributionLogic_ResourceAware:
		return models.ResourceAwareDistributorName
	}

	return models.RoundRobinDistributorName
}
