package distribution_test

import (
	"github.com/golang/mock/gomock"
	"github.com/resource-aware-jds/resource-aware-jds/generated/mock_proto"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/distribution"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestRoundRobinDistributorTestSuite(t *testing.T) {
	suite.Run(t, new(RoundRobinDistributorTestSuite))
}

type RoundRobinDistributorTestSuite struct {
	BaseDistributionTest
}

func (s *RoundRobinDistributorTestSuite) SetupSubTest() {
	s.BaseDistributionTest.SetupSubTest()
	s.underTest = distribution.ProvideRoundRobinDistributor(s.meter)
}

func (s *RoundRobinDistributorTestSuite) TestDistribution() {
	s.Run("Should distribute to correct node", func() {
		mockNodeMapperTest := []nodeMapperTestData{
			s.newNodeMapper("Node 1", models.AvailableResource{
				CpuCores:               1,
				AvailableCpuPercentage: 100,
				AvailableMemory: models.MemorySize{
					Size: 500,
					Unit: "MiB",
				},
			}),
			s.newNodeMapper("Node 2", models.AvailableResource{
				CpuCores:               1,
				AvailableCpuPercentage: 0,
				AvailableMemory: models.MemorySize{
					Size: 0,
					Unit: "MiB",
				},
			}),
			s.newNodeMapper("Node 3", models.AvailableResource{
				CpuCores:               1,
				AvailableCpuPercentage: 100,
				AvailableMemory: models.MemorySize{
					Size: 0,
					Unit: "MiB",
				},
			}),
			s.newNodeMapper("Node 4", models.AvailableResource{
				CpuCores:               1,
				AvailableCpuPercentage: 0,
				AvailableMemory: models.MemorySize{
					Size: 100,
					Unit: "MiB",
				},
			}),
		}
		nodeMapper := toNodeMapper(mockNodeMapperTest)
		mockJobID := primitive.NewObjectID()
		taskToDistribute := []models.Task{
			{
				ID:       newObjectID(),
				JobID:    &mockJobID,
				ImageUrl: "task-a-image",
			},
			{
				ID:       newObjectID(),
				JobID:    &mockJobID,
				ImageUrl: "task-b-image",
			},
			{
				ID:       newObjectID(),
				JobID:    &mockJobID,
				ImageUrl: "task-c-image",
			},
			{
				ID:       newObjectID(),
				JobID:    &mockJobID,
				ImageUrl: "task-d-image",
			},
			{
				ID:       newObjectID(),
				JobID:    &mockJobID,
				ImageUrl: "task-e-image",
			},
		}

		mockCall := make([]*gomock.Call, 0, 5)
		for i := 0; i < 5; i++ {
			mockCall = append(mockCall, mockNodeMapperTest[i%4].MockWorkerNodeGRPC.EXPECT().SendTask(s.ctx, &proto.RecievedTask{
				ID:             taskToDistribute[i].ID.Hex(),
				TaskAttributes: nil,
				DockerImage:    taskToDistribute[i].ImageUrl,
			}))
		}

		gomock.InOrder(mockCall...)

		successTask, failureTask, err := s.underTest.Distribute(s.ctx, nodeMapper, taskToDistribute)
		s.NoError(err)
		s.Len(failureTask, 0)
		s.Len(successTask, 5)
	})
}

func (s *RoundRobinDistributorTestSuite) newNodeMapper(nodeID string, mockAvailableResource models.AvailableResource) nodeMapperTestData {
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
