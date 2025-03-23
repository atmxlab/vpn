package tunnel

import (
	"context"
	"time"

	"github.com/atmxlab/vpn/internal/protocol"
	"github.com/atmxlab/vpn/internal/server"
	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/sirupsen/logrus"
)

type KPAHandler struct {
	peerManager  server.PeerManager
	keepAliveTTL time.Duration
}

func (h *KPAHandler) Handle(ctx context.Context, packet *protocol.TunnelPacket) error {
	_, exists, err := h.peerManager.GetByAddrAndExtend(ctx, packet.Addr(), h.keepAliveTTL)
	if err != nil {
		return errors.Wrap(err, "get peer by addr")
	}
	if !exists {
		return errors.Wrap(errors.ErrNotFound, "peer not found")
	}

	logrus.Debugf("Keep alive: peer addr=[%s]", packet.Addr())

	return nil
}
