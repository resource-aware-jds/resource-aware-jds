package distribution_test

import (
	"github.com/resource-aware-jds/resource-aware-jds/config"
	"github.com/resource-aware-jds/resource-aware-jds/generated/proto/github.com/resource-aware-jds/resource-aware-jds/generated/proto"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/distribution"
	"github.com/resource-aware-jds/resource-aware-jds/service/mock_service"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestResourceAwareDistributorTestSuite(t *testing.T) {
	suite.Run(t, new(ResourceAwareDistributorTestSuite))
}

type ResourceAwareDistributorTestSuite struct {
	BaseDistributionTest

	mockTaskService *mock_service.MockTask
}

func (s *ResourceAwareDistributorTestSuite) SetupSubTest() {
	s.BaseDistributionTest.SetupSubTest()
	s.mockTaskService = mock_service.NewMockTask(s.ctrl)
}

func (s *ResourceAwareDistributorTestSuite) PrepareDistribution(config config.ResourceAwareDistributorConfigModel) {
	s.underTest = distribution.ProvideResourceAwareDistributor(config, s.meter, s.mockTaskService)
}

func (s *ResourceAwareDistributorTestSuite) TestDistribution() {
	s.Run("Should distribute to correct node", func() {
		s.PrepareDistribution(config.ResourceAwareDistributorConfigModel{
			AvailableResourceClearanceThreshold: 80,
		})

		mockNodeMapperTest := []nodeMapperTestData{
			s.newNodeMapper("Node 1", models.AvailableResource{
				CpuCores:               1,
				AvailableCpuPercentage: 100,
				AvailableMemory: models.MemorySize{
					Size: 500,
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
		}

		s.mockTaskService.EXPECT().GetAverageResourceUsage(s.ctx, &mockJobID).Return(&models.TaskResourceUsage{
			Memory: 100,
			CPU:    10,
		}, nil)

		for i := 0; i < 3; i++ {
			mockNodeMapperTest[0].MockWorkerNodeGRPC.EXPECT().SendTask(gomock.Any(), &proto.RecievedTask{
				ID:             taskToDistribute[i].ID.Hex(),
				TaskAttributes: nil,
				DockerImage:    taskToDistribute[i].ImageUrl,
			}).Times(1)
		}

		successTask, failureTask, err := s.underTest.Distribute(s.ctx, nodeMapper, taskToDistribute)
		s.NoError(err)
		s.Len(failureTask, 0)

		s.Len(successTask, len(taskToDistribute))
	})
	s.Run("Should avoid the clearance threshold", func() {
		s.PrepareDistribution(config.ResourceAwareDistributorConfigModel{
			AvailableResourceClearanceThreshold: 50,
		})

		mockNodeMapperTest := []nodeMapperTestData{
			s.newNodeMapper("Node 1", models.AvailableResource{
				CpuCores:               1,
				AvailableCpuPercentage: 100,
				AvailableMemory: models.MemorySize{
					Size: 500,
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
		}

		s.mockTaskService.EXPECT().GetAverageResourceUsage(s.ctx, &mockJobID).Return(&models.TaskResourceUsage{
			Memory: 100,
			CPU:    20,
		}, nil)

		for i := 0; i < 2; i++ {
			mockNodeMapperTest[0].MockWorkerNodeGRPC.EXPECT().SendTask(gomock.Any(), &proto.RecievedTask{
				ID:             taskToDistribute[i].ID.Hex(),
				TaskAttributes: nil,
				DockerImage:    taskToDistribute[i].ImageUrl,
			}).Times(1)
		}

		successTask, failureTask, err := s.underTest.Distribute(s.ctx, nodeMapper, taskToDistribute)
		s.NoError(err)
		s.Len(failureTask, 1)
		s.Len(successTask, 2)

		s.Len(failureTask[0].Task.Logs, 1)
		s.Equal("Fail to distribute task to node", failureTask[0].Task.Logs[0].Message)
		s.Equal(distribution.ErrNotEnoughResource.Error(), failureTask[0].Task.Logs[0].Parameters["error"])
		s.Equal(models.WarnLogSeverity, failureTask[0].Task.Logs[0].Severity)
	})
	s.Run("Should avoid the clearance threshold and distribute to the next node", func() {
		s.PrepareDistribution(config.ResourceAwareDistributorConfigModel{
			AvailableResourceClearanceThreshold: 50,
		})

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
				AvailableCpuPercentage: 100,
				AvailableMemory: models.MemorySize{
					Size: 500,
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
		}

		s.mockTaskService.EXPECT().GetAverageResourceUsage(s.ctx, &mockJobID).Return(&models.TaskResourceUsage{
			Memory: 100,
			CPU:    20,
		}, nil)

		for i := 0; i < 2; i++ {
			mockNodeMapperTest[0].MockWorkerNodeGRPC.EXPECT().SendTask(gomock.Any(), &proto.RecievedTask{
				ID:             taskToDistribute[i].ID.Hex(),
				TaskAttributes: nil,
				DockerImage:    taskToDistribute[i].ImageUrl,
			})
		}
		mockNodeMapperTest[1].MockWorkerNodeGRPC.EXPECT().SendTask(gomock.Any(), &proto.RecievedTask{
			ID:             taskToDistribute[2].ID.Hex(),
			TaskAttributes: nil,
			DockerImage:    taskToDistribute[2].ImageUrl,
		})

		successTask, failureTask, err := s.underTest.Distribute(s.ctx, nodeMapper, taskToDistribute)
		s.NoError(err)
		s.Len(failureTask, 0)
		s.Len(successTask, 3)
	})
	s.Run("Should avoid the clearance threshold and distribute to the next available node", func() {
		s.PrepareDistribution(config.ResourceAwareDistributorConfigModel{
			AvailableResourceClearanceThreshold: 50,
		})

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
				AvailableCpuPercentage: 100,
				AvailableMemory: models.MemorySize{
					Size: 10,
					Unit: "MiB",
				},
			}),
			s.newNodeMapper("Node 3", models.AvailableResource{
				CpuCores:               1,
				AvailableCpuPercentage: 40,
				AvailableMemory: models.MemorySize{
					Size: 500,
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
		}

		s.mockTaskService.EXPECT().GetAverageResourceUsage(s.ctx, &mockJobID).Return(&models.TaskResourceUsage{
			Memory: 100,
			CPU:    20,
		}, nil)

		for i := 0; i < 2; i++ {
			mockNodeMapperTest[0].MockWorkerNodeGRPC.EXPECT().SendTask(gomock.Any(), &proto.RecievedTask{
				ID:             taskToDistribute[i].ID.Hex(),
				TaskAttributes: nil,
				DockerImage:    taskToDistribute[i].ImageUrl,
			})
		}
		mockNodeMapperTest[2].MockWorkerNodeGRPC.EXPECT().SendTask(gomock.Any(), &proto.RecievedTask{
			ID:             taskToDistribute[2].ID.Hex(),
			TaskAttributes: nil,
			DockerImage:    taskToDistribute[2].ImageUrl,
		})

		successTask, failureTask, err := s.underTest.Distribute(s.ctx, nodeMapper, taskToDistribute)
		s.NoError(err)
		s.Len(failureTask, 0)
		s.Len(successTask, 3)
	})
}
