// Code generated by MockGen. DO NOT EDIT.
// Source: ./task.go
//
// Generated by this command:
//
//	mockgen -source=./task.go -destination=./mock_service/mock_task.go -package=mock_service
//

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	models "github.com/resource-aware-jds/resource-aware-jds/models"
	primitive "go.mongodb.org/mongo-driver/bson/primitive"
	gomock "go.uber.org/mock/gomock"
)

// MockTask is a mock of Task interface.
type MockTask struct {
	ctrl     *gomock.Controller
	recorder *MockTaskMockRecorder
}

// MockTaskMockRecorder is the mock recorder for MockTask.
type MockTaskMockRecorder struct {
	mock *MockTask
}

// NewMockTask creates a new mock instance.
func NewMockTask(ctrl *gomock.Controller) *MockTask {
	mock := &MockTask{ctrl: ctrl}
	mock.recorder = &MockTaskMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTask) EXPECT() *MockTaskMockRecorder {
	return m.recorder
}

// CountUnfinishedTaskByJobID mocks base method.
func (m *MockTask) CountUnfinishedTaskByJobID(ctx context.Context, jobID *primitive.ObjectID) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountUnfinishedTaskByJobID", ctx, jobID)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountUnfinishedTaskByJobID indicates an expected call of CountUnfinishedTaskByJobID.
func (mr *MockTaskMockRecorder) CountUnfinishedTaskByJobID(ctx, jobID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountUnfinishedTaskByJobID", reflect.TypeOf((*MockTask)(nil).CountUnfinishedTaskByJobID), ctx, jobID)
}

// CreateTask mocks base method.
func (m *MockTask) CreateTask(ctx context.Context, job *models.Job, taskAttributes [][]byte, isExperiment bool) ([]models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTask", ctx, job, taskAttributes, isExperiment)
	ret0, _ := ret[0].([]models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTask indicates an expected call of CreateTask.
func (mr *MockTaskMockRecorder) CreateTask(ctx, job, taskAttributes, isExperiment any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTask", reflect.TypeOf((*MockTask)(nil).CreateTask), ctx, job, taskAttributes, isExperiment)
}

// GetAllDistributedTask mocks base method.
func (m *MockTask) GetAllDistributedTask(ctx context.Context) ([]models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllDistributedTask", ctx)
	ret0, _ := ret[0].([]models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllDistributedTask indicates an expected call of GetAllDistributedTask.
func (mr *MockTaskMockRecorder) GetAllDistributedTask(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllDistributedTask", reflect.TypeOf((*MockTask)(nil).GetAllDistributedTask), ctx)
}

// GetAvailableTask mocks base method.
func (m *MockTask) GetAvailableTask(ctx context.Context, jobIDs []models.Job) (*models.Job, []models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAvailableTask", ctx, jobIDs)
	ret0, _ := ret[0].(*models.Job)
	ret1, _ := ret[1].([]models.Task)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetAvailableTask indicates an expected call of GetAvailableTask.
func (mr *MockTaskMockRecorder) GetAvailableTask(ctx, jobIDs any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAvailableTask", reflect.TypeOf((*MockTask)(nil).GetAvailableTask), ctx, jobIDs)
}

// GetAverageResourceUsage mocks base method.
func (m *MockTask) GetAverageResourceUsage(ctx context.Context, jobID *primitive.ObjectID) (*models.TaskResourceUsage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAverageResourceUsage", ctx, jobID)
	ret0, _ := ret[0].(*models.TaskResourceUsage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAverageResourceUsage indicates an expected call of GetAverageResourceUsage.
func (mr *MockTaskMockRecorder) GetAverageResourceUsage(ctx, jobID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAverageResourceUsage", reflect.TypeOf((*MockTask)(nil).GetAverageResourceUsage), ctx, jobID)
}

// GetTaskByID mocks base method.
func (m *MockTask) GetTaskByID(ctx context.Context, taskID primitive.ObjectID) (*models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTaskByID", ctx, taskID)
	ret0, _ := ret[0].(*models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTaskByID indicates an expected call of GetTaskByID.
func (mr *MockTaskMockRecorder) GetTaskByID(ctx, taskID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTaskByID", reflect.TypeOf((*MockTask)(nil).GetTaskByID), ctx, taskID)
}

// GetTaskByJob mocks base method.
func (m *MockTask) GetTaskByJob(ctx context.Context, job *models.Job) ([]models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTaskByJob", ctx, job)
	ret0, _ := ret[0].([]models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTaskByJob indicates an expected call of GetTaskByJob.
func (mr *MockTaskMockRecorder) GetTaskByJob(ctx, job any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTaskByJob", reflect.TypeOf((*MockTask)(nil).GetTaskByJob), ctx, job)
}

// UpdateAllTaskToWorkOnFailure mocks base method.
func (m *MockTask) UpdateAllTaskToWorkOnFailure(ctx context.Context, job *models.Job, jobErrorMessage string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAllTaskToWorkOnFailure", ctx, job, jobErrorMessage)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateAllTaskToWorkOnFailure indicates an expected call of UpdateAllTaskToWorkOnFailure.
func (mr *MockTaskMockRecorder) UpdateAllTaskToWorkOnFailure(ctx, job, jobErrorMessage any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAllTaskToWorkOnFailure", reflect.TypeOf((*MockTask)(nil).UpdateAllTaskToWorkOnFailure), ctx, job, jobErrorMessage)
}

// UpdateTaskAfterDistribution mocks base method.
func (m *MockTask) UpdateTaskAfterDistribution(ctx context.Context, successTasks []models.Task, errorTasks []models.DistributeError) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTaskAfterDistribution", ctx, successTasks, errorTasks)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTaskAfterDistribution indicates an expected call of UpdateTaskAfterDistribution.
func (mr *MockTaskMockRecorder) UpdateTaskAfterDistribution(ctx, successTasks, errorTasks any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTaskAfterDistribution", reflect.TypeOf((*MockTask)(nil).UpdateTaskAfterDistribution), ctx, successTasks, errorTasks)
}

// UpdateTaskSuccess mocks base method.
func (m *MockTask) UpdateTaskSuccess(ctx context.Context, taskID primitive.ObjectID, nodeID string, result []byte, averageCPUUsage float32, averageMemoryUsage float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTaskSuccess", ctx, taskID, nodeID, result, averageCPUUsage, averageMemoryUsage)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTaskSuccess indicates an expected call of UpdateTaskSuccess.
func (mr *MockTaskMockRecorder) UpdateTaskSuccess(ctx, taskID, nodeID, result, averageCPUUsage, averageMemoryUsage any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTaskSuccess", reflect.TypeOf((*MockTask)(nil).UpdateTaskSuccess), ctx, taskID, nodeID, result, averageCPUUsage, averageMemoryUsage)
}

// UpdateTaskToBeReadyToBeDistributed mocks base method.
func (m *MockTask) UpdateTaskToBeReadyToBeDistributed(ctx context.Context, jobID *primitive.ObjectID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTaskToBeReadyToBeDistributed", ctx, jobID)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTaskToBeReadyToBeDistributed indicates an expected call of UpdateTaskToBeReadyToBeDistributed.
func (mr *MockTaskMockRecorder) UpdateTaskToBeReadyToBeDistributed(ctx, jobID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTaskToBeReadyToBeDistributed", reflect.TypeOf((*MockTask)(nil).UpdateTaskToBeReadyToBeDistributed), ctx, jobID)
}

// UpdateTaskWaitTimeout mocks base method.
func (m *MockTask) UpdateTaskWaitTimeout(ctx context.Context, taskID primitive.ObjectID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTaskWaitTimeout", ctx, taskID)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTaskWaitTimeout indicates an expected call of UpdateTaskWaitTimeout.
func (mr *MockTaskMockRecorder) UpdateTaskWaitTimeout(ctx, taskID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTaskWaitTimeout", reflect.TypeOf((*MockTask)(nil).UpdateTaskWaitTimeout), ctx, taskID)
}

// UpdateTaskWorkOnFailure mocks base method.
func (m *MockTask) UpdateTaskWorkOnFailure(ctx context.Context, taskID primitive.ObjectID, nodeID, errMessage string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTaskWorkOnFailure", ctx, taskID, nodeID, errMessage)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTaskWorkOnFailure indicates an expected call of UpdateTaskWorkOnFailure.
func (mr *MockTaskMockRecorder) UpdateTaskWorkOnFailure(ctx, taskID, nodeID, errMessage any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTaskWorkOnFailure", reflect.TypeOf((*MockTask)(nil).UpdateTaskWorkOnFailure), ctx, taskID, nodeID, errMessage)
}
