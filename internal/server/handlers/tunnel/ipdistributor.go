package tunnel

import "net"

// IpDistributor - распределитель IP адресов
// Из выделенной подсети выделяет и освобождает IP адреса
type IpDistributor interface {
	AcquireIP() (net.IP, error)
	ReleaseIP(net.IP) error
}
