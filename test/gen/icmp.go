package gen

import (
	"github.com/atmxlab/vpn/test"
	"github.com/atmxlab/vpn/test/stub"
)

func RandICMPReq(hook ...func(p *stub.ICMPPacket)) *stub.ICMPPacket {
	p := &stub.ICMPPacket{}

	test.ApplyHooks(p, hook)

	return p
}
