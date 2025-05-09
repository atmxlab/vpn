package actions

import "net"

//go:generate mock Tunnel
type Tunnel interface {
	SYN(addr net.Addr, payload []byte) (int, error)
	KPA(addr net.Addr, payload []byte) (int, error)
}
