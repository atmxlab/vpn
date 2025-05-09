package tunnel

import (
	"context"
	"time"

	"github.com/atmxlab/vpn/internal/protocol"
	"github.com/atmxlab/vpn/internal/server"
	"github.com/atmxlab/vpn/pkg/errors"
)

type SYNHandler struct {
	peerManager   PeerManager
	tunnel        Tunnel
	ipDistributor IpDistributor
	keepAliveTTL  time.Duration
}

func NewSYNHandler(
	peerManager PeerManager,
	tunnel Tunnel,
	ipDistributor IpDistributor,
	keepAliveTTL time.Duration,
) *SYNHandler {
	return &SYNHandler{
		peerManager:   peerManager,
		tunnel:        tunnel,
		ipDistributor: ipDistributor,
		keepAliveTTL:  keepAliveTTL,
	}
}

func (h *SYNHandler) Handle(ctx context.Context, packet *protocol.TunnelPacket) error {
	l := log(packet)

	l.Debug("Handle packet")

	has, err := h.peerManager.HasPeer(ctx, packet.Addr())
	if err != nil {
		return errors.Wrap(err, "peerManager.GetByAddr")
	}
	if has {
		return errors.AlreadyExistsf("peer already exists: addr=[%s]", packet.Addr())
	}

	acquiredIP, err := h.ipDistributor.AcquireIP()
	if err != nil {
		return errors.Wrap(err, "ipDistributor.AcquireIP")
	}

	peer := server.NewPeer(acquiredIP, packet.Addr())

	err = h.peerManager.Add(ctx, peer, h.keepAliveTTL, func(p *server.Peer) error {
		if localErr := h.ipDistributor.ReleaseIP(p.DedicatedIP()); localErr != nil {
			return errors.Wrap(localErr, "ipDistributor.ReleaseIP")
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "peerManager.Add")
	}

	l.
		WithField("DedicatedIP", peer.DedicatedIP()).
		Info("Created new peer")

	if _, err = h.tunnel.ACK(peer.Addr(), peer.DedicatedIP().To4()); err != nil {
		return errors.Wrap(err, "tunnel.ACK")
	}

	return nil
}
