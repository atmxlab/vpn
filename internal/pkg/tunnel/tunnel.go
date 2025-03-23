package tunnel

import (
	"context"
	"net"
	"sync"

	"github.com/atmxlab/vpn/internal/protocol"
	"github.com/atmxlab/vpn/internal/tunnel"
	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Tunnel struct {
	tunnel tunnel.Tunnel
}

func NewTunnel(tunnel tunnel.Tunnel) *Tunnel {
	return &Tunnel{tunnel: tunnel}
}

func (t *Tunnel) PSH(addr net.Addr, payload []byte) (int, error) {
	tunnelPacket := protocol.NewTunnelPacket(
		protocol.NewHeader(protocol.FlagPSH),
		payload,
		addr,
	)

	return t.Write(tunnelPacket)
}

func (t *Tunnel) ACK(addr net.Addr, payload []byte) (int, error) {
	tunnelPacket := protocol.NewTunnelPacket(
		protocol.NewHeader(protocol.FlagACK),
		payload,
		addr,
	)

	return t.Write(tunnelPacket)
}

func (t *Tunnel) Write(tunnelPacket *protocol.TunnelPacket) (int, error) {
	n, err := t.tunnel.WriteTo(tunnelPacket.Marshal(), tunnelPacket.Addr())
	if err != nil {
		return 0, errors.Wrap(err, "tunnel.WriteTo")
	}

	logrus.Debugf("Write to TUNNEL %d bytes", n)

	return n, nil
}

// ReadFromWithContext - необходим, чтобы учитывать отмену контекста при чтении из тоннеля
func (t *Tunnel) ReadFromWithContext(ctx context.Context, buf []byte) (int, net.Addr, error) {
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

		n, addr, err := t.tunnel.ReadFrom(buf)
		resultChan <- result{n, addr, err}
	}()

	select {
	case <-ctx.Done():
		logrus.Warnf("Context canceled: %v", ctx.Err())
		if err := t.tunnel.Close(); err != nil {
			return 0, nil, errors.Join(err, ctx.Err())
		}

		logrus.Warn("Waiting ending read from Tunnel...")
		wg.Wait()

		return 0, nil, ctx.Err()
	case res := <-resultChan:
		return res.n, res.addr, res.err
	}
}

func (t *Tunnel) Close() error {
	return t.tunnel.Close()
}
