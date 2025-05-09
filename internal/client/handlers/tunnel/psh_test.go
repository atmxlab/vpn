package tunnel_test

import (
	"context"
	"testing"

	"github.com/atmxlab/vpn/internal/client/handlers/tunnel"
	"github.com/atmxlab/vpn/internal/client/handlers/tunnel/mocks"
	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/atmxlab/vpn/test/gen"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestPSHHandler(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		tp := gen.RandTunnelPacket()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tun := mocks.NewMockTun(ctrl)

		tun.EXPECT().Write(tp.Payload()).Return(len(tp.Payload()), nil)

		h := tunnel.NewPSHHandler(tun)

		err := h.Handle(context.Background(), tp)
		require.NoError(t, err)
	})

	t.Run("tun error", func(t *testing.T) {
		t.Parallel()

		tp := gen.RandTunnelPacket()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tun := mocks.NewMockTun(ctrl)

		tunErr := errors.New("tun error")
		tun.EXPECT().Write(tp.Payload()).Return(0, tunErr)

		h := tunnel.NewPSHHandler(tun)

		err := h.Handle(context.Background(), tp)
		require.ErrorIs(t, err, tunErr)
	})
}
