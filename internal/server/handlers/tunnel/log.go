package tunnel

import (
	"github.com/atmxlab/vpn/internal/protocol"
	"github.com/sirupsen/logrus"
)

func log(packet *protocol.TunnelPacket) *logrus.Entry {
	return logrus.
		WithField("Namespace", "TUNNEL|HANDLER").
		WithField("Flag", packet.Header().Flag()).
		WithField("Len", packet.Payload().Len()).
		WithField("PacketAddr", packet.Addr())
}
