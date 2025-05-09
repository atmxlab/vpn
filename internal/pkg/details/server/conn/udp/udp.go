package udp

import (
	"net"

	"github.com/atmxlab/vpn/pkg/errors"
)

type Tunnel struct {
	*net.UDPConn
}

func New(udpAddr *net.UDPAddr) (*Tunnel, error) {
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, errors.Wrap(err, "net.ListenUDP")
	}

	return &Tunnel{
		UDPConn: conn,
	}, nil
}
