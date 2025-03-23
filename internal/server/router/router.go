package router

import (
	"context"
	"io"
	"net"

	"github.com/atmxlab/vpn/internal/protocol"
	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

// RouteConfigurator - конфигурирует сеть на сервере под VPN сервер
type RouteConfigurator interface {
	// EnableIPForward - включает транзит IP пакетов на сервере
	// e.g. sysctl -w net.ipv4.ip_forward=1
	EnableIPForward() error
	// ConfigureFirewall - конфигурирует сетевой фильтр
	// e.g. netfilter - iptables, nftables
	// TODO: в идеале эта штука должна создать отдельную цепочку
	ConfigureFirewall(subnet net.IPNet) error
	// SetDefaultRoute - указывает шлюз по умолчанию для подсети
	SetDefaultRoute(subnet net.IPNet) error
}

// Tun - TUN интерфейс - из него читаем и отдаем пакеты обработчику
type Tun interface {
	ReadWithContext(ctx context.Context, data []byte) (int, error)
	io.Closer
}

// Tunnel - тоннель - из нее читаем и отдаем пакеты обработчику
type Tunnel interface {
	ReadFromWithContext(ctx context.Context, p []byte) (n int, addr net.Addr, err error)
	io.Closer
}

// TunHandler - обрабатывает пакеты из TUN интерфейса
type TunHandler interface {
	Handle(
		ctx context.Context,
		packet *protocol.TunPacket,
	) error
}

// TunnelHandler - обрабатывает пакеты из тоннеля
type TunnelHandler interface {
	Handle(
		ctx context.Context,
		packet *protocol.TunnelPacket,
	) error
}

type Router struct {
	// Туннель между сервером и клиентом
	tunnel Tunnel
	// Виртуальный L3 сетевой интерфейс - выход в интернет
	tun Tun

	// Канал, в который попадают пакеты из тоннеля - от клиента
	// Пакеты из этого канала отправляются в TUN интерфейс - в интернет
	tunnelPackets chan *protocol.TunnelPacket
	// Канал, в который попадают пакеты из TUN интерфейса - из интернета
	// Пакеты из этого канала отправляются в тоннель - клиенту
	tunPackets chan *protocol.TunPacket

	cfg               *config
	routeConfigurator RouteConfigurator

	tunHandler          TunHandler
	tunnelHandlerByFlag map[protocol.Flag]TunnelHandler

	eg *errgroup.Group
	// Необходим для грейс фул остановки сервера
	cancel context.CancelFunc
}

func (r *Router) Run(ctx context.Context) error {
	defer r.tun.Close()
	defer r.tunnel.Close()
	defer close(r.tunPackets)
	defer close(r.tunnelPackets)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	r.cancel = cancel

	if err := r.setup(); err != nil {
		return errors.Wrap(err, "setup")
	}

	r.eg, ctx = errgroup.WithContext(ctx)

	r.eg.Go(func() error {
		if err := r.listenTun(ctx); err != nil {
			return errors.Wrap(err, "listen TUN")
		}
		return nil
	})

	r.eg.Go(func() error {
		if err := r.listenTunnel(ctx); err != nil {
			return errors.Wrap(err, "listen Tunnel")
		}
		return nil
	})

	r.eg.Go(func() error {
		if err := r.consumeTun(ctx); err != nil {
			return errors.Wrap(err, "tun consumer")
		}
		return nil
	})

	r.eg.Go(func() error {
		if err := r.consumeTunnel(ctx); err != nil {
			return errors.Wrap(err, "tunnel consumer")
		}
		return nil
	})

	if err := r.eg.Wait(); err != nil {
		return errors.Wrap(err, "error group wait")
	}

	return nil
}

func (r *Router) Stop() {
	if r.cancel == nil {
		return
	}

	r.cancel()

	if err := r.eg.Wait(); err != nil {
		logrus.Errorf("error group wait: %v", err)
	}
}
