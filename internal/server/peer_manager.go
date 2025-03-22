package server

import "net"

// PeerManager - управляющий пирами
type PeerManager interface {
	Add(peer *Peer) error
	Remove(peer *Peer) error
	FindByAddr(ip net.IP) (peer *Peer, exists bool, err error)
	FindByDedicatedIP(ip net.IP) (peer *Peer, exists bool, err error)
}
