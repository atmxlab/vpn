package tunnel

import (
	"context"
	"io"

	"github.com/atmxlab/vpn/internal/pkg/ip"
	"github.com/atmxlab/vpn/internal/protocol"
	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/sirupsen/logrus"
)

//go:generate mock Tun
type Tun interface {
	io.Writer
}

type PSHHandler struct {
	tun Tun
}

func NewPSHHandler(tun Tun) *PSHHandler {
	return &PSHHandler{tun: tun}
}

func (h *PSHHandler) Handle(_ context.Context, packet *protocol.TunnelPacket) error {
	ip.LogHeader(packet.Payload())

	n, err := h.tun.Write(packet.Payload())
	if err != nil {
		return errors.Wrap(err, "tun.Write")
	}

	logrus.Debugf("Write bytes to TUN: len=[%d]", n)

	return nil
}
