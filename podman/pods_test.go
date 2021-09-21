package podman_test

import (
	"fmt"
	"testing"

	"github.com/Dadard29/podman-proxy/podman"
)

func TestPods(t *testing.T) {
	r, err := podman.NewPodmanRuntime()
	if err != nil {
		t.Error(err)
	}

	// testing list of pods
	pods, err := r.ListPods()
	if err != nil {
		t.Error(err)
	}

	for _, pod := range pods {
		fmt.Println(*pod)
	}

	if len(pods) > 0 {
		p := pods[0]

		// testing get pod by name
		podByName, err := r.GetPodFromName(p.Name)
		if err != nil {
			t.Error(err)
		}

		// testing get pod by id
		podById, err := r.GetPodFromID(p.Id)
		if err != nil {
			t.Error(err)
		}

		if podByName.Id != podById.Id {
			t.Errorf("mismatch between pods: %v != %v", *podById, *podByName)
		}
	}
}

func TestPodsError(t *testing.T) {
	r, err := podman.NewPodmanRuntime()
	if err != nil {
		t.Error(err)
	}

	// get pod with invalid id
	c, err := r.GetPodFromID("")
	if err == nil {
		t.Errorf("expected error while getting pod: %v", *c)
	}

	// get pod with invalid name
	c, err = r.GetPodFromName("")
	if err == nil {
		t.Errorf("expected error while getting pod: %v", *c)
	}
}
