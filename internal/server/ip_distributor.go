package server

import "net"

type IpDistributorFactory interface {
	Create(ipNet net.IPNet) (IpDistributor, error)
}

type IpDistributor interface {
	AcquireIP() (net.IP, error)
	ReleaseIP(net.IP) error
}
