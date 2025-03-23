package peermanager_test

import (
	"context"
	"testing"

	"github.com/atmxlab/vpn/internal/pkg/peermanager"
	"github.com/atmxlab/vpn/test/gen"
	"github.com/stretchr/testify/require"
)

func TestManager(t *testing.T) {
	t.Parallel()

	t.Run("get by addr", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		peer := gen.RandPeer()

		pm := peermanager.New()

		require.NoError(t, pm.Add(ctx, peer, gen.RandDuration()))

		actualPeer, exists, err := pm.GetByAddr(ctx, peer.Addr())
		require.NoError(t, err)
		require.True(t, exists)

		require.Equal(t, peer, actualPeer)
	})

	t.Run("get by dedicated ip", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		peer := gen.RandPeer()

		pm := peermanager.New()

		require.NoError(t, pm.Add(ctx, peer, gen.RandDuration()))

		actualPeer, exists, err := pm.GetByDedicatedIP(ctx, peer.DedicatedIP())
		require.NoError(t, err)
		require.True(t, exists)

		require.Equal(t, peer, actualPeer)
	})

	t.Run("double add", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		peer := gen.RandPeer()

		pm := peermanager.New()

		require.NoError(t, pm.Add(ctx, peer, gen.RandDuration()))
		require.NoError(t, pm.Add(ctx, peer, gen.RandDuration()))
	})

	t.Run("get by addr with many peers", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		peer := gen.RandPeer()

		pm := peermanager.New()

		require.NoError(t, pm.Add(ctx, gen.RandPeer(), gen.RandDuration()))
		require.NoError(t, pm.Add(ctx, gen.RandPeer(), gen.RandDuration()))
		require.NoError(t, pm.Add(ctx, peer, gen.RandDuration()))
		require.NoError(t, pm.Add(ctx, gen.RandPeer(), gen.RandDuration()))
		require.NoError(t, pm.Add(ctx, gen.RandPeer(), gen.RandDuration()))

		actualPeer, exists, err := pm.GetByAddr(ctx, peer.Addr())
		require.NoError(t, err)
		require.True(t, exists)
		require.Equal(t, peer, actualPeer)
	})

	t.Run("get by dedicated ip with many peers", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		peer := gen.RandPeer()

		pm := peermanager.New()

		require.NoError(t, pm.Add(ctx, gen.RandPeer(), gen.RandDuration()))
		require.NoError(t, pm.Add(ctx, gen.RandPeer(), gen.RandDuration()))
		require.NoError(t, pm.Add(ctx, peer, gen.RandDuration()))
		require.NoError(t, pm.Add(ctx, gen.RandPeer(), gen.RandDuration()))
		require.NoError(t, pm.Add(ctx, gen.RandPeer(), gen.RandDuration()))

		actualPeer, exists, err := pm.GetByDedicatedIP(ctx, peer.DedicatedIP())
		require.NoError(t, err)
		require.True(t, exists)
		require.Equal(t, peer, actualPeer)
	})

	t.Run("get by addr and extend", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		peer := gen.RandPeer()

		pm := peermanager.New()

		require.NoError(t, pm.Add(ctx, gen.RandPeer(), gen.RandDuration()))
		require.NoError(t, pm.Add(ctx, gen.RandPeer(), gen.RandDuration()))
		require.NoError(t, pm.Add(ctx, peer, gen.RandDuration()))
		require.NoError(t, pm.Add(ctx, gen.RandPeer(), gen.RandDuration()))
		require.NoError(t, pm.Add(ctx, gen.RandPeer(), gen.RandDuration()))

		actualPeer, exists, err := pm.GetByAddrAndExtend(ctx, peer.Addr(), gen.RandDuration())
		require.NoError(t, err)
		require.True(t, exists)
		require.Equal(t, peer, actualPeer)
	})

	t.Run("has without peer needs peer", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		pm := peermanager.New()

		require.NoError(t, pm.Add(ctx, gen.RandPeer(), gen.RandDuration()))
		require.NoError(t, pm.Add(ctx, gen.RandPeer(), gen.RandDuration()))
		require.NoError(t, pm.Add(ctx, gen.RandPeer(), gen.RandDuration()))
		require.NoError(t, pm.Add(ctx, gen.RandPeer(), gen.RandDuration()))

		has, err := pm.HasPeer(ctx, gen.RandAddr())
		require.NoError(t, err)
		require.False(t, has)
	})

	t.Run("has with peer needs peer", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		peer := gen.RandPeer()

		pm := peermanager.New()

		require.NoError(t, pm.Add(ctx, gen.RandPeer(), gen.RandDuration()))
		require.NoError(t, pm.Add(ctx, gen.RandPeer(), gen.RandDuration()))
		require.NoError(t, pm.Add(ctx, gen.RandPeer(), gen.RandDuration()))
		require.NoError(t, pm.Add(ctx, peer, gen.RandDuration()))
		require.NoError(t, pm.Add(ctx, gen.RandPeer(), gen.RandDuration()))

		has, err := pm.HasPeer(ctx, peer.Addr())
		require.NoError(t, err)
		require.True(t, has)
	})
}
