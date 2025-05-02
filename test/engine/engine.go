package engine

import (
	"context"
	"net"
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
	"github.com/atmxlab/vpn/test/gen"
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
			ServerAddr:       gen.RandAddr(),
			BufferSize:       1500,
			PeerKeepAliveTTL: 10 * time.Second,
			Tun: config.ServerTun{
				Subnet: net.IPNet{
					IP:   net.ParseIP("10.0.0.0"),
					Mask: net.CIDRMask(24, 32),
				},
				MTU:            1500,
				TunChanSize:    1000,
				TunnelChanSize: 1000,
			},
		},
	}

	test.ApplyHooks(cfg, hook)

	embeddedTunStub := stub.NewEmbeddedTun("TestTun", int(cfg.serverConfig.Tun.TunChanSize))
	tn := tun.NewTun(embeddedTunStub)

	tunnelConnStub := stub.NewTunnelConnection(cfg.serverConfig.ServerAddr, int(cfg.serverConfig.Tun.TunnelChanSize))
	tunl := tunnel.New(tunnelConnStub)

	pm := peermanager.New()

	ipDistributor, err := ipdistributor.New(cfg.serverConfig.Tun.Subnet)
	require.NoError(t, err, "ipdistributor.New")

	routerBuilder := router.NewBuilder()

	routerBuilder.
		Config(func(b *router.ConfigBuilder) {
			b.
				BufferSize(cfg.serverConfig.BufferSize).
				TunMtu(cfg.serverConfig.Tun.MTU).
				TunSubnet(cfg.serverConfig.Tun.Subnet).
				TunChanSize(cfg.serverConfig.Tun.TunChanSize).
				TunnelChanSize(cfg.serverConfig.Tun.TunnelChanSize)
		}).
		Tun(tn).
		Tunnel(tunl).
		TunHandler(tunhandler.NewHandler(tunl, pm)).
		TunnelHandler(func(build *router.TunnelHandlerBuilder) {
			build.SYN(tunnelhandler.NewSYNHandler(pm, tunl, ipDistributor, cfg.serverConfig.PeerKeepAliveTTL))
			build.FIN(tunnelhandler.NewFINHandler(pm, ipDistributor))
			build.PSH(tunnelhandler.NewPSHHandler(pm, tn, tunl))
			build.KPA(tunnelhandler.NewKPAHandler(pm, cfg.serverConfig.PeerKeepAliveTTL))
		})

	rt := routerBuilder.Build()

	go func() {
		require.NoError(t, rt.Run(ctx), "router.Run")
	}()

	return &engine{
		app:        newApp(t, ctx, tunnelConnStub, embeddedTunStub, pm),
		cancelFunc: cancel,
		cfg:        cfg,
	}
}

func (e *engine) REPLAY(actions ...test.Action) {
	defer e.cancelFunc()

	for _, action := range actions {
		action.Handle(e.app)
		time.Sleep(e.cfg.actionDelay)
	}
}
