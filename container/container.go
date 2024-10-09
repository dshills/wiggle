package container

/*
Package container provides an interface for interacting with containers in a generic way.
It abstracts common container operations such as pulling images, starting and stopping containers,
injecting local code, cloning repositories, running commands inside containers, and waiting for their completion.

This package is designed to work with various container backends (like Docker) by implementing the
Container and Instance interfaces. It also supports the use of configurable options for controlling
mounts, environment variables, and other container settings, making it suitable for automated testing,
code sandboxing, and CI/CD pipelines.
*/

// Container defines an interface for managing container lifecycle operations such as pulling images,
// starting, stopping, and removing containers.
type Container interface {
	// PullImage pulls the specified container image from a remote registry.
	// The image argument specifies the image name and tag (e.g., "golang:1.20").
	PullImage(image string) error

	// StartContainer starts a new container based on the specified image and options.
	// It returns an Instance representing the running container, or an error if the operation fails.
	StartContainer(image string, options Options) (Instance, error)

	// StopContainer stops a running container identified by containerID.
	// This method does not remove the container from the system.
	StopContainer(containerID string) error

	// RemoveContainer removes the container identified by containerID from the system.
	// The container must be stopped before it can be removed.
	RemoveContainer(containerID string) error
}

// Instance defines an interface representing a running container, allowing operations such as
// executing commands, injecting local code, cloning repositories, and waiting for the container to exit.
type Instance interface {
	// ContainerID returns the unique identifier of the running container.
	ContainerID() string

	// Options returns the configuration options with which the container was started.
	Options() Options

	// ExecCommand runs the specified command inside the container and returns the output or an error.
	// The command argument is a slice of strings, with the first element as the command and the rest as its arguments.
	ExecCommand(command []string) (string, error)

	// InjectLocalCode copies or mounts a local path into the container at the specified container path.
	// The localPath must be a valid directory within SafeLocalPaths defined in the Options.
	InjectLocalCode(localPath, containerPath string) error

	// CloneRepo clones a git repository from repoURL into the specified target directory inside the container.
	CloneRepo(repoURL, targetDir string) error

	// Wait blocks until the container exits and returns its exit code and an error, if any occurred.
	Wait() (int, error)
}

// Options holds the configuration options for starting a container, such as volume mounts,
// environment variables, working directory, and security constraints for injecting local code.
type Options struct {
	// Mounts defines a list of volume or bind mounts to be attached to the container.
	Mounts []Mount

	// Environment defines environment variables to be set inside the container.
	Environment map[string]string

	// WorkingDir specifies the working directory inside the container.
	WorkingDir string

	// SafeLocalPaths defines a list of safe local directories from which local code can be injected into the container.
	// Only directories within SafeLocalPaths can be mounted or copied into the container.
	SafeLocalPaths []string

	// Name assigns a specific name to the container for easier identification.
	Name string

	// KeepAlive determines whether the container should remain running after its primary task completes.
	KeepAlive bool
}

// Mount defines a source-target pair for mounting local directories or volumes inside the container.
type Mount struct {
	// Source is the path to the directory on the local host.
	Source string

	// Target is the path where the directory should be mounted inside the container.
	Target string
}
