package engine

import (
	"github.com/atmxlab/vpn/internal/protocol"
	"github.com/atmxlab/vpn/test"
	"github.com/stretchr/testify/require"
)

func TUN(packet protocol.TunPacket) test.Action {
	return newSimpleAction(func(a test.App) {
		_, err := a.Tun().WriteToInput(packet.Payload())
		require.NoError(a.T(), err)
	})
}
