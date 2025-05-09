package main

import (
	"fmt"
	"net"

	"github.com/atmxlab/vpn/cmd"
	"github.com/atmxlab/vpn/internal/config"
	"github.com/atmxlab/vpn/internal/pkg/details/client/conn/udp"
	"github.com/atmxlab/vpn/internal/pkg/tunnel"
	"github.com/atmxlab/vpn/pkg/errors"
)

func setupTunnelConn(cfg config.ClientConfig) tunnel.Connection {
	switch cfg.Tunnel.Network {
	case "udp":
		serverUDPAddr := resolveServerAddr(cfg)
		var clientUDPAddr *net.UDPAddr
		if cfg.Tunnel.IP != "" && cfg.Tunnel.Port != 0 {
			clientUDPAddr = resolveClientAddr(cfg)
		}

		conn, err := udp.New(clientUDPAddr, serverUDPAddr)
		cmd.Exitf(err, "udp.New")

		return conn
	default:
		cmd.Exitf(errors.New("invalid tunnel network"), "unexpected network: %s", cfg.Tunnel.Network)
		return nil
	}
}

func resolveServerAddr(cfg config.ClientConfig) *net.UDPAddr {
	switch cfg.Tunnel.Network {
	case "udp":
		updAddr, err := net.ResolveUDPAddr(
			cfg.Tunnel.Network,
			fmt.Sprintf("%s:%d", cfg.Tunnel.ServerIP, cfg.Tunnel.ServerPort),
		)
		cmd.Exitf(err, "net.ResolveUDPAddr")

		return updAddr
	default:
		cmd.Exitf(errors.New("invalid tunnel network"), "unexpected network: %s", cfg.Tunnel.Network)
		return nil
	}
}

func resolveClientAddr(cfg config.ClientConfig) *net.UDPAddr {
	switch cfg.Tunnel.Network {
	case "udp":
		updAddr, err := net.ResolveUDPAddr(
			cfg.Tunnel.Network,
			fmt.Sprintf("%s:%d", cfg.Tunnel.IP, cfg.Tunnel.Port),
		)
		cmd.Exitf(err, "net.ResolveUDPAddr")

		return updAddr
	default:
		cmd.Exitf(errors.New("invalid tunnel network"), "unexpected network: %s", cfg.Tunnel.Network)
		return nil
	}
}
