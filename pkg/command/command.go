package command

import (
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/atmxlab/vpn/pkg/errors"
)

type Command struct {
	stdout io.Writer
	stderr io.Writer
	before func(cmd *Command) error
	name   string
	args   []string
}

func (c *Command) Stdout() io.Writer {
	return c.stdout
}

func (c *Command) Stderr() io.Writer {
	return c.stderr
}

func (c *Command) Name() string {
	return c.name
}

func (c *Command) Args() []string {
	return c.args
}

func (c *Command) String() string {
	return fmt.Sprintf("%s %s", c.name, strings.Join(c.args, " "))
}

func (c *Command) Run() error {
	if c.before != nil {
		if err := c.before(c); err != nil {
			return errors.Wrap(err, "before command")
		}
	}

	cmd := exec.Command(c.name, c.args...)
	cmd.Stdout = c.stdout
	cmd.Stderr = c.stderr
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "failed to run command")
	}

	return nil
}
