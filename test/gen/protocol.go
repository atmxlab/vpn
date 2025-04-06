package gen

import (
	"github.com/atmxlab/vpn/internal/protocol"
)

func RandTunnelPacket() *protocol.TunnelPacket {
	return protocol.NewTunnelPacket(
		RandHeader(),
		RandPayload(),
		RandAddr(),
	)
}

func RandTunnelPacketWithFlag(flag protocol.Flag) *protocol.TunnelPacket {
	return protocol.NewTunnelPacket(
		protocol.NewHeader(flag),
		RandPayload(),
		RandAddr(),
	)
}

func RandTunnelPSHPacket() *protocol.TunnelPacket {
	return RandTunnelPacketWithFlag(protocol.FlagPSH)
}

func RandTunnelSYNPacket() *protocol.TunnelPacket {
	return RandTunnelPacketWithFlag(protocol.FlagSYN)
}

func RandTunnelACKPacket() *protocol.TunnelPacket {
	return RandTunnelPacketWithFlag(protocol.FlagACK)
}

func RandHeader() *protocol.Header {
	return protocol.NewHeader(
		RandElement(protocol.Flags()...),
	)
}

func RandPayload() protocol.Payload {
	return protocol.Payload(RandString())
}
