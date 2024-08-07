// Code generated by MockGen. DO NOT EDIT.
// Source: ./parser/json.go
//
// Generated by this command:
//
//	mockgen -source ./parser/json.go -destination ./tests/mock/parser/json.go
//

// Package mock_parser is a generated GoMock package.
package mock_parser

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockJSONInterface is a mock of JSONInterface interface.
type MockJSONInterface struct {
	ctrl     *gomock.Controller
	recorder *MockJSONInterfaceMockRecorder
}

// MockJSONInterfaceMockRecorder is the mock recorder for MockJSONInterface.
type MockJSONInterfaceMockRecorder struct {
	mock *MockJSONInterface
}

// NewMockJSONInterface creates a new mock instance.
func NewMockJSONInterface(ctrl *gomock.Controller) *MockJSONInterface {
	mock := &MockJSONInterface{ctrl: ctrl}
	mock.recorder = &MockJSONInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockJSONInterface) EXPECT() *MockJSONInterfaceMockRecorder {
	return m.recorder
}

// Marshal mocks base method.
func (m *MockJSONInterface) Marshal(orig any) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Marshal", orig)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Marshal indicates an expected call of Marshal.
func (mr *MockJSONInterfaceMockRecorder) Marshal(orig any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Marshal", reflect.TypeOf((*MockJSONInterface)(nil).Marshal), orig)
}

// Unmarshal mocks base method.
func (m *MockJSONInterface) Unmarshal(blob []byte, dest any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unmarshal", blob, dest)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unmarshal indicates an expected call of Unmarshal.
func (mr *MockJSONInterfaceMockRecorder) Unmarshal(blob, dest any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unmarshal", reflect.TypeOf((*MockJSONInterface)(nil).Unmarshal), blob, dest)
}
