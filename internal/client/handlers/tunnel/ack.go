package tunnel

import (
	"context"
	"net"

	"github.com/atmxlab/vpn/internal/protocol"
	"github.com/atmxlab/vpn/pkg/errors"
)

//go:generate mock TunConfigurator
type TunConfigurator interface {
	ChangeAddr(ctx context.Context, subnet net.IPNet) error
}

//go:generate mock NetConfigurator
type NetConfigurator interface {
	ConfigureRouting(ctx context.Context, subnet net.IPNet) error
}

//go:generate mock Signaller
type Signaller interface {
	Signal(ctx context.Context) error
}

type ACKHandler struct {
	tunConfigurator  TunConfigurator
	netConfigurator  NetConfigurator
	connectSignaller Signaller
	ipMasc           net.IPMask
}

func NewACKHandler(
	tunConfigurator TunConfigurator,
	netConfigurator NetConfigurator,
	ipMasc net.IPMask,
	connectSignaller Signaller,
) *ACKHandler {
	return &ACKHandler{
		tunConfigurator:  tunConfigurator,
		netConfigurator:  netConfigurator,
		ipMasc:           ipMasc,
		connectSignaller: connectSignaller,
	}
}

func (h *ACKHandler) Handle(ctx context.Context, packet *protocol.TunnelPacket) error {
	dedicatedIP, err := packet.Payload().IP()
	if err != nil {
		return errors.Wrap(err, "failed to decode ip address from payload")
	}

	subnet := net.IPNet{
		IP:   dedicatedIP,
		Mask: h.ipMasc,
	}

	if err = h.tunConfigurator.ChangeAddr(ctx, subnet); err != nil {
		return errors.Wrap(err, "failed to change tun subnet")
	}

	if err = h.netConfigurator.ConfigureRouting(ctx, subnet); err != nil {
		return errors.Wrap(err, "failed to configure routing")
	}

	if err = h.connectSignaller.Signal(ctx); err != nil {
		return errors.Wrap(err, "failed to signal connect")
	}

	return nil
}
