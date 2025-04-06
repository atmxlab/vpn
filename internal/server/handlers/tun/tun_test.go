package tun_test

import (
	"context"
	"net"
	"testing"

	"github.com/atmxlab/vpn/internal/protocol"
	"github.com/atmxlab/vpn/internal/server"
	"github.com/atmxlab/vpn/internal/server/handlers/tun"
	"github.com/atmxlab/vpn/internal/server/handlers/tun/mocks"
	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/atmxlab/vpn/test"
	"github.com/atmxlab/vpn/test/gen"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestHandle(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		var (
			ctx      = context.Background()
			ipPacket = test.NewIPPacketBuilder(t).
					SrcIP(gen.RandIP()).
					DstIP(gen.RandIP()).
					Payload(gen.RandPayload()).
					TCP().
					Build()
			peer      = server.NewPeer(ipPacket.DstIP(), gen.RandAddr())
			tunPacket = protocol.NewTunPacket(ipPacket.Bytes())
		)

		ctrl := gomock.NewController(t)
		peerManager := mocks.NewMockPeerManager(ctrl)
		peerManager.
			EXPECT().
			GetByDedicatedIP(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, dedicatedIP net.IP) (*server.Peer, bool, error) {
				require.True(t, ipPacket.DstIP().Equal(dedicatedIP))
				return peer, true, nil
			})

		tunl := mocks.NewMockTunnel(ctrl)
		tunl.EXPECT().PSH(peer.Addr(), tunPacket.Payload()).Return(len(tunPacket.Payload()), nil)

		handler := tun.NewHandler(tunl, peerManager)

		err := handler.Handle(ctx, tunPacket)
		require.NoError(t, err)
	})

	t.Run("get peer by dedicated ip error", func(t *testing.T) {
		t.Parallel()

		var (
			ctx      = context.Background()
			ipPacket = test.NewIPPacketBuilder(t).
					SrcIP(gen.RandIP()).
					DstIP(gen.RandIP()).
					Payload(gen.RandPayload()).
					TCP().
					Build()
			tunPacket = protocol.NewTunPacket(ipPacket.Bytes())
		)

		ctrl := gomock.NewController(t)
		peerManager := mocks.NewMockPeerManager(ctrl)
		peerManagerError := errors.New("test error")
		peerManager.
			EXPECT().
			GetByDedicatedIP(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, dedicatedIP net.IP) (*server.Peer, bool, error) {
				require.True(t, ipPacket.DstIP().Equal(dedicatedIP))
				return nil, false, peerManagerError
			})

		tunl := mocks.NewMockTunnel(ctrl)

		handler := tun.NewHandler(tunl, peerManager)

		err := handler.Handle(ctx, tunPacket)
		require.ErrorIs(t, err, peerManagerError)
	})

	t.Run("peer not found", func(t *testing.T) {
		t.Parallel()

		var (
			ctx      = context.Background()
			ipPacket = test.NewIPPacketBuilder(t).
					SrcIP(gen.RandIP()).
					DstIP(gen.RandIP()).
					Payload(gen.RandPayload()).
					TCP().
					Build()
			tunPacket = protocol.NewTunPacket(ipPacket.Bytes())
		)

		ctrl := gomock.NewController(t)
		peerManager := mocks.NewMockPeerManager(ctrl)
		peerManager.
			EXPECT().
			GetByDedicatedIP(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, dedicatedIP net.IP) (*server.Peer, bool, error) {
				require.True(t, ipPacket.DstIP().Equal(dedicatedIP))
				return nil, false, nil
			})

		tunl := mocks.NewMockTunnel(ctrl)

		handler := tun.NewHandler(tunl, peerManager)

		err := handler.Handle(ctx, tunPacket)
		require.ErrorIs(t, err, errors.ErrNotFound)
	})

	t.Run("write to tunnel error", func(t *testing.T) {
		t.Parallel()

		var (
			ctx      = context.Background()
			ipPacket = test.NewIPPacketBuilder(t).
					SrcIP(gen.RandIP()).
					DstIP(gen.RandIP()).
					Payload(gen.RandPayload()).
					TCP().
					Build()
			peer      = server.NewPeer(ipPacket.DstIP(), gen.RandAddr())
			tunPacket = protocol.NewTunPacket(ipPacket.Bytes())
		)

		ctrl := gomock.NewController(t)
		peerManager := mocks.NewMockPeerManager(ctrl)
		peerManager.
			EXPECT().
			GetByDedicatedIP(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, dedicatedIP net.IP) (*server.Peer, bool, error) {
				require.True(t, ipPacket.DstIP().Equal(dedicatedIP))
				return peer, true, nil
			})

		tunl := mocks.NewMockTunnel(ctrl)
		tunlError := errors.New("test error")
		tunl.EXPECT().PSH(peer.Addr(), tunPacket.Payload()).Return(0, tunlError)

		handler := tun.NewHandler(tunl, peerManager)

		err := handler.Handle(ctx, tunPacket)
		require.ErrorIs(t, err, tunlError)
	})
}
