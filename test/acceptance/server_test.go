package acceptance

import (
	"net"
	"testing"
	"time"

	"github.com/atmxlab/vpn/internal/protocol"
	"github.com/atmxlab/vpn/test/engine"
	"github.com/atmxlab/vpn/test/gen"
	"github.com/atmxlab/vpn/test/stub"
)

func TestServer(t *testing.T) {
	t.Parallel()

	t.Run("syn single client", func(t *testing.T) {
		t.Parallel()

		client := gen.RandAddr()
		eng := engine.New(t)

		eng.REPLAY(
			engine.SYN(client),
			engine.CHECKPOINT(
				engine.ExpectPeer(client),
				engine.ExpectEmptyTun(),
				engine.ExpectTunnelACK(client),
			),
		)
	})

	t.Run("syn multi client", func(t *testing.T) {
		t.Parallel()

		firstClient := gen.RandAddr()
		secondClient := gen.RandAddr()
		eng := engine.New(t)

		eng.REPLAY(
			engine.SYN(firstClient),
			engine.CHECKPOINT(
				engine.ExpectPeer(firstClient),
				engine.ExpectEmptyTun(),
				engine.ExpectTunnelACK(firstClient),
			),

			engine.SYN(secondClient),
			engine.CHECKPOINT(
				engine.ExpectPeer(firstClient),
				engine.ExpectPeer(secondClient),
				engine.ExpectEmptyTun(),
				engine.ExpectTunnelACK(secondClient),
			),
		)
	})

	t.Run("keepalive deadline", func(t *testing.T) {
		t.Parallel()

		client := gen.RandAddr()

		eng := engine.New(
			t,
			engine.WithActionDelay(10*time.Millisecond),
			engine.WithPeerKeepAliveTTL(100*time.Millisecond),
		)

		eng.REPLAY(
			engine.SYN(client),
			engine.CHECKPOINT(
				engine.ExpectPeer(client),
				engine.ExpectBusyDedicatedIP(),
			),
			engine.WAIT(100*time.Millisecond),
			engine.CHECKPOINT(
				engine.UnexpectPeer(client),
				engine.ExpectFreeAllDedicatedIPs(),
			),
		)
	})

	t.Run("fin", func(t *testing.T) {
		t.Parallel()

		client := gen.RandAddr()
		eng := engine.New(t)

		eng.REPLAY(
			engine.SYN(client),
			engine.CHECKPOINT(
				engine.ExpectPeer(client),
				engine.ExpectEmptyTun(),
				engine.ExpectTunnelACK(client),
			),
			engine.FIN(client),
			engine.CHECKPOINT(
				engine.UnexpectPeer(client),
				engine.ExpectEmptyTun(),
			),
		)
	})

	t.Run("push data", func(t *testing.T) {
		t.Parallel()

		client := gen.RandAddr()
		payload := gen.RandPayload()
		eng := engine.New(t)

		eng.REPLAY(
			engine.SYN(client),
			engine.CHECKPOINT(
				engine.ExpectPeer(client),
				engine.ExpectEmptyTun(),
				engine.ExpectTunnelACK(client),
			),
			engine.PSH(client, payload),
			engine.CHECKPOINT(
				engine.ExpectTun(protocol.NewTunPacket(payload)),
			),
		)
	})

	t.Run("receive data", func(t *testing.T) {
		t.Parallel()

		client := gen.RandAddr()
		tunPacket := gen.RandTunICMPReq(t, func(h *stub.IPHeader) {
			h.Dst = net.IPv4(10, 0, 0, 0)
		})
		eng := engine.New(t)

		eng.REPLAY(
			engine.SYN(client),
			engine.CHECKPOINT(
				engine.ExpectPeer(client),
				engine.ExpectEmptyTun(),
				engine.ExpectTunnelACK(client),
			),
			engine.TUN(tunPacket),
			engine.CHECKPOINT(
				engine.ExpectTunnelPSH(client, tunPacket.Payload()),
			),
		)
	})
}
