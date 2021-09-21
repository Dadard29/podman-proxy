package podman

import (
	"context"

	"github.com/containers/libpod/libpod"
)

type PodmanRuntime struct {
	runtime *libpod.Runtime
	stop    context.CancelFunc
}

func NewPodmanRuntime() (*PodmanRuntime, error) {
	ctx, stop := context.WithCancel(context.Background())
	r, err := libpod.NewRuntime(ctx)
	if err != nil {
		stop()
		return nil, err
	}

	return &PodmanRuntime{
		runtime: r,
		stop:    stop,
	}, nil
}
