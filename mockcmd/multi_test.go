package mockcmd_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/cirrusdata/cdsexec"
	"github.com/cirrusdata/cdsexec/mockcmd"
)

func TestMultiCmdMockCmd(t *testing.T) {
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
		{
			Name: "rm",
			Args: []string{"file2"},
			Err:  errors.New("permission denied"),
		},
	}

	cmdConstructor := mockcmd.MultiCmdMock(configs...)

	tests := []struct {
		name           string
		cmdName        string
		args           []string
		expectedOutput string
		expectedErr    error
	}{
		{"LS Command", "ls", []string{"-l"}, "file1\nfile2\n", nil},
		{"Cat Command", "cat", []string{"file1"}, "contents of file1", nil},
		{"RM Command (Error)", "rm", []string{"file2"}, "", errors.New("permission denied")},
		{"Unmatched Command", "unknown", []string{}, "", mockcmd.ErrNoMatchingCommand},
		{"Partial Match", "ls", []string{"-l", "-a"}, "", mockcmd.ErrNoMatchingCommand},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := cmdConstructor(context.Background(), tt.cmdName, tt.args...)
			output, err := cmd.Output()

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
				} else {
					// No need to check output if an error is expected
					return
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if string(output) != tt.expectedOutput {
				t.Errorf("Expected output %q, got %q", tt.expectedOutput, string(output))
			}

			// Test String() method
			multiCmd, ok := cmd.(*mockcmd.MultiCmdMockCmd)
			if !ok {
				t.Fatalf("Expected *mockcmd.MultiCmdMockCmd, got %T", cmd)
			}
			stringOutput := multiCmd.String()
			if tt.expectedErr == nil {
				if !strings.Contains(stringOutput, "Matched command") {
					t.Errorf("Expected String() to contain 'Matched command', got %q", stringOutput)
				}
			} else {
				if !strings.Contains(stringOutput, "No matching command found") {
					t.Errorf("Expected String() to contain 'No matching command found', got %q", stringOutput)
				}
			}
		})
	}
}

// MockService represents a service that uses command execution
type MockService struct {
	commandContext cdsexec.CommandConstructor
}

func NewMockService(commandContext cdsexec.CommandConstructor) *MockService {
	return &MockService{commandContext: commandContext}
}

func (s *MockService) ListFiles(ctx context.Context) (string, error) {
	cmd := s.commandContext(ctx, "ls", "-l")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func (s *MockService) ReadFile(ctx context.Context, filename string) (string, error) {
	cmd := s.commandContext(ctx, "cat", filename)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func TestMockServiceWithMultiCmdMock(t *testing.T) {
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
	service := NewMockService(mockCommandContext)

	t.Run("ListFiles", func(t *testing.T) {
		result, err := service.ListFiles(context.Background())
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if result != "file1\nfile2\n" {
			t.Errorf("Unexpected result: %s", result)
		}
	})

	t.Run("ReadFile", func(t *testing.T) {
		result, err := service.ReadFile(context.Background(), "file1")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if result != "contents of file1" {
			t.Errorf("Unexpected result: %s", result)
		}
	})

	t.Run("ReadNonexistentFile", func(t *testing.T) {
		_, err := service.ReadFile(context.Background(), "nonexistent")
		if !errors.Is(err, mockcmd.ErrNoMatchingCommand) {
			t.Fatalf("Expected ErrNoMatchingCommand, got: %v", err)
		}
	})
}

func TestMultiCmdMockCmdCombinedOutput(t *testing.T) {
	configs := []mockcmd.CommandConfig{
		{
			Name:   "echo",
			Args:   []string{"hello"},
			Stdout: []byte("hello"),
			Stderr: []byte("warning: echo"),
		},
	}

	mockCommandContext := mockcmd.MultiCmdMock(configs...)

	cmd := mockCommandContext(context.Background(), "echo", "hello")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedOutput := "hellowarning: echo"
	if string(output) != expectedOutput {
		t.Errorf("Expected combined output %q, got %q", expectedOutput, string(output))
	}
}
