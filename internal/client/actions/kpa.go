package actions

import (
	"net"

	"github.com/atmxlab/vpn/pkg/errors"
)

type KPAAction struct {
	tunnel     Tunnel
	serverAddr net.Addr
}

func NewKPAAction(tunnel Tunnel, serverAddr net.Addr) *KPAAction {
	return &KPAAction{tunnel: tunnel, serverAddr: serverAddr}
}

func (a *KPAAction) Run() error {
	_, err := a.tunnel.KPA(a.serverAddr, nil)
	if err != nil {
		return errors.Wrap(err, "tunnel.KPA")
	}

	return nil
}
