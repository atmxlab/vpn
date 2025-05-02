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

func TestSYNHandler(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		kpaTTL := time.Minute
		acquiredIP := gen.RandIP()

		tp := gen.RandTunnelPacket()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		pm := mocks.NewMockPeerManager(ctrl)
		pm.EXPECT().HasPeer(gomock.Any(), tp.Addr()).Return(false, nil)

		pm.EXPECT().
			Add(gomock.Any(), gomock.Any(), kpaTTL, gomock.Any()).
			DoAndReturn(func(
				ctx context.Context,
				peer *server.Peer,
				kpa time.Duration,
				afterTTL ...func(*server.Peer) error,
			) error {
				require.Equal(t, tp.Addr(), peer.Addr())
				require.Equal(t, kpa, kpaTTL)
				require.Equal(t, acquiredIP, peer.DedicatedIP())

				return nil
			})

		ipd := mocks.NewMockIpDistributor(ctrl)
		ipd.EXPECT().AcquireIP().Return(acquiredIP, nil)

		tn := mocks.NewMockTunnel(ctrl)
		tn.EXPECT().ACK(tp.Addr(), acquiredIP).Return(0, nil)

		h := tunnel.NewSYNHandler(
			pm,
			tn,
			ipd,
			kpaTTL,
		)

		err := h.Handle(context.Background(), tp)
		require.NoError(t, err)
	})

	t.Run("peer already exists", func(t *testing.T) {
		t.Parallel()

		kpaTTL := time.Minute

		tp := gen.RandTunnelPacket()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		pm := mocks.NewMockPeerManager(ctrl)
		pm.EXPECT().HasPeer(gomock.Any(), tp.Addr()).Return(true, nil)

		ipd := mocks.NewMockIpDistributor(ctrl)
		tn := mocks.NewMockTunnel(ctrl)

		h := tunnel.NewSYNHandler(
			pm,
			tn,
			ipd,
			kpaTTL,
		)

		err := h.Handle(context.Background(), tp)
		require.ErrorIs(t, err, errors.ErrAlreadyExists)
	})
}
