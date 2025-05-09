package configurator

import (
	"net"
)

// TODO: impl

type Configurator struct {
}

func NewConfigurator() *Configurator {
	return &Configurator{}
}

func (c Configurator) EnableIPForward() error {
	return nil
}

func (c Configurator) ConfigureFirewall(subnet net.IPNet) error {
	return nil
}

func (c Configurator) SetDefaultRoute(subnet net.IPNet) error {
	return nil
}
