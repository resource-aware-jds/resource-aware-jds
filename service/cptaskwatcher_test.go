package service_test

import (
	"context"
	"github.com/resource-aware-jds/resource-aware-jds/config"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/service"
	"github.com/resource-aware-jds/resource-aware-jds/service/mock_service"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestCPTaskWatcherTestSuite(t *testing.T) {
	suite.Run(t, new(CPTaskWatcherTestSuite))
}

type CPTaskWatcherTestSuite struct {
	suite.Suite

	ctx             context.Context
	ctrl            *gomock.Controller
	mockTaskService *mock_service.MockTask

	underTest service.CPTaskWatcher
}

func (s *CPTaskWatcherTestSuite) SetupSubTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockTaskService = mock_service.NewMockTask(s.ctrl)
	s.ctx = context.Background()
}

func (s *CPTaskWatcherTestSuite) Prepare(config config.TaskWatcherConfigModel) {
	s.underTest = service.ProvideCPTaskWatcher(s.mockTaskService, config)
}

func (s *CPTaskWatcherTestSuite) TestWatcherLoop() {
	s.Run("Call the Timeout for the passed deadline task", func() {
		s.Prepare(config.TaskWatcherConfigModel{
			Timeout: 0 * time.Second,
		})

		mockedObjectID := primitive.NewObjectID()
		s.underTest.AddTaskToWatch(mockedObjectID)

		time.Sleep(1 * time.Second)

		s.mockTaskService.EXPECT().UpdateTaskWaitTimeout(s.ctx, mockedObjectID).Return(nil)
		s.underTest.WatcherLoop(s.ctx)

		time.Sleep(1 * time.Second)
	})
	s.Run("Should remove the objectID once updated", func() {
		s.Prepare(config.TaskWatcherConfigModel{
			Timeout: 20 * time.Second,
		})

		mockedObjectID := primitive.NewObjectID()
		s.underTest.AddTaskToWatch(mockedObjectID)
		mockedObjectID2 := primitive.NewObjectID()
		s.underTest.AddTaskToWatch(mockedObjectID2)

		err := s.underTest.OnEvent(s.ctx, models.TaskEventBus{
			TaskID:    mockedObjectID,
			EventType: models.SuccessTaskEventType,
		})
		s.NoError(err)

		s.underTest.WatcherLoop(s.ctx)

		time.Sleep(1 * time.Second)
	})
}
