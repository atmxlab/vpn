package tunnel

import (
	"context"

	"github.com/atmxlab/vpn/internal/pkg/ip"
	"github.com/atmxlab/vpn/internal/protocol"
	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/sirupsen/logrus"
)

type PSHHandler struct {
	tun Tun
}

func (h *PSHHandler) Handle(_ context.Context, packet *protocol.TunnelPacket) error {
	ip.LogHeader(packet.Payload())

	// TODO: проверять есть ли пир.
	//  Если нет, отправлять флаг с ошибкой, чтобы клиент знал, что надо переподключиться

	n, err := h.tun.Write(packet.Payload())
	if err != nil {
		return errors.Wrap(err, "tun.Write")
	}

	logrus.Debugf("Write bytes to TUN: %d", n)

	return nil
}
