package tun

import (
	"context"
	"net"

	"github.com/atmxlab/vpn/internal/protocol"
	"github.com/atmxlab/vpn/internal/server"
	"github.com/atmxlab/vpn/pkg/errors"
	"golang.org/x/net/ipv4"
)

type Tunnel interface {
	PSH(addr net.Addr, payload []byte) (int, error)
}

type Handler struct {
	tunnel      Tunnel
	peerManager server.PeerManager
}

func NewHandler(tunnel Tunnel, peerManager server.PeerManager) *Handler {
	return &Handler{tunnel: tunnel, peerManager: peerManager}
}

func (h *Handler) Handle(
	ctx context.Context,
	packet *protocol.TunPacket,
) error {
	header, err := ipv4.ParseHeader(packet.Payload())
	if err != nil {
		return errors.Wrap(err, "ipv4.ParseHeader")
	}

	peer, exists, err := h.peerManager.FindByDedicatedIP(ctx, header.Dst)
	if err != nil {
		return errors.Wrap(err, "peerManager.FindByDedicatedIP")
	}
	if !exists {
		return errors.Wrap(errors.ErrNotFound, "peerManager.FindByDedicatedIP not found")
	}

	_, err = h.tunnel.PSH(peer.Addr(), packet.Payload())
	if err != nil {
		return errors.Wrap(err, "tunnel.PSH")
	}

	return nil
}
