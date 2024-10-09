package docker

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func sanitizeLocalPathStrict(localPath string, safePaths ...string) (string, error) {
	cleanedPath, err := filepath.Abs(filepath.Clean(localPath))
	if err != nil {
		return "", fmt.Errorf("invalid local path: %v", err)
	}

	info, err := os.Stat(cleanedPath)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("local path not found")
	}
	if !info.IsDir() {
		return "", fmt.Errorf("local path must be a directory")
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
	if absBasePath == absTargetPath {
		return true
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
