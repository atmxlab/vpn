package protocol_test

import (
	"bytes"
	"testing"

	"github.com/atmxlab/vpn/internal/protocol"
	"github.com/atmxlab/vpn/test/gen"
	"github.com/stretchr/testify/require"
)

func TestTunnelPacket_Marshall(t *testing.T) {
	t.Parallel()

	tunnelPacket := gen.RandTunnelPacket()

	var buf []byte

	buffer := bytes.NewBuffer(buf)

	buffer.WriteByte(tunnelPacket.Header().Flag().Byte())
	buffer.Write(tunnelPacket.Payload())

	expected := buffer.Bytes()

	actual := tunnelPacket.Marshal()

	require.Equal(t, expected, actual)
}

func TestTunnelPacket_UnMarshall(t *testing.T) {
	t.Parallel()

	tunnelPacket := gen.RandTunnelPacket()

	unmarshalledTunnelPacket := protocol.UnmarshalTunnelPacket(tunnelPacket.Addr(), tunnelPacket.Marshal())

	require.Equal(t, tunnelPacket, unmarshalledTunnelPacket)
}
