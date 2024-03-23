package daemon_test

import (
	"github.com/resource-aware-jds/resource-aware-jds/config"
	"github.com/resource-aware-jds/resource-aware-jds/daemon"
	"github.com/resource-aware-jds/resource-aware-jds/handlerservice/mock_handlerservice"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/pool/mock_pool"
	"github.com/resource-aware-jds/resource-aware-jds/service/mock_service"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestControlPlaneTestSuite(t *testing.T) {
	suite.Run(t, new(ControlPlaneTestSuite))
}

type ControlPlaneTestSuite struct {
	suite.Suite

	ctrl                           *gomock.Controller
	mockNodePool                   *mock_pool.MockWorkerNode
	mockControlPlaneHandlerService *mock_handlerservice.MockIControlPlane
	mockTaskService                *mock_service.MockTask
	mockJobService                 *mock_service.MockJob
	mockCPTaskWatcher              *mock_service.MockCPTaskWatcher

	underTest daemon.IControlPlane
}

func (s *ControlPlaneTestSuite) SetupSubTest() {

	s.ctrl = gomock.NewController(s.T())
	s.mockNodePool = mock_pool.NewMockWorkerNode(s.ctrl)
	s.mockControlPlaneHandlerService = mock_handlerservice.NewMockIControlPlane(s.ctrl)
	s.mockTaskService = mock_service.NewMockTask(s.ctrl)
	s.mockJobService = mock_service.NewMockJob(s.ctrl)
	s.mockCPTaskWatcher = mock_service.NewMockCPTaskWatcher(s.ctrl)

	underTest, _ := daemon.ProvideControlPlaneDaemon(
		s.mockNodePool,
		s.mockControlPlaneHandlerService,
		s.mockTaskService,
		s.mockJobService,
		s.mockCPTaskWatcher,
		config.ControlPlaneConfigModel{},
	)

	s.underTest = underTest
}

func (s *ControlPlaneTestSuite) TestCheckTheDistributedTask() {
	s.Run("No Task to be recovered", func() {

		workerNodeResponseTasks := map[primitive.ObjectID]bool{}
		distributedTasks := []models.Task{}

		gomock.InOrder(
			s.mockNodePool.EXPECT().CheckRunningTaskInEachWorkerNode(gomock.Any()).Return(workerNodeResponseTasks),
			s.mockTaskService.EXPECT().GetAllDistributedTask(gomock.Any()).Return(distributedTasks, nil),
		)

		s.underTest.CheckTheDistributedTask()
	})
	s.Run("Worker Node response some of the running task", func() {
		taskAID := primitive.NewObjectID()
		workerNodeResponseTasks := map[primitive.ObjectID]bool{
			taskAID: true,
		}
		distributedTasks := []models.Task{
			{
				ID: &taskAID,
			},
		}

		gomock.InOrder(
			s.mockNodePool.EXPECT().CheckRunningTaskInEachWorkerNode(gomock.Any()).Return(workerNodeResponseTasks),
			s.mockTaskService.EXPECT().GetAllDistributedTask(gomock.Any()).Return(distributedTasks, nil),
			s.mockCPTaskWatcher.EXPECT().AddTaskToWatch(taskAID),
		)

		s.underTest.CheckTheDistributedTask()
	})
	s.Run("Worker Node not response the working task", func() {
		taskAID := primitive.NewObjectID()
		workerNodeResponseTasks := map[primitive.ObjectID]bool{}
		distributedTasks := []models.Task{
			{
				ID: &taskAID,
			},
		}

		gomock.InOrder(
			s.mockNodePool.EXPECT().CheckRunningTaskInEachWorkerNode(gomock.Any()).Return(workerNodeResponseTasks),
			s.mockTaskService.EXPECT().GetAllDistributedTask(gomock.Any()).Return(distributedTasks, nil),
			s.mockTaskService.EXPECT().UpdateTaskWorkOnFailure(gomock.Any(), taskAID, "", "Control Plane startup process detect no worker in the pool running this task").Return(nil),
		)

		s.underTest.CheckTheDistributedTask()
	})
	s.Run("Mix-case", func() {
		taskAID := primitive.NewObjectID()
		taskBID := primitive.NewObjectID()
		workerNodeResponseTasks := map[primitive.ObjectID]bool{
			taskBID: true,
		}
		distributedTasks := []models.Task{
			{
				ID: &taskAID,
			},
			{
				ID: &taskBID,
			},
		}

		gomock.InOrder(
			s.mockNodePool.EXPECT().CheckRunningTaskInEachWorkerNode(gomock.Any()).Return(workerNodeResponseTasks),
			s.mockTaskService.EXPECT().GetAllDistributedTask(gomock.Any()).Return(distributedTasks, nil),
			s.mockTaskService.EXPECT().UpdateTaskWorkOnFailure(gomock.Any(), taskAID, "", "Control Plane startup process detect no worker in the pool running this task").Return(nil),
			s.mockCPTaskWatcher.EXPECT().AddTaskToWatch(taskBID),
		)

		s.underTest.CheckTheDistributedTask()
	})
}
