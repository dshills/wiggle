package docker

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/dshills/wiggle/container"
)

// Compile check
var _ container.Container = (*Docker)(nil)

type Docker struct {
}

func NewDocker() *Docker {
	return &Docker{}
}

func (d *Docker) PullImage(image string) error {
	cmd := exec.Command("docker", "pull", image)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to pull image: %s\n%s", err, output)
	}
	return nil
}

func (d *Docker) StartContainer(image string, options container.Options) (container.Instance, error) {
	args := []string{"run", "-d", "--rm"}

	// Name the container
	if options.Name != "" {
		args = append(args, "--name", options.Name)
	}

	// Add volume mounts
	for _, mount := range options.Mounts {
		args = append(args, "-v", fmt.Sprintf("%s:%s", mount.Source, mount.Target))
	}

	// Set working directory
	if options.WorkingDir != "" {
		args = append(args, "-w", options.WorkingDir)
	}

	// Add environment variables
	for key, value := range options.Environment {
		args = append(args, "-e", fmt.Sprintf("%s=%s", key, value))
	}

	args = append(args, image) // Docker image

	// Keep the container running
	if options.KeepAlive {
		args = append(args, []string{"tail", "-f", "/dev/null"}...)
	}

	fmt.Printf("%+v\n", args)

	cmd := exec.Command("docker", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %s\n%s", err, output)
	}

	containerID := strings.TrimSpace(string(output))
	return NewInstance(containerID, options), nil
}

func (d *Docker) StopContainer(containerID string) error {
	cmd := exec.Command("docker", "stop", containerID)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to stop container: %s\n%s", err, output)
	}
	return nil
}

func (d *Docker) RemoveContainer(containerID string) error {
	cmd := exec.Command("docker", "rm", "-f", containerID)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to remove container: %s\n%s", err, output)
	}
	return nil
}
