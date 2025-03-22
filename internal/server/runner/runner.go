package runner

import (
	"github.com/atmxlab/vpn/internal/protocol"
	"github.com/atmxlab/vpn/internal/server"
	"github.com/atmxlab/vpn/internal/tunnel"
	"github.com/atmxlab/vpn/internal/tuntap"
)

type state struct {
	// Конфигуратор сети на сервере
	routeConfigurator server.RouteConfigurator
	// Управляющий пирами
	peerManager server.PeerManager
	// Распределитель IP адресов
	ipDistributor server.IpDistributor

	net struct {
		// Туннель между сервером и клиентом
		tunnel tunnel.Tunnel
		// Виртуальный L3 сетевой интерфейс - выход в интернет
		tun tuntap.Tun

		// Канал, в который попадают пакеты из тоннеля - от клиента
		// Пакеты из этого канала отправляются в TUN интерфейс - в интернет
		tunnelPackets chan *protocol.TunnelPacket
		// Канал, в который попадают пакеты из TUN интерфейса - из интернета
		// Пакеты из этого канала отправляются в тоннель - клиенту
		tunPackets chan *protocol.TunPacket
	}
}

type Runner struct {
	cfg                  server.Config
	ipDistributorFactory server.IpDistributorFactory
	tunnelFactory        tunnel.Factory
	tunFactory           tuntap.TunFactory
}
