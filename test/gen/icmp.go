package gen

import (
	"testing"

	"github.com/atmxlab/vpn/test"
	"github.com/atmxlab/vpn/test/stub"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func RandICMPPacket(t *testing.T, hook ...func(p *stub.ICMPPacket)) []byte {
	msg := &stub.ICMPPacket{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   1,
			Seq:  1,
			Data: []byte("hello icmp request"),
		},
	}

	test.ApplyHooks(msg, hook)

	packet, err := msg.Marshal(nil)
	require.NoError(t, err)

	return packet
}
