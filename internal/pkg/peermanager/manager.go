package peermanager

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/atmxlab/vpn/internal/server"
	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/sirupsen/logrus"
)

type peer struct {
	peer       *server.Peer
	ttl        time.Duration
	extendChan chan struct{}
	doneChan   chan struct{}

	afterTTL []func(*server.Peer) error
}

func newPeer(p *server.Peer, ttl time.Duration, afterTTL []func(*server.Peer) error) *peer {
	return &peer{
		peer:       p,
		ttl:        ttl,
		afterTTL:   afterTTL,
		extendChan: make(chan struct{}, 1),
		doneChan:   make(chan struct{}, 1),
	}
}

type Manager struct {
	mu                 sync.RWMutex
	indexByDedicatedIP map[string]*peer
	indexByAddress     map[string]*peer
}

func New() *Manager {
	return &Manager{
		indexByDedicatedIP: make(map[string]*peer),
		indexByAddress:     make(map[string]*peer),
	}
}

func (pm *Manager) Add(
	ctx context.Context,
	peer *server.Peer,
	ttl time.Duration,
	afterTTL ...func(p *server.Peer) error,
) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	p := newPeer(peer, ttl, afterTTL)
	pm.indexByDedicatedIP[peer.DedicatedIP().String()] = p
	pm.indexByAddress[peer.Addr().String()] = p

	go func() {
		if err := pm.monitorPeer(ctx, p); err != nil {
			logrus.Errorf("monitor peer error: addr=[%s], err=[%s]", p.peer.Addr(), err.Error())
		}
	}()

	return nil
}

func (pm *Manager) Remove(ctx context.Context, peer *server.Peer) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	p, ok, err := pm.getByAddrLocked(peer.Addr())
	if err != nil {
		return errors.Wrap(err, "get peer by addr locked")
	}
	if !ok {
		return nil
	}

	p.doneChan <- struct{}{}

	delete(pm.indexByDedicatedIP, peer.DedicatedIP().String())
	delete(pm.indexByAddress, peer.Addr().String())

	return nil
}

func (pm *Manager) GetByDedicatedIP(_ context.Context, ip net.IP) (*server.Peer, bool, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	p, ok := pm.indexByDedicatedIP[ip.String()]
	return p.peer, ok, nil
}

func (pm *Manager) GetByAddr(_ context.Context, addr net.Addr) (*server.Peer, bool, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	if p, ok, err := pm.getByAddrLocked(addr); ok {
		return p.peer, ok, err
	}

	return nil, false, nil
}

func (pm *Manager) GetByAddrAndExtend(_ context.Context, addr net.Addr, _ time.Duration) (
	*server.Peer,
	bool,
	error,
) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	if p, ok, err := pm.getByAddrLocked(addr); ok {
		p.extendChan <- struct{}{}
		return p.peer, ok, err
	}

	return nil, false, nil
}

func (pm *Manager) HasPeer(_ context.Context, addr net.Addr) (bool, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	_, ok, err := pm.getByAddrLocked(addr)

	return ok, err
}

func (pm *Manager) getByAddrLocked(addr net.Addr) (*peer, bool, error) {
	p, ok := pm.indexByAddress[addr.String()]
	return p, ok, nil
}

func (pm *Manager) monitorPeer(ctx context.Context, p *peer) error {
	timer := time.NewTimer(p.ttl)
	defer timer.Stop()

	for {
		select {
		case <-p.extendChan:
			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(p.ttl)
		case <-timer.C:
			if err := pm.Remove(ctx, p.peer); err != nil {
				return errors.Wrap(err, "failed to remove peer")
			}

			for _, hook := range p.afterTTL {
				if err := hook(p.peer); err != nil {
					return errors.Wrap(err, "failed to execute after TTL hook")
				}
			}

			return nil
		case <-p.doneChan:
			return nil
		}
	}
}
