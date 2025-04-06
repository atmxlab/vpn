// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/atmxlab/vpn/internal/server/handlers/tun (interfaces: PeerManager)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	net "net"
	reflect "reflect"

	server "github.com/atmxlab/vpn/internal/server"
	gomock "github.com/golang/mock/gomock"
)

// MockPeerManager is a mock of PeerManager interface.
type MockPeerManager struct {
	ctrl     *gomock.Controller
	recorder *MockPeerManagerMockRecorder
}

// MockPeerManagerMockRecorder is the mock recorder for MockPeerManager.
type MockPeerManagerMockRecorder struct {
	mock *MockPeerManager
}

// NewMockPeerManager creates a new mock instance.
func NewMockPeerManager(ctrl *gomock.Controller) *MockPeerManager {
	mock := &MockPeerManager{ctrl: ctrl}
	mock.recorder = &MockPeerManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPeerManager) EXPECT() *MockPeerManagerMockRecorder {
	return m.recorder
}

// GetByDedicatedIP mocks base method.
func (m *MockPeerManager) GetByDedicatedIP(arg0 context.Context, arg1 net.IP) (*server.Peer, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByDedicatedIP", arg0, arg1)
	ret0, _ := ret[0].(*server.Peer)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetByDedicatedIP indicates an expected call of GetByDedicatedIP.
func (mr *MockPeerManagerMockRecorder) GetByDedicatedIP(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByDedicatedIP", reflect.TypeOf((*MockPeerManager)(nil).GetByDedicatedIP), arg0, arg1)
}
