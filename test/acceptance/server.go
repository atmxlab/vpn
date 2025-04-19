package acceptance

import (
	"testing"
	"time"

	"github.com/atmxlab/vpn/test/gen"
)

// TODO: так как некоторые вещи завязаны на время
//  можно использовать новую фичу гошки с ускорением времени

func TestServer(t *testing.T) {
	t.Parallel()

	serverIP := gen.RandIP()
	clientIP := gen.RandIP()
	targetIP := gen.RandIP()
	// TODO: движок должен использовать реальные реализации
	//  подменять реализации нужно только если есть сайд эффекты
	eng := engine.NewServer(t, serverIP)

	eng.Replay(
		// Wrote SYN packet into tunnel
		engine.SYN(clientIP),
		// Waiting 10 mls before running next action
		engine.WAIT(10*time.Millisecond),
		// Checkpoint
		engine.ASSERT(
			eng.ExpectPeer(clientIP),
			eng.UnexpectTun(),
			eng.ExpectTunnelACK(),
			func() {
				// custom checks
			},
		),
		engine.PSH(engine.ICMPReq(func(packet *packets.ICMP) {
			// change packet
			packet.SrcIP = eng.DedicatedIP(clientIP)
			packet.DstIP = targetIP
		})),
		// Checkpoint
		engine.ASSERT(
			eng.ExpectPeer(clientIP),
			eng.ExpectTun(),
			eng.ExpectTunnelACK(),

			// rename
			eng.ExpectTunEqualsTunnel(),
		),
		engine.TUN(engine.ICMPResp(func(packet *packets.ICMP) {
			// change packet
			packet.SrcIP = targetIP
			packet.DstIP = clientIP // after masquerading must be client IP
		})),
		// Checkpoint
		engine.ASSERT(
			eng.ExpectPeer(clientIP),
			eng.ExpectTun(),
			eng.ExpectTunnelPSH(),

			// rename
			eng.ExpectTunEqualsTunnel(),
		),
	)
}
