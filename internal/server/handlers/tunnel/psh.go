package tunnel

import (
	"context"

	"github.com/atmxlab/vpn/internal/pkg/ip"
	"github.com/atmxlab/vpn/internal/protocol"
	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/sirupsen/logrus"
)

type PSHHandler struct {
	peerManager PeerManager
	tun         Tun
	tunnel      Tunnel
}

func (h *PSHHandler) Handle(ctx context.Context, packet *protocol.TunnelPacket) error {
	ip.LogHeader(packet.Payload())

	has, err := h.peerManager.HasPeer(ctx, packet.Addr())
	if err != nil {
		return errors.Wrap(err, "peerManager.HasPeer")
	}

	if !has {
		logrus.Warnf("Peer not found: addr=[%s]", packet.Addr())

		_, err = h.tunnel.SYN(packet.Addr(), nil)
		if err != nil {
			return errors.Wrap(err, "tunnel.SYN")
		}

		return nil
	}

	n, err := h.tun.Write(packet.Payload())
	if err != nil {
		return errors.Wrap(err, "tun.Write")
	}

	logrus.Debugf("Write bytes to TUN: len=[%d]", n)

	return nil
}
