package ipdistributor

import (
	"net"
	"sync"

	"github.com/atmxlab/vpn/internal/pkg/ip"
	"github.com/atmxlab/vpn/pkg/errors"
)

type Distributor struct {
	subnet       net.IPNet
	m            sync.Mutex
	ipPool       map[string]bool
	sortedIPPool []string
}

func New(subnet net.IPNet) (*Distributor, error) {
	d := &Distributor{
		subnet:       subnet,
		ipPool:       make(map[string]bool, ip.CountInMask(subnet.Mask)),
		sortedIPPool: make([]string, 0, ip.CountInMask(subnet.Mask)),
	}

	d.generateIpPool()

	return d, nil
}

func (ipd *Distributor) AcquireIP() (net.IP, error) {
	ipd.m.Lock()
	defer ipd.m.Unlock()

	for _, ipp := range ipd.sortedIPPool {
		isBusy, _ := ipd.ipPool[ipp]
		if !isBusy {
			ipd.ipPool[ipp] = true
			return net.ParseIP(ipp), nil
		}
	}

	return nil, errors.New("IP pool is empty")
}

func (ipd *Distributor) ReleaseIP(ip net.IP) error {
	ipd.m.Lock()
	defer ipd.m.Unlock()

	if _, ok := ipd.ipPool[ip.String()]; !ok {
		return errors.New("ip not found")
	}

	ipd.ipPool[ip.String()] = false

	return nil
}

func (ipd *Distributor) generateIpPool() {
	ipd.m.Lock()
	defer ipd.m.Unlock()

	for ipp := ipd.subnet.IP.Mask(ipd.subnet.Mask); ipd.subnet.Contains(ipp); ipd.incIP(ipp) {
		ipd.sortedIPPool = append(ipd.sortedIPPool, ipp.String())
		ipd.ipPool[ipp.String()] = false
	}
}

func (ipd *Distributor) incIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func (ipd *Distributor) HasBusy() bool {
	ipd.m.Lock()
	defer ipd.m.Unlock()

	for _, isBusy := range ipd.ipPool {
		if isBusy {
			return true
		}
	}

	return false
}
