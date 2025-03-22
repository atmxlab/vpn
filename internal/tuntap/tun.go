package tuntap

import (
	"io"
	"net"
)

type TunFactory interface {
	Create(subnet net.IPNet, mtu uint16) (Tun, error)
}

// Tun - виртуальный сетевой интерфейс
// В Linux можно создавать два интерфейса - TUN и TAP
// Нам необходим только TUN, так как работаем на L3 уровне, а не на L2
type Tun interface {
	io.ReadWriteCloser
	Name() string
}
