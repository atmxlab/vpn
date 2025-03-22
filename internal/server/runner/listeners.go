package runner

import (
	"context"
	"net"
	"sync"

	"github.com/atmxlab/vpn/internal/protocol"
	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/ipv4"
)

func (r *Runner) listenTun(ctx context.Context) (err error) {
	for {
		buf, n := make([]byte, r.cfg.BufferSize), 0

		select {
		case <-ctx.Done():
			logrus.Warnf("Context canceled: %v", ctx.Err())
			return ctx.Err()
		default:
			n, err = r.readTunWithContext(ctx, buf)
			if err != nil {
				return errors.Wrap(err, "read with context from TUN interface")
			}

			logrus.Debugf("Readed bytes from TUN; len=[%d]", n)

			if err = r.logIPHeader(buf[:n]); err != nil {
				return errors.Wrap(err, "log IP header")
			}
		}

		select {
		case <-ctx.Done():
			logrus.Warnf("Context canceled: %v", ctx.Err())
			return ctx.Err()
		case r.state.net.tunPackets <- protocol.NewTunPacket(buf[:n]):
		}
	}
}

// readTunWithContext - необходим, чтобы учитывать отмену контекста при чтении из потока
func (r *Runner) readTunWithContext(ctx context.Context, buf []byte) (int, error) {
	type result struct {
		n   int
		err error
	}

	resultChan := make(chan result, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		defer close(resultChan)

		n, err := r.state.net.tun.Read(buf)
		resultChan <- result{n, err}
	}()

	select {
	case <-ctx.Done():
		logrus.Warnf("Context canceled: %v", ctx.Err())
		if err := r.state.net.tun.Close(); err != nil {
			return 0, errors.Join(err, ctx.Err())
		}

		logrus.Warn("Waiting ending read from Tunnel...")
		wg.Wait()

		return 0, ctx.Err()
	case res := <-resultChan:
		return res.n, res.err
	}
}

func (r *Runner) listenTunnel(ctx context.Context) (err error) {
	for {
		buf, n := make([]byte, r.cfg.BufferSize), 0
		var addr net.Addr

		select {
		case <-ctx.Done():
			logrus.Warnf("Context canceled: %v", ctx.Err())
			return ctx.Err()
		default:

			n, addr, err = r.readTunnelWithContext(ctx, buf)
			if err != nil {
				return errors.Wrap(err, "read with context from TUN interface")
			}

			logrus.Debugf("Readed bytes from TUN; len=[%d]", n)

			if err = r.logIPHeader(buf[:n]); err != nil {
				return errors.Wrap(err, "log IP header")
			}
		}

		select {
		case <-ctx.Done():
			logrus.Warnf("Context canceled: %v", ctx.Err())
			return ctx.Err()
		case r.state.net.tunnelPackets <- protocol.UnmarshalTunnelPacket(addr, buf[:n]):
		}
	}
}

// readTunnelWithContext - необходим, чтобы учитывать отмену контекста при чтении из тоннеля
func (r *Runner) readTunnelWithContext(ctx context.Context, buf []byte) (int, net.Addr, error) {
	type result struct {
		n    int
		addr net.Addr
		err  error
	}

	resultChan := make(chan result, 1)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		defer close(resultChan)

		n, addr, err := r.state.net.tunnel.ReadFrom(buf)
		resultChan <- result{n, addr, err}
	}()

	select {
	case <-ctx.Done():
		logrus.Warnf("Context canceled: %v", ctx.Err())
		if err := r.state.net.tunnel.Close(); err != nil {
			return 0, nil, errors.Join(err, ctx.Err())
		}

		logrus.Warn("Waiting ending read from Tunnel...")
		wg.Wait()

		return 0, nil, ctx.Err()
	case res := <-resultChan:
		return res.n, res.addr, res.err
	}
}

func (r *Runner) logIPHeader(frame []byte) error {
	header, err := ipv4.ParseHeader(frame)
	if err != nil {
		return errors.Wrap(err, "parse header")
	}

	logrus.Debug(header.String())

	return nil
}
