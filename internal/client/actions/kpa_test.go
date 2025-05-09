package actions_test

import (
	"context"
	"testing"
	"time"

	"github.com/atmxlab/vpn/internal/client/actions"
	"github.com/atmxlab/vpn/internal/client/actions/mocks"
	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/atmxlab/vpn/test/gen"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestKPAAction(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		serverAddr := gen.RandAddr()
		tickDuration := 100 * time.Millisecond

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tunl := mocks.NewMockTunnel(ctrl)

		tunl.EXPECT().KPA(serverAddr, nil).Return(0, nil)

		h := actions.NewKPAAction(tunl, serverAddr, tickDuration)

		ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
		defer cancel()

		err := h.Run(ctx)
		require.ErrorIs(t, err, context.DeadlineExceeded)
	})

	t.Run("double tick", func(t *testing.T) {
		t.Parallel()

		serverAddr := gen.RandAddr()
		tickDuration := 100 * time.Millisecond

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tunl := mocks.NewMockTunnel(ctrl)

		tunl.EXPECT().KPA(serverAddr, nil).Return(0, nil).Times(2)

		h := actions.NewKPAAction(tunl, serverAddr, tickDuration)

		ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
		defer cancel()

		err := h.Run(ctx)
		require.ErrorIs(t, err, context.DeadlineExceeded)
	})

	t.Run("without tick", func(t *testing.T) {
		t.Parallel()

		serverAddr := gen.RandAddr()
		tickDuration := 100 * time.Millisecond

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tunl := mocks.NewMockTunnel(ctrl)

		h := actions.NewKPAAction(tunl, serverAddr, tickDuration)

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		err := h.Run(ctx)
		require.ErrorIs(t, err, context.DeadlineExceeded)
	})

	t.Run("tunnel error", func(t *testing.T) {
		t.Parallel()

		serverAddr := gen.RandAddr()
		tickDuration := 100 * time.Millisecond

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tunl := mocks.NewMockTunnel(ctrl)

		tunlErr := errors.New("tunl error")
		tunl.EXPECT().KPA(serverAddr, nil).Return(0, tunlErr)

		h := actions.NewKPAAction(tunl, serverAddr, tickDuration)

		ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
		defer cancel()

		err := h.Run(ctx)
		require.ErrorIs(t, err, tunlErr)
	})
}
