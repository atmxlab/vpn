package tuntap

import (
	"io"
	"net"
)

type TunFactory interface {
	Create(subnet net.IPNet, mtu uint16) (Tun, error)
}

type Tun interface {
	io.ReadWriteCloser
	Name() string
}
