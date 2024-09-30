package cdsexec

import (
	"context"
	"io"
	"os"
	"os/exec"
)

var _ Commander = (*Cmd)(nil)

func CommandContext(ctx context.Context, name string, arg ...string) Commander {
	return &Cmd{
		Cmd: exec.CommandContext(ctx, name, arg...),
	}
}

// Cmd is a wrapper around exec.Cmd.
type Cmd struct {
	*exec.Cmd
}

// SetDir sets the working directory of the command.
func (c *Cmd) SetDir(dir string) {
	c.Cmd.Dir = dir
}

// SetEnv sets the environment variables for the command.
func (c *Cmd) SetEnv(env []string) {
	c.Cmd.Env = env
}

// SetStdin sets the standard input for the command.
func (c *Cmd) SetStdin(in io.Reader) {
	c.Cmd.Stdin = in
}

// SetStdout sets the standard output for the command.
func (c *Cmd) SetStdout(out io.Writer) {
	c.Cmd.Stdout = out
}

// SetStderr sets the standard error for the command.
func (c *Cmd) SetStderr(out io.Writer) {
	c.Cmd.Stderr = out
}

// Process returns the process.
func (c *Cmd) Process() *os.Process {
	return c.Cmd.Process
}

// ProcessState returns the process state.
func (c *Cmd) ProcessState() *os.ProcessState {
	return c.Cmd.ProcessState
}
