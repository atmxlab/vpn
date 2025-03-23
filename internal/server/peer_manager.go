package server

import (
	"context"
	"net"
)

// PeerManager - управляющий пирами
type PeerManager interface {
	Add(ctx context.Context, peer *Peer) error
	Remove(ctx context.Context, peer *Peer) error
	FindByAddr(ctx context.Context, addr net.Addr) (peer *Peer, exists bool, err error)
	FindByDedicatedIP(ctx context.Context, ip net.IP) (peer *Peer, exists bool, err error)
}
