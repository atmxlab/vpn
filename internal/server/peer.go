package server

import (
	"net"
)

// Peer - информация о подключенном клиенте
type Peer struct {
	// Адрес пира - это уникальное поле, по которому можно находить пир
	addr net.Addr
	// Выделенный для пира IP адрес
	dedicatedIP net.IP
}

func NewPeer(dedicatedIP net.IP, addr net.Addr) *Peer {
	return &Peer{
		addr:        addr,
		dedicatedIP: dedicatedIP,
	}
}

func (p *Peer) DedicatedIP() net.IP {
	return p.dedicatedIP
}

func (p *Peer) Addr() net.Addr {
	return p.addr
}
