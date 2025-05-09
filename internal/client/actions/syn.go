package actions

import (
	"context"
	"net"

	"github.com/atmxlab/vpn/pkg/errors"
)

type SYNAction struct {
	tunnel     Tunnel
	serverAddr net.Addr
}

func NewSYNAction(tunnel Tunnel, serverAddr net.Addr) *SYNAction {
	return &SYNAction{tunnel: tunnel, serverAddr: serverAddr}
}

func (a *SYNAction) Run(_ context.Context) error {
	_, err := a.tunnel.SYN(a.serverAddr, nil)
	if err != nil {
		return errors.Wrap(err, "tunnel.SYN")
	}

	return nil
}
