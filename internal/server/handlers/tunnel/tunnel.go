package tunnel

import (
	"net"
)

type Tunnel interface {
	SYN(addr net.Addr, payload []byte) (int, error)
	ACK(addr net.Addr, payload []byte) (int, error)
}
