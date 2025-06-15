package client

import (
	"context"
	"time"

	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

//go:generate mock Router
type Router interface {
	Run(ctx context.Context) error
}

//go:generate mock KPAAction
type KPAAction interface {
	Run(ctx context.Context) error
}

//go:generate mock SynAction
type SynAction interface {
	Run(ctx context.Context) error
}

//go:generate mock Signaller
type Signaller interface {
	WaitWithTimeout(timeout time.Duration) error
}

type Client struct {
	router     Router
	synAction  SynAction
	kpaAction  KPAAction
	connSignal Signaller

	eg *errgroup.Group
}

func NewClient(
	router Router,
	synAction SynAction,
	kpaAction KPAAction,
	connSignal Signaller,
) *Client {
	return &Client{
		router:     router,
		synAction:  synAction,
		kpaAction:  kpaAction,
		connSignal: connSignal,
	}
}

func (c *Client) Run(ctx context.Context, connTimeout time.Duration) (err error) {
	l := logrus.WithField("Namespace", "CLIENT")

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	c.eg, ctx = errgroup.WithContext(ctx)

	defer func() {
		if err != nil {
			cancel()
			l.
				WithError(err).
				Error("Client app error")
		}

		l.Debug("Waiting app closing...")
		if egErr := c.eg.Wait(); egErr != nil {
			err = errors.Join(err, errors.Wrap(egErr, "error group wait"))
			if !errors.Is(err, context.Canceled) {
				l.
					WithError(egErr).
					Error("Client app error after wait")
			}
		}
	}()

	c.eg.Go(func() error {
		defer l.Debug("Stopped router")

		l.Debug("Starting router...")
		if err = c.router.Run(ctx); err != nil {
			return errors.Wrap(err, "failed to start router")
		}
		return nil
	})

	l.Debug("Init connection")
	if err = c.synAction.Run(ctx); err != nil {
		return errors.Wrap(err, "failed to start syn action")
	}

	l.Debug("Waiting connection signal...")
	if err = c.connSignal.WaitWithTimeout(connTimeout); err != nil {
		return errors.Wrap(err, "failed to wait for connection signal")
	}
	l.Debug("Connected to server!")

	c.eg.Go(func() error {
		defer l.Debug("Stopped KPA action")

		l.Debug("Starting KPA action...")
		if err = c.kpaAction.Run(ctx); err != nil {
			return errors.Wrap(err, "failed to start kpa action")
		}
		return nil
	})

	return nil
}
