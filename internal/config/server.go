package config

import (
	"net"

	"github.com/atmxlab/vpn/pkg/errors"
)

type ServerConfig struct {
	// Кол-во байт, которые будут читаться из tun интерфейса и туннеля
	BufferSize uint16 `json:"bufferSize,omitempty"`
	// Время жизни пира до следующего keepalive запроса
	PeerKeepAliveTTL Duration `json:"peerKeepAliveTTL,omitempty"`

	Tun    ServerTun    `json:"tun"`
	Tunnel ServerTunnel `json:"tunnel"`
}

type ServerTun struct {
	// 10.1.1.0/24 Подсеть TUN интерфейса
	// Из этой подсети будут выдаваться IP адреса клиентам
	SubnetCIDR string `json:"subnetCIDR,omitempty"`
	// Maximum Transition Unit  - максимальная длина неделимого пакета
	// MTU = 1500 - 20 (IP) - 8 (UDP) - 1 (VPN) = 1471
	MTU uint16 `json:"mtu,omitempty"`
	// Размер буфера канала,
	// в который будут складываться пакеты из TUN интерфейса
	TunChanSize uint `json:"tunChanSize,omitempty"`
}

func (st ServerTun) GetCIDR() (net.IP, net.IPNet, error) {
	tunIP, subnet, err := net.ParseCIDR(st.SubnetCIDR)
	if err != nil {
		return net.IP{}, net.IPNet{}, errors.Wrap(err, "parsing subnet GetCIDR")
	}

	return tunIP, *subnet, nil
}

type ServerTunnel struct {
	// Размер буфера канала,
	// в который будут складываться пакеты из тоннеля
	TunnelChanSize uint `json:"tunnelChanSize,omitempty"`
	// Протокол, в котором будет работать тоннель
	Network string `json:"network,omitempty"`
	// IP адрес, который будет слушать тоннель
	IP string `json:"ip,omitempty"`
	// Port, который будет слушать тоннель
	Port uint16 `json:"port,omitempty"`
}
