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

	containers, err := r.ListContainers()
	if err != nil {
		t.Error(err)
	}

	for _, c := range containers {
		fmt.Println(c.String())
	}

}
