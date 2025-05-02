package engine

import (
	"net"

	"github.com/atmxlab/vpn/internal/protocol"
	"github.com/atmxlab/vpn/test"
	"github.com/stretchr/testify/require"
)

func FIN(clientAddr net.Addr) test.Action {
	return newSimpleAction(func(a test.App) {
		tunnelPacket := protocol.NewTunnelPacket(
			protocol.NewHeader(protocol.FlagFIN),
			nil,
			clientAddr,
		)

		_, err := a.Tunnel().WriteToInput(tunnelPacket.Marshal(), tunnelPacket.Addr())
		require.NoError(a.T(), err)
	})
}

func SYN(clientAddr net.Addr) test.Action {
	return newSimpleAction(func(a test.App) {
		tunnelPacket := protocol.NewTunnelPacket(
			protocol.NewHeader(protocol.FlagSYN),
			nil,
			clientAddr,
		)

		_, err := a.Tunnel().WriteToInput(tunnelPacket.Marshal(), tunnelPacket.Addr())
		require.NoError(a.T(), err)
	})
}

func PSH(clientAddr net.Addr, payload protocol.Payload) test.Action {
	return newSimpleAction(func(a test.App) {
		tunnelPacket := protocol.NewTunnelPacket(
			protocol.NewHeader(protocol.FlagPSH),
			payload,
			clientAddr,
		)

		_, err := a.Tunnel().WriteToInput(tunnelPacket.Marshal(), tunnelPacket.Addr())
		require.NoError(a.T(), err)
	})
}
