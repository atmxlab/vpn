package tun

import (
	"context"
	"net"

	"github.com/atmxlab/vpn/internal/server"
)

// PeerManager - управляющий пирами
type PeerManager interface {
	GetByDedicatedIP(ctx context.Context, ip net.IP) (peer *server.Peer, exists bool, err error)
}
