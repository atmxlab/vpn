package tunnel

import (
	"context"
	"time"

	"github.com/atmxlab/vpn/internal/protocol"
	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/sirupsen/logrus"
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
	peer, err := h.peerManager.GetByAddr(ctx, packet.Addr())
	if err != nil {
		return errors.Wrap(err, "peerManager.GetByAddr")
	}

	if err = h.peerManager.Extend(ctx, peer, h.keepAliveTTL); err != nil {
		return errors.Wrap(err, "peerManager.Extend")
	}

	logrus.Debugf("Keep alive: peer addr=[%s]", packet.Addr())

	return nil
}
