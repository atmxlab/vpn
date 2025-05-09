package router

import (
	"github.com/atmxlab/vpn/internal/protocol"
)

type Builder struct {
	cfg                 *config
	tunnel              Tunnel
	tun                 Tun
	tunHandler          TunHandler
	tunnelHandlerByFlag map[protocol.Flag]TunnelHandler
}

func NewBuilder() *Builder {
	return &Builder{
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

func (b *Builder) Build() *Router {
	return &Router{
		tunnel:              b.tunnel,
		tun:                 b.tun,
		tunnelPackets:       make(chan *protocol.TunnelPacket, b.cfg.tunnelChanSize),
		tunPackets:          make(chan *protocol.TunPacket, b.cfg.tunChanSize),
		cfg:                 b.cfg,
		tunHandler:          b.tunHandler,
		tunnelHandlerByFlag: b.tunnelHandlerByFlag,
	}
}

type ConfigBuilder struct {
	bufferSize     uint16
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
		bufferSize:     c.bufferSize,
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
