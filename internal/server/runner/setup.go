package runner

import (
	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (r *Runner) setup() error {
	if err := r.setupTUN(); err != nil {
		return errors.Wrap(err, "setupTUN")
	}

	if err := r.setupTunnel(); err != nil {
		return errors.Wrap(err, "setupTunnel")
	}

	if err := r.setupIPDistributor(); err != nil {
		return errors.Wrap(err, "setupIPDistributor")
	}

	if err := r.configureRoute(); err != nil {
		return errors.Wrap(err, "configureRoute")
	}

	return nil
}

func (r *Runner) setupTUN() error {
	tun, err := r.tunFactory.Create(r.cfg.Tun.Subnet, r.cfg.Tun.MTU)
	if err != nil {
		return errors.Wrap(err, "tunFactory.Create")
	}
	logrus.Debugf("Created TUN interface: name=[%s]", tun.Name())

	r.state.net.tun = tun

	return nil
}

func (r *Runner) setupTunnel() error {
	tunl, err := r.tunnelFactory.Create(r.cfg.ServerAddr)
	if err != nil {
		return errors.Wrap(err, "tunnelFactory.Create")
	}
	logrus.Debugf("Created tunnel: addr=[%s]", tunl.LocalAddr())

	r.state.net.tunnel = tunl

	return nil
}

func (r *Runner) setupIPDistributor() error {
	ipDistributor, err := r.ipDistributorFactory.Create(r.cfg.Tun.Subnet)
	if err != nil {
		return errors.Wrap(err, "ipDistributor.Create")
	}
	logrus.Debug("Created ip distributor")

	r.state.ipDistributor = ipDistributor

	return nil
}

func (r *Runner) configureRoute() error {
	if err := r.routeConfigurator.EnableIPForward(); err != nil {
		return errors.Wrap(err, "routeConfigurator.EnableIPForward")
	}
	logrus.Debug("Route configurator configured IP forwarding")

	if err := r.routeConfigurator.ConfigureFirewall(r.cfg.Tun.Subnet); err != nil {
		return errors.Wrap(err, "routeConfigurator.ConfigureFirewall")
	}
	logrus.Debug("Route configurator configured firewall")

	if err := r.routeConfigurator.SetDefaultRoute(r.cfg.Tun.Subnet); err != nil {
		return errors.Wrap(err, "routeConfigurator.SetDefaultRoute")
	}
	logrus.Debug("Route configurator configured default route")

	return nil
}
