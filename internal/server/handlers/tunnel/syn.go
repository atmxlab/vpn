package tunnel

import (
	"context"
	"net"

	"github.com/atmxlab/vpn/internal/protocol"
	"github.com/atmxlab/vpn/internal/server"
	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/sirupsen/logrus"
)

// IpDistributor - распределитель IP адресов
// Из выделенной подсети выделяет и освобождает IP адреса
type IpDistributor interface {
	AcquireIP() (net.IP, error)
	ReleaseIP(net.IP) error
}

type SynHandler struct {
	tunnel        Tunnel
	peerManager   server.PeerManager
	ipDistributor IpDistributor
}

func (h *SynHandler) Handle(ctx context.Context, packet *protocol.TunnelPacket) error {
	_, exists, err := h.peerManager.FindByAddr(ctx, packet.Addr())
	if err != nil {
		return errors.Wrap(err, "peerManager.FindByAddr")
	}
	if exists {
		return errors.Wrap(errors.ErrNotFound, "peerManager.FindByAddr not found")
	}

	acquiredIP, err := h.ipDistributor.AcquireIP()
	if err != nil {
		return errors.Wrap(err, "ipDistributor.AcquireIP")
	}

	peer := server.NewPeer(acquiredIP, packet.Addr())

	logrus.Infof("Created new peer with addr: %s and dedicated ip: %s", peer.Addr(), peer.DedicatedIP())

	err = h.peerManager.Add(ctx, peer)
	if err != nil {
		return errors.Wrap(err, "peerManager.Add")
	}

	_, err = h.tunnel.ACK(peer.Addr(), peer.DedicatedIP().To4())
	if err != nil {
		return errors.Wrap(err, "tunnel.ACK")
	}

	return nil
}
