// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/asatraitis/mangrove/internal/bll (interfaces: BLL)
//
// Generated by this command:
//
//	mockgen -destination=./mocks/mock_bll.go -package=mocks github.com/asatraitis/mangrove/internal/bll BLL
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	bll "github.com/asatraitis/mangrove/internal/bll"
	gomock "go.uber.org/mock/gomock"
)

// MockBLL is a mock of BLL interface.
type MockBLL struct {
	ctrl     *gomock.Controller
	recorder *MockBLLMockRecorder
	isgomock struct{}
}

// MockBLLMockRecorder is the mock recorder for MockBLL.
type MockBLLMockRecorder struct {
	mock *MockBLL
}

// NewMockBLL creates a new mock instance.
func NewMockBLL(ctrl *gomock.Controller) *MockBLL {
	mock := &MockBLL{ctrl: ctrl}
	mock.recorder = &MockBLLMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBLL) EXPECT() *MockBLLMockRecorder {
	return m.recorder
}

// Client mocks base method.
func (m *MockBLL) Client(arg0 context.Context) bll.ClientBLL {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Client", arg0)
	ret0, _ := ret[0].(bll.ClientBLL)
	return ret0
}

// Client indicates an expected call of Client.
func (mr *MockBLLMockRecorder) Client(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Client", reflect.TypeOf((*MockBLL)(nil).Client), arg0)
}

// Config mocks base method.
func (m *MockBLL) Config(arg0 context.Context) bll.ConfigBLL {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Config", arg0)
	ret0, _ := ret[0].(bll.ConfigBLL)
	return ret0
}

// Config indicates an expected call of Config.
func (mr *MockBLLMockRecorder) Config(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Config", reflect.TypeOf((*MockBLL)(nil).Config), arg0)
}

// User mocks base method.
func (m *MockBLL) User(arg0 context.Context) bll.UserBLL {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "User", arg0)
	ret0, _ := ret[0].(bll.UserBLL)
	return ret0
}

// User indicates an expected call of User.
func (mr *MockBLLMockRecorder) User(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "User", reflect.TypeOf((*MockBLL)(nil).User), arg0)
}
