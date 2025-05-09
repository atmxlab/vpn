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

func TestSYNHandler(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		tp := gen.RandTunnelPacket()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tunl := mocks.NewMockTunnel(ctrl)

		tunl.EXPECT().SYN(tp.Addr(), nil).Return(0, nil)

		h := tunnel.NewSYNHandler(tunl)

		err := h.Handle(context.Background(), tp)
		require.NoError(t, err)
	})

	t.Run("tunnel error", func(t *testing.T) {
		t.Parallel()

		tp := gen.RandTunnelPacket()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tunl := mocks.NewMockTunnel(ctrl)

		tunlErr := errors.New("tunl error")
		tunl.EXPECT().SYN(tp.Addr(), nil).Return(0, tunlErr)

		h := tunnel.NewSYNHandler(tunl)

		err := h.Handle(context.Background(), tp)
		require.ErrorIs(t, err, tunlErr)
	})
}
