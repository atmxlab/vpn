package tun_test

import (
	"context"
	"net"
	"testing"

	"github.com/atmxlab/vpn/internal/server"
	"github.com/atmxlab/vpn/internal/server/handlers/tun"
	"github.com/atmxlab/vpn/internal/server/handlers/tun/mocks"
	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/atmxlab/vpn/test/gen"
	"github.com/atmxlab/vpn/test/stub"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHandle(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		var (
			ctx         = context.Background()
			peer        = gen.RandPeer()
			dedicatedIP = gen.RandIP()
			tunPacket   = gen.RandTunICMPReq(t, func(h *stub.IPHeader) {
				h.Dst = dedicatedIP
			})
		)

		ctrl := gomock.NewController(t)
		peerManager := mocks.NewMockPeerManager(ctrl)
		peerManager.
			EXPECT().
			GetByDedicatedIP(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, dIP net.IP) (*server.Peer, bool, error) {
				require.True(t, dedicatedIP.Equal(dIP))
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
			ctx         = context.Background()
			dedicatedIP = gen.RandIP()
			tunPacket   = gen.RandTunICMPReq(t, func(h *stub.IPHeader) {
				h.Dst = dedicatedIP
			})
		)

		ctrl := gomock.NewController(t)
		peerManager := mocks.NewMockPeerManager(ctrl)
		peerManagerError := errors.New("test error")
		peerManager.
			EXPECT().
			GetByDedicatedIP(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, dIP net.IP) (*server.Peer, bool, error) {
				require.True(t, dedicatedIP.Equal(dIP))
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
			ctx         = context.Background()
			dedicatedIP = gen.RandIP()
			tunPacket   = gen.RandTunICMPReq(t, func(h *stub.IPHeader) {
				h.Dst = dedicatedIP
			})
		)

		ctrl := gomock.NewController(t)
		peerManager := mocks.NewMockPeerManager(ctrl)
		peerManager.
			EXPECT().
			GetByDedicatedIP(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, dIP net.IP) (*server.Peer, bool, error) {
				require.True(t, dedicatedIP.Equal(dIP))
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
			ctx         = context.Background()
			peer        = gen.RandPeer()
			dedicatedIP = gen.RandIP()
			tunPacket   = gen.RandTunICMPReq(t, func(h *stub.IPHeader) {
				h.Dst = dedicatedIP
			})
		)

		ctrl := gomock.NewController(t)
		peerManager := mocks.NewMockPeerManager(ctrl)
		peerManager.
			EXPECT().
			GetByDedicatedIP(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, dIP net.IP) (*server.Peer, bool, error) {
				require.True(t, dedicatedIP.Equal(dedicatedIP))
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
