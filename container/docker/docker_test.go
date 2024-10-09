//go:build docker
// +build docker

package docker_test

import (
	"fmt"
	"testing"

	"github.com/dshills/wiggle/container"
	"github.com/dshills/wiggle/container/docker"
)

func TestPullImage(t *testing.T) {
	image := "golang:latest"
	dock := docker.NewDocker()
	if err := dock.PullImage(image); err != nil {
		t.Error(err)
	}
}

func TestStartContainer(t *testing.T) {
	image := "golang:latest"
	options := container.Options{
		Mounts: []container.Mount{
			{Source: "/Users/dshills/Development/projects/goschema", Target: "/go/src/app"},
		},
		WorkingDir:     "/go/src/app",
		SafeLocalPaths: []string{"/Users/dshills/Development/projects"},
		Name:           "TestStartContainer",
		KeepAlive:      true,
	}
	dock := docker.NewDocker()
	inst, err := dock.StartContainer(image, options)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := dock.StopContainer(inst.ContainerID()); err != nil {
			t.Fatal(err)
		}
	}()
	output, err := inst.ExecCommand([]string{"go", "test", "./..."})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(output)
}
