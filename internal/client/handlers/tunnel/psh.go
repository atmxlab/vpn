package tunnel

import (
	"context"
	"io"

	"github.com/atmxlab/vpn/internal/pkg/ip"
	"github.com/atmxlab/vpn/internal/protocol"
	"github.com/atmxlab/vpn/pkg/errors"
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
	l := log(packet)

	l.Debug("Handle packet")

	ip.LogHeader(packet.Payload())

	_, err := h.tun.Write(packet.Payload())
	if err != nil {
		return errors.Wrap(err, "tun.Write")
	}

	return nil
}
