package tun

import "net"

type Tunnel interface {
	PSH(addr net.Addr, payload []byte) (int, error)
}
