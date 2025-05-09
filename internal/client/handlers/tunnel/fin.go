package tunnel

import (
	"context"

	"github.com/atmxlab/vpn/internal/protocol"
	"github.com/atmxlab/vpn/pkg/errors"
)

//go:generate mock Stopper
type Stopper interface {
	Stop(ctx context.Context) error
}

type FINHandler struct {
	stopper Stopper
}

func (h *FINHandler) Handle(ctx context.Context, _ *protocol.TunnelPacket) error {
	if err := h.stopper.Stop(ctx); err != nil {
		return errors.Wrap(err, "failed to stop client connection")
	}

	return nil
}
