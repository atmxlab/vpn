package tunnel_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/atmxlab/vpn/internal/pkg/tunnel"
	"github.com/atmxlab/vpn/internal/pkg/tunnel/mocks"
	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/atmxlab/vpn/test/gen"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestTunnel_Write(t *testing.T) {
	t.Parallel()

	t.Run("without error", func(t *testing.T) {
		t.Parallel()

		tp := gen.RandTunnelPacket()
		ctrl := gomock.NewController(t)
		conn := mocks.NewMockConnection(ctrl)

		conn.EXPECT().WriteTo(tp.Marshal(), tp.Addr()).Return(len(tp.Marshal()), nil)

		tn := tunnel.New(conn)

		n, err := tn.Write(tp)
		require.NoError(t, err)
		require.Equal(t, len(tp.Marshal()), n)
	})

	t.Run("write error", func(t *testing.T) {
		t.Parallel()

		tp := gen.RandTunnelPacket()
		ctrl := gomock.NewController(t)
		conn := mocks.NewMockConnection(ctrl)

		writeError := errors.New("write error")
		conn.EXPECT().WriteTo(tp.Marshal(), tp.Addr()).Return(len(tp.Marshal()), writeError)

		tn := tunnel.New(conn)

		n, err := tn.Write(tp)
		require.ErrorIs(t, err, writeError)
		require.Zero(t, n)
	})

	t.Run("PSH", func(t *testing.T) {
		t.Parallel()

		tp := gen.RandTunnelPSHPacket()
		ctrl := gomock.NewController(t)
		conn := mocks.NewMockConnection(ctrl)

		conn.EXPECT().WriteTo(tp.Marshal(), tp.Addr()).Return(len(tp.Marshal()), nil)

		tn := tunnel.New(conn)

		n, err := tn.PSH(tp.Addr(), tp.Payload())
		require.NoError(t, err)
		require.Equal(t, len(tp.Marshal()), n)
	})

	t.Run("SYN", func(t *testing.T) {
		t.Parallel()

		tp := gen.RandTunnelSYNPacket()
		ctrl := gomock.NewController(t)
		conn := mocks.NewMockConnection(ctrl)

		conn.EXPECT().WriteTo(tp.Marshal(), tp.Addr()).Return(len(tp.Marshal()), nil)

		tn := tunnel.New(conn)

		n, err := tn.SYN(tp.Addr(), tp.Payload())
		require.NoError(t, err)
		require.Equal(t, len(tp.Marshal()), n)
	})

	t.Run("ACK", func(t *testing.T) {
		t.Parallel()

		tp := gen.RandTunnelACKPacket()
		ctrl := gomock.NewController(t)
		conn := mocks.NewMockConnection(ctrl)

		conn.EXPECT().WriteTo(tp.Marshal(), tp.Addr()).Return(len(tp.Marshal()), nil)

		tn := tunnel.New(conn)

		n, err := tn.ACK(tp.Addr(), tp.Payload())
		require.NoError(t, err)
		require.Equal(t, len(tp.Marshal()), n)
	})
}

func TestTunnel_Read(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		addr := gen.RandAddr()

		ctrl := gomock.NewController(t)
		conn := mocks.NewMockConnection(ctrl)

		buffer := make([]byte, 1024)

		conn.EXPECT().ReadFrom(buffer).Return(len(buffer), addr, nil)

		tn := tunnel.New(conn)

		n, a, err := tn.ReadFromWithContext(ctx, buffer)
		require.NoError(t, err)
		require.Equal(t, len(buffer), n)
		require.Equal(t, addr, a)
	})

	t.Run("cancel ctx", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		addr := gen.RandAddr()

		ctrl := gomock.NewController(t)
		conn := mocks.NewMockConnection(ctrl)

		buffer := make([]byte, 1024)

		conn.EXPECT().ReadFrom(buffer).DoAndReturn(func(buf []byte) (int, net.Addr, error) {
			// Имитируем блокировку чтения
			time.Sleep(50 * time.Millisecond)
			return len(buf), addr, nil
		})

		conn.EXPECT().Close().Return(nil)

		go func() {
			// Отменяем контекст до того, как завершит чтение
			time.Sleep(10 * time.Millisecond)
			cancel()
		}()

		tn := tunnel.New(conn)

		n, a, err := tn.ReadFromWithContext(ctx, buffer)
		require.ErrorIs(t, err, context.Canceled)
		require.Zero(t, n)
		require.Zero(t, a)
	})

	t.Run("cancel ctx close error", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		addr := gen.RandAddr()

		ctrl := gomock.NewController(t)
		conn := mocks.NewMockConnection(ctrl)

		buffer := make([]byte, 1024)

		conn.EXPECT().ReadFrom(buffer).DoAndReturn(func(buf []byte) (int, net.Addr, error) {
			// Имитируем блокировку чтения
			time.Sleep(50 * time.Millisecond)
			return len(buf), addr, nil
		})

		closeError := errors.New("close error")
		conn.EXPECT().Close().Return(closeError)

		go func() {
			// Отменяем контекст до того, как завершит чтение
			time.Sleep(10 * time.Millisecond)
			cancel()
		}()

		tn := tunnel.New(conn)

		n, a, err := tn.ReadFromWithContext(ctx, buffer)
		require.ErrorIs(t, err, closeError)
		require.Zero(t, n)
		require.Zero(t, a)
	})
}

func TestTunnel_Close(t *testing.T) {
	t.Parallel()

	t.Run("without error", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		conn := mocks.NewMockConnection(ctrl)
		conn.EXPECT().Close().Return(nil)

		tn := tunnel.New(conn)

		require.NoError(t, tn.Close())
	})

	t.Run("without error", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		conn := mocks.NewMockConnection(ctrl)

		closeError := errors.New("close error")
		conn.EXPECT().Close().Return(closeError)

		tn := tunnel.New(conn)

		require.ErrorIs(t, tn.Close(), closeError)
	})
}
