package docker

import (
	"fmt"
	"os/exec"

	"github.com/dshills/wiggle/container"
)

// Compile check
var _ container.Instance = (*Instance)(nil)

// Define a type for the function that runs commands
type CommandRunner func(name string, arg ...string) *exec.Cmd

type Instance struct {
	RunCommand  CommandRunner
	containerID string
	options     container.Options
}

func NewInstance(id string, options container.Options) *Instance {
	return &Instance{containerID: id, options: options, RunCommand: exec.Command}
}

func (in *Instance) ContainerID() string {
	return in.containerID
}

func (in *Instance) Options() container.Options {
	return in.options
}

func (in *Instance) Wait() (int, error) {
	// Run the `docker wait` command, which blocks until the container exits
	// nolint
	cmd := exec.Command("docker", "wait", in.containerID)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return -1, fmt.Errorf("failed to wait for container: %v\n%s", err, output)
	}

	// Parse the exit code
	var exitCode int
	_, err = fmt.Sscanf(string(output), "%d", &exitCode)
	if err != nil {
		return -1, fmt.Errorf("failed to parse exit code: %v", err)
	}

	return exitCode, nil
}

func (in *Instance) ExecCommand(command []string) (string, error) {
	args := append([]string{"exec", in.containerID}, command...)
	cmd := exec.Command("docker", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to exec command: %s\n%s", err, output)
	}
	return string(output), nil
}

func (in *Instance) InjectLocalCode(localPath, containerPath string) error {
	// Sanitize inputs

	cleanLocalPath, err := sanitizeLocalPathStrict(localPath, in.options.SafeLocalPaths...)
	if err != nil {
		return fmt.Errorf("failed to sanitize local path: %v", err)
	}

	cleanContainerID, err := sanitizeContainerID(in.containerID)
	if err != nil {
		return fmt.Errorf("failed to sanitize container ID: %v", err)
	}

	cleanContainerPath, err := sanitizeContainerPath(containerPath)
	if err != nil {
		return fmt.Errorf("failed to sanitize container path: %v", err)
	}

	cmd := in.RunCommand("docker", "cp", cleanLocalPath, fmt.Sprintf("%s:%s", cleanContainerID, cleanContainerPath))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to inject local code: %s\n%s", err, output)
	}

	return nil
}

func (in *Instance) CloneRepo(repoURL, targetDir string) error {
	cmd := in.RunCommand("docker", "exec", in.containerID, "git", "clone", repoURL, targetDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to clone repository: %s\n%s", err, output)
	}
	return nil
}
