package udp

import (
	"fmt"
	"net"

	"github.com/atmxlab/vpn/pkg/errors"
)

type Tunnel struct {
	*net.UDPConn
}

func New(localUDPAddr, serverUDPAddr *net.UDPAddr) (*Tunnel, error) {
	conn, err := net.DialUDP("udp", localUDPAddr, serverUDPAddr)
	if err != nil {
		return nil, errors.Wrap(err, "net.DialUDP")
	}

	return &Tunnel{
		UDPConn: conn,
	}, nil
}

// WriteTo не поддерживается при UDP соединении
func (t *Tunnel) WriteTo(p []byte, _ net.Addr) (n int, err error) {
	return t.UDPConn.Write(p)
}

// ReadFrom не поддерживается при UDP соединении
func (t *Tunnel) ReadFrom(p []byte) (int, net.Addr, error) {
	n, err := t.UDPConn.Read(p)

	if err != nil {
		return n, t.UDPConn.RemoteAddr(), fmt.Errorf("unable read from tunnel: %w", err)
	}

	return n, t.UDPConn.RemoteAddr(), nil
}
