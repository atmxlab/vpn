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

		closer := mocks.NewMockCloser(ctrl)

		closer.EXPECT().Close().Return(nil)

		h := tunnel.NewFINHandler(closer)

		err := h.Handle(context.Background(), nil)
		require.NoError(t, err)
	})

	t.Run("closer error", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		closer := mocks.NewMockCloser(ctrl)

		closeErr := errors.New("closer error")
		closer.EXPECT().Close().Return(closeErr)

		h := tunnel.NewFINHandler(closer)

		err := h.Handle(context.Background(), nil)
		require.ErrorIs(t, err, closeErr)
	})
}
