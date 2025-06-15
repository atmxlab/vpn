package tunnel

import (
	"github.com/atmxlab/vpn/internal/protocol"
	"github.com/sirupsen/logrus"
)

func log(packet *protocol.TunnelPacket) *logrus.Entry {
	l := logrus.
		WithField("Namespace", "TUNNEL|HANDLER")

	if packet != nil {
		l.
			WithField("Flag", packet.Header().Flag()).
			WithField("Len", packet.Payload().Len()).
			WithField("PacketAddr", packet.Addr())
	}

	return l
}
