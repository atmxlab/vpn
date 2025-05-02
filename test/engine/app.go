package engine

import (
	"context"
	"testing"

	"github.com/atmxlab/vpn/test"
)

type app struct {
	t             *testing.T
	ctx           context.Context
	tunnel        test.Tunnel
	tun           test.Tun
	peerManager   test.PeerManager
	ipDistributor test.IPDistributor
}

func newApp(
	t *testing.T,
	ctx context.Context,
	tunnel test.Tunnel,
	tun test.Tun,
	peerManager test.PeerManager,
	ipDistributor test.IPDistributor,
) *app {
	return &app{
		t:             t,
		ctx:           ctx,
		tunnel:        tunnel,
		tun:           tun,
		peerManager:   peerManager,
		ipDistributor: ipDistributor,
	}
}

func (a *app) T() *testing.T {
	return a.t
}

func (a *app) Ctx() context.Context {
	return a.ctx
}

func (a *app) Tunnel() test.Tunnel {
	return a.tunnel
}

func (a *app) Tun() test.Tun {
	return a.tun
}

func (a *app) PeerManager() test.PeerManager {
	return a.peerManager
}

func (a *app) IPDistributor() test.IPDistributor {
	return a.ipDistributor
}
