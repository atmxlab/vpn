package main

import (
	"fmt"
	"net"

	"github.com/atmxlab/vpn/cmd"
	"github.com/atmxlab/vpn/internal/config"
	"github.com/atmxlab/vpn/internal/pkg/details/server/conn/udp"
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
		cmd.Exitf(err, "net.ResolveUDPAddr")
		conn, err := udp.New(udpAddr)
		cmd.Exitf(err, "udp.New")

		return conn
	default:
		cmd.Exitf(errors.New("invalid tunnel network"), "unexpected network: %s", cfg.Tunnel.Network)
		return nil
	}
}
