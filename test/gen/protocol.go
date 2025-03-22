package gen

import "github.com/atmxlab/vpn/internal/protocol"

func RandTunnelPacket() *protocol.TunnelPacket {
	return protocol.NewTunnelPacket(
		RandHeader(),
		RandPayload(),
		RandAddr(),
	)
}

func RandHeader() *protocol.Header {
	return protocol.NewHeader(
		RandElement(protocol.Flags()...),
	)
}

func RandPayload() protocol.Payload {
	return protocol.Payload(RandString())
}
