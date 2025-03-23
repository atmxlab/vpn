package router

import (
	"net"

	"github.com/atmxlab/vpn/internal/protocol"
)

type Builder struct {
	cfg                 *config
	tunnel              Tunnel
	tun                 Tun
	tunnelChanSize      uint
	tunChanSize         uint
	routeConfigurator   RouteConfigurator
	tunHandler          TunHandler
	tunnelHandlerByFlag map[protocol.Flag]TunnelHandler
}

func NewBuilder() *Builder {
	return &Builder{
		tunChanSize:         1024,
		tunnelChanSize:      1024,
		tunnelHandlerByFlag: make(map[protocol.Flag]TunnelHandler, len(protocol.Flags())),
	}
}

func (b *Builder) Config(fn func(b *ConfigBuilder)) *Builder {
	builder := newConfigBuilder()
	fn(builder)

	b.cfg = builder.build()

	return b
}

func (b *Builder) Tunnel(tunnel Tunnel) *Builder {
	b.tunnel = tunnel
	return b
}

func (b *Builder) Tun(tun Tun) *Builder {
	b.tun = tun
	return b
}

func (b *Builder) TunnelChanSize(size uint) *Builder {
	b.tunnelChanSize = size
	return b
}

func (b *Builder) TunChanSize(size uint) *Builder {
	b.tunChanSize = size
	return b
}

func (b *Builder) RouteConfigurator(routeConfigurator RouteConfigurator) *Builder {
	b.routeConfigurator = routeConfigurator
	return b
}

func (b *Builder) TunHandler(tunHandler TunHandler) *Builder {
	b.tunHandler = tunHandler
	return b
}

func (b *Builder) TunnelHandler(fn func(build *TunnelHandlerBuilder)) *Builder {
	builder := newTunnelHandlerBuilder()
	fn(builder)

	b.tunnelHandlerByFlag = builder.Build()

	return b
}

type ConfigBuilder struct {
	bufferSize     uint16
	subnet         net.IPNet
	mtu            uint16
	tunnelChanSize uint
	tunChanSize    uint
}

func newConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{}
}

func (c *ConfigBuilder) BufferSize(bufferSize uint16) *ConfigBuilder {
	c.bufferSize = bufferSize
	return c
}

func (c *ConfigBuilder) TunSubnet(subnet net.IPNet) *ConfigBuilder {
	c.subnet = subnet
	return c
}

func (c *ConfigBuilder) TunMtu(mtu uint16) *ConfigBuilder {
	c.mtu = mtu
	return c
}

func (c *ConfigBuilder) TunnelChanSize(tunnelChanSize uint) *ConfigBuilder {
	c.tunnelChanSize = tunnelChanSize
	return c
}

func (c *ConfigBuilder) TunChanSize(tunChanSize uint) *ConfigBuilder {
	c.tunChanSize = tunChanSize
	return c
}

func (c *ConfigBuilder) build() *config {
	return &config{
		bufferSize: c.bufferSize,
		tun: struct {
			subnet net.IPNet
			mtu    uint16
		}{
			subnet: c.subnet,
			mtu:    c.mtu,
		},
		tunnelChanSize: c.tunnelChanSize,
		tunChanSize:    c.tunChanSize,
	}
}

type TunnelHandlerBuilder struct {
	tunnelHandlerByFlag map[protocol.Flag]TunnelHandler
}

func newTunnelHandlerBuilder() *TunnelHandlerBuilder {
	return &TunnelHandlerBuilder{tunnelHandlerByFlag: make(map[protocol.Flag]TunnelHandler, len(protocol.Flags()))}
}

func (b *TunnelHandlerBuilder) SYN(handler TunnelHandler) *TunnelHandlerBuilder {
	b.tunnelHandlerByFlag[protocol.FlagSYN] = handler
	return b
}
func (b *TunnelHandlerBuilder) ACK(handler TunnelHandler) *TunnelHandlerBuilder {
	b.tunnelHandlerByFlag[protocol.FlagACK] = handler
	return b
}
func (b *TunnelHandlerBuilder) FIN(handler TunnelHandler) *TunnelHandlerBuilder {
	b.tunnelHandlerByFlag[protocol.FlagFIN] = handler
	return b
}
func (b *TunnelHandlerBuilder) PSH(handler TunnelHandler) *TunnelHandlerBuilder {
	b.tunnelHandlerByFlag[protocol.FlagPSH] = handler
	return b
}
func (b *TunnelHandlerBuilder) KPA(handler TunnelHandler) *TunnelHandlerBuilder {
	b.tunnelHandlerByFlag[protocol.FlagKPA] = handler
	return b
}

func (b *TunnelHandlerBuilder) Build() map[protocol.Flag]TunnelHandler {
	return b.tunnelHandlerByFlag
}
