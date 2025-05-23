// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/atmxlab/vpn/internal/client/handlers/tunnel (interfaces: Tun)
//
// Generated by this command:
//
//	mockgen -destination=/Users/timur_abdurashidov/Desktop/other/vpn/internal/client/handlers/tunnel/mocks/Tun_mock.go -package=mocks . Tun
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockTun is a mock of Tun interface.
type MockTun struct {
	ctrl     *gomock.Controller
	recorder *MockTunMockRecorder
}

// MockTunMockRecorder is the mock recorder for MockTun.
type MockTunMockRecorder struct {
	mock *MockTun
}

// NewMockTun creates a new mock instance.
func NewMockTun(ctrl *gomock.Controller) *MockTun {
	mock := &MockTun{ctrl: ctrl}
	mock.recorder = &MockTunMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTun) EXPECT() *MockTunMockRecorder {
	return m.recorder
}

// Write mocks base method.
func (m *MockTun) Write(arg0 []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Write indicates an expected call of Write.
func (mr *MockTunMockRecorder) Write(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockTun)(nil).Write), arg0)
}
