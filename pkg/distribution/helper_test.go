package distribution_test

import (
	"context"
	"github.com/resource-aware-jds/resource-aware-jds/generated/mock_proto"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/distribution"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/metrics"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/mock/gomock"
)

type nodeMapperTestData struct {
	NodeMapper         distribution.NodeMapper
	MockWorkerNodeGRPC *mock_proto.MockWorkerNodeClient
}

func toNodeMapper(data []nodeMapperTestData) []distribution.NodeMapper {
	result := make([]distribution.NodeMapper, len(data))
	for i, d := range data {
		result[i] = d.NodeMapper
	}

	return result
}

func newObjectID() *primitive.ObjectID {
	id := primitive.NewObjectID()
	return &id
}

type BaseDistributionTest struct {
	suite.Suite

	ctrl      *gomock.Controller
	meter     metric.Meter
	ctx       context.Context
	underTest distribution.Distributor
}

func (s *BaseDistributionTest) SetupSubTest() {
	s.ctrl = gomock.NewController(s.T())
	meter, err := metrics.ProvideMeter()
	s.NoError(err)
	s.meter = meter
	s.ctx = context.Background()
}

func (s *BaseDistributionTest) newNodeMapper(nodeID string, mockAvailableResource models.AvailableResource) nodeMapperTestData {
	mockWorkerNodeGRPC := mock_proto.NewMockWorkerNodeClient(s.ctrl)
	nodeEntry := models.NodeEntry{
		NodeID: nodeID,
	}

	result := distribution.NodeMapper{
		NodeEntry:         nodeEntry,
		GRPCConnection:    mockWorkerNodeGRPC,
		AvailableResource: mockAvailableResource,
		Logger:            logrus.WithField("type", "unittest"),
	}

	return nodeMapperTestData{
		NodeMapper:         result,
		MockWorkerNodeGRPC: mockWorkerNodeGRPC,
	}
}
