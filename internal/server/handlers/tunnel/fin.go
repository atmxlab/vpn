package tunnel

import (
	"context"

	"github.com/atmxlab/vpn/internal/protocol"
	"github.com/atmxlab/vpn/pkg/errors"
)

type FINHandler struct {
	peerManager   PeerManager
	ipDistributor IpDistributor
}

func NewFINHandler(peerManager PeerManager, ipDistributor IpDistributor) *FINHandler {
	return &FINHandler{peerManager: peerManager, ipDistributor: ipDistributor}
}

func (h *FINHandler) Handle(ctx context.Context, packet *protocol.TunnelPacket) error {
	l := log(packet)

	l.Debug("Handle packet")

	peer, err := h.peerManager.GetByAddr(ctx, packet.Addr())
	if err != nil {
		return errors.Wrap(err, "peerManager.GetByAddr")
	}

	if err = h.peerManager.Remove(ctx, peer); err != nil {
		return errors.Wrap(err, "peerManager.Remove")
	}

	if err = h.ipDistributor.ReleaseIP(peer.DedicatedIP()); err != nil {
		return errors.Wrap(err, "ipDistributor.ReleaseIP")
	}

	l.
		WithField("PeerAddr", peer.Addr()).
		WithField("DedicatedIP", peer.DedicatedIP()).
		Info("Removed peer and release dedicated ip")

	return nil
}
