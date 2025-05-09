package gen

import (
	"fmt"
	"net"

	"github.com/atmxlab/vpn/test/stub"
)

func RandAddr() net.Addr {
	return stub.NewAddr(
		RandNetwork(),
		RandAddress(),
	)
}

func RandNetwork() string {
	return RandElement("udp", "tcp", "ip")
}

func RandAddress() string {
	return fmt.Sprintf("%s:%d", RandIP().String(), RandPort())
}

func RandPort() uint16 {
	return RandUInt16()
}

func RandIP() net.IP {
	return net.IPv4(
		RandByte(),
		RandByte(),
		RandByte(),
		RandByte(),
	).
		To4()
}

func RandIPMask() net.IPMask {
	return net.IPv4Mask(
		RandByte(),
		RandByte(),
		RandByte(),
		RandByte(),
	)
}
