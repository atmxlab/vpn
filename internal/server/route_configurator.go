package server

import "net"

// RouteConfigurator - конфигурирует сеть на сервере под VPN сервер
type RouteConfigurator interface {
	// EnableIPForward - включает транзит IP пакетов на сервере
	// e.g. sysctl -w net.ipv4.ip_forward=1
	EnableIPForward() error
	// ConfigureFirewall - конфигурирует сетевой фильтр
	// e.g. netfilter - iptables, nftables
	// TODO: в идеале эта штука должна создать отдельную цепочку
	ConfigureFirewall(subnet net.IPNet) error
	// SetDefaultRoute - указывает шлюз по умолчанию для подсети
	SetDefaultRoute(subnet net.IPNet) error
}
