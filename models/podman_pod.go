package models

import (
	"time"

	"github.com/containers/libpod/libpod"
)

type PodmanPod struct {
	Id             string             `json:"id"`
	Name           string             `json:"name"`
	CreatedAt      time.Time          `json:"created_at"`
	InfraContainer *PodmanContainer   `json:"infra_container"`
	Containers     []*PodmanContainer `json:"containers"`
}

func NewPodmanPod(pod *libpod.Pod, infra *PodmanContainer) (*PodmanPod, error) {

	containers, err := pod.AllContainers()
	if err != nil {
		return nil, err
	}

	podmanContainers := make([]*PodmanContainer, 0)
	for _, container := range containers {
		podmanContainer, err := NewPodmanContainer(container)
		if err != nil {
			return nil, err
		}
		podmanContainers = append(podmanContainers, podmanContainer)
	}

	return &PodmanPod{
		Id:             pod.ID(),
		Name:           pod.Name(),
		CreatedAt:      pod.CreatedTime(),
		InfraContainer: infra,
		Containers:     podmanContainers,
	}, nil
}
