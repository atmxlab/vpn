package configurator

import (
	"bytes"
	"net"

	"github.com/atmxlab/vpn/pkg/command"
	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Configurator struct {
	stderr *bytes.Buffer
	stdout *bytes.Buffer
}

func NewConfigurator() *Configurator {
	return &Configurator{
		stderr: bytes.NewBuffer(nil),
		stdout: bytes.NewBuffer(nil),
	}
}

func (c *Configurator) EnableIPForward() error {
	b := c.createCommandBuilder().
		Add(func(b *command.Builder) {
			b.Before(func(cmd string) {
				logrus.Infof("Enable IP forwarding")
				logrus.Infof("Run cmd: [%s]", cmd)
			})
			b.Cmd("sysctl", "-w", "net.ipv4.ip_forward=1")
		})

	if err := c.runCommands(b); err != nil {
		return errors.Wrap(err, "run commands")
	}

	return nil
}

// TODO: делаем маскарадинг для второго интерфейса + можно еще форвард проверить на него

func (c *Configurator) ConfigureFirewall() error {
	b := c.createCommandBuilder().
		Add(func(b *command.Builder) {
			b.Before(func(cmd string) {
				logrus.Infof("Delete all NAT rules from firewall")
				logrus.Infof("Run cmd: [%s]", cmd)
			})
			b.Cmd("iptables", "-t", "nat", "-F")
		}).
		Add(func(b *command.Builder) {
			b.Before(func(cmd string) {
				logrus.Infof("Enable MASQUERADE for eth0 network interface")
				logrus.Infof("Run cmd: [%s]", cmd)
			})
			b.Cmd("iptables", "-t", "nat", "-A", "POSTROUTING", "-o", "eth0", "-j", "MASQUERADE")
		}).
		Add(func(b *command.Builder) {
			b.Before(func(cmd string) {
				logrus.Infof("Enable MASQUERADE for eth1 network interface")
				logrus.Infof("Run cmd: [%s]", cmd)
			})
			b.Cmd("iptables", "-t", "nat", "-A", "POSTROUTING", "-o", "eth1", "-j", "MASQUERADE")
		}).
		Add(func(b *command.Builder) {
			b.Before(func(cmd string) {
				logrus.Infof("Accept traffic forwarding from tun0 to eth0")
				logrus.Infof("Run cmd: [%s]", cmd)
			})
			b.Cmd("iptables", "-t", "filter", "-A", "FORWARD", "-i", "tun0", "-o", "eth0", "-j", "ACCEPT")
		}).
		Add(func(b *command.Builder) {
			b.Before(func(cmd string) {
				logrus.Infof("Accept traffic forwarding from eth0 to tun0")
				logrus.Infof("Run cmd: [%s]", cmd)
			})
			b.Cmd("iptables", "-t", "filter", "-A", "FORWARD", "-i", "eth0", "-o", "tun0", "-m", "state", "--state", "RELATED,ESTABLISHED", "-j", "ACCEPT")
		}).
		Add(func(b *command.Builder) {
			b.Before(func(cmd string) {
				logrus.Infof("Accept traffic forwarding from tun0 to eth1")
				logrus.Infof("Run cmd: [%s]", cmd)
			})
			b.Cmd("iptables", "-t", "filter", "-A", "FORWARD", "-i", "tun0", "-o", "eth1", "-j", "ACCEPT")
		}).
		Add(func(b *command.Builder) {
			b.Before(func(cmd string) {
				logrus.Infof("Accept traffic forwarding from eth1 to tun0")
				logrus.Infof("Run cmd: [%s]", cmd)
			})
			b.Cmd("iptables", "-t", "filter", "-A", "FORWARD", "-i", "eth1", "-o", "tun0", "-m", "state", "--state", "RELATED,ESTABLISHED", "-j", "ACCEPT")
		})

	if err := c.runCommands(b); err != nil {
		return errors.Wrap(err, "run commands")
	}

	return nil
}

func (c *Configurator) SetDefaultRoute(subnet net.IPNet) error {
	return nil
}

func (c *Configurator) createCommandBuilder() *command.CommandsBuilder {
	commandBuilder := command.NewCommandsBuilder()

	return commandBuilder.
		Stdout(c.stdout).
		Stderr(c.stderr)
}

func (c *Configurator) runCommands(builder *command.CommandsBuilder) error {
	if err := builder.BuildAndRun(); err != nil {
		logrus.Errorf("Stdout: %s", c.stdout.String())
		logrus.Errorf("Stderr: %s", c.stderr.String())
		return errors.Wrap(err, "failed to build commands")
	}

	return nil
}
