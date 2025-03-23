package server

import (
	"github.com/atmxlab/vpn/internal/config"
	"github.com/atmxlab/vpn/internal/server/router"
)

func main() {
	cfg := config.ServerConfig{
		ServerAddr:                  nil,
		BufferSize:                  0,
		PeerKeepAliveMissingTimeout: 0,
	}

	routerBuilder := router.NewBuilder()

	routerBuilder.
		Config(func(b *router.ConfigBuilder) {
			b.
				TunMtu(cfg.Tun.MTU).
				TunSubnet(cfg.Tun.Subnet).
				TunChanSize(cfg.Tun.TunChanSize).
				TunnelChanSize(cfg.Tun.TunnelChanSize)
		})
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
