package main

import (
	"fmt"
	"net"

	"github.com/atmxlab/vpn/cmd"
	"github.com/atmxlab/vpn/internal/config"
	"github.com/atmxlab/vpn/internal/pkg/details/server/conn/udp"
	"github.com/atmxlab/vpn/internal/pkg/tunnel"
	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/sirupsen/logrus"
)

func setupTunnelConn(cfg config.ServerConfig) tunnel.Connection {
	switch cfg.Tunnel.Network {
	case "udp":
		serverAddr := fmt.Sprintf("%s:%d", cfg.Tunnel.IP, cfg.Tunnel.Port)

		logrus.
			WithField("ServerAddr", serverAddr).
			Info("[TUNNEL] Setup tunnel UDP connection")

		udpAddr, err := net.ResolveUDPAddr(
			cfg.Tunnel.Network,
			serverAddr,
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
