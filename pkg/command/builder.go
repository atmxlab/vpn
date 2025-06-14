package command

import (
	"io"
	"os"

	"github.com/atmxlab/vpn/pkg/errors"
)

type CommandsBuilder struct {
	stdout io.Writer
	stderr io.Writer

	commands []*Command
}

func NewCommandsBuilder() *CommandsBuilder {
	return &CommandsBuilder{
		stdout: os.Stdout,
		stderr: os.Stderr,
	}
}

func (b *CommandsBuilder) Stdout(stdout io.Writer) *CommandsBuilder {
	b.stdout = stdout
	return b
}

func (b *CommandsBuilder) Stderr(stderr io.Writer) *CommandsBuilder {
	b.stderr = stderr
	return b
}

func (b *CommandsBuilder) Add(fn func(b *Builder)) *CommandsBuilder {
	builder := NewBuilder()
	builder.stdout = b.stdout
	builder.stderr = b.stderr
	fn(builder)

	b.commands = append(b.commands, builder.Build())
	return b
}

func (b *CommandsBuilder) Build() []*Command {
	return b.commands
}

func (b *CommandsBuilder) BuildAndRun() error {
	for _, cmd := range b.Build() {
		if err := cmd.Run(); err != nil {
			return errors.Wrap(err, "failed to execute command")
		}
	}

	return nil
}

type Builder struct {
	stdout io.Writer
	stderr io.Writer

	before func(cmd string)
	name   string
	argv   []string
}

func NewBuilder() *Builder {
	return &Builder{
		stdout: os.Stdout,
		stderr: os.Stderr,
	}
}

func (b *Builder) Stderr(stderr io.Writer) *Builder {
	b.stderr = stderr
	return b
}

func (b *Builder) Stdout(stdout io.Writer) *Builder {
	b.stdout = stdout
	return b
}

func (b *Builder) Before(fn func(cmd string)) *Builder {
	b.before = fn
	return b
}

func (b *Builder) Cmd(name string, argv ...string) *Builder {
	b.name = name
	b.argv = argv
	return b
}

func (b *Builder) Build() *Command {
	return &Command{
		stdout: b.stdout,
		stderr: b.stderr,
		name:   b.name,
		args:   b.argv,
		before: b.before,
	}
}
