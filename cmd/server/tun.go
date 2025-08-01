package main

import (
	"bytes"
	"net"
	"strconv"

	"github.com/atmxlab/vpn/internal/pkg/ip"
	"github.com/atmxlab/vpn/pkg/command"
	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/songgao/water"
)

func setupTun(tunIP net.IP, tunSubnet net.IPNet, mtu uint16) (*water.Interface, error) {
	cfg := water.Config{
		DeviceType: water.TUN,
	}

	iface, err := water.New(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create water")
	}

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	commandBuilder := command.NewCommandsBuilder()

	tunCIDR := ip.BuildCIDR(tunIP, tunSubnet.Mask)

	commandBuilder.
		Stdout(stdout).
		Stderr(stderr).
		Add(func(b *command.Builder) {
			b.Before(func(cmd string) {
				logrus.Debug("Set MTU for tun interface")
				logrus.Infof("Run cmd: [%s]", cmd)
			})
			b.Cmd("ip", "link", "set", "dev", iface.Name(), "mtu", strconv.Itoa(int(mtu)))
		}).
		Add(func(b *command.Builder) {
			b.Before(func(cmd string) {
				logrus.Debug("Set IP address for tun interface")
				logrus.Infof("Run cmd: [%s]", cmd)
			})
			b.Cmd("ip", "addr", "add", tunCIDR, "dev", iface.Name())
		}).
		Add(func(b *command.Builder) {
			b.Before(func(cmd string) {
				logrus.Debug("Up tun interface")
				logrus.Infof("Run cmd: [%s]", cmd)
			})
			b.Cmd("ip", "link", "set", "dev", iface.Name(), "up")
		})

	if err = commandBuilder.BuildAndRun(); err != nil {
		logrus.Errorf("Stdout: %s", stdout.String())
		logrus.Errorf("Stderr: %s", stderr.String())
		return nil, errors.Wrap(err, "failed to build commands")
	}

	return iface, nil
}
