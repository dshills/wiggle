package container

type Container interface {
	PullImage(image string) error
	StartContainer(image string, options Options) (Instance, error)
	StopContainer(containerID string) error
	RemoveContainer(containerID string) error
}

type Instance interface {
	ContainerID() string
	Options() Options
	ExecCommand(command []string) (string, error)          // Runs a command inside the container
	InjectLocalCode(localPath, containerPath string) error // Injects local code by copying or mounting
	CloneRepo(repoURL, targetDir string) error             // Clones a git repository inside the container
	Wait() (int, error)                                    // Wait for container to exit
}

// Options would hold any configurable options like mounts, environment variables, etc.
type Options struct {
	Mounts      []Mount
	Environment map[string]string
	WorkingDir  string
	// InjectLocalCode localPath must be in SafeLocalPaths or a sub directory
	SafeLocalPaths []string // Required for local code injection.
	Name           string   // Container name
	KeepAlive      bool     // Keep the container running
}

// Mount defines the source and target paths for a volume or bind mount.
type Mount struct {
	Source string
	Target string
}
