package actions

import (
	"context"
	"net"
	"time"

	"github.com/atmxlab/vpn/pkg/errors"
)

type KPAAction struct {
	tunnel     Tunnel
	serverAddr net.Addr
	tick       time.Duration
}

func NewKPAAction(tunnel Tunnel, serverAddr net.Addr, tick time.Duration) *KPAAction {
	return &KPAAction{tunnel: tunnel, serverAddr: serverAddr, tick: tick}
}

func (a *KPAAction) Run(ctx context.Context) error {
	ticker := time.NewTicker(a.tick)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			_, err := a.tunnel.KPA(a.serverAddr, nil)
			if err != nil {
				return errors.Wrap(err, "tunnel.KPA")
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
