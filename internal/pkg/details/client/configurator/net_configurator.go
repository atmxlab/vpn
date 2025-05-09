package configurator

import (
	"context"
	"net"
)

type NetConfigurator struct {
}

func NewNetConfigurator() *NetConfigurator {
	return &NetConfigurator{}
}

func (c NetConfigurator) ConfigureRouting(ctx context.Context, subnet net.IPNet) error {
	// TODO implement me
	return nil
}
