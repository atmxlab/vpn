package tunnel_test

import (
	"context"
	"net"
	"testing"

	"github.com/atmxlab/vpn/internal/client/handlers/tunnel"
	"github.com/atmxlab/vpn/internal/client/handlers/tunnel/mocks"
	"github.com/atmxlab/vpn/internal/protocol"
	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/atmxlab/vpn/test/gen"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestACKHandler(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		tp := gen.RandTunnelPacket()
		ip, err := tp.Payload().IP()
		require.NoError(t, err)
		subnet := net.IPNet{
			IP:   ip,
			Mask: gen.RandIPMask(),
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tunCfg := mocks.NewMockTunConfigurator(ctrl)
		tunCfg.EXPECT().ChangeTunAddr(gomock.Any(), subnet).Return(nil)

		signaller := mocks.NewMockSignaller(ctrl)
		signaller.EXPECT().Signal(gomock.Any()).Return(nil)

		h := tunnel.NewACKHandler(tunCfg, signaller, subnet.Mask)

		err = h.Handle(context.Background(), tp)
		require.NoError(t, err)
	})

	t.Run("tun configurator error", func(t *testing.T) {
		t.Parallel()

		tp := gen.RandTunnelPacket()
		ip, err := tp.Payload().IP()
		require.NoError(t, err)
		subnet := net.IPNet{
			IP:   ip,
			Mask: gen.RandIPMask(),
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tunCfg := mocks.NewMockTunConfigurator(ctrl)

		tunCfgErr := errors.New("tun config err")
		tunCfg.EXPECT().ChangeTunAddr(gomock.Any(), subnet).Return(tunCfgErr)

		signaller := mocks.NewMockSignaller(ctrl)

		h := tunnel.NewACKHandler(tunCfg, signaller, subnet.Mask)

		err = h.Handle(context.Background(), tp)
		require.ErrorIs(t, err, tunCfgErr)
	})

	t.Run("signaller error", func(t *testing.T) {
		t.Parallel()

		tp := gen.RandTunnelPacket()
		ip, err := tp.Payload().IP()
		require.NoError(t, err)
		subnet := net.IPNet{
			IP:   ip,
			Mask: gen.RandIPMask(),
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tunCfg := mocks.NewMockTunConfigurator(ctrl)
		tunCfg.EXPECT().ChangeTunAddr(gomock.Any(), subnet).Return(nil)

		signallerErr := errors.New("signal err")
		signaller := mocks.NewMockSignaller(ctrl)
		signaller.EXPECT().Signal(gomock.Any()).Return(signallerErr)

		h := tunnel.NewACKHandler(tunCfg, signaller, subnet.Mask)

		err = h.Handle(context.Background(), tp)
		require.ErrorIs(t, err, signallerErr)
	})

	t.Run("signaller already exists error", func(t *testing.T) {
		t.Parallel()

		tp := gen.RandTunnelPacket()
		ip, err := tp.Payload().IP()
		require.NoError(t, err)
		subnet := net.IPNet{
			IP:   ip,
			Mask: gen.RandIPMask(),
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tunCfg := mocks.NewMockTunConfigurator(ctrl)
		tunCfg.EXPECT().ChangeTunAddr(gomock.Any(), subnet).Return(nil)

		signallerErr := errors.AlreadyExists("test error")
		signaller := mocks.NewMockSignaller(ctrl)
		signaller.EXPECT().Signal(gomock.Any()).Return(signallerErr)

		h := tunnel.NewACKHandler(tunCfg, signaller, subnet.Mask)

		err = h.Handle(context.Background(), tp)
		require.NoError(t, err)
	})

	t.Run("ip from payload error", func(t *testing.T) {
		t.Parallel()

		tp := protocol.NewTunnelPacket(
			gen.RandHeader(),
			nil,
			gen.RandAddr(),
		)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tunCfg := mocks.NewMockTunConfigurator(ctrl)
		signaller := mocks.NewMockSignaller(ctrl)

		h := tunnel.NewACKHandler(tunCfg, signaller, gen.RandIPMask())

		err := h.Handle(context.Background(), tp)
		require.Error(t, err)
	})
}
