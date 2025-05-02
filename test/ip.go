package test

import (
	"encoding/binary"
	"net"
	"testing"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/stretchr/testify/require"
)

type IPPacket struct {
	srcIP   net.IP
	dstIP   net.IP
	payload []byte

	bytes []byte
}

func (p *IPPacket) DstIP() net.IP {
	return p.dstIP
}

func (p *IPPacket) SrcIP() net.IP {
	return p.srcIP
}

func (p *IPPacket) Payload() []byte {
	return p.payload
}

func (p *IPPacket) Bytes() []byte {
	return p.bytes
}

type IPPacketBuilder struct {
	t       *testing.T
	version uint8
	srcIP   net.IP
	dstIP   net.IP
	payload []byte
	proto   layers.IPProtocol
}

func NewIPPacketBuilder(t *testing.T) *IPPacketBuilder {
	return &IPPacketBuilder{version: 4, t: t}
}

func (i *IPPacketBuilder) Version(version uint8) *IPPacketBuilder {
	i.version = version
	return i
}

func (i *IPPacketBuilder) SrcIP(srcIP net.IP) *IPPacketBuilder {
	i.srcIP = srcIP
	return i
}

func (i *IPPacketBuilder) DstIP(dstIP net.IP) *IPPacketBuilder {
	i.dstIP = dstIP
	return i
}

func (i *IPPacketBuilder) Payload(payload []byte) *IPPacketBuilder {
	i.payload = payload
	return i
}

func (i *IPPacketBuilder) TCP() *IPPacketBuilder {
	i.proto = layers.IPProtocolTCP
	return i
}

func (i *IPPacketBuilder) UDP() *IPPacketBuilder {
	i.proto = layers.IPProtocolUDP
	return i
}

func (i *IPPacketBuilder) ICMPv4() *IPPacketBuilder {
	i.proto = layers.IPProtocolICMPv4
	return i
}

func (i *IPPacketBuilder) Build() *IPPacket {
	switch i.proto {
	case layers.IPProtocolTCP:
		return i.buildTCP()
	}

	i.t.Fatalf("Unsupported IP protocol: %v", i.proto)
	return nil
}

func (i *IPPacketBuilder) buildTCP() *IPPacket {
	ip := &layers.IPv4{
		Version:  i.version,
		TTL:      64,
		SrcIP:    i.srcIP,
		DstIP:    i.dstIP,
		Protocol: i.proto,
	}

	// Создание полезной нагрузки (например, TCP-сегмента)
	tcp := &layers.TCP{
		SrcPort: layers.TCPPort(12345),
		DstPort: layers.TCPPort(80),
		SYN:     true, // Установка флага SYN (для TCP)
	}

	// Сборка пакета
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	payload := gopacket.Payload(i.payload)

	// Установка контрольной суммы TCP (если нужно)
	require.NoError(i.t, tcp.SetNetworkLayerForChecksum(ip))

	// Сериализация пакета
	require.NoError(i.t, gopacket.SerializeLayers(buf, opts, ip, tcp, payload))

	return &IPPacket{
		srcIP:   i.srcIP,
		dstIP:   i.dstIP,
		payload: i.payload,
		bytes:   buf.Bytes(),
	}
}

func IPHeaderChecksum(b []byte) uint16 {
	csum := uint32(0)
	for i := 0; i < len(b); i += 2 {
		csum += uint32(binary.BigEndian.Uint16(b[i:min(i+2, len(b))]))
	}
	for csum > 0xffff {
		csum = (csum >> 16) + (csum & 0xffff)
	}
	return ^uint16(csum)
}
