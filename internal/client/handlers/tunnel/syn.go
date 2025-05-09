package tunnel

import (
	"context"

	"github.com/atmxlab/vpn/internal/protocol"
	"github.com/atmxlab/vpn/pkg/errors"
)

type SYNHandler struct {
	tunnel Tunnel
}

func NewSYNHandler(tunnel Tunnel) *SYNHandler {
	return &SYNHandler{tunnel: tunnel}
}

func (h *SYNHandler) Handle(_ context.Context, packet *protocol.TunnelPacket) error {
	l := log(packet)

	l.Debug("Handle packet")

	_, err := h.tunnel.SYN(packet.Addr(), nil)
	if err != nil {
		return errors.Wrap(err, "failed to syn")
	}

	return nil
}
