package tunnel

import (
	"context"
	"net"

	"github.com/atmxlab/vpn/internal/protocol"
	"github.com/atmxlab/vpn/pkg/errors"
)

//go:generate mock TunConfigurator
type TunConfigurator interface {
	ChangeTunAddr(ctx context.Context, subnet net.IPNet) error
}

//go:generate mock Signaller
type Signaller interface {
	Signal(ctx context.Context) error
}

type ACKHandler struct {
	tunConfigurator  TunConfigurator
	connectSignaller Signaller
	ipMasc           net.IPMask
}

func NewACKHandler(
	tunConfigurator TunConfigurator,
	connectSignaller Signaller,
	ipMasc net.IPMask,
) *ACKHandler {
	return &ACKHandler{
		tunConfigurator:  tunConfigurator,
		connectSignaller: connectSignaller,
		ipMasc:           ipMasc,
	}
}

func (h *ACKHandler) Handle(ctx context.Context, packet *protocol.TunnelPacket) error {
	l := log(packet)

	l.Debug("Handle packet")

	dedicatedIP, err := packet.Payload().IP()
	if err != nil {
		return errors.Wrap(err, "failed to decode ip address from payload")
	}

	l.
		WithField("DedicatedIP", dedicatedIP).
		Debug("Got dedicated IP from payload")

	subnet := net.IPNet{
		IP:   dedicatedIP,
		Mask: h.ipMasc,
	}

	if err = h.tunConfigurator.ChangeTunAddr(ctx, subnet); err != nil {
		return errors.Wrap(err, "failed to change tun subnet")
	}

	if err = h.connectSignaller.Signal(ctx); errors.IsSomeBut(err, errors.ErrAlreadyExists) {
		return errors.Wrap(err, "failed to signal connect")
	}

	return nil
}
