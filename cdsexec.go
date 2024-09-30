package cdsexec

import (
	"context"
	"io"
	"os"
	"os/exec"
)

// Commander is an interface that abstracts the exec.Cmd functionality.
type Commander interface {
	CommandRunner
	SetDir(dir string)
	SetEnv(env []string)
	SetStdin(in io.Reader)
	SetStdout(out io.Writer)
	SetStderr(out io.Writer)
	Process() *os.Process
	ProcessState() *os.ProcessState
}

// CommandRunner is an interface that abstracts the exec.Cmd functionality.
type CommandRunner interface {
	Run() error
	Output() ([]byte, error)
	CombinedOutput() ([]byte, error)
	Start() error
	Wait() error
	StdinPipe() (io.WriteCloser, error)
	StdoutPipe() (io.ReadCloser, error)
	StderrPipe() (io.ReadCloser, error)
}

// make sure that this interface satisfies exec.Cmd
var _ CommandRunner = (*exec.Cmd)(nil)

// CommandConstructor is a function that creates a new Commander instance. os/exec.CommandContext is the default implementation.
// only command with context is supported.
// this constructor allows module to replace the default implementation with their own implementation.
type CommandConstructor func(ctx context.Context, name string, arg ...string) Commander
