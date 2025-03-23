package server

import (
	"github.com/atmxlab/vpn/internal/config"
	"github.com/atmxlab/vpn/internal/pkg/tun"
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

	t := tun.NewTun(tunIface)

	routerBuilder := router.NewBuilder()

	routerBuilder.
		Config(func(b *router.ConfigBuilder) {
			b.
				TunMtu(cfg.Tun.MTU).
				TunSubnet(cfg.Tun.Subnet).
				TunChanSize(cfg.Tun.TunChanSize).
				TunnelChanSize(cfg.Tun.TunnelChanSize)
		}).
		Tun(t)
}
