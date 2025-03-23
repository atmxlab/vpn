package peermanager

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/atmxlab/vpn/internal/server"
)

type Manager struct {
	mu                 sync.RWMutex
	indexByDedicatedIP map[string]*server.Peer
	indexByAddress     map[string]*server.Peer
}

func New() *Manager {
	return &Manager{
		indexByDedicatedIP: make(map[string]*server.Peer),
		indexByAddress:     make(map[string]*server.Peer),
	}
}

func (pm *Manager) Add(_ context.Context, peer *server.Peer, _ time.Duration) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.indexByDedicatedIP[peer.DedicatedIP().String()] = peer
	pm.indexByAddress[peer.Addr().String()] = peer

	return nil
}

func (pm *Manager) Remove(_ context.Context, peer *server.Peer) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	delete(pm.indexByDedicatedIP, peer.DedicatedIP().String())
	delete(pm.indexByAddress, peer.Addr().String())

	return nil
}

func (pm *Manager) GetByDedicatedIP(_ context.Context, ip net.IP) (*server.Peer, bool, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	peer, ok := pm.indexByDedicatedIP[ip.String()]
	return peer, ok, nil
}

func (pm *Manager) GetByAddr(_ context.Context, addr net.Addr) (*server.Peer, bool, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	return pm.getByAddrLocked(addr)
}

func (pm *Manager) GetByAddrAndExtend(_ context.Context, addr net.Addr, _ time.Duration) (
	*server.Peer,
	bool,
	error,
) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	return pm.getByAddrLocked(addr)
}

func (pm *Manager) HasPeer(_ context.Context, addr net.Addr) (bool, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	_, ok, err := pm.getByAddrLocked(addr)

	return ok, err
}

func (pm *Manager) getByAddrLocked(addr net.Addr) (*server.Peer, bool, error) {
	peer, ok := pm.indexByAddress[addr.String()]
	return peer, ok, nil
}
