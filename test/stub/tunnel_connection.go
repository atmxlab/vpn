package stub

import (
	"io"
	"net"
	"sync"

	"github.com/atmxlab/vpn/internal/protocol"
)

type TunnelConnection struct {
	mu        sync.Mutex
	closeOnce sync.Once

	input chan *protocol.TunnelPacket

	output        chan *protocol.TunnelPacket
	outputPackets []*protocol.TunnelPacket

	addr net.Addr
}

func NewTunnelConnection(addr net.Addr, dataChanSize int) *TunnelConnection {
	return &TunnelConnection{
		input:  make(chan *protocol.TunnelPacket, dataChanSize),
		output: make(chan *protocol.TunnelPacket, dataChanSize),
		addr:   addr,
	}
}

func (t *TunnelConnection) ReadFrom(p []byte) (int, net.Addr, error) {
	tp, ok := <-t.input
	// Канал может быть закрыт (e.g. Отмена контекста)
	// В таком случае отдаем EOF
	if !ok {
		return 0, nil, io.EOF
	}

	payload := tp.Marshal()
	n := copy(p, payload)

	return n, tp.Addr(), nil
}

func (t *TunnelConnection) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	tp := protocol.UnmarshalTunnelPacket(addr, p)

	t.output <- tp
	t.mu.Lock()
	defer t.mu.Unlock()
	t.outputPackets = append(t.outputPackets, tp)

	return len(tp.Marshal()), nil
}

func (t *TunnelConnection) WriteToInput(p []byte, addr net.Addr) (n int, err error) {
	tp := protocol.UnmarshalTunnelPacket(addr, p)
	t.input <- tp
	return len(tp.Marshal()), nil
}

func (t *TunnelConnection) LocalAddr() net.Addr {
	return t.addr
}

func (t *TunnelConnection) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	close(t.input)
	close(t.output)

	return nil
}

func (t *TunnelConnection) GetLastPacket() (*protocol.TunnelPacket, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if len(t.outputPackets) == 0 {
		return nil, false
	}

	return t.outputPackets[len(t.outputPackets)-1], true
}
