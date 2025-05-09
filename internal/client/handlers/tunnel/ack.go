package tunnel

import (
	"context"
	"net"

	"github.com/atmxlab/vpn/internal/protocol"
	"github.com/atmxlab/vpn/pkg/errors"
)

type TunConfigurator interface {
	ChangeAddr(ctx context.Context, subnet net.IPNet) error
}

type NetConfigurator interface {
	ConfigureRouting(ctx context.Context, subnet net.IPNet) error
}

type ACKHandler struct {
	tunConfigurator TunConfigurator
	netConfigurator NetConfigurator
	ipMasc          net.IPMask
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

	return nil
}
