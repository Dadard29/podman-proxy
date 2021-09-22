package models

import (
	"time"

	"github.com/containers/libpod/libpod"
)

type PodmanPod struct {
	Id             string       `json:"id"`
	Name           string       `json:"name"`
	CreatedAt      time.Time    `json:"created_at"`
	InfraContainer *Container   `json:"infra_container"`
	Containers     []*Container `json:"containers"`
}

func NewPodmanPod(pod *libpod.Pod, infra *Container) (*PodmanPod, error) {

	containers, err := pod.AllContainers()
	if err != nil {
		return nil, err
	}

	Containers := make([]*Container, 0)
	for _, container := range containers {
		Container, err := NewContainer(container)
		if err != nil {
			return nil, err
		}
		Containers = append(Containers, Container)
	}

	return &PodmanPod{
		Id:             pod.ID(),
		Name:           pod.Name(),
		CreatedAt:      pod.CreatedTime(),
		InfraContainer: infra,
		Containers:     Containers,
	}, nil
}
