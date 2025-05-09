package tun_test

import (
	"context"
	"testing"

	"github.com/atmxlab/vpn/internal/client/handlers/tun"
	"github.com/atmxlab/vpn/internal/server/handlers/tun/mocks"
	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/atmxlab/vpn/test/gen"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHandle(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		tp := gen.RandTunPacket()
		serverAddr := gen.RandAddr()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tunl := mocks.NewMockTunnel(ctrl)

		tunl.EXPECT().PSH(serverAddr, tp.Payload()).Return(len(tp.Payload()), nil)

		h := tun.NewHandler(tunl, serverAddr)

		err := h.Handle(context.Background(), tp)
		require.NoError(t, err)
	})

	t.Run("tunnel error", func(t *testing.T) {
		t.Parallel()

		tp := gen.RandTunPacket()
		serverAddr := gen.RandAddr()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tunl := mocks.NewMockTunnel(ctrl)

		tunlErr := errors.New("test error")
		tunl.EXPECT().PSH(serverAddr, tp.Payload()).Return(0, tunlErr)

		h := tun.NewHandler(tunl, serverAddr)

		err := h.Handle(context.Background(), tp)
		require.ErrorIs(t, err, tunlErr)
	})
}
