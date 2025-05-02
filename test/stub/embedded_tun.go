package stub

import "github.com/atmxlab/vpn/internal/protocol"

type EmbeddedTun struct {
	input chan []byte

	output        chan []byte
	outputPackets []*protocol.TunPacket

	name string
}

func NewEmbeddedTun(name string, dataChanSize int) *EmbeddedTun {
	return &EmbeddedTun{
		input:  make(chan []byte, dataChanSize),
		output: make(chan []byte, dataChanSize),
		name:   name,
	}
}

func (e *EmbeddedTun) Read(p []byte) (n int, err error) {
	bytes := <-e.input
	return copy(p, bytes), nil
}

func (e *EmbeddedTun) ReadFromOutput(p []byte) (n int, err error) {
	bytes := <-e.output
	return copy(p, bytes), nil
}

func (e *EmbeddedTun) Write(p []byte) (n int, err error) {
	e.output <- p
	e.outputPackets = append(e.outputPackets, protocol.NewTunPacket(p))

	return len(p), nil
}

func (e *EmbeddedTun) WriteToInput(p []byte) (n int, err error) {
	e.input <- p
	return len(p), nil
}

func (e *EmbeddedTun) Close() error {
	close(e.input)
	close(e.output)
	return nil
}

func (e *EmbeddedTun) Name() string {
	return e.name
}

func (e *EmbeddedTun) GetLastPacket() (*protocol.TunPacket, bool) {
	if len(e.outputPackets) == 0 {
		return nil, false
	}

	return e.outputPackets[len(e.outputPackets)-1], true
}
