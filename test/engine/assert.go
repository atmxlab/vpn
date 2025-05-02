package engine

import (
	"net"

	"github.com/atmxlab/vpn/internal/protocol"
	"github.com/atmxlab/vpn/test"
	"github.com/stretchr/testify/require"
)

func CHECKPOINT(actions ...test.Action) test.Action {
	return newSimpleAction(func(a test.App) {
		a.T().Helper()

		for _, action := range actions {
			action.Handle(a)
		}
	})
}

func ExpectPeer(addr net.Addr) test.Action {
	return newSimpleAction(func(a test.App) {
		_, exists, err := a.PeerManager().GetByAddr(a.Ctx(), addr)
		require.NoError(a.T(), err)
		require.Truef(a.T(), exists, "peer with addr=[%v] must be exists", addr.String())
	})
}

func UnexpectPeer(addr net.Addr) test.Action {
	return newSimpleAction(func(a test.App) {
		_, exists, err := a.PeerManager().GetByAddr(a.Ctx(), addr)
		require.NoError(a.T(), err)
		require.Falsef(a.T(), exists, "peer with addr=[%v] cannot be exists", addr.String())
	})
}

func ExpectFreeAllDedicatedIPs() test.Action {
	return newSimpleAction(func(a test.App) {
		require.False(a.T(), a.IPDistributor().HasBusy(), "cannot exists busy ips")
	})
}

func ExpectBusyDedicatedIP() test.Action {
	return newSimpleAction(func(a test.App) {
		require.True(a.T(), a.IPDistributor().HasBusy(), "busy ip must be exists")
	})
}

func ExpectTun(packet *protocol.TunPacket) test.Action {
	return newSimpleAction(func(a test.App) {
		lastPacket, ok := a.Tun().GetLastPacket()
		require.Truef(a.T(), ok, "tun cannot be empty")
		require.Equal(a.T(), packet, lastPacket)
	})
}

func ExpectEmptyTun() test.Action {
	return newSimpleAction(func(a test.App) {
		_, ok := a.Tun().GetLastPacket()
		require.Falsef(a.T(), ok, "tun must be empty")
	})
}

func ExpectTunnelACK(dst net.Addr) test.Action {
	return newSimpleAction(func(a test.App) {
		lastPacket, ok := a.Tunnel().GetLastPacket()

		require.Truef(a.T(), ok, "tunnel cannot be empty")
		require.Equalf(a.T(), protocol.FlagACK, lastPacket.Header().Flag(), "invalid last tunnel packet: expected [%s], actual [%s]", protocol.FlagACK, lastPacket.Header().Flag())
		require.Equal(a.T(), dst, lastPacket.Addr(), "invalid destination addr")
	})
}

func ExpectTunnelPSH(dst net.Addr, payload protocol.Payload) test.Action {
	return newSimpleAction(func(a test.App) {
		a.T().Helper()

		lastPacket, ok := a.Tunnel().GetLastPacket()

		require.Truef(a.T(), ok, "tunnel cannot be empty")
		require.Equalf(a.T(), protocol.FlagPSH, lastPacket.Header().Flag(), "invalid last tunnel packet: expected [%s], actual [%s]", protocol.FlagPSH, lastPacket.Header().Flag())
		require.Equal(a.T(), dst, lastPacket.Addr(), "invalid destination addr")
		require.Equal(a.T(), payload, lastPacket.Payload(), "invalid payload")
	})
}
