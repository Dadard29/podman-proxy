package podman_test

import (
	"fmt"
	"testing"

	"github.com/Dadard29/podman-proxy/podman"
)

func TestContainers(t *testing.T) {
	r, err := podman.NewPodmanRuntime()
	if err != nil {
		t.Error(err)
	}

	// testing list of containers
	containers, err := r.ListContainers()
	if err != nil {
		t.Error(err)
	}

	for _, c := range containers {
		fmt.Println(*c)
	}

	if len(containers) > 0 {
		c := containers[0]

		// testing get container from name
		containerByName, err := r.GetContainerFromID(c.Id)
		if err != nil {
			t.Error(err)
		}
		fmt.Println(*containerByName)

		// testing get container from id
		containerById, err := r.GetContainerFromName(c.Name)
		if err != nil {
			t.Error(err)
		}

		if *containerByName != *containerById {
			t.Errorf("mismatch between containers: %v != %v", *containerByName, *containerById)
		}
	}
}

func TestContainersError(t *testing.T) {
	r, err := podman.NewPodmanRuntime()
	if err != nil {
		t.Error(err)
	}

	// get container with invalid id
	c, err := r.GetContainerFromID("")
	if err == nil {
		t.Errorf("expected error while getting container: %v", *c)
	}

	// get container with invalid name
	c, err = r.GetContainerFromName("")
	if err == nil {
		t.Errorf("expected error while getting container: %v", *c)
	}
}
