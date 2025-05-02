package test

import (
	"context"
	"net"
	"testing"

	"github.com/atmxlab/vpn/internal/protocol"
	"github.com/atmxlab/vpn/internal/server"
)

type Engine interface {
	REPLAY(...Action)
}

type Action interface {
	Handle(a App)
}

type App interface {
	T() *testing.T
	Ctx() context.Context
	Tunnel() Tunnel
	Tun() Tun
	PeerManager() PeerManager
	IPDistributor() IPDistributor
}

type Tunnel interface {
	WriteToInput(p []byte, addr net.Addr) (n int, err error)
	GetLastPacket() (*protocol.TunnelPacket, bool)
}

type Tun interface {
	WriteToInput(p []byte) (n int, err error)
	GetLastPacket() (*protocol.TunPacket, bool)
}

type PeerManager interface {
	GetByAddr(ctx context.Context, addr net.Addr) (peer *server.Peer, exists bool, err error)
}

type IPDistributor interface {
	HasBusy() bool
}
