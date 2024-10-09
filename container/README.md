# Container Package

The `container` package provides an abstraction layer for interacting with containerized environments in a generic way. It supports common container operations like pulling images, starting and stopping containers, injecting local code, cloning repositories, running commands, and waiting for container execution to complete. The package is designed to work with various container backends (e.g., Docker) by implementing the `Container` and `Instance` interfaces.

## Features

- **Container Lifecycle Management**: Pull images, start, stop, and remove containers.
- **Running Commands**: Execute commands inside containers and retrieve their output.
- **Local Code Injection**: Safely inject local code into containers by copying or mounting directories.
- **Repository Cloning**: Clone git repositories inside containers for testing or development purposes.
- **Wait for Containers**: Wait for containers to finish their tasks and retrieve their exit codes.
- **Flexible Options**: Configure container options, such as environment variables, mounts, and working directories.

Usage

Defining a Container Interface

The container.Container interface defines the basic operations for managing containers:

```go
type Container interface {
    PullImage(image string) error
    StartContainer(image string, options Options) (Instance, error)
    StopContainer(containerID string) error
    RemoveContainer(containerID string) error
}
```

- PullImage: Pull a container image from a remote registry.
- StartContainer: Start a container with specified options, returning an Instance of the running container.
- StopContainer: Stop a running container.
- RemoveContainer: Remove a container from the system.

Working with a Running Container

The Instance interface represents a running container and allows you to perform operations like executing commands, injecting local code, cloning repositories, and waiting for the container to exit:

```go
type Instance interface {
    ContainerID() string
    Options() Options
    ExecCommand(command []string) (string, error)
    InjectLocalCode(localPath, containerPath string) error
    CloneRepo(repoURL, targetDir string) error
    Wait() (int, error)
}
```

- ExecCommand: Run commands inside the container and retrieve their output.
- InjectLocalCode: Inject local code into the container by copying or mounting directories.
- CloneRepo: Clone a git repository into the container.
- Wait: Block until the container exits and return its exit code.

Example

Here’s an example of how you might use the container package to start a container, inject code, run tests, and clean up the container:

```go
package main

import (
    "fmt"
    "log"
    "container"
)

func main() {
    d := container.NewDocker() // Assumes an implementation for Docker exists

    // Pull the Go image
    err := d.PullImage("golang:1.20")
    if err != nil {
        log.Fatalf("Failed to pull Go image: %v", err)
    }

    // Set up container options
    opts := container.Options{
        Mounts: []container.Mount{
            {Source: "/path/to/your/local/go/project", Target: "/go/src/app"},
        },
        WorkingDir: "/go/src/app",
        Name:       "my-go-test-container",
        KeepAlive:  false, // Automatically stop after tests
    }

    // Start the container
    containerInstance, err := d.StartContainer("golang:1.20", opts)
    if err != nil {
        log.Fatalf("Failed to start container: %v", err)
    }

    // Inject local code into the container (if needed)
    err = containerInstance.InjectLocalCode("/path/to/your/local/go/project", "/go/src/app")
    if err != nil {
        log.Fatalf("Failed to inject local code: %v", err)
    }

    // Run tests inside the container
    output, err := containerInstance.ExecCommand([]string{"go", "test", "./..."})
    if err != nil {
        log.Fatalf("Failed to run tests: %v", err)
    }
    fmt.Printf("Test Output:\n%s\n", output)

    // Wait for the container to finish
    exitCode, err := containerInstance.Wait()
    if err != nil {
        log.Fatalf("Failed to wait for container: %v", err)
    }
    fmt.Printf("Container exit code: %d\n", exitCode)

    // Stop and remove the container
    err = d.StopContainer(containerInstance.ContainerID())
    if err != nil {
        log.Fatalf("Failed to stop container: %v", err)
    }

    err = d.RemoveContainer(containerInstance.ContainerID())
    if err != nil {
        log.Fatalf("Failed to remove container: %v", err)
    }
}
```

Options Configuration

The Options struct allows you to specify various settings for starting containers:

```go
type Options struct {
    Mounts         []Mount             // Defines volume or bind mounts
    Environment    map[string]string   // Environment variables to set inside the container
    WorkingDir     string              // Working directory inside the container
    SafeLocalPaths []string            // Defines safe local directories for code injection
    Name           string              // Optional name for the container
    KeepAlive      bool                // Determines whether the container should stay running after execution
}
```

Mounting Volumes

Use the Mount struct to define source and target paths for mounting directories:

```go
type Mount struct {
    Source string // Local directory to mount
    Target string // Path inside the container where the directory will be mounted
}
```

Waiting for Containers

You can wait for a container to finish execution by using the Wait() method, which returns the exit code:

```go
exitCode, err := instance.Wait()
```
This ensures that you can block execution until the container task is complete and retrieve the final exit status.
