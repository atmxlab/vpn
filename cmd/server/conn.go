package main

import (
	"fmt"
	"net"

	"github.com/atmxlab/vpn/internal/config"
	"github.com/atmxlab/vpn/internal/pkg/details/conn/server/udp"
	"github.com/atmxlab/vpn/internal/pkg/tunnel"
	"github.com/atmxlab/vpn/pkg/errors"
)

func setupTunnelConn(cfg config.ServerConfig) tunnel.Connection {
	switch cfg.Tunnel.Network {
	case "udp":
		udpAddr, err := net.ResolveUDPAddr(
			cfg.Tunnel.Network,
			fmt.Sprintf("%s:%d", cfg.Tunnel.IP, cfg.Tunnel.Port),
		)
		exitf(err, "net.ResolveUDPAddr")
		conn, err := udp.New(udpAddr)
		exitf(err, "udp.New")

		return conn
	default:
		exitf(errors.New("invalid tunnel network"), "unexpected network: %s", cfg.Tunnel.Network)
		return nil
	}
}
