package main

import (
	"context"
	"net"

	"github.com/atmxlab/vpn/cmd"
	"github.com/atmxlab/vpn/internal/client"
	"github.com/atmxlab/vpn/internal/client/actions"
	tunhandler "github.com/atmxlab/vpn/internal/client/handlers/tun"
	tunnelhandler "github.com/atmxlab/vpn/internal/client/handlers/tunnel"
	"github.com/atmxlab/vpn/internal/config"
	"github.com/atmxlab/vpn/internal/pkg/closer"
	"github.com/atmxlab/vpn/internal/pkg/details/client/configurator"
	"github.com/atmxlab/vpn/internal/pkg/ip"
	_ "github.com/atmxlab/vpn/internal/pkg/logger"
	"github.com/atmxlab/vpn/internal/pkg/tun"
	"github.com/atmxlab/vpn/internal/pkg/tunnel"
	"github.com/atmxlab/vpn/internal/router"
	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/atmxlab/vpn/pkg/jsonconfig"
	"github.com/atmxlab/vpn/pkg/signal"
)

func main() {
	defer cmd.Recover()

	ctx, cancel := cmd.SignalCtx()
	defer cancel()

	const configPath = "./config/client.json"

	cfg, err := jsonconfig.Load[config.ClientConfig](configPath)
	cmd.Exitf(err, "jsonconfig.Load")

	serverAddr := resolveServerAddr(cfg)

	netConfigurator := configurator.NewConfigurator()

	tunIPMask := ip.MaskFromIP(net.ParseIP(cfg.Tun.IPMask))

	signaller := signal.NewSignaller()

	embeddedTun, err := setupTun(cfg.Tun.MTU)
	cmd.Exitf(err, "setupTun")

	err = setupOS(cfg.Tunnel.GetServerIP(), cfg.GetGatewayIP())
	cmd.Exitf(err, "setupOS")

	tn := tun.NewTun(embeddedTun)
	tunl := tunnel.New(setupTunnelConn(cfg))

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
		TunHandler(tunhandler.NewHandler(tunl, serverAddr)).
		TunnelHandler(func(build *router.TunnelHandlerBuilder) {
			build.ACK(tunnelhandler.NewACKHandler(netConfigurator, signaller, tunIPMask))
			build.SYN(tunnelhandler.NewSYNHandler(tunl))
			build.FIN(tunnelhandler.NewFINHandler(closer.NewCloser(cancel)))
			build.PSH(tunnelhandler.NewPSHHandler(tn))
		})

	rt := routerBuilder.Build()

	synAction := actions.NewSYNAction(tunl, serverAddr)
	kpaAction := actions.NewKPAAction(tunl, serverAddr, cfg.KeepAliveTickDuration.ToDuration())

	c := client.NewClient(rt, synAction, kpaAction, signaller)

	if err = c.Run(ctx, cfg.Tunnel.ServerConnectionTimeout.ToDuration()); errors.IsSomeBut(err, context.Canceled) {
		cmd.Exitf(err, "client.Run")
	}
}
