package tun_test

import (
	"context"
	"testing"
	"time"

	"github.com/atmxlab/vpn/internal/pkg/tun"
	"github.com/atmxlab/vpn/internal/pkg/tun/mocks"
	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestTun_Write(t *testing.T) {
	t.Parallel()

	t.Run("without error", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		embeddedTun := mocks.NewMockEmbeddedTun(ctrl)

		payload := make([]byte, 1024)

		embeddedTun.EXPECT().Write(payload).Return(len(payload), nil)

		tn := tun.NewTun(embeddedTun)

		n, err := tn.Write(payload)
		require.NoError(t, err)
		require.Equal(t, len(payload), n)
	})

	t.Run("write error", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		embeddedTun := mocks.NewMockEmbeddedTun(ctrl)

		payload := make([]byte, 1024)

		writeError := errors.New("write error")
		embeddedTun.EXPECT().Write(payload).Return(0, writeError)

		tn := tun.NewTun(embeddedTun)

		n, err := tn.Write(payload)
		require.ErrorIs(t, err, writeError)
		require.Zero(t, n)
	})
}

func TestTun_Read(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		ctrl := gomock.NewController(t)
		embeddedTun := mocks.NewMockEmbeddedTun(ctrl)

		buffer := make([]byte, 1024)

		embeddedTun.EXPECT().Read(buffer).Return(len(buffer), nil)

		tn := tun.NewTun(embeddedTun)

		n, err := tn.ReadWithContext(ctx, buffer)
		require.NoError(t, err)
		require.Equal(t, len(buffer), n)
	})

	t.Run("cancel ctx", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		ctrl := gomock.NewController(t)
		embeddedTun := mocks.NewMockEmbeddedTun(ctrl)

		buffer := make([]byte, 1024)

		embeddedTun.EXPECT().Read(buffer).DoAndReturn(func(buf []byte) (n int, err error) {
			// Имитируем блокировку чтения
			time.Sleep(50 * time.Millisecond)
			return len(buf), nil
		})

		embeddedTun.EXPECT().Close().Return(nil)

		tn := tun.NewTun(embeddedTun)

		go func() {
			// Отменяем контекст до того, как завершит чтение
			time.Sleep(10 * time.Millisecond)
			cancel()
		}()

		n, err := tn.ReadWithContext(ctx, buffer)
		require.ErrorIs(t, err, context.Canceled)
		require.Zero(t, n)
	})

	t.Run("cancel ctx and close reader error", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		ctrl := gomock.NewController(t)
		embeddedTun := mocks.NewMockEmbeddedTun(ctrl)

		buffer := make([]byte, 1024)

		embeddedTun.EXPECT().Read(buffer).DoAndReturn(func(buf []byte) (n int, err error) {
			// Имитируем блокировку чтения
			time.Sleep(50 * time.Millisecond)
			return len(buf), nil
		})

		closeError := errors.New("close error")
		embeddedTun.EXPECT().Close().Return(closeError)

		tn := tun.NewTun(embeddedTun)

		go func() {
			// Отменяем контекст до того, как завершит чтение
			time.Sleep(10 * time.Millisecond)
			cancel()
		}()

		n, err := tn.ReadWithContext(ctx, buffer)
		require.ErrorIs(t, err, closeError)
		require.Zero(t, n)
	})
}

func TestTun_Close(t *testing.T) {
	t.Parallel()

	t.Run("without error", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		embeddedTun := mocks.NewMockEmbeddedTun(ctrl)
		embeddedTun.EXPECT().Close().Return(nil)

		tn := tun.NewTun(embeddedTun)

		require.NoError(t, tn.Close())
	})

	t.Run("without error", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		embeddedTun := mocks.NewMockEmbeddedTun(ctrl)
		embeddedTun.EXPECT().Close().Return(errors.New("test error"))

		tn := tun.NewTun(embeddedTun)

		require.Error(t, tn.Close())
	})
}
