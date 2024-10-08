package container

type Container interface {
	PullImage(image string) error
	StartContainer(image string, options Options) (Instance, error)
	StopContainer(containerID string) error
	ExecCommand(containerID string, command []string) (string, error) // Runs a command inside the container
	RemoveContainer(containerID string) error
	InjectLocalCode(containerID, localPath, containerPath string) error // Injects local code by copying or mounting
	CloneRepo(containerID, repoURL, targetDir string) error             // Clones a git repository inside the container
}

// Options would hold any configurable options like mounts, environment variables, etc.
type Options struct {
	Mounts      []Mount
	Environment map[string]string
	WorkingDir  string
	Cmd         []string
	// InjectLocalCode localPath must be in SafeLocalPaths or a sub directory
	SafeLocalPaths []string // Required for local code injection.
}

// Mount defines the source and target paths for a volume or bind mount.
type Mount struct {
	Source string
	Target string
}

// ContainerInstance represents a running container.
type Instance struct {
	ID string
}
