package main

import (
	"bytes"
	"strconv"

	"github.com/atmxlab/vpn/pkg/command"
	"github.com/atmxlab/vpn/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/songgao/water"
)

func setupTun(mtu uint16) (*water.Interface, error) {
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

	commandBuilder.
		Stdout(stdout).
		Stderr(stderr).
		Add(func(b *command.Builder) {
			b.Before(func(cmd command.Command) error {
				logrus.Infof("Назначаем размер MTU: [%d], для созданного интерфейса: [%s]", mtu, iface.Name())
				logrus.Infof("Run cmd: [%s]", cmd.String())
				return nil
			})
			b.Cmd("ip", "link", "set", "dev", iface.Name(), "mtu", strconv.Itoa(int(mtu)))
		}).
		Add(func(b *command.Builder) {
			b.Before(func(cmd command.Command) error {
				logrus.Infof("Включаем созданный интерфейс")
				logrus.Infof("Run cmd: [%s]", cmd.String())
				return nil
			})
			b.Cmd("ip", "link", "set", "dev", iface.Name(), "up")
		})

	if err = commandBuilder.BuildAndRun(); err != nil {
		logrus.Errorf("Stdout: %s", stdout.String())
		return nil, errors.Wrap(err, "failed to build commands")
	}

	return iface, nil
}
