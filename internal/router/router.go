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

	cfg *config

	tunHandler          TunHandler
	tunnelHandlerByFlag map[protocol.Flag]TunnelHandler

	eg *errgroup.Group
	// Необходим для грейс фул остановки сервера
	cancel context.CancelFunc
}

func (r *Router) Run(ctx context.Context) error {
	defer r.tun.Close()
	defer r.tunnel.Close()

	log := logrus.WithField("Namespace", "ROUTER")

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	r.cancel = cancel

	r.eg, ctx = errgroup.WithContext(ctx)

	r.eg.Go(func() error {
		defer close(r.tunPackets)
		defer log.Debug("Stop tun listening")

		log.Debug("Start tun listening")
		if err := r.listenTun(ctx); err != nil {
			return errors.Wrap(err, "listen TUN")
		}

		return nil
	})

	r.eg.Go(func() error {
		defer close(r.tunnelPackets)
		defer log.Debug("Stop tunnel listening")

		log.Debug("Start tunnel listening")
		if err := r.listenTunnel(ctx); err != nil {
			return errors.Wrap(err, "listen Tunnel")
		}
		return nil
	})

	r.eg.Go(func() error {
		defer log.Debug("Stop tun consuming")

		log.Debug("Start tun consuming.")
		if err := r.consumeTun(ctx); err != nil {
			return errors.Wrap(err, "tun consumer")
		}
		return nil
	})

	r.eg.Go(func() error {
		defer log.Debug("Stop tunnel consuming")

		log.Debug("Start tunnel consuming")
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

func (r *Router) Close() error {
	if r.cancel == nil {
		return nil
	}

	r.cancel()

	if err := r.eg.Wait(); err != nil {
		return errors.Wrap(err, "error group wait")
	}

	return nil
}
