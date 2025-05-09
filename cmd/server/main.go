package main

import (
	"context"
	"runtime/debug"

	"github.com/atmxlab/vpn/internal/config"
	"github.com/atmxlab/vpn/internal/pkg/details/route"
	"github.com/atmxlab/vpn/internal/pkg/ipdistributor"
	_ "github.com/atmxlab/vpn/internal/pkg/logger"
	"github.com/atmxlab/vpn/internal/pkg/peermanager"
	"github.com/atmxlab/vpn/internal/pkg/tun"
	"github.com/atmxlab/vpn/internal/pkg/tunnel"
	"github.com/atmxlab/vpn/internal/router"
	tunhandler "github.com/atmxlab/vpn/internal/server/handlers/tun"
	tunnelhandler "github.com/atmxlab/vpn/internal/server/handlers/tunnel"
	"github.com/atmxlab/vpn/pkg/jsonconfig"
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

	const configPath = "./config/server.json"

	cfg, err := jsonconfig.Load[config.ServerConfig](configPath)
	exitf(err, "jsonconfig.Load")

	tunSubnet, err := cfg.Tun.Subnet()
	exitf(err, "cfg.Tun.Subnet")

	embeddedTun, err := setupTun(tunSubnet, cfg.Tun.MTU)
	exitf(err, "setupTun")

	tn := tun.NewTun(embeddedTun)
	tunl := tunnel.New(setupTunnelConn(cfg))
	pm := peermanager.New()

	ipDistributor, err := ipdistributor.New(tunSubnet)
	exitf(err, "ipdistributor.New")

	exitf(setupOS(route.NewConfigurator(), cfg), "setupOS")

	routerBuilder := router.NewBuilder()

	routerBuilder.
		Config(func(b *router.ConfigBuilder) {
			b.
				BufferSize(cfg.BufferSize).
				TunChanSize(cfg.Tun.TunChanSize).
				TunnelChanSize(cfg.Tunnel.TunnelChanSize)
		}).
		Tun(tn).
		Tunnel(tunl).
		TunHandler(tunhandler.NewHandler(tunl, pm)).
		TunnelHandler(func(build *router.TunnelHandlerBuilder) {
			build.SYN(tunnelhandler.NewSYNHandler(pm, tunl, ipDistributor, cfg.PeerKeepAliveTTL.ToDuration()))
			build.FIN(tunnelhandler.NewFINHandler(pm, ipDistributor))
			build.PSH(tunnelhandler.NewPSHHandler(pm, tn, tunl))
			build.KPA(tunnelhandler.NewKPAHandler(pm, cfg.PeerKeepAliveTTL.ToDuration()))
		})

	rt := routerBuilder.Build()

	exitf(rt.Run(ctx), "router.Run")
}
