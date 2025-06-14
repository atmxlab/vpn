package main

import (
	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/sirupsen/logrus"
)

// RouteConfigurator - конфигурирует сеть на сервере под VPN сервер
type RouteConfigurator interface {
	// EnableIPForward - включает транзит IP пакетов на сервере
	// e.g. sysctl -w net.ipv4.ip_forward=1
	EnableIPForward() error
	// ConfigureFirewall - конфигурирует сетевой фильтр
	// e.g. netfilter - iptables, nftables
	// TODO: в идеале эта штука должна создать отдельную цепочку
	ConfigureFirewall() error
}

func setupOS(rc RouteConfigurator) error {
	if err := rc.EnableIPForward(); err != nil {
		return errors.Wrap(err, "routeConfigurator.EnableIPForward")
	}
	logrus.Debug("Route configurator configured IP forwarding")

	if err := rc.ConfigureFirewall(); err != nil {
		return errors.Wrap(err, "routeConfigurator.ConfigureFirewall")
	}
	logrus.Debug("Route configurator configured firewall")

	return nil
}
