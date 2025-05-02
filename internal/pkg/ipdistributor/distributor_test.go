package ipdistributor_test

import (
	"net"
	"testing"

	"github.com/atmxlab/vpn/internal/pkg/ipdistributor"
	"github.com/stretchr/testify/require"
)

func TestDistributor(t *testing.T) {
	t.Parallel()

	t.Run("32 bits mask", func(t *testing.T) {
		t.Parallel()

		subnet := net.IPNet{
			IP:   net.IPv4(1, 1, 1, 2),
			Mask: net.IPv4Mask(255, 255, 255, 255),
		}

		d, err := ipdistributor.New(subnet)
		require.NoError(t, err)

		_, err = d.AcquireIP()
		require.Error(t, err)
	})

	t.Run("24 bits mask", func(t *testing.T) {
		t.Parallel()

		subnet := net.IPNet{
			IP:   net.IPv4(1, 1, 1, 2),
			Mask: net.IPv4Mask(255, 255, 255, 0),
		}

		d, err := ipdistributor.New(subnet)
		require.NoError(t, err)

		acquiredIP, err := d.AcquireIP()
		require.NoError(t, err)
		require.Equal(t, net.IPv4(1, 1, 1, 2).To4(), acquiredIP.To4())

		acquiredIP, err = d.AcquireIP()
		require.NoError(t, err)
		require.Equal(t, net.IPv4(1, 1, 1, 3).To4(), acquiredIP.To4())

		acquiredIP, err = d.AcquireIP()
		require.NoError(t, err)
		require.Equal(t, net.IPv4(1, 1, 1, 4).To4(), acquiredIP.To4())
	})

	t.Run("24 bits mask - check last ip", func(t *testing.T) {
		t.Parallel()

		subnet := net.IPNet{
			IP:   net.IPv4(1, 1, 1, 2),
			Mask: net.IPv4Mask(255, 255, 255, 0),
		}

		d, err := ipdistributor.New(subnet)
		require.NoError(t, err)

		var lastIP net.IP
		var acquiredIP net.IP

		for {
			acquiredIP, err = d.AcquireIP()
			if err != nil {
				break
			} else {
				lastIP = acquiredIP
			}
		}

		require.Equal(t, net.IPv4(1, 1, 1, 254).To4(), lastIP.To4())
	})
}
