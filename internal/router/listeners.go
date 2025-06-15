package router

import (
	"context"

	"github.com/atmxlab/vpn/internal/protocol"
	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (r *Router) listenTun(ctx context.Context) error {
	log := logrus.WithField("Namespace", "TUN")

	for {
		buf := make([]byte, r.cfg.bufferSize)

		n, err := r.tun.ReadWithContext(ctx, buf)
		if err != nil {
			return errors.Wrap(err, "failed to read from tun")
		}

		payload := make([]byte, n)
		copy(payload, buf[:n])

		select {
		case <-ctx.Done():
			log.
				WithError(ctx.Err()).
				Debug("Stop listening because context canceled")
			return ctx.Err()
		case r.tunPackets <- protocol.NewTunPacket(buf[:n]):
		}
	}
}

func (r *Router) listenTunnel(ctx context.Context) error {
	log := logrus.WithField("Namespace", "TUNNEL")
	for {
		buf := make([]byte, r.cfg.bufferSize)

		n, addr, err := r.tunnel.ReadFromWithContext(ctx, buf)
		if err != nil {
			return errors.Wrap(err, "failed to read from tunnel")
		}

		payload := make([]byte, n)
		copy(payload, buf[:n])

		select {
		case <-ctx.Done():
			log.
				WithError(ctx.Err()).
				Debug("Stop listening because context canceled")
			return ctx.Err()
		case r.tunnelPackets <- protocol.UnmarshalTunnelPacket(addr, payload):
		}
	}
}
