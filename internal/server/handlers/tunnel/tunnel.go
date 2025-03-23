package tunnel

import (
	"net"
)

type Tunnel interface {
	PSH(addr net.Addr, payload []byte) (int, error)
	ACK(addr net.Addr, payload []byte) (int, error)
}
