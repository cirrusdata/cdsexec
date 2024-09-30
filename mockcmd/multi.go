package mockcmd

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/cirrusdata/cdsexec"
)

var ErrNoMatchingCommand = errors.New("no matching command found in this mock")

// CommandConfig represents a single command configuration
type CommandConfig struct {
	Name   string
	Args   []string
	Stdout []byte
	Stderr []byte
	Err    error
}

// MultiCmdMockCmd is a mock that can handle multiple command configurations
type MultiCmdMockCmd struct {
	MockCmd
	configs        []CommandConfig
	lastMatchedCmd *CommandConfig
}

// matchCommand checks if the given command matches any of the configured commands
func (m *MultiCmdMockCmd) matchCommand() error {
	m.lastMatchedCmd = nil
	for _, config := range m.configs {
		if m.Name == config.Name && reflect.DeepEqual(m.Args, config.Args) {
			m.Stdout = config.Stdout
			m.Stderr = config.Stderr
			m.Err = config.Err
			m.lastMatchedCmd = &config
			return nil
		}
	}
	m.Stderr = nil
	m.Err = ErrNoMatchingCommand
	return nil
}

// Run implements the Commander interface
func (m *MultiCmdMockCmd) Run() error {
	if err := m.matchCommand(); err != nil {
		return err
	}
	return m.Err
}

// Output implements the Commander interface
func (m *MultiCmdMockCmd) Output() ([]byte, error) {
	if err := m.matchCommand(); err != nil {
		return nil, err
	}
	return m.Stdout, m.Err
}

// CombinedOutput implements the Commander interface
func (m *MultiCmdMockCmd) CombinedOutput() ([]byte, error) {
	if err := m.matchCommand(); err != nil {
		return nil, err
	}
	return append(m.Stdout, m.Stderr...), m.Err
}

// String returns a string representation of the last matched command
func (m *MultiCmdMockCmd) String() string {
	if m.lastMatchedCmd == nil {
		return fmt.Sprintf("No matching command found for: %s %s", m.Name, strings.Join(m.Args, " "))
	}
	return fmt.Sprintf("Matched command: %s %s", m.lastMatchedCmd.Name, strings.Join(m.lastMatchedCmd.Args, " "))
}

// MultiCmdMock creates a CommandConstructor that returns a MultiCmdMockCmd
func MultiCmdMock(configs ...CommandConfig) cdsexec.CommandConstructor {
	return func(ctx context.Context, name string, arg ...string) cdsexec.Commander {
		cmd := &MultiCmdMockCmd{
			configs: configs,
		}
		cmd.Ctx = ctx
		cmd.Name = name
		cmd.Args = arg
		return cmd
	}
}
