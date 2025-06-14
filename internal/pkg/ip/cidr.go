package ip

import (
	"fmt"
	"net"
)

func BuildCIDR(ip net.IP, mask net.IPMask) string {
	once, _ := mask.Size()
	return fmt.Sprintf("%s/%d", ip.String(), once)
}
