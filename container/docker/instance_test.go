package docker_test

import (
	"bytes"
	"fmt"
	"os/exec"
	"testing"

	"github.com/dshills/wiggle/container"
	"github.com/dshills/wiggle/container/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCmd structure to simulate the exec.Cmd behavior
type MockCmd struct {
	mock.Mock
	Stdout bytes.Buffer
	Stderr bytes.Buffer
}

// CombinedOutput method mock for exec.Cmd
func (m *MockCmd) CombinedOutput() ([]byte, error) {
	args := m.Called()
	return args.Get(0).([]byte), args.Error(1)
}

// Helper to create an exec.Cmd-like structure
func fakeExecCommand(command string, args ...string) *exec.Cmd {
	// Create a real exec.Cmd, but override its behavior
	cmd := exec.Command("echo") // Just use a dummy command

	// We will replace its stdout/stderr and the Run method.
	mockCmd := &MockCmd{}

	// Success case for docker cp
	if command == "docker" && args[0] == "cp" {
		mockCmd.On("CombinedOutput").Return([]byte("Success"), nil)
	} else {
		mockCmd.On("CombinedOutput").Return([]byte{}, fmt.Errorf("docker cp failed"))
	}
	return cmd
}

// Test case for successful injection
func TestInjectLocalCode_Success(t *testing.T) {
	codePath := "/tmp"
	options := container.Options{
		SafeLocalPaths: []string{codePath},
	}
	inst := docker.NewInstance("TEST", options)
	inst.RunCommand = fakeExecCommand

	err := inst.InjectLocalCode(codePath, "/container/path")

	assert.NoError(t, err)
}

// Test case for injection failure
func TestInjectLocalCode_Failure(t *testing.T) {
	codePath := "/fake/local/path"
	options := container.Options{
		SafeLocalPaths: []string{codePath},
	}
	inst := docker.NewInstance("TEST", options)
	inst.RunCommand = fakeExecCommand

	err := inst.InjectLocalCode(codePath, "/container/path")

	assert.Error(t, err)
}
