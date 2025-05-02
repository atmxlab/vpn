package main

import (
	"net"

	"github.com/atmxlab/vpn/internal/config"
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
	ConfigureFirewall(subnet net.IPNet) error
	// SetDefaultRoute - указывает шлюз по умолчанию для подсети
	SetDefaultRoute(subnet net.IPNet) error
}

// TODO: эта часть не сделана совсем

func setupOS(rc RouteConfigurator, cfg config.ServerConfig) error {
	tunSubnet, err := cfg.Tun.Subnet()
	if err != nil {
		return errors.Wrap(err, "cfg.Tun.Subnet")
	}

	if err := rc.EnableIPForward(); err != nil {
		return errors.Wrap(err, "routeConfigurator.EnableIPForward")
	}
	logrus.Debug("Route configurator configured IP forwarding")

	if err := rc.ConfigureFirewall(tunSubnet); err != nil {
		return errors.Wrap(err, "routeConfigurator.ConfigureFirewall")
	}
	logrus.Debug("Route configurator configured firewall")

	if err := rc.SetDefaultRoute(tunSubnet); err != nil {
		return errors.Wrap(err, "routeConfigurator.SetDefaultRoute")
	}
	logrus.Debug("Route configurator configured default route")

	return nil
}
