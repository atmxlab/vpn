package protocol

import (
	"bytes"
	"net"
)

// TunnelPacket - пакет, который передается через тоннель
type TunnelPacket struct {
	header  *Header
	payload Payload
	// Адрес назначения пакета
	addr net.Addr
}

func NewTunnelPacket(header *Header, payload Payload, addr net.Addr) *TunnelPacket {
	return &TunnelPacket{header: header, payload: payload, addr: addr}
}

func (p *TunnelPacket) Header() *Header {
	return p.header
}

func (p *TunnelPacket) Payload() Payload {
	return p.payload
}

func (p *TunnelPacket) Addr() net.Addr {
	return p.addr
}

func (p *TunnelPacket) Marshal() []byte {
	var buf []byte

	buffer := bytes.NewBuffer(buf)

	buffer.WriteByte(p.Header().Flag().Byte())
	buffer.Write(p.Payload())

	return buffer.Bytes()
}

func UnmarshalTunnelPacket(addr net.Addr, bytes []byte) *TunnelPacket {
	return NewTunnelPacket(
		// Первый байт это флаг
		NewHeader(Flag(bytes[0])),
		// Все что после заголовка - это полезные данные
		bytes[HeaderSize:],
		addr,
	)
}
