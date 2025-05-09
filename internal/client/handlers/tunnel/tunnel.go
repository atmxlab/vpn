package tunnel

import (
	"net"
)

//go:generate mock Tunnel
type Tunnel interface {
	SYN(addr net.Addr, payload []byte) (int, error)
	ACK(addr net.Addr, payload []byte) (int, error)
}
