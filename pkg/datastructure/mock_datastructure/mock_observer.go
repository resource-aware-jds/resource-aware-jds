// Code generated by MockGen. DO NOT EDIT.
// Source: ./observer.go
//
// Generated by this command:
//
//	mockgen -source=./observer.go -destination=./mock_datastructure/mock_observer.go -package=mock_datastructure
//

// Package mock_datastructure is a generated GoMock package.
package mock_datastructure

import (
	context "context"
	reflect "reflect"

	datastructure "github.com/resource-aware-jds/resource-aware-jds/pkg/datastructure"
	gomock "go.uber.org/mock/gomock"
)

// MockObserver is a mock of Observer interface.
type MockObserver[T any] struct {
	ctrl     *gomock.Controller
	recorder *MockObserverMockRecorder[T]
}

// MockObserverMockRecorder is the mock recorder for MockObserver.
type MockObserverMockRecorder[T any] struct {
	mock *MockObserver[T]
}

// NewMockObserver creates a new mock instance.
func NewMockObserver[T any](ctrl *gomock.Controller) *MockObserver[T] {
	mock := &MockObserver[T]{ctrl: ctrl}
	mock.recorder = &MockObserverMockRecorder[T]{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockObserver[T]) EXPECT() *MockObserverMockRecorder[T] {
	return m.recorder
}

// OnEvent mocks base method.
func (m *MockObserver[T]) OnEvent(arg0 context.Context, arg1 T) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OnEvent", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// OnEvent indicates an expected call of OnEvent.
func (mr *MockObserverMockRecorder[T]) OnEvent(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnEvent", reflect.TypeOf((*MockObserver[T])(nil).OnEvent), arg0, arg1)
}

// MockObservable is a mock of Observable interface.
type MockObservable[T any] struct {
	ctrl     *gomock.Controller
	recorder *MockObservableMockRecorder[T]
}

// MockObservableMockRecorder is the mock recorder for MockObservable.
type MockObservableMockRecorder[T any] struct {
	mock *MockObservable[T]
}

// NewMockObservable creates a new mock instance.
func NewMockObservable[T any](ctrl *gomock.Controller) *MockObservable[T] {
	mock := &MockObservable[T]{ctrl: ctrl}
	mock.recorder = &MockObservableMockRecorder[T]{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockObservable[T]) EXPECT() *MockObservableMockRecorder[T] {
	return m.recorder
}

// AddObserver mocks base method.
func (m *MockObservable[T]) AddObserver(arg0 datastructure.Observer[T]) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddObserver", arg0)
}

// AddObserver indicates an expected call of AddObserver.
func (mr *MockObservableMockRecorder[T]) AddObserver(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddObserver", reflect.TypeOf((*MockObservable[T])(nil).AddObserver), arg0)
}

// NotifyObserver mocks base method.
func (m *MockObservable[T]) NotifyObserver(arg0 context.Context, arg1 T) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NotifyObserver", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// NotifyObserver indicates an expected call of NotifyObserver.
func (mr *MockObservableMockRecorder[T]) NotifyObserver(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NotifyObserver", reflect.TypeOf((*MockObservable[T])(nil).NotifyObserver), arg0, arg1)
}

// RemoveObserver mocks base method.
func (m *MockObservable[T]) RemoveObserver(arg0 datastructure.Observer[T]) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RemoveObserver", arg0)
}

// RemoveObserver indicates an expected call of RemoveObserver.
func (mr *MockObservableMockRecorder[T]) RemoveObserver(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveObserver", reflect.TypeOf((*MockObservable[T])(nil).RemoveObserver), arg0)
}
