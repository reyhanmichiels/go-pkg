// Code generated by MockGen. DO NOT EDIT.
// Source: ./auth/auth.go
//
// Generated by this command:
//
//	mockgen -source ./auth/auth.go -destination ./tests/mock/auth/auth.go
//

// Package mock_auth is a generated GoMock package.
package mock_auth

import (
	context "context"
	reflect "reflect"

	auth "github.com/reyhanmichiels/go-pkg/auth"
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

// CreateAccessToken mocks base method.
func (m *MockInterface) CreateAccessToken(userID int64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccessToken", userID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAccessToken indicates an expected call of CreateAccessToken.
func (mr *MockInterfaceMockRecorder) CreateAccessToken(userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccessToken", reflect.TypeOf((*MockInterface)(nil).CreateAccessToken), userID)
}

// CreateRefreshToken mocks base method.
func (m *MockInterface) CreateRefreshToken(userID int64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRefreshToken", userID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRefreshToken indicates an expected call of CreateRefreshToken.
func (mr *MockInterfaceMockRecorder) CreateRefreshToken(userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRefreshToken", reflect.TypeOf((*MockInterface)(nil).CreateRefreshToken), userID)
}

// GetUserAuthInfo mocks base method.
func (m *MockInterface) GetUserAuthInfo(ctx context.Context) (auth.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserAuthInfo", ctx)
	ret0, _ := ret[0].(auth.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserAuthInfo indicates an expected call of GetUserAuthInfo.
func (mr *MockInterfaceMockRecorder) GetUserAuthInfo(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserAuthInfo", reflect.TypeOf((*MockInterface)(nil).GetUserAuthInfo), ctx)
}

// SetUserAuthInfo mocks base method.
func (m *MockInterface) SetUserAuthInfo(ctx context.Context, user auth.User) context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetUserAuthInfo", ctx, user)
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// SetUserAuthInfo indicates an expected call of SetUserAuthInfo.
func (mr *MockInterfaceMockRecorder) SetUserAuthInfo(ctx, user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetUserAuthInfo", reflect.TypeOf((*MockInterface)(nil).SetUserAuthInfo), ctx, user)
}

// ValidateAccessToken mocks base method.
func (m *MockInterface) ValidateAccessToken(token string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateAccessToken", token)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateAccessToken indicates an expected call of ValidateAccessToken.
func (mr *MockInterfaceMockRecorder) ValidateAccessToken(token any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateAccessToken", reflect.TypeOf((*MockInterface)(nil).ValidateAccessToken), token)
}

// ValidateRefreshToken mocks base method.
func (m *MockInterface) ValidateRefreshToken(token string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateRefreshToken", token)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateRefreshToken indicates an expected call of ValidateRefreshToken.
func (mr *MockInterfaceMockRecorder) ValidateRefreshToken(token any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateRefreshToken", reflect.TypeOf((*MockInterface)(nil).ValidateRefreshToken), token)
}
