package tunnel_test

import (
	"context"
	"testing"

	"github.com/atmxlab/vpn/internal/server/handlers/tunnel"
	"github.com/atmxlab/vpn/internal/server/handlers/tunnel/mocks"
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

		pm := mocks.NewMockPeerManager(ctrl)
		pm.EXPECT().HasPeer(gomock.Any(), tp.Addr()).Return(true, nil)

		tnl := mocks.NewMockTunnel(ctrl)

		tn := mocks.NewMockTun(ctrl)
		tn.EXPECT().Write(tp.Payload()).Return(0, nil)

		h := tunnel.NewPSHHandler(pm, tn, tnl)

		err := h.Handle(context.Background(), tp)
		require.NoError(t, err)
	})

	t.Run("peer not found", func(t *testing.T) {
		t.Parallel()

		tp := gen.RandTunnelPacket()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		pm := mocks.NewMockPeerManager(ctrl)
		pm.EXPECT().HasPeer(gomock.Any(), tp.Addr()).Return(false, nil)

		tnl := mocks.NewMockTunnel(ctrl)
		tnl.EXPECT().SYN(tp.Addr(), nil).Return(0, nil)

		tn := mocks.NewMockTun(ctrl)

		h := tunnel.NewPSHHandler(pm, tn, tnl)

		err := h.Handle(context.Background(), tp)
		require.NoError(t, err)
	})
}
