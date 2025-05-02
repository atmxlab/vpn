package tunnel_test

import (
	"context"
	"testing"
	"time"

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

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		pm := mocks.NewMockPeerManager(ctrl)

		pm.EXPECT().
			GetByAddrAndExtend(gomock.Any(), tp.Addr(), kpaTTL).
			Return(nil, true, nil)

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
			GetByAddrAndExtend(gomock.Any(), tp.Addr(), kpaTTL).
			Return(nil, false, nil)

		h := tunnel.NewKPAHandler(pm, kpaTTL)

		err := h.Handle(context.Background(), tp)
		require.ErrorIs(t, err, errors.ErrNotFound)
	})

	t.Run("peer manager error", func(t *testing.T) {
		t.Parallel()

		kpaTTL := 1 * time.Minute
		tp := gen.RandTunnelPacket()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		pm := mocks.NewMockPeerManager(ctrl)

		pmErr := errors.New("peer manager error")
		pm.EXPECT().
			GetByAddrAndExtend(gomock.Any(), tp.Addr(), kpaTTL).
			Return(nil, false, pmErr)

		h := tunnel.NewKPAHandler(pm, kpaTTL)

		err := h.Handle(context.Background(), tp)
		require.ErrorIs(t, err, pmErr)
	})
}
