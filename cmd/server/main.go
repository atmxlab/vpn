package main

import (
	"context"
	"net"
	"runtime/debug"

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
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("Panic recovered: %v", err)
			logrus.Fatalf("Stack trace:\n%s", debug.Stack())
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// TODO: вынести в конфиг
	updAddr, err := net.ResolveUDPAddr("udp", ":6000")
	exitf(err, "net.ResolveUDPAddr")

	cfg := config.ServerConfig{
		ServerAddr:       updAddr,
		BufferSize:       0,
		PeerKeepAliveTTL: 0,
		Tun: config.ServerTun{
			Subnet: net.IPNet{
				IP:   net.ParseIP("10.0.0.0"),
				Mask: net.CIDRMask(24, 32),
			},
			MTU:            1500, // TODO: подобрать
			TunChanSize:    1000,
			TunnelChanSize: 1000,
		},
	}

	embeddedTun, err := setupTun(cfg.Tun.Subnet, cfg.Tun.MTU)
	exitf(err, "setupTun")

	tn := tun.NewTun(embeddedTun)

	conn, err := udp.New(cfg.ServerAddr)
	exitf(err, "udp.New")

	tunl := tunnel.New(conn)

	pm := peermanager.New()

	ipDistributor, err := ipdistributor.New(cfg.Tun.Subnet)
	exitf(err, "ipdistributor.New")

	exitf(setupOS(route.NewConfigurator(), cfg), "setupOS")

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

	exitf(rt.Run(ctx), "router.Run")
}
