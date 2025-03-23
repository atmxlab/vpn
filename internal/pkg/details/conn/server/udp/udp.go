package udp

import (
	"net"

	"github.com/atmxlab/vpn/pkg/errors"
)

type Tunnel struct {
	*net.UDPConn
}

func New(addr net.Addr) (*Tunnel, error) {
	udpAddr, err := net.ResolveUDPAddr(addr.Network(), addr.String())
	if err != nil {
		return nil, errors.Wrap(err, "net.ResolveUDPAddr")
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, errors.Wrap(err, "net.ListenUDP")
	}

	return &Tunnel{
		UDPConn: conn,
	}, nil
}
