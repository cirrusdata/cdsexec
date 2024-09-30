# cdsexec

`cdsexec` is a Go package that provides a flexible abstraction over the standard `os/exec` package, allowing for easier testing and mocking of command execution.

## Features

- Abstracts `exec.Cmd` functionality through interfaces
- Provides a real implementation that wraps `exec.Cmd`
- Includes a mock implementation for testing in a separate `mockcmd` subpackage
- Supports context-based command creation
- Offers utility functions for creating mock commands with various behaviors

## Limitations
- Only a subset of `exec.Cmd` methods are currently supported in the mock implementation. 
- only CommandContext is supported. Command without context is not supported.
- Streaming output is not supported in the mock implementation. You can only get the output after the command has finished.

If you need additional functionality, you can build other mock implementations that satisfy the `Commander` interface.

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

To use the mock implementation in your tests, import the `mockcmd` subpackage:

```go
import (
    "context"
    "testing"

    "github.com/cirrusdata/cdsexec"
    "github.com/cirrusdata/cdsexec/mockcmd"
)

func TestSomeFunction(t *testing.T) {
    // Create a mock command with fixed output
    mockCommandContext := mockcmd.MakeMockCmdWithOutput("mocked output", func(cmd *mockcmd.MockCmd) error {
        if cmd.Name != "expected-command" {
            return fmt.Errorf("unexpected command: %s", cmd.Name)
        }
        return nil
    })

    // Use the mock constructor in your code
    cmd := mockCommandContext(context.Background(), "expected-command", "-arg1", "-arg2")
    
    output, err := cmd.Output()
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }
    if string(output) != "mocked output" {
        t.Errorf("Unexpected output: %s", string(output))
    }

    // You can also assert on the MockCmd's properties if needed
    mockCmd := cmd.(*mockcmd.MockCmd)
    if !mockCmd.startCalled {
        t.Error("Start was not called")
    }
    if !mockCmd.waitCalled {
        t.Error("Wait was not called")
    }
}
```

## Mock Features

The `mockcmd` subpackage provides several utility functions for creating mock commands:

- `MakeMockCmdWithOutput`: Creates a mock command with fixed output
- `MakeMockCmdWithOutputGenericError`: Creates a mock command that returns a generic error
- `MakeMockCmdWithOutputSpecificError`: Creates a mock command with fixed output and a specific error
- `MakeMockCmd`: Creates a mock command from a pre-configured `MockCmd` struct

The `MockCmd` struct provides several features to help with testing:

- `Stdout` and `Stderr`: Set these to provide predefined output
- `Err`: Set this to simulate command errors
- `CheckFunc`: Use this to verify if the command was constructed correctly
- `startCalled` and `waitCalled`: These flags track whether `Start()` and `Wait()` were called

## Example: Using Mock in a Service

Here's an example of how to use the mock in a service that depends on command execution:

```go
type MyService struct {
    commandContext cdsexec.CommandConstructor
}

func NewMyService(commandContext cdsexec.CommandConstructor) *MyService {
    return &MyService{commandContext: commandContext}
}

func (s *MyService) DoSomething(ctx context.Context) (string, error) {
    cmd := s.commandContext(ctx, "some-command", "-arg1", "-arg2")
    output, err := cmd.Output()
    if err != nil {
        return "", err
    }
    return string(output), nil
}

// In your tests:
func TestMyService(t *testing.T) {
    mockCommandContext := mockcmd.MakeMockCmdWithOutput("expected output", nil)
    service := NewMyService(mockCommandContext)
    
    result, err := service.DoSomething(context.Background())
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }
    if result != "expected output" {
        t.Errorf("Unexpected result: %s", result)
    }
}

// In your main application:
func main() {
    service := NewMyService(cdsexec.CommandContext)
    // Use the service...
}
```