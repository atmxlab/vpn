package main

import (
	"context"

	"github.com/atmxlab/vpn/cmd"
	"github.com/atmxlab/vpn/internal/config"
	"github.com/atmxlab/vpn/internal/pkg/ipdistributor"
	_ "github.com/atmxlab/vpn/internal/pkg/logger"
	"github.com/atmxlab/vpn/internal/pkg/peermanager"
	"github.com/atmxlab/vpn/internal/pkg/tun"
	"github.com/atmxlab/vpn/internal/pkg/tunnel"
	"github.com/atmxlab/vpn/internal/router"
	tunhandler "github.com/atmxlab/vpn/internal/server/handlers/tun"
	tunnelhandler "github.com/atmxlab/vpn/internal/server/handlers/tunnel"
	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/atmxlab/vpn/pkg/jsonconfig"
)

func main() {
	defer cmd.Recover()

	ctx, cancel := cmd.SignalCtx()
	defer cancel()

	const configPath = "./config/server.json"

	cfg, err := jsonconfig.Load[config.ServerConfig](configPath)
	cmd.Exitf(err, "jsonconfig.Load")

	tunIP, tunSubnet, err := cfg.Tun.GetCIDR()
	cmd.Exitf(err, "cfg.Tun.GetCIDR")

	embeddedTun, err := setupTun(tunIP, tunSubnet, cfg.Tun.MTU)
	cmd.Exitf(err, "setupTun")

	tn := tun.NewTun(embeddedTun)
	tunl := tunnel.New(setupTunnelConn(cfg))
	pm := peermanager.New()

	ipDistributor, err := ipdistributor.New(tunSubnet)
	cmd.Exitf(err, "ipdistributor.New")

	cmd.Exitf(setupOS(), "setupOS")

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

	if err = rt.Run(ctx); errors.IsSomeBut(err, context.Canceled) {
		cmd.Exitf(err, "router.Run")
	}
}
