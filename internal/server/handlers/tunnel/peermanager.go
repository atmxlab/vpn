package tunnel

import (
	"context"
	"net"
	"time"

	"github.com/atmxlab/vpn/internal/server"
)

// PeerManager - управляющий пирами
type PeerManager interface {
	Add(ctx context.Context, peer *server.Peer, ttl time.Duration) error
	Remove(ctx context.Context, peer *server.Peer) error
	GetByAddrAndExtend(ctx context.Context, addr net.Addr, ttl time.Duration) (
		peer *server.Peer,
		exists bool,
		err error,
	)
	GetByAddr(ctx context.Context, addr net.Addr) (peer *server.Peer, exists bool, err error)
	HasPeer(ctx context.Context, addr net.Addr) (bool, error)
}
