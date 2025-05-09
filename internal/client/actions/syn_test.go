package actions_test

import (
	"context"
	"testing"

	"github.com/atmxlab/vpn/internal/client/actions"
	"github.com/atmxlab/vpn/internal/client/actions/mocks"
	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/atmxlab/vpn/test/gen"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestSYNAction(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		serverAddr := gen.RandAddr()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tunl := mocks.NewMockTunnel(ctrl)

		tunl.EXPECT().SYN(serverAddr, nil).Return(0, nil)

		h := actions.NewSYNAction(tunl, serverAddr)

		err := h.Run(context.Background())
		require.NoError(t, err)
	})

	t.Run("tunnel error", func(t *testing.T) {
		t.Parallel()

		serverAddr := gen.RandAddr()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tunl := mocks.NewMockTunnel(ctrl)

		tunlErr := errors.New("tunnel error")
		tunl.EXPECT().SYN(serverAddr, nil).Return(0, tunlErr)

		h := actions.NewSYNAction(tunl, serverAddr)

		err := h.Run(context.Background())
		require.ErrorIs(t, err, tunlErr)
	})
}
