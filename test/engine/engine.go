package engine

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/atmxlab/vpn/internal/config"
	"github.com/atmxlab/vpn/internal/pkg/ipdistributor"
	"github.com/atmxlab/vpn/internal/pkg/peermanager"
	"github.com/atmxlab/vpn/internal/pkg/tun"
	"github.com/atmxlab/vpn/internal/pkg/tunnel"
	tunhandler "github.com/atmxlab/vpn/internal/server/handlers/tun"
	tunnelhandler "github.com/atmxlab/vpn/internal/server/handlers/tunnel"
	"github.com/atmxlab/vpn/internal/server/router"
	"github.com/atmxlab/vpn/test"
	"github.com/atmxlab/vpn/test/stub"
	"github.com/stretchr/testify/require"
)

type engine struct {
	app        *app
	cancelFunc context.CancelFunc
	cfg        *Config
}

func New(
	t *testing.T,
	hook ...func(c *Config),
) test.Engine {
	ctx, cancel := context.WithCancel(context.Background())

	cfg := &Config{
		actionDelay: 10 * time.Millisecond,
		serverConfig: &config.ServerConfig{
			BufferSize:       1500,
			PeerKeepAliveTTL: config.Duration(10 * time.Second),
			Tun: config.ServerTun{
				SubnetCIDR:  "10.0.0.0/24",
				MTU:         1500,
				TunChanSize: 1000,
			},
			Tunnel: config.ServerTunnel{
				TunnelChanSize: 1000,
				Network:        "udp",
				IP:             "",
				Port:           6000,
			},
		},
	}

	test.ApplyHooks(cfg, hook)

	serverAddr := stub.NewAddr(
		cfg.serverConfig.Tunnel.Network,
		fmt.Sprintf("%s:%d", cfg.serverConfig.Tunnel.IP, cfg.serverConfig.Tunnel.Port),
	)

	tunSubnet, err := cfg.serverConfig.Tun.Subnet()
	require.NoError(t, err)

	embeddedTunStub := stub.NewEmbeddedTun("TestTun", int(cfg.serverConfig.Tun.TunChanSize))
	tn := tun.NewTun(embeddedTunStub)

	tunnelConnStub := stub.NewTunnelConnection(serverAddr, int(cfg.serverConfig.Tunnel.TunnelChanSize))
	tunl := tunnel.New(tunnelConnStub)

	pm := peermanager.New()

	ipDistributor, err := ipdistributor.New(tunSubnet)
	require.NoError(t, err, "ipdistributor.New")

	routerBuilder := router.NewBuilder()

	routerBuilder.
		Config(func(b *router.ConfigBuilder) {
			b.
				BufferSize(cfg.serverConfig.BufferSize).
				TunChanSize(cfg.serverConfig.Tun.TunChanSize).
				TunnelChanSize(cfg.serverConfig.Tunnel.TunnelChanSize)
		}).
		Tun(tn).
		Tunnel(tunl).
		TunHandler(tunhandler.NewHandler(tunl, pm)).
		TunnelHandler(func(build *router.TunnelHandlerBuilder) {
			build.SYN(tunnelhandler.NewSYNHandler(pm, tunl, ipDistributor, cfg.serverConfig.PeerKeepAliveTTL.ToDuration()))
			build.FIN(tunnelhandler.NewFINHandler(pm, ipDistributor))
			build.PSH(tunnelhandler.NewPSHHandler(pm, tn, tunl))
			build.KPA(tunnelhandler.NewKPAHandler(pm, cfg.serverConfig.PeerKeepAliveTTL.ToDuration()))
		})

	rt := routerBuilder.Build()

	go func() {
		require.NoError(t, rt.Run(ctx), "router.Run")
	}()

	return &engine{
		app:        newApp(t, ctx, tunnelConnStub, embeddedTunStub, pm, ipDistributor),
		cancelFunc: cancel,
		cfg:        cfg,
	}
}

func (e *engine) REPLAY(actions ...test.Action) {
	e.app.T().Helper()

	defer e.cancelFunc()

	for _, action := range actions {
		action.Handle(e.app)
		time.Sleep(e.cfg.actionDelay)
	}
}
