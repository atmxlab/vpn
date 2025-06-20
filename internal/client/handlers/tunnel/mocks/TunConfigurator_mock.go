// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/atmxlab/vpn/internal/client/handlers/tunnel (interfaces: TunConfigurator)
//
// Generated by this command:
//
//	mockgen -destination=/Users/timur_abdurashidov/Desktop/other/vpn/internal/client/handlers/tunnel/mocks/TunConfigurator_mock.go -package=mocks . TunConfigurator
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	net "net"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockTunConfigurator is a mock of TunConfigurator interface.
type MockTunConfigurator struct {
	ctrl     *gomock.Controller
	recorder *MockTunConfiguratorMockRecorder
}

// MockTunConfiguratorMockRecorder is the mock recorder for MockTunConfigurator.
type MockTunConfiguratorMockRecorder struct {
	mock *MockTunConfigurator
}

// NewMockTunConfigurator creates a new mock instance.
func NewMockTunConfigurator(ctrl *gomock.Controller) *MockTunConfigurator {
	mock := &MockTunConfigurator{ctrl: ctrl}
	mock.recorder = &MockTunConfiguratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTunConfigurator) EXPECT() *MockTunConfiguratorMockRecorder {
	return m.recorder
}

// ChangeTunAddr mocks base method.
func (m *MockTunConfigurator) ChangeTunAddr(arg0 context.Context, arg1 net.IPNet) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeTunAddr", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeTunAddr indicates an expected call of ChangeTunAddr.
func (mr *MockTunConfiguratorMockRecorder) ChangeTunAddr(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeTunAddr", reflect.TypeOf((*MockTunConfigurator)(nil).ChangeTunAddr), arg0, arg1)
}
