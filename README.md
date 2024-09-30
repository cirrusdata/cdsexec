# cdsexec

`cdsexec` is a Go package that provides a flexible abstraction over the standard `os/exec` package, allowing for easier testing and mocking of command execution.

## Features

- Abstracts `exec.Cmd` functionality through interfaces
- Provides a real implementation that wraps `exec.Cmd`
- Includes mock implementations for testing in a separate `mockcmd` subpackage:
    - Single command mock for simple scenarios
    - Multi-command mock for complex testing situations
- Supports context-based command creation

## Installation

To install `cdsexec`, use `go get`:

```bash
go get github.com/cirrusdata/cdsexec
```

## Usage

### Real Command Execution

To use the real command execution in your code:

```go
import "github.com/cirrusdata/cdsexec"

func main() {
    ctx := context.Background()
    cmd := cdsexec.CommandContext(ctx, "ls", "-l")
    output, err := cmd.Output()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(output))
}
```

### Mocking in Tests

The `mockcmd` subpackage provides two types of mocks: single command mock and multi-command mock.

#### Single Command Mock

For simple scenarios where you need to mock a single command:

```go
import (
    "context"
    "testing"

    "github.com/cirrusdata/cdsexec"
    "github.com/cirrusdata/cdsexec/mockcmd"
)

func TestWithSingleCommandMock(t *testing.T) {
    mockCommandContext := mockcmd.MakeMockCmdWithOutput("mocked output", func(cmd *mockcmd.MockCmd) error {
        if cmd.Name != "expected-command" {
            return fmt.Errorf("unexpected command: %s", cmd.Name)
        }
        return nil
    })

    cmd := mockCommandContext(context.Background(), "expected-command", "-arg1", "-arg2")
    
    output, err := cmd.Output()
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }
    if string(output) != "mocked output" {
        t.Errorf("Unexpected output: %s", string(output))
    }
}
```

#### Multi-Command Mock

For more complex scenarios where you need to mock multiple commands:

```go
import (
    "context"
    "testing"

    "github.com/cirrusdata/cdsexec"
    "github.com/cirrusdata/cdsexec/mockcmd"
)

func TestWithMultiCommandMock(t *testing.T) {
    configs := []mockcmd.CommandConfig{
        {
            Name:   "ls",
            Args:   []string{"-l"},
            Stdout: []byte("file1\nfile2\n"),
        },
        {
            Name:   "cat",
            Args:   []string{"file1"},
            Stdout: []byte("contents of file1"),
        },
    }

    mockCommandContext := mockcmd.MultiCmdMock(configs...)

    // Test matched command
    cmd := mockCommandContext(context.Background(), "ls", "-l")
    output, err := cmd.Output()
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }
    if string(output) != "file1\nfile2\n" {
        t.Errorf("Unexpected output: %s", string(output))
    }

    // Test unmatched command
    cmd = mockCommandContext(context.Background(), "unknown", "command")
    _, err = cmd.Output()
    if err != mockcmd.ErrNoMatchingCommand {
        t.Errorf("Expected ErrNoMatchingCommand, got: %v", err)
    }
}
```

## Mock Features

### Single Command Mock

The `mockcmd` subpackage provides several utility functions for creating single command mocks:

- `MakeMockCmdWithOutput`: Creates a mock command with fixed output
- `MakeMockCmdWithOutputGenericError`: Creates a mock command that returns a generic error
- `MakeMockCmdWithOutputSpecificError`: Creates a mock command with fixed output and a specific error
- `MakeMockCmd`: Creates a mock command from a pre-configured `MockCmd` struct

### Multi-Command Mock

The `MultiCmdMock` function creates a mock that can handle multiple command configurations:

```go
func MultiCmdMock(configs ...CommandConfig) cdsexec.CommandConstructor
```

The `CommandConfig` struct allows you to specify:

- `Name`: The name of the command
- `Args`: The arguments for the command
- `Stdout`: The simulated standard output
- `Stderr`: The simulated standard error
- `Err`: Any error that should be returned

When an unmatched command is executed, the mock returns `ErrNoMatchingCommand`.

## Example: Using Mock in a Service

Here's an example of how to use the multi-command mock in a service that depends on command execution:

```go
type MyService struct {
    commandContext cdsexec.CommandConstructor
}

func NewMyService(commandContext cdsexec.CommandConstructor) *MyService {
    return &MyService{commandContext: commandContext}
}

func (s *MyService) ListFiles(ctx context.Context) (string, error) {
    cmd := s.commandContext(ctx, "ls", "-l")
    output, err := cmd.Output()
    if err != nil {
        return "", err
    }
    return string(output), nil
}

// In your tests:
func TestMyService(t *testing.T) {
    configs := []mockcmd.CommandConfig{
        {
            Name:   "ls",
            Args:   []string{"-l"},
            Stdout: []byte("file1\nfile2\n"),
        },
    }
    mockCommandContext := mockcmd.MultiCmdMock(configs...)
    service := NewMyService(mockCommandContext)
    
    result, err := service.ListFiles(context.Background())
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }
    if result != "file1\nfile2\n" {
        t.Errorf("Unexpected result: %s", result)
    }
}

// In your main application:
func main() {
    service := NewMyService(cdsexec.CommandContext)
    // Use the service...
}
```

## Error Handling

The `MultiCmdMock` returns `ErrNoMatchingCommand` when an unmatched command is executed:

```go
var ErrNoMatchingCommand = errors.New("no matching command found in this mock")
```

You can check for this error in your tests to verify that an unexpected command was not executed.