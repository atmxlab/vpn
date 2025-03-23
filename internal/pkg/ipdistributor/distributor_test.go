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

		allocatedIP, err := d.AllocateIP()
		require.NoError(t, err)
		require.Equal(t, subnet.IP, allocatedIP)

		_, err = d.AllocateIP()
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

		allocatedIP, err := d.AllocateIP()
		require.NoError(t, err)
		require.Equal(t, net.IPv4(1, 1, 1, 0).To4(), allocatedIP.To4())

		allocatedIP, err = d.AllocateIP()
		require.NoError(t, err)
		require.Equal(t, net.IPv4(1, 1, 1, 1).To4(), allocatedIP.To4())

		allocatedIP, err = d.AllocateIP()
		require.NoError(t, err)
		require.Equal(t, net.IPv4(1, 1, 1, 2).To4(), allocatedIP.To4())
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
		var allocatedIP net.IP

		for {
			allocatedIP, err = d.AllocateIP()
			if err != nil {
				break
			} else {
				lastIP = allocatedIP
			}
		}

		require.Equal(t, net.IPv4(1, 1, 1, 255).To4(), lastIP.To4())
	})
}
