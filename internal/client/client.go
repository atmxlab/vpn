package client

import (
	"context"
	"time"

	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/atmxlab/vpn/pkg/signal"
	"golang.org/x/sync/errgroup"
)

type Router interface {
	Run(ctx context.Context) error
}

type KPAAction interface {
	Run(ctx context.Context) error
}

type SynAction interface {
	Run(ctx context.Context) error
}

type Client struct {
	router     Router
	synAction  SynAction
	kpaAction  KPAAction
	connSignal *signal.Signal

	cancel context.CancelFunc
	eg     *errgroup.Group
}

func NewClient(
	router Router,
	synAction SynAction,
	kpaAction KPAAction,
	connSignal *signal.Signal,
) *Client {
	return &Client{
		router:     router,
		synAction:  synAction,
		kpaAction:  kpaAction,
		connSignal: connSignal,
	}
}

func (c *Client) Run(ctx context.Context, connTimeout time.Duration) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	c.cancel = cancel

	c.eg, ctx = errgroup.WithContext(ctx)

	defer func() {
		if err = c.eg.Wait(); err != nil {
			err = errors.Wrap(err, "error group wait")
		}
	}()
	defer cancel()

	c.eg.Go(func() error {
		if err = c.router.Run(ctx); err != nil {
			return errors.Wrap(err, "failed to start router")
		}
		return nil
	})

	if err = c.synAction.Run(ctx); err != nil {
		return errors.Wrap(err, "failed to start syn action")
	}

	if err = c.connSignal.WaitWithTimeout(ctx, connTimeout); err != nil {
		return errors.Wrap(err, "failed to wait for connection signal")
	}

	c.eg.Go(func() error {
		if err = c.kpaAction.Run(ctx); err != nil {
			return errors.Wrap(err, "failed to start kpa action")
		}
		return nil
	})

	return nil
}

func (c *Client) Close() error {
	if c.cancel == nil {
		return nil
	}

	c.cancel()

	if err := c.eg.Wait(); err != nil {
		return errors.Wrap(err, "error group wait")
	}

	return nil
}
