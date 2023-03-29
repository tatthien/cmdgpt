package executils

import (
	"io"
	"os/exec"
)

type Option struct {
	Cmd *exec.Cmd
}

type OptionFns func(*Option)

func WithArgs(args ...string) OptionFns {
	return func(o *Option) {
		o.Cmd.Args = append([]string{o.Cmd.String()}, args...)
	}
}

func WithStdOut(stdout io.Writer) OptionFns {
	return func(o *Option) {
		o.Cmd.Stdout = stdout
	}
}

func Run(cmd string, options ...OptionFns) error {
	c := exec.Command(cmd)

	cmdOption := &Option{
		Cmd: c,
	}

	for _, opt := range options {
		opt(cmdOption)
	}

	if err := c.Run(); err != nil {
		return err
	}

	return nil
}
