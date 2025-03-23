package tunnel

import (
	"context"
	"io"
	"net"
	"sync"

	"github.com/atmxlab/vpn/internal/protocol"
	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Connection - тоннель через который проходит трафик
// от VPN-сервера к VPN-клиент и обратно
type Connection interface {
	ReadFrom(p []byte) (n int, addr net.Addr, err error)
	WriteTo(p []byte, addr net.Addr) (n int, err error)
	LocalAddr() net.Addr
	io.Closer
}

type Tunnel struct {
	conn Connection
}

func New(conn Connection) *Tunnel {
	return &Tunnel{conn: conn}
}

func (t *Tunnel) PSH(addr net.Addr, payload []byte) (int, error) {
	tunnelPacket := protocol.NewTunnelPacket(
		protocol.NewHeader(protocol.FlagPSH),
		payload,
		addr,
	)

	return t.Write(tunnelPacket)
}

func (t *Tunnel) SYN(addr net.Addr, payload []byte) (int, error) {
	tunnelPacket := protocol.NewTunnelPacket(
		protocol.NewHeader(protocol.FlagSYN),
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
	n, err := t.conn.WriteTo(tunnelPacket.Marshal(), tunnelPacket.Addr())
	if err != nil {
		return 0, errors.Wrap(err, "conn.WriteTo")
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

		n, addr, err := t.conn.ReadFrom(buf)
		resultChan <- result{n, addr, err}
	}()

	select {
	case <-ctx.Done():
		logrus.Warnf("Context canceled: %v", ctx.Err())
		if err := t.conn.Close(); err != nil {
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
	return t.conn.Close()
}
