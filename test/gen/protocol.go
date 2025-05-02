package gen

import (
	"encoding/binary"
	"testing"

	"github.com/atmxlab/vpn/internal/protocol"
	"github.com/atmxlab/vpn/test"
	"github.com/atmxlab/vpn/test/stub"
	"github.com/stretchr/testify/require"
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

func RandTunPacket() *protocol.TunPacket {
	return protocol.NewTunPacket(RandPayload())
}

func RandTunICMPReq(t *testing.T, hook ...func(h *stub.IPHeader)) *protocol.TunPacket {
	icmpBytes := RandICMPPacket(t)

	// 3. Создаем IP-заголовок
	ipHeader := &stub.IPHeader{
		Version:  4,
		Len:      20, // Без опций
		TOS:      0,
		TotalLen: 20 + len(icmpBytes), // IP заголовок + ICMP
		ID:       54321,
		Flags:    0,
		FragOff:  0,
		TTL:      64,
		Protocol: 1, // ICMP
		Dst:      RandIP(),
		Src:      RandIP(),
	}

	test.ApplyHooks(ipHeader, hook)

	ipBytes, err := ipHeader.Marshal()
	require.NoError(t, err)

	ipChecksum := test.IPHeaderChecksum(ipBytes)
	binary.BigEndian.PutUint16(ipBytes[10:12], ipChecksum)

	// 6. Объединяем IP + ICMP
	fullPacket := append(ipBytes, icmpBytes...)

	return protocol.NewTunPacket(fullPacket)
}
