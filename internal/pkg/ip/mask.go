package ip

import "net"

func MaskFromIP(ip net.IP) net.IPMask {
	return net.IPv4Mask(ip[0], ip[1], ip[2], ip[3])
}
