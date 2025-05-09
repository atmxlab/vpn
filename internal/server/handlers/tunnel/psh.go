package tunnel

import (
	"context"

	"github.com/atmxlab/vpn/internal/pkg/ip"
	"github.com/atmxlab/vpn/internal/protocol"
	"github.com/atmxlab/vpn/pkg/errors"
)

type PSHHandler struct {
	peerManager PeerManager
	tun         Tun
	tunnel      Tunnel
}

func NewPSHHandler(peerManager PeerManager, tun Tun, tunnel Tunnel) *PSHHandler {
	return &PSHHandler{peerManager: peerManager, tun: tun, tunnel: tunnel}
}

func (h *PSHHandler) Handle(ctx context.Context, packet *protocol.TunnelPacket) error {
	l := log(packet)

	l.Debug("Handle packet")

	ip.LogHeader(packet.Payload())

	has, err := h.peerManager.HasPeer(ctx, packet.Addr())
	if err != nil {
		return errors.Wrap(err, "peerManager.HasPeer")
	}

	if !has {
		l.Warn("Peer not found")

		// After server syn, client must init syn
		_, err = h.tunnel.SYN(packet.Addr(), nil)
		if err != nil {
			return errors.Wrap(err, "tunnel.SYN")
		}

		return nil
	}

	_, err = h.tun.Write(packet.Payload())
	if err != nil {
		return errors.Wrap(err, "tun.Write")
	}

	return nil
}
