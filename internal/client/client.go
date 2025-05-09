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

	cancel context.CancelFunc
	eg     *errgroup.Group
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
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	c.cancel = cancel

	c.eg, ctx = errgroup.WithContext(ctx)

	defer func() {
		if err != nil {
			cancel()
			logrus.Warnf("Client app error: %v", err)
		}

		logrus.Info("Waiting app closing...")
		if egErr := c.eg.Wait(); egErr != nil {
			err = errors.Join(err, errors.Wrap(egErr, "error group wait"))
		}
	}()

	c.eg.Go(func() error {
		logrus.Debug("Starting router")
		defer logrus.Warn("Stopped router")
		if err = c.router.Run(ctx); err != nil {
			return errors.Wrap(err, "failed to start router")
		}
		return nil
	})

	logrus.Info("Init connection")
	if err = c.synAction.Run(ctx); err != nil {
		return errors.Wrap(err, "failed to start syn action")
	}

	logrus.Info("Waiting connection signal...")
	if err = c.connSignal.WaitWithTimeout(connTimeout); err != nil {
		return errors.Wrap(err, "failed to wait for connection signal")
	}
	logrus.Info("Connected to server!")

	c.eg.Go(func() error {
		logrus.Debug("Starting KPA action")
		defer logrus.Warn("stopped KPA action")
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
