package server

import (
	"context"
	"net"
	"time"
)

// PeerManager - управляющий пирами
type PeerManager interface {
	Add(ctx context.Context, peer *Peer, ttl time.Duration) error
	Remove(ctx context.Context, peer *Peer) error
	GetByAddrAndExtend(ctx context.Context, addr net.Addr, ttl time.Duration) (peer *Peer, exists bool, err error)
	GetByAddr(ctx context.Context, addr net.Addr) (peer *Peer, exists bool, err error)
	GetByDedicatedIP(ctx context.Context, ip net.IP) (peer *Peer, exists bool, err error)
	HasPeer(ctx context.Context, addr net.Addr) (bool, error)
}
