package configurator

import (
	"bytes"
	"context"
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

func (c *Configurator) ChangeTunAddr(_ context.Context, subnet net.IPNet) error {
	b := c.createCommandBuilder().
		Add(func(b *command.Builder) {
			b.Before(func(cmd string) {
				logrus.Debug("Add ip addr for tun0 interface")
				logrus.Infof("Run cmd: [%s]", cmd)
			})
			b.Cmd("ip", "addr", "add", subnet.String(), "dev", "tun0")
		})

	if err := c.runCommands(b); err != nil {
		return errors.Wrap(err, "run commands")
	}

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
