package docker

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/dshills/wiggle/container"
)

// Compile check
var _ container.Container = (*Docker)(nil)

type Docker struct {
	options container.Options
}

func (d *Docker) PullImage(image string) error {
	cmd := exec.Command("docker", "pull", image)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to pull image: %s\n%s", err, output)
	}
	return nil
}

func (d *Docker) StartContainer(image string, options container.Options) (container.Instance, error) {
	d.options = options
	args := []string{"run", "-d", "--rm"}

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

	args = append(args, image)          // Docker image
	args = append(args, options.Cmd...) // Command to run inside the container

	cmd := exec.Command("docker", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return container.Instance{}, fmt.Errorf("failed to start container: %s\n%s", err, output)
	}

	containerID := strings.TrimSpace(string(output))
	return container.Instance{ID: containerID}, nil
}

func (d *Docker) StopContainer(containerID string) error {
	cmd := exec.Command("docker", "stop", containerID)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to stop container: %s\n%s", err, output)
	}
	return nil
}

func (d *Docker) ExecCommand(containerID string, command []string) (string, error) {
	args := append([]string{"exec", containerID}, command...)
	cmd := exec.Command("docker", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to exec command: %s\n%s", err, output)
	}
	return string(output), nil
}

func (d *Docker) RemoveContainer(containerID string) error {
	cmd := exec.Command("docker", "rm", "-f", containerID)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to remove container: %s\n%s", err, output)
	}
	return nil
}

func (d *Docker) CloneRepo(containerID, repoURL, targetDir string) error {
	cmd := exec.Command("docker", "exec", containerID, "git", "clone", repoURL, targetDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to clone repository: %s\n%s", err, output)
	}
	return nil
}

func sanitizeLocalPathStrict(localPath string, safePaths ...string) (string, error) {
	cleanedPath, err := filepath.Abs(filepath.Clean(localPath))
	if err != nil {
		return "", fmt.Errorf("invalid local path: %v", err)
	}

	for _, baseDir := range safePaths {
		if isSubdirectory(baseDir, cleanedPath) {
			if err := checkFilePermissions(cleanedPath); err == nil {
				return cleanedPath, nil
			}
		}
	}

	return "", fmt.Errorf("path must be within a defined safe directory or sub-directory: %v", safePaths)
}

// isSubdirectory checks if targetPath is equal to or a subdirectory of basePath.
func isSubdirectory(basePath, targetPath string) bool {
	// Clean and get absolute versions of the paths
	absBasePath, err := filepath.Abs(filepath.Clean(basePath))
	if err != nil {
		return false
	}

	absTargetPath, err := filepath.Abs(filepath.Clean(targetPath))
	if err != nil {
		return false
	}

	// Ensure that basePath ends with a trailing slash for proper prefix matching
	if !strings.HasSuffix(absBasePath, string(filepath.Separator)) {
		absBasePath += string(filepath.Separator)
	}

	// Check if the target path is the same as or a subdirectory of the base path
	return strings.HasPrefix(absTargetPath, absBasePath)
}

func checkFilePermissions(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to stat path: %v", err)
	}

	if info.Mode()&0444 == 0 { // Check read permissions
		return fmt.Errorf("insufficient read permissions for path: %s", path)
	}

	return nil
}

// Validate container ID (Docker's container IDs are alphanumeric)
func sanitizeContainerID(containerID string) (string, error) {
	validID := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if !validID.MatchString(containerID) {
		return "", errors.New("invalid container ID")
	}
	return containerID, nil
}

// Ensure container path doesn't contain special characters
func sanitizeContainerPath(containerPath string) (string, error) {
	// Check for any potentially harmful special characters
	if matched, _ := regexp.MatchString(`[;|&]`, containerPath); matched {
		return "", errors.New("invalid container path: contains special characters")
	}
	return containerPath, nil
}

func (d *Docker) InjectLocalCode(containerID, localPath, containerPath string) error {
	// Sanitize inputs
	cleanLocalPath, err := sanitizeLocalPathStrict(localPath, d.options.SafeLocalPaths...)
	if err != nil {
		return fmt.Errorf("failed to sanitize local path: %v", err)
	}

	cleanContainerID, err := sanitizeContainerID(containerID)
	if err != nil {
		return fmt.Errorf("failed to sanitize container ID: %v", err)
	}

	cleanContainerPath, err := sanitizeContainerPath(containerPath)
	if err != nil {
		return fmt.Errorf("failed to sanitize container path: %v", err)
	}

	// Use subshell with sanitized input
	// nolint - this will be flaged us unsafe because of the localPath.
	// the sanitize function requires that it is in the listed SafeLocalPaths
	cmd := exec.Command("sh", "-c", fmt.Sprintf("docker cp %s %s:%s", cleanLocalPath, cleanContainerID, cleanContainerPath))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to inject local code: %s\n%s", err, output)
	}

	return nil
}
