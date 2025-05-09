package tunnel

import (
	"context"
	"net"
	"time"

	"github.com/atmxlab/vpn/internal/server"
)

// PeerManager - управляющий пирами
//
//go:generate mock PeerManager
type PeerManager interface {
	Add(
		ctx context.Context,
		peer *server.Peer,
		ttl time.Duration,
		afterTTL ...func(p *server.Peer) error,
	) error
	Remove(ctx context.Context, peer *server.Peer) error
	Extend(ctx context.Context, peer *server.Peer, ttl time.Duration) (err error)
	GetByAddr(ctx context.Context, addr net.Addr) (peer *server.Peer, err error)
	HasPeer(ctx context.Context, addr net.Addr) (bool, error)
}
