package stub

type Addr struct {
	network string
	address string
}

func NewAddr(network string, address string) *Addr {
	return &Addr{network: network, address: address}
}

func (a *Addr) Network() string {
	return a.network
}

func (a *Addr) String() string {
	return a.address
}
