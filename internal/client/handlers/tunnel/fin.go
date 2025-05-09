package tunnel

import (
	"context"

	"github.com/atmxlab/vpn/internal/protocol"
	"github.com/atmxlab/vpn/pkg/errors"
)

//go:generate mock Closer
type Closer interface {
	Close() error
}

type FINHandler struct {
	closer Closer
}

func NewFINHandler(closer Closer) *FINHandler {
	return &FINHandler{closer: closer}
}

func (h *FINHandler) Handle(_ context.Context, _ *protocol.TunnelPacket) error {
	if err := h.closer.Close(); err != nil {
		return errors.Wrap(err, "failed to close client connection")
	}

	return nil
}
