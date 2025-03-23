package main

import (
	"context"

	"github.com/atmxlab/vpn/internal/config"
	"github.com/atmxlab/vpn/internal/pkg/details/conn/server/udp"
	"github.com/atmxlab/vpn/internal/pkg/details/route"
	"github.com/atmxlab/vpn/internal/pkg/ipdistributor"
	"github.com/atmxlab/vpn/internal/pkg/peermanager"
	"github.com/atmxlab/vpn/internal/pkg/tun"
	"github.com/atmxlab/vpn/internal/pkg/tunnel"
	tunhandler "github.com/atmxlab/vpn/internal/server/handlers/tun"
	tunnelhandler "github.com/atmxlab/vpn/internal/server/handlers/tunnel"
	"github.com/atmxlab/vpn/internal/server/router"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.ServerConfig{
		ServerAddr:       nil,
		BufferSize:       0,
		PeerKeepAliveTTL: 0,
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

	pm := peermanager.New()

	ipDistributor, err := ipdistributor.New(cfg.Tun.Subnet)
	if err != nil {
		logrus.Fatal(err)
	}

	rc := route.NewConfigurator()

	routerBuilder := router.NewBuilder()

	routerBuilder.
		Config(func(b *router.ConfigBuilder) {
			b.
				BufferSize(cfg.BufferSize).
				TunMtu(cfg.Tun.MTU).
				TunSubnet(cfg.Tun.Subnet).
				TunChanSize(cfg.Tun.TunChanSize).
				TunnelChanSize(cfg.Tun.TunnelChanSize)
		}).
		RouteConfigurator(rc).
		Tun(tn).
		Tunnel(tunl).
		TunHandler(tunhandler.NewHandler(tunl, pm)).
		TunnelHandler(func(build *router.TunnelHandlerBuilder) {
			build.SYN(tunnelhandler.NewSYNHandler(pm, tunl, ipDistributor, cfg.PeerKeepAliveTTL))
			build.FIN(tunnelhandler.NewFINHandler(pm, ipDistributor))
			build.PSH(tunnelhandler.NewPSHHandler(pm, tn, tunl))
			build.KPA(tunnelhandler.NewKPAHandler(pm, cfg.PeerKeepAliveTTL))
		})

	rt := routerBuilder.Build()

	if err = rt.Run(ctx); err != nil {
		logrus.Fatal(err)
	}
}
