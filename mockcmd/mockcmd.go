package mockcmd

import (
	"bytes"
	"context"
	"github.com/cirrusdata/cdsexec"
	"io"
	"os"
	"os/exec"
)

// MockCmd is a simplified mock implementation of the Commander interface.
// It allows for predefined outputs and command construction checking.
type MockCmd struct {
	// Predefined outputs and error
	Stdout []byte
	Stderr []byte
	Err    error

	// Command construction details
	Ctx  context.Context
	Name string
	Args []string
	Dir  string
	Env  []string

	// Function to check if the command was constructed correctly
	CheckFunc func(*MockCmd) error

	// Flags to track method calls
	startCalled bool
	waitCalled  bool
}

// mockCommandContext creates a new MockCmd with the given context, name, and arguments.
func mockCommandContext(ctx context.Context, name string, arg ...string) *MockCmd {
	return &MockCmd{
		Ctx:  ctx,
		Name: name,
		Args: arg,
	}
}

// Run simulates running the command and returns any predefined error.
// It also executes the CheckFunc if defined.
func (m *MockCmd) Run() error {
	if m.CheckFunc != nil {
		if err := m.CheckFunc(m); err != nil {
			return err
		}
	}
	return m.Err
}

// Output returns the predefined stdout and any error.
// It also executes the CheckFunc if defined.
func (m *MockCmd) Output() ([]byte, error) {
	if m.CheckFunc != nil {
		if err := m.CheckFunc(m); err != nil {
			return nil, err
		}
	}

	return m.Stdout, m.Err
}

// CombinedOutput returns the combined predefined stdout and stderr, and any error.
// It also executes the CheckFunc if defined.
func (m *MockCmd) CombinedOutput() ([]byte, error) {
	if m.CheckFunc != nil {
		if err := m.CheckFunc(m); err != nil {
			return nil, err
		}
	}
	return append(m.Stdout, m.Stderr...), m.Err
}

// Start simulates starting the command and marks it as started.
// It executes the CheckFunc if defined.
func (m *MockCmd) Start() error {
	m.startCalled = true
	if m.CheckFunc != nil {
		return m.CheckFunc(m)
	}
	return m.Err
}

// Wait simulates waiting for the command to complete and marks it as waited.
func (m *MockCmd) Wait() error {
	m.waitCalled = true
	return m.Err
}

// StdinPipe returns a mock WriteCloser for stdin.
func (m *MockCmd) StdinPipe() (io.WriteCloser, error) {
	return &mockWriteCloser{}, nil
}

// StdoutPipe returns a ReadCloser with the predefined stdout.
func (m *MockCmd) StdoutPipe() (io.ReadCloser, error) {
	return io.NopCloser(bytes.NewBuffer(m.Stdout)), nil
}

// StderrPipe returns a ReadCloser with the predefined stderr.
func (m *MockCmd) StderrPipe() (io.ReadCloser, error) {
	return io.NopCloser(bytes.NewBuffer(m.Stderr)), nil
}

// SetDir sets the working directory for the mock command.
func (m *MockCmd) SetDir(dir string) {
	m.Dir = dir
}

// SetEnv sets the environment variables for the mock command.
func (m *MockCmd) SetEnv(env []string) {
	m.Env = env
}

// SetStdin, SetStdout, and SetStderr are no-op implementations to satisfy the interface.

func (m *MockCmd) SetStdin(in io.Reader)   {}
func (m *MockCmd) SetStdout(out io.Writer) {}
func (m *MockCmd) SetStderr(out io.Writer) {}

// Process and ProcessState return nil to satisfy the interface.
func (m *MockCmd) Process() *os.Process           { return nil }
func (m *MockCmd) ProcessState() *os.ProcessState { return nil }

// mockWriteCloser is a simple implementation of io.WriteCloser.
type mockWriteCloser struct {
	bytes.Buffer
}

// Close is a no-op implementation to satisfy the io.Closer interface.
func (mwc *mockWriteCloser) Close() error {
	return nil
}

func MakeMockCmdWithOutput(fixedOutput string, checkFunc func(*MockCmd) error) cdsexec.CommandConstructor {
	return func(ctx context.Context, name string, arg ...string) cdsexec.Commander {
		c := mockCommandContext(ctx, name, arg...)
		c.Stdout = []byte(fixedOutput)
		c.CheckFunc = checkFunc
		return c
	}
}

func MakeMockCmdWithOutputGenericError(checkFunc func(*MockCmd) error) cdsexec.CommandConstructor {
	return func(ctx context.Context, name string, arg ...string) cdsexec.Commander {
		c := mockCommandContext(ctx, name, arg...)
		c.CheckFunc = checkFunc
		c.Err = exec.ErrNotFound
		return c
	}
}

func MakeMockCmdWithOutputSpecificError(fixedOutput string, specificError error, checkFunc func(*MockCmd) error) cdsexec.CommandConstructor {
	return func(ctx context.Context, name string, arg ...string) cdsexec.Commander {
		c := mockCommandContext(ctx, name, arg...)
		c.Stdout = []byte(fixedOutput)
		c.CheckFunc = checkFunc
		c.Err = specificError
		return c
	}
}

func MakeMockCmd(c *MockCmd) cdsexec.CommandConstructor {
	return func(ctx context.Context, name string, arg ...string) cdsexec.Commander {
		return c
	}
}
