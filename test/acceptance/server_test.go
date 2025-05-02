package acceptance

import (
	"testing"
	"time"

	"github.com/atmxlab/vpn/test/engine"
	"github.com/atmxlab/vpn/test/gen"
)

// TODO: так как некоторые вещи завязаны на время
//  можно использовать новую фичу гошки с ускорением времени

// func TestServer(t *testing.T) {
// 	t.Parallel()
//
// 	serverIP := gen.RandIP()
// 	clientIP := gen.RandIP()
// 	targetIP := gen.RandIP()
// 	// TODO: движок должен использовать реальные реализации
// 	//  подменять реализации нужно только если есть сайд эффекты
// 	eng := engine.NewServer(t, serverIP)
//
// 	eng.Replay(
// 		// Wrote SYN packet into tunnel
// 		engine.SYN(clientIP),
// 		// Waiting 10 mls before running next action
// 		engine.WAIT(10*time.Millisecond),
// 		// Checkpoint
// 		engine.ASSERT(
// 			eng.ExpectPeer(clientIP),
// 			eng.UnexpectTun(),
// 			eng.ExpectTunnelACK(),
// 			func() {
// 				// custom checks
// 			},
// 		),
// 		engine.PSH(engine.ICMPReq(func(packet *packets.ICMP) {
// 			// change packet
// 			packet.SrcIP = eng.DedicatedIP(clientIP)
// 			packet.DstIP = targetIP
// 		})),
// 		// Checkpoint
// 		engine.ASSERT(
// 			eng.ExpectPeer(clientIP),
// 			eng.ExpectTun(),
// 			eng.ExpectTunnelACK(),
//
// 			// rename
// 			eng.ExpectTunEqualsTunnel(),
// 		),
// 		engine.TUN(engine.ICMPResp(func(packet *packets.ICMP) {
// 			// change packet
// 			packet.SrcIP = targetIP
// 			packet.DstIP = clientIP // after masquerading must be client IP
// 		})),
// 		// Checkpoint
// 		engine.ASSERT(
// 			eng.ExpectPeer(clientIP),
// 			eng.ExpectTun(),
// 			eng.ExpectTunnelPSH(),
//
// 			// rename
// 			eng.ExpectTunEqualsTunnel(),
// 		),
// 	)
// }

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
		t.Skip("waiting developed keepalive logic")
		t.Parallel()

		client := gen.RandAddr()

		eng := engine.New(
			t,
			engine.WithPeerKeepAliveTTL(20*time.Millisecond),
		)

		eng.REPLAY(
			engine.SYN(client),
			engine.CHECKPOINT(
				engine.ExpectPeer(client),
			),
			engine.WAIT(40*time.Millisecond),
			engine.CHECKPOINT(
				engine.UnexpectPeer(client),
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
}
