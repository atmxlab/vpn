package tunnel

import (
	"context"
	"time"

	"github.com/atmxlab/vpn/internal/protocol"
	"github.com/atmxlab/vpn/pkg/errors"
)

type KPAHandler struct {
	peerManager   PeerManager
	ipDistributor IpDistributor
	keepAliveTTL  time.Duration
}

func NewKPAHandler(peerManager PeerManager, keepAliveTTL time.Duration) *KPAHandler {
	return &KPAHandler{peerManager: peerManager, keepAliveTTL: keepAliveTTL}
}

func (h *KPAHandler) Handle(ctx context.Context, packet *protocol.TunnelPacket) error {
	l := log(packet)

	l.Debug("Handle packet")

	peer, err := h.peerManager.GetByAddr(ctx, packet.Addr())
	if err != nil {
		return errors.Wrap(err, "peerManager.GetByAddr")
	}

	if err = h.peerManager.Extend(ctx, peer, h.keepAliveTTL); err != nil {
		return errors.Wrap(err, "peerManager.Extend")
	}

	return nil
}
