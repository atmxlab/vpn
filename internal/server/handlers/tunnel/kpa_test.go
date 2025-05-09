package tunnel_test

import (
	"context"
	"testing"
	"time"

	"github.com/atmxlab/vpn/internal/server"
	"github.com/atmxlab/vpn/internal/server/handlers/tunnel"
	"github.com/atmxlab/vpn/internal/server/handlers/tunnel/mocks"
	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/atmxlab/vpn/test/gen"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestKPAHandler(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		kpaTTL := 1 * time.Minute
		tp := gen.RandTunnelPacket()
		peer := server.NewPeer(gen.RandIP(), tp.Addr())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		pm := mocks.NewMockPeerManager(ctrl)

		pm.EXPECT().
			GetByAddr(gomock.Any(), tp.Addr()).
			Return(peer, nil)

		pm.EXPECT().
			Extend(gomock.Any(), peer, kpaTTL).
			Return(nil)

		h := tunnel.NewKPAHandler(pm, kpaTTL)

		err := h.Handle(context.Background(), tp)
		require.NoError(t, err)
	})

	t.Run("peer not found", func(t *testing.T) {
		t.Parallel()

		kpaTTL := 1 * time.Minute
		tp := gen.RandTunnelPacket()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		pm := mocks.NewMockPeerManager(ctrl)

		pm.EXPECT().
			GetByAddr(gomock.Any(), tp.Addr()).
			Return(nil, errors.ErrNotFound)

		h := tunnel.NewKPAHandler(pm, kpaTTL)

		err := h.Handle(context.Background(), tp)
		require.ErrorIs(t, err, errors.ErrNotFound)
	})

	t.Run("extend error", func(t *testing.T) {
		t.Parallel()

		kpaTTL := 1 * time.Minute
		tp := gen.RandTunnelPacket()
		peer := server.NewPeer(gen.RandIP(), tp.Addr())

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		pm := mocks.NewMockPeerManager(ctrl)

		pm.EXPECT().
			GetByAddr(gomock.Any(), tp.Addr()).
			Return(peer, nil)

		pmErr := errors.New("peer manager error")

		pm.EXPECT().
			Extend(gomock.Any(), peer, kpaTTL).
			Return(pmErr)

		h := tunnel.NewKPAHandler(pm, kpaTTL)

		err := h.Handle(context.Background(), tp)
		require.ErrorIs(t, err, pmErr)
	})
}
