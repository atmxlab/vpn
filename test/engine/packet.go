package engine

import (
	"github.com/atmxlab/vpn/internal/protocol"
	"github.com/atmxlab/vpn/test/gen"
	"github.com/atmxlab/vpn/test/stub"
)

func ICMPReq(hook ...func(p *stub.ICMPPacket)) *protocol.TunPacket {
	p := gen.RandICMPReq(hook...)

	bytes, err := p.Marshal()
	if err != nil {
		panic(err)
	}

	return protocol.NewTunPacket(bytes)
}
