package tun

import (
	"context"
	"net"

	"github.com/atmxlab/vpn/internal/pkg/ip"
	"github.com/atmxlab/vpn/internal/protocol"
	"github.com/atmxlab/vpn/pkg/errors"
)

type Handler struct {
	tunnel     Tunnel
	serverAddr net.Addr
}

func NewHandler(tunnel Tunnel, serverAddr net.Addr) *Handler {
	return &Handler{tunnel: tunnel, serverAddr: serverAddr}
}

func (h *Handler) Handle(
	_ context.Context,
	packet *protocol.TunPacket,
) error {
	ip.LogHeader(packet.Payload())

	_, err := h.tunnel.PSH(h.serverAddr, packet.Payload())
	if err != nil {
		return errors.Wrap(err, "tunnel.PSH")
	}

	return nil
}
