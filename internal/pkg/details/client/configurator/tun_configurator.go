package configurator

import (
	"context"
	"net"
)

type TunConfigurator struct {
}

func NewTunConfigurator() *TunConfigurator {
	return &TunConfigurator{}
}

func (t TunConfigurator) ChangeAddr(ctx context.Context, subnet net.IPNet) error {
	// TODO implement me
	return nil
}
