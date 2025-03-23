package main

import (
	"github.com/atmxlab/vpn/internal/config"
	"github.com/atmxlab/vpn/internal/pkg/details/conn/server/udp"
	"github.com/atmxlab/vpn/internal/pkg/tun"
	"github.com/atmxlab/vpn/internal/pkg/tunnel"
	"github.com/atmxlab/vpn/internal/server/router"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.ServerConfig{
		ServerAddr:                  nil,
		BufferSize:                  0,
		PeerKeepAliveMissingTimeout: 0,
	}

	tunIface, err := setupTun(cfg.Tun.Subnet, cfg.Tun.MTU)
	if err != nil {
		logrus.Fatal(err)
	}

	tn := tun.NewTun(tunIface)

	conn, err := udp.New(cfg.ServerAddr)
	if err != nil {
		logrus.Fatal(err)
	}

	tunl := tunnel.New(conn)

	routerBuilder := router.NewBuilder()

	routerBuilder.
		Config(func(b *router.ConfigBuilder) {
			b.
				TunMtu(cfg.Tun.MTU).
				TunSubnet(cfg.Tun.Subnet).
				TunChanSize(cfg.Tun.TunChanSize).
				TunnelChanSize(cfg.Tun.TunnelChanSize)
		}).
		Tun(tn).
		Tunnel(tunl)
}
