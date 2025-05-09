package tunnel_test

import (
	"context"
	"testing"

	"github.com/atmxlab/vpn/internal/server/handlers/tunnel"
	"github.com/atmxlab/vpn/internal/server/handlers/tunnel/mocks"
	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/atmxlab/vpn/test/gen"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestFINHandler(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		tp := gen.RandTunnelPacket()
		peer := gen.RandPeer()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		pm := mocks.NewMockPeerManager(ctrl)
		pm.EXPECT().GetByAddr(gomock.Any(), tp.Addr()).Return(peer, nil)
		pm.EXPECT().Remove(gomock.Any(), peer).Return(nil)

		ipd := mocks.NewMockIpDistributor(ctrl)
		ipd.EXPECT().ReleaseIP(peer.DedicatedIP()).Return(nil)

		h := tunnel.NewFINHandler(pm, ipd)

		err := h.Handle(context.Background(), tp)
		require.NoError(t, err)
	})

	t.Run("peer not found", func(t *testing.T) {
		t.Parallel()

		tp := gen.RandTunnelPacket()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		pm := mocks.NewMockPeerManager(ctrl)
		pm.EXPECT().GetByAddr(gomock.Any(), tp.Addr()).Return(nil, errors.ErrNotFound)

		ipd := mocks.NewMockIpDistributor(ctrl)

		h := tunnel.NewFINHandler(pm, ipd)

		err := h.Handle(context.Background(), tp)
		require.ErrorIs(t, err, errors.ErrNotFound)
	})
}
