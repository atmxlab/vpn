package tunnel

import (
	"context"
	"time"

	"github.com/atmxlab/vpn/internal/protocol"
	"github.com/atmxlab/vpn/internal/server"
	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/sirupsen/logrus"
)

type SYNHandler struct {
	tunnel        Tunnel
	peerManager   server.PeerManager
	ipDistributor IpDistributor
	keepAliveTTL  time.Duration
}

func (h *SYNHandler) Handle(ctx context.Context, packet *protocol.TunnelPacket) error {
	_, exists, err := h.peerManager.GetByAddr(ctx, packet.Addr())
	if err != nil {
		return errors.Wrap(err, "peerManager.GetByAddr")
	}
	if exists {
		return errors.Wrap(errors.ErrNotFound, "peerManager.GetByAddr not found")
	}

	acquiredIP, err := h.ipDistributor.AcquireIP()
	if err != nil {
		return errors.Wrap(err, "ipDistributor.AcquireIP")
	}

	peer := server.NewPeer(acquiredIP, packet.Addr())

	logrus.Infof("Created new peer with addr: %s and dedicated ip: %s", peer.Addr(), peer.DedicatedIP())

	if err = h.peerManager.Add(ctx, peer, h.keepAliveTTL); err != nil {
		return errors.Wrap(err, "peerManager.Add")
	}

	if _, err = h.tunnel.ACK(peer.Addr(), peer.DedicatedIP().To4()); err != nil {
		return errors.Wrap(err, "tunnel.ACK")
	}

	return nil
}
