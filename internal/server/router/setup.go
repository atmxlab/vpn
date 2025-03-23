package router

import (
	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (r *Router) setup() error {
	if err := r.configureRoute(); err != nil {
		return errors.Wrap(err, "configureRoute")
	}

	return nil
}

// TODO: возможно все-таки стоит от этого избавиться отсюда.
//  Можно сделать это выше

func (r *Router) configureRoute() error {
	if err := r.routeConfigurator.EnableIPForward(); err != nil {
		return errors.Wrap(err, "routeConfigurator.EnableIPForward")
	}
	logrus.Debug("Route configurator configured IP forwarding")

	if err := r.routeConfigurator.ConfigureFirewall(r.cfg.tun.subnet); err != nil {
		return errors.Wrap(err, "routeConfigurator.ConfigureFirewall")
	}
	logrus.Debug("Route configurator configured firewall")

	if err := r.routeConfigurator.SetDefaultRoute(r.cfg.tun.subnet); err != nil {
		return errors.Wrap(err, "routeConfigurator.SetDefaultRoute")
	}
	logrus.Debug("Route configurator configured default route")

	return nil
}
