package runner

import (
	"context"

	"github.com/atmxlab/vpn/internal/protocol"
	"github.com/atmxlab/vpn/internal/server"
	"github.com/atmxlab/vpn/internal/tunnel"
	"github.com/atmxlab/vpn/internal/tuntap"
	"github.com/atmxlab/vpn/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type state struct {
	ctx struct {
		// Необходим для грейс фул остановки сервера
		cancel context.CancelFunc
	}

	// Управляющий пирами
	peerManager server.PeerManager
	// Распределитель IP адресов
	ipDistributor server.IpDistributor

	net struct {
		// Туннель между сервером и клиентом
		tunnel tunnel.Tunnel
		// Виртуальный L3 сетевой интерфейс - выход в интернет
		tun tuntap.Tun

		// Канал, в который попадают пакеты из тоннеля - от клиента
		// Пакеты из этого канала отправляются в TUN интерфейс - в интернет
		tunnelPackets chan *protocol.TunnelPacket
		// Канал, в который попадают пакеты из TUN интерфейса - из интернета
		// Пакеты из этого канала отправляются в тоннель - клиенту
		tunPackets chan *protocol.TunPacket
	}
}

type Runner struct {
	cfg                  server.Config
	ipDistributorFactory server.IpDistributorFactory
	tunnelFactory        tunnel.Factory
	tunFactory           tuntap.TunFactory
	routeConfigurator    server.RouteConfigurator

	state state
}

func NewRunner(
	cfg server.Config,
	ipDistributorFactory server.IpDistributorFactory,
	tunnelFactory tunnel.Factory,
	tunFactory tuntap.TunFactory,
	routeConfigurator server.RouteConfigurator,
	peerManager server.PeerManager,
) *Runner {
	return &Runner{
		cfg:                  cfg,
		ipDistributorFactory: ipDistributorFactory,
		tunnelFactory:        tunnelFactory,
		tunFactory:           tunFactory,
		routeConfigurator:    routeConfigurator,
		state: state{
			peerManager: peerManager,
			net: struct {
				tunnel        tunnel.Tunnel
				tun           tuntap.Tun
				tunnelPackets chan *protocol.TunnelPacket
				tunPackets    chan *protocol.TunPacket
			}{
				tunPackets:    make(chan *protocol.TunPacket, 1024),
				tunnelPackets: make(chan *protocol.TunnelPacket, 1024),
			},
		},
	}
}

func (r *Runner) Run(ctx context.Context) error {
	defer r.state.net.tun.Close()
	defer r.state.net.tunnel.Close()
	defer close(r.state.net.tunPackets)
	defer close(r.state.net.tunnelPackets)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	r.state.ctx.cancel = cancel

	if err := r.setup(); err != nil {
		return errors.Wrap(err, "setup")
	}

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		if err := r.listenTun(ctx); err != nil {
			return errors.Wrap(err, "listen TUN")
		}
		return nil
	})

	eg.Go(func() error {
		if err := r.listenTunnel(ctx); err != nil {
			return errors.Wrap(err, "listen Tunnel")
		}
		return nil
	})

	// TODO: read channels

	if err := eg.Wait(); err != nil {
		return errors.Wrap(err, "error group wait")
	}

	return nil
}

func (r *Runner) Stop() {
	if r.state.ctx.cancel != nil {
		r.state.ctx.cancel()
	}
}
