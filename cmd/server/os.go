package main

import (
	"bytes"

	"github.com/atmxlab/vpn/pkg/command"
	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/sirupsen/logrus"
)

func setupOS() error {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	commandBuilder := command.NewCommandsBuilder()

	commandBuilder.
		Stdout(stdout).
		Stderr(stderr).
		Add(func(b *command.Builder) {
			b.Before(func(cmd string) {
				logrus.Debug("Enable IP forwarding")
				logrus.Infof("Run cmd: [%s]", cmd)
			})
			b.Cmd("sysctl", "-w", "net.ipv4.ip_forward=1")
		}).
		Add(func(b *command.Builder) {
			b.Before(func(cmd string) {
				logrus.Debug("Delete all NAT rules from firewall")
				logrus.Infof("Run cmd: [%s]", cmd)
			})
			b.Cmd("iptables", "-t", "nat", "-F")
		}).
		Add(func(b *command.Builder) {
			b.Before(func(cmd string) {
				logrus.Debug("Enable MASQUERADE for eth0 network interface")
				logrus.Infof("Run cmd: [%s]", cmd)
			})
			b.Cmd("iptables", "-t", "nat", "-A", "POSTROUTING", "-o", "eth0", "-j", "MASQUERADE")
		}).
		Add(func(b *command.Builder) {
			b.Before(func(cmd string) {
				logrus.Debug("Enable MASQUERADE for eth1 network interface")
				logrus.Infof("Run cmd: [%s]", cmd)
			})
			b.Cmd("iptables", "-t", "nat", "-A", "POSTROUTING", "-o", "eth1", "-j", "MASQUERADE")
		}).
		Add(func(b *command.Builder) {
			b.Before(func(cmd string) {
				logrus.Debug("Accept traffic forwarding from tun0 to eth0")
				logrus.Infof("Run cmd: [%s]", cmd)
			})
			b.Cmd("iptables", "-t", "filter", "-A", "FORWARD", "-i", "tun0", "-o", "eth0", "-j", "ACCEPT")
		}).
		Add(func(b *command.Builder) {
			b.Before(func(cmd string) {
				logrus.Debug("Accept traffic forwarding from eth0 to tun0")
				logrus.Infof("Run cmd: [%s]", cmd)
			})
			b.Cmd("iptables", "-t", "filter", "-A", "FORWARD", "-i", "eth0", "-o", "tun0", "-m", "state", "--state", "RELATED,ESTABLISHED", "-j", "ACCEPT")
		}).
		Add(func(b *command.Builder) {
			b.Before(func(cmd string) {
				logrus.Debug("Accept traffic forwarding from tun0 to eth1")
				logrus.Infof("Run cmd: [%s]", cmd)
			})
			b.Cmd("iptables", "-t", "filter", "-A", "FORWARD", "-i", "tun0", "-o", "eth1", "-j", "ACCEPT")
		}).
		Add(func(b *command.Builder) {
			b.Before(func(cmd string) {
				logrus.Debug("Accept traffic forwarding from eth1 to tun0")
				logrus.Infof("Run cmd: [%s]", cmd)
			})
			b.Cmd("iptables", "-t", "filter", "-A", "FORWARD", "-i", "eth1", "-o", "tun0", "-m", "state", "--state", "RELATED,ESTABLISHED", "-j", "ACCEPT")
		})

	if err := commandBuilder.BuildAndRun(); err != nil {
		logrus.Errorf("Stdout: %s", stdout.String())
		logrus.Errorf("Stderr: %s", stderr.String())
		return errors.Wrap(err, "failed to build commands")
	}

	return nil
}
