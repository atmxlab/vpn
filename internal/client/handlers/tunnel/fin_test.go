package tunnel_test

import (
	"context"
	"testing"

	"github.com/atmxlab/vpn/internal/client/handlers/tunnel"
	"github.com/atmxlab/vpn/internal/client/handlers/tunnel/mocks"
	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestFINHandler(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		stopper := mocks.NewMockStopper(ctrl)

		stopper.EXPECT().Stop(gomock.Any()).Return(nil)

		h := tunnel.NewFINHandler(stopper)

		err := h.Handle(context.Background(), nil)
		require.NoError(t, err)
	})

	t.Run("stopper error", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		stopper := mocks.NewMockStopper(ctrl)

		stopErr := errors.New("stopper error")
		stopper.EXPECT().Stop(gomock.Any()).Return(stopErr)

		h := tunnel.NewFINHandler(stopper)

		err := h.Handle(context.Background(), nil)
		require.ErrorIs(t, err, stopErr)
	})
}
