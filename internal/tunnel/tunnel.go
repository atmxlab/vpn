package tunnel

import (
	"io"
	"net"
)

type Factory interface {
	Create(addr net.Addr) (Tunnel, error)
}

// Tunnel - тоннель через который проходит трафик
// от VPN-сервера к VPN-клиент и обратно
type Tunnel interface {
	ReadFrom(p []byte) (n int, addr net.Addr, err error)
	WriteTo(p []byte, addr net.Addr) (n int, err error)
	LocalAddr() net.Addr
	io.Closer
}
