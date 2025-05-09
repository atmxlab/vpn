package tun

import "net"

//go:generate mock Tunnel
type Tunnel interface {
	PSH(addr net.Addr, payload []byte) (int, error)
}
