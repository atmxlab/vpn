package tunnel

import (
	"context"

	"github.com/atmxlab/vpn/internal/protocol"
	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/sirupsen/logrus"
)

type FINHandler struct {
	peerManager   PeerManager
	ipDistributor IpDistributor
}

func NewFINHandler(peerManager PeerManager, ipDistributor IpDistributor) *FINHandler {
	return &FINHandler{peerManager: peerManager, ipDistributor: ipDistributor}
}

func (h *FINHandler) Handle(ctx context.Context, packet *protocol.TunnelPacket) error {
	peer, err := h.peerManager.GetByAddr(ctx, packet.Addr())
	if err != nil {
		return errors.Wrap(err, "peerManager.GetByAddr")
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
