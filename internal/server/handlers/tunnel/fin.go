package tunnel

import (
	"context"

	"github.com/atmxlab/vpn/internal/protocol"
	"github.com/atmxlab/vpn/internal/server"
	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/sirupsen/logrus"
)

type FINHandler struct {
	peerManager   server.PeerManager
	ipDistributor IpDistributor
}

func (h *FINHandler) Handle(ctx context.Context, packet *protocol.TunnelPacket) error {
	peer, exists, err := h.peerManager.FindByAddr(ctx, packet.Addr())
	if err != nil {
		return errors.Wrap(err, "peerManager.FindByAddr")
	}
	if !exists {
		return errors.Wrap(errors.ErrNotFound, "peer not found")
	}

	if err = h.peerManager.Remove(ctx, peer); err != nil {
		return errors.Wrap(err, "peerManager.Remove")
	}

	logrus.Infof("Removed peer: addr=[%s], dedicated ip=[%s]", peer.Addr(), peer.DedicatedIP())

	if err = h.ipDistributor.ReleaseIP(peer.DedicatedIP()); err != nil {
		return errors.Wrap(err, "ipDistributor.ReleaseIP")
	}

	return nil
}
