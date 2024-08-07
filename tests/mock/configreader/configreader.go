// Code generated by MockGen. DO NOT EDIT.
// Source: ./configreader/configreader.go
//
// Generated by this command:
//
//	mockgen -source ./configreader/configreader.go -destination ./tests/mock/configreader/configreader.go
//

// Package mock_configreader is a generated GoMock package.
package mock_configreader

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockInterface is a mock of Interface interface.
type MockInterface struct {
	ctrl     *gomock.Controller
	recorder *MockInterfaceMockRecorder
}

// MockInterfaceMockRecorder is the mock recorder for MockInterface.
type MockInterfaceMockRecorder struct {
	mock *MockInterface
}

// NewMockInterface creates a new mock instance.
func NewMockInterface(ctrl *gomock.Controller) *MockInterface {
	mock := &MockInterface{ctrl: ctrl}
	mock.recorder = &MockInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInterface) EXPECT() *MockInterfaceMockRecorder {
	return m.recorder
}

// ReadConfig mocks base method.
func (m *MockInterface) ReadConfig(cfg any) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ReadConfig", cfg)
}

// ReadConfig indicates an expected call of ReadConfig.
func (mr *MockInterfaceMockRecorder) ReadConfig(cfg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadConfig", reflect.TypeOf((*MockInterface)(nil).ReadConfig), cfg)
}
