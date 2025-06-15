package main

import (
	"bytes"
	"net"

	"github.com/atmxlab/vpn/pkg/command"
	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/sirupsen/logrus"
)

func setupOS(serverIP, gatewayIP net.IP) error {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	commandBuilder := command.NewCommandsBuilder()

	commandBuilder.
		Stdout(stdout).
		Stderr(stderr).
		Add(func(b *command.Builder) {
			b.Before(func(cmd string) {
				logrus.Debug("Replace default route")
				logrus.Infof("Run cmd: [%s]", cmd)
			})
			b.Cmd("ip", "route", "replace", "default", "dev", "tun0")
		}).
		Add(func(b *command.Builder) {
			b.Before(func(cmd string) {
				logrus.Debug("Add route for tunnel")
				logrus.Infof("Run cmd: [%s]", cmd)
			})
			b.Cmd("ip", "route", "replace", serverIP.String(), "via", gatewayIP.String(), "dev", "eth0")
		})

	if err := commandBuilder.BuildAndRun(); err != nil {
		logrus.Errorf("Stdout: %s", stdout.String())
		logrus.Errorf("Stderr: %s", stderr.String())
		return errors.Wrap(err, "failed to build commands")
	}

	return nil
}
